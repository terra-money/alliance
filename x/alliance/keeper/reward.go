package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/terra-money/alliance/x/alliance/types"
	"time"
)

type RewardsKeeper interface {
}

var (
	_ RewardsKeeper = Keeper{}
)

const (
	// YEAR_IN_NANOS is used to calculate the pro-rated take rate when it is periodically applied with DeductAssetsHook
	YEAR_IN_NANOS int64 = 31_557_000_000_000_000
)

// ClaimValidatorRewards claims the validator rewards (minus commission) from the distribution module
// This should be called everytime validator delegation changes (e.g. [un/re]delegation) to update the reward claim history
func (k Keeper) ClaimValidatorRewards(ctx sdk.Context, val types.AllianceValidator) (sdk.Coins, error) {
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

// ClaimDelegationRewards claims delegation rewards and transfers to the delegator account
// This method updates the delegation so you will need to re-query an updated version from the database
func (k Keeper) ClaimDelegationRewards(ctx sdk.Context, delAddr sdk.AccAddress, val types.AllianceValidator, denom string) (sdk.Coins, error) {
	asset, found := k.GetAssetByDenom(ctx, denom)
	if !found {
		return nil, types.ErrUnknownAsset
	}
	delegation, found := k.GetDelegation(ctx, delAddr, val, denom)
	if !found {
		return sdk.Coins{}, stakingtypes.ErrNoDelegatorForAddress
	}

	_, err := k.ClaimValidatorRewards(ctx, val)
	if err != nil {
		return nil, err
	}

	coins, newIndices, err := k.CalculateDelegationRewards(ctx, delegation, val, asset)
	if err != nil {
		return nil, err
	}

	delegation.RewardHistory = newIndices
	delegation.LastRewardClaimHeight = uint64(ctx.BlockHeight())
	k.SetDelegation(ctx, delAddr, val.GetOperator(), denom, delegation)

	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.RewardsPoolName, delAddr, coins)
	if err != nil {
		return nil, err
	}

	return coins, nil
}

// CalculateDelegationRewards calculates the rewards that can be claimed for a delegation
// It takes past reward_rate changes into account by using the RewardRateChangeSnapshot entry
func (k Keeper) CalculateDelegationRewards(ctx sdk.Context, delegation types.Delegation, val types.AllianceValidator, asset types.AllianceAsset) (sdk.Coins, types.RewardHistories, error) {
	var totalRewards sdk.Coins
	currentRewardHistory := types.NewRewardHistories(val.GlobalRewardHistory)
	delegationRewardHistories := types.NewRewardHistories(delegation.RewardHistory)
	// If there are reward rate changes between last and current claim, sequentially claim with the help of the snapshots
	snapshotIter := k.IterateWeightChangeSnapshot(ctx, asset.Denom, val.GetOperator(), delegation.LastRewardClaimHeight)
	for ; snapshotIter.Valid(); snapshotIter.Next() {
		var snapshot types.RewardWeightChangeSnapshot
		b := snapshotIter.Value()
		k.cdc.MustUnmarshal(b, &snapshot)
		var rewards sdk.Coins
		rewards, delegationRewardHistories = accumulateRewards(types.NewRewardHistories(snapshot.RewardHistories), delegationRewardHistories, asset, snapshot.PrevRewardWeight, delegation, val)
		totalRewards = totalRewards.Add(rewards...)
	}
	rewards, _ := accumulateRewards(currentRewardHistory, delegationRewardHistories, asset, asset.RewardWeight, delegation, val)
	totalRewards = totalRewards.Add(rewards...)
	return totalRewards, currentRewardHistory, nil
}

