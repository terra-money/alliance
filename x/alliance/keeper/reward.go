package keeper

import (
	"alliance/x/alliance/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"golang.org/x/exp/slices"
)

type RewardsKeeper interface {
	ClaimDelegationRewards(ctx sdk.Context, delAddr sdk.AccAddress, val stakingtypes.Validator, denom string) (sdk.Coins, error)
	CalculateDelegationRewards(ctx sdk.Context, delegation types.Delegation, asset types.AllianceAsset) (sdk.Coins, types.RewardIndices, error)
	AddAssetsToRewardPool(ctx sdk.Context, from sdk.AccAddress, coins sdk.Coins) error
}

var (
	_ RewardsKeeper = Keeper{}
)

func (k Keeper) ClaimDelegationRewards(ctx sdk.Context, delAddr sdk.AccAddress, val stakingtypes.Validator, denom string) (sdk.Coins, error) {
	asset := k.GetAssetByDenom(ctx, denom)
	delegation, found := k.GetDelegation(ctx, delAddr, val, denom)
	if !found {
		return sdk.Coins{}, stakingtypes.ErrNoDelegatorForAddress
	}

	coins, newIndices, err := k.CalculateDelegationRewards(ctx, delegation, asset)
	if err != nil {
		return nil, err
	}

	delegation.RewardIndices = newIndices
	k.SetDelegation(ctx, delAddr, val, denom, delegation)

	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.RewardsPoolName, delAddr, coins)
	if err != nil {
		return nil, err
	}

	return coins, nil
}

func (k Keeper) CalculateDelegationRewards(ctx sdk.Context, delegation types.Delegation, asset types.AllianceAsset) (sdk.Coins, types.RewardIndices, error) {
	// TODO: check if there is a rewards rate change
	var rewards sdk.Coins
	globalIndices := k.GlobalRewardIndices(ctx)
	for _, index := range globalIndices {
		idx := slices.IndexFunc(delegation.RewardIndices, func(r types.RewardIndex) bool {
			return r.Denom == index.Denom
		})
		if idx >= 0 && delegation.RewardIndices[idx].Index == index.Index {
			continue
		}
		var localIndex sdk.Dec
		if idx < 0 {
			localIndex = sdk.ZeroDec()
		} else {
			localIndex = delegation.RewardIndices[idx].Index
		}

		claimWeight := delegation.Shares.MulInt(asset.TotalTokens).Quo(asset.TotalShares).Mul(asset.RewardWeight)
		totalClaimable := (index.Index.Sub(localIndex)).Mul(claimWeight)
		rewards = append(rewards, sdk.NewCoin(index.Denom, totalClaimable.TruncateInt()))
	}
	return rewards, globalIndices, nil
}

func (k Keeper) AddAssetsToRewardPool(ctx sdk.Context, from sdk.AccAddress, coins sdk.Coins) error {
	globalIndices := k.GlobalRewardIndices(ctx)
	totalAssetWeight := k.totalAssetWeight(ctx)
	// We need some delegations before we can split rewards. Else rewards belong to no one
	if totalAssetWeight.IsZero() {
		return types.ErrZeroDelegations
	}

	for _, c := range coins {
		index, found := globalIndices.GetIndexByDenom(c.Denom)
		if !found {
			globalIndices = append(globalIndices, types.RewardIndex{
				Denom: c.Denom,
				Index: sdk.NewDecFromInt(c.Amount).Quo(totalAssetWeight),
			})
		} else {
			index.Index = index.Index.Add(sdk.NewDecFromInt(c.Amount).Quo(totalAssetWeight))
		}
	}
	k.SetGlobalRewardIndex(ctx, globalIndices)

	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, from, types.RewardsPoolName, coins)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) totalAssetWeight(ctx sdk.Context) sdk.Dec {
	assets := k.GetAllAssets(ctx)
	total := sdk.ZeroDec()
	for _, asset := range assets {
		total = total.Add(asset.RewardWeight.MulInt(asset.TotalTokens))
	}
	return total
}
