package keeper

import (
	"alliance/x/alliance/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"golang.org/x/exp/slices"
	"time"
)

type RewardsKeeper interface {
}

var (
	_ RewardsKeeper = Keeper{}
)

const (
	YEAR_IN_NANOS int64 = 31_557_000_000_000_000
)

// ClaimDistributionRewards to be called right before any reward claims so that we get
// the latest rewards
func (k Keeper) ClaimDistributionRewards(ctx sdk.Context, val types.AllianceValidator) (sdk.Coins, error) {
	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)

	_, found := k.stakingKeeper.GetDelegation(ctx, moduleAddr, val.GetOperator())
	if !found {
		return sdk.NewCoins(), nil
	}

	coins, err := k.distributionKeeper.WithdrawDelegationRewards(ctx, moduleAddr, val.GetOperator())
	if err != nil || coins.IsZero() {
		return nil, err
	}
	err = k.AddAssetsToRewardPool(ctx, moduleAddr, val, coins)
	if err != nil {
		return nil, err
	}
	return coins, nil
}

func (k Keeper) ClaimDelegationRewards(ctx sdk.Context, delAddr sdk.AccAddress, val types.AllianceValidator, denom string) (sdk.Coins, error) {
	asset, found := k.GetAssetByDenom(ctx, denom)
	if !found {
		return nil, types.ErrUnknownAsset
	}
	delegation, found := k.GetDelegation(ctx, delAddr, val, denom)
	if !found {
		return sdk.Coins{}, stakingtypes.ErrNoDelegatorForAddress
	}

	_, err := k.ClaimDistributionRewards(ctx, val)
	if err != nil {
		return nil, err
	}

	coins, newIndices, err := k.CalculateDelegationRewards(ctx, delegation, val, asset)
	if err != nil {
		return nil, err
	}

	delegation.RewardHistory = newIndices
	delegation.LastRewardClaimHeight = uint64(ctx.BlockHeight())
	k.SetDelegation(ctx, delAddr, val, denom, delegation)

	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.RewardsPoolName, delAddr, coins)
	if err != nil {
		return nil, err
	}

	return coins, nil
}

func (k Keeper) CalculateDelegationRewards(ctx sdk.Context, delegation types.Delegation, val types.AllianceValidator, asset types.AllianceAsset) (sdk.Coins, types.RewardHistories, error) {
	var rewards sdk.Coins
	currentRewardHistory := types.NewRewardHistories(val.GlobalRewardHistory)
	delegationRewardHistories := delegation.RewardHistory
	// If there are reward rate changes between last and current claim, claim that first using the snapshots
	snapshotIter := k.IterateRewardRatesChangeSnapshot(ctx, asset.Denom, val.GetOperator(), delegation.LastRewardClaimHeight)
	for ; snapshotIter.Valid(); snapshotIter.Next() {
		var snapshot types.RewardRateChangeSnapshot
		b := snapshotIter.Value()
		k.cdc.MustUnmarshal(b, &snapshot)
		// Go through each reward denom and accumulate rewards
		for _, history := range snapshot.RewardHistories {
			idx := slices.IndexFunc(delegationRewardHistories, func(r types.RewardHistory) bool {
				return r.Denom == history.Denom
			})

			// If local history == global history, it means that user has already claimed
			// Index should never be more than global unless some rewards are withdrawn from the pool
			if idx >= 0 && delegationRewardHistories[idx].Index.GTE(history.Index) {
				continue
			}
			if idx < 0 {
				idx = len(delegationRewardHistories)
				delegationRewardHistories = append(delegationRewardHistories, types.RewardHistory{
					Denom: history.Denom,
					Index: sdk.ZeroDec(),
				})
			}
			delegationTokens := sdk.NewDecFromInt(types.GetDelegationTokens(delegation, val, asset).Amount)

			claimWeight := delegationTokens.Mul(snapshot.PrevRewardWeight)
			totalClaimable := (history.Index.Sub(delegationRewardHistories[idx].Index)).Mul(claimWeight)
			delegationRewardHistories[idx].Index = history.Index
			rewards = rewards.Add(sdk.NewCoin(history.Denom, totalClaimable.TruncateInt()))
		}
	}

	// Go through each reward denom and accumulate rewards
	for _, history := range currentRewardHistory {
		idx := slices.IndexFunc(delegationRewardHistories, func(r types.RewardHistory) bool {
			return r.Denom == history.Denom
		})

		// If local history == global history, it means that user has already claimed
		// Index should never be more than global unless some rewards are withdrawn from the pool
		if idx >= 0 && delegationRewardHistories[idx].Index.GTE(history.Index) {
			continue
		}
		var localRewardHistory sdk.Dec
		if idx < 0 {
			localRewardHistory = sdk.ZeroDec()
		} else {
			localRewardHistory = delegationRewardHistories[idx].Index
		}
		delegationTokens := sdk.NewDecFromInt(types.GetDelegationTokens(delegation, val, asset).Amount)

		claimWeight := delegationTokens.Mul(asset.RewardWeight)
		totalClaimable := (history.Index.Sub(localRewardHistory)).Mul(claimWeight)
		rewards = rewards.Add(sdk.NewCoin(history.Denom, totalClaimable.TruncateInt()))
	}
	return rewards, currentRewardHistory, nil
}