// accumulateRewards compares the latest reward history with the delegation's reward history
// It takes the difference and calculates how much can be claimed
func accumulateRewards(latestRewardHistories types.RewardHistories, rewardHistories types.RewardHistories, asset types.AllianceAsset, rewardWeight sdk.Dec, delegation types.Delegation, validator types.AllianceValidator) (sdk.Coins, types.RewardHistories) {
	// Go through each reward denom and accumulate rewards
	var rewards sdk.Coins
	for _, history := range latestRewardHistories {
		rewardHistory, found := rewardHistories.GetIndexByDenom(history.Denom)
		if !found {
			rewardHistory.Denom = history.Denom
			rewardHistory.Index = sdk.ZeroDec()
		}
		if rewardHistory.Index.GTE(history.Index) {
			continue
		}
		delegationTokens := sdk.NewDecFromInt(types.GetDelegationTokens(delegation, validator, asset).Amount)

		claimWeight := delegationTokens.Mul(rewardWeight)
		totalClaimable := (history.Index.Sub(rewardHistory.Index)).Mul(claimWeight)
		rewardHistory.Index = history.Index
		rewards = rewards.Add(sdk.NewCoin(history.Denom, totalClaimable.TruncateInt()))
		if !found {
			rewardHistories = append(rewardHistories, *rewardHistory)
		}
	}
	return rewards, rewardHistories
}

// AddAssetsToRewardPool increments a reward history array. A reward history stores the average reward per token/reward_weight.
// To calculate the number of rewards claimable, take reward_history * alliance_token_amount * reward_weight
func (k Keeper) AddAssetsToRewardPool(ctx sdk.Context, from sdk.AccAddress, val types.AllianceValidator, coins sdk.Coins) error {
	rewardHistories := types.NewRewardHistories(val.GlobalRewardHistory)
	totalAssetWeight := k.totalAssetWeight(ctx, val)
	// We need some delegations before we can split rewards. Else rewards belong to no one
	if totalAssetWeight.IsZero() {
		fmt.Println(val.ValidatorShares, totalAssetWeight.String())
		return types.ErrZeroDelegations
	}

	for _, c := range coins {
		rewardHistory, found := rewardHistories.GetIndexByDenom(c.Denom)
		if !found {
			rewardHistories = append(rewardHistories, types.RewardHistory{
				Denom: c.Denom,
				Index: sdk.NewDecFromInt(c.Amount).Quo(totalAssetWeight),
			})
		} else {
			rewardHistory.Index = rewardHistory.Index.Add(sdk.NewDecFromInt(c.Amount).Quo(totalAssetWeight))
		}
	}

	val.GlobalRewardHistory = rewardHistories
	k.SetValidator(ctx, val)
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, from, types.RewardsPoolName, coins)
	if err != nil {
		return err
	}

	return nil
}

// DeductAssetsHook is called periodically to deduct from an alliance asset (calculated by take_rate).
// The interval in which assets are deducted is set in module params
func (k Keeper) DeductAssetsHook(ctx sdk.Context) (sdk.Coins, error) {
	last := k.LastRewardClaimTime(ctx)
	interval := k.RewardClaimInterval(ctx)
	next := last.Add(interval)
	if ctx.BlockTime().After(next) {
		return k.DeductAssetsWithTakeRate(ctx, last)
	}
	return nil, nil
}

// DeductAssetsWithTakeRate Deducts an alliance asset using the take_rate
// The deducted asset is distributed to the fee_collector module account to be redistributed to stakers
func (k Keeper) DeductAssetsWithTakeRate(ctx sdk.Context, lastClaim time.Time) (sdk.Coins, error) {
	assets := k.GetAllAssets(ctx)
	durationSinceLastClaim := ctx.BlockTime().Sub(lastClaim)
	prorate := sdk.NewDec(durationSinceLastClaim.Nanoseconds()).Quo(sdk.NewDec(YEAR_IN_NANOS))

	var coins sdk.Coins
	for _, asset := range assets {
		if asset.TotalTokens.IsPositive() && asset.TakeRate.IsPositive() {
			reward := asset.TakeRate.Mul(prorate).MulInt(asset.TotalTokens).TruncateInt()
			asset.TotalTokens = asset.TotalTokens.Sub(reward)
			coins = coins.Add(sdk.NewCoin(asset.Denom, reward))
			k.SetAsset(ctx, *asset)
		}
	}

	if !coins.Empty() && !coins.IsZero() {
		err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, authtypes.FeeCollectorName, coins)
		if err != nil {
			return nil, err
		}
		// Only update if there was a token transfer to prevent < 1 amounts to be ignored
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
		totalValTokens := val.TotalDecTokensWithAsset(asset)
		total = total.Add(asset.RewardWeight.Mul(totalValTokens))
	}
	return total
}