func (k Keeper) AddAssetsToRewardPool(ctx sdk.Context, from sdk.AccAddress, val types.AllianceValidator, coins sdk.Coins) error {
	globalIndices := types.NewRewardHistories(val.GlobalRewardHistory)
	totalAssetWeight := k.totalAssetWeight(ctx, val)
	// We need some delegations before we can split rewards. Else rewards belong to no one
	if totalAssetWeight.IsZero() {
		return types.ErrZeroDelegations
	}

	for _, c := range coins {
		index, found := globalIndices.GetIndexByDenom(c.Denom)
		if !found {
			globalIndices = append(globalIndices, types.RewardHistory{
				Denom: c.Denom,
				Index: sdk.NewDecFromInt(c.Amount).Quo(totalAssetWeight),
			})
		} else {
			index.Index = index.Index.Add(sdk.NewDecFromInt(c.Amount).Quo(totalAssetWeight))
		}
	}

	val.GlobalRewardHistory = globalIndices
	k.SetValidator(ctx, val)

	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, from, types.RewardsPoolName, coins)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) ClaimAssetsWithTakeRateRateLimited(ctx sdk.Context) (sdk.Coins, error) {
	last := k.LastRewardClaimTime(ctx)
	interval := k.RewardClaimInterval(ctx)
	next := last.Add(interval)
	if ctx.BlockTime().After(next) {
		return k.ClaimAssetsWithTakeRate(ctx, last)
	}
	return nil, nil
}

func (k Keeper) ClaimAssetsWithTakeRate(ctx sdk.Context, lastClaim time.Time) (sdk.Coins, error) {
	assets := k.GetAllAssets(ctx)
	durationSinceLastClaim := ctx.BlockTime().Sub(lastClaim)
	prorate := sdk.NewDec(durationSinceLastClaim.Nanoseconds()).Quo(sdk.NewDec(YEAR_IN_NANOS))

	var coins sdk.Coins
	for _, asset := range assets {
		if asset.TotalTokens.IsPositive() && asset.TakeRate.IsPositive() {
			reward := asset.TakeRate.Mul(prorate).MulInt(asset.TotalTokens).TruncateInt()
			asset.TotalTokens = asset.TotalTokens.Sub(reward)
			coins = append(coins, sdk.NewCoin(asset.Denom, reward))
			k.SetAsset(ctx, *asset)
		}
	}

	if !coins.Empty() && !coins.IsZero() {
		err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, authtypes.FeeCollectorName, coins)
		if err != nil {
			return nil, err
		}
		// Only update if there was a token transfer to prevent < 1 amounts to be totally ignored
		// TODO: Look into how to deal with rounding issues if claim interval is too short
		k.SetLastRewardClaimTime(ctx, ctx.BlockTime())
	}
	return coins, nil
}

func (k Keeper) totalAssetWeight(ctx sdk.Context, val types.AllianceValidator) sdk.Dec {
	total := sdk.ZeroDec()
	for _, token := range val.TotalDelegatorShares {
		asset, found := k.GetAssetByDenom(ctx, token.Denom)
		if !found {
			continue
		}
		totalValTokens := val.TotalTokensWithAsset(asset)
		total = total.Add(asset.RewardWeight.MulInt(totalValTokens))
	}
	return total
}
