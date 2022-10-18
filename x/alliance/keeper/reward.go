package keeper

import (
	"alliance/x/alliance/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type RewardsKeeper interface {
	CalculateDelegationRewards(ctx sdk.Context, delAddr sdk.AccAddress, val stakingtypes.Validator, denom string) (sdk.Coins, error)
}

var (
	_ RewardsKeeper = Keeper{}
)

func (k Keeper) CalculateDelegationRewards(ctx sdk.Context, delAddr sdk.AccAddress, val stakingtypes.Validator, denom string) (sdk.Coins, error) {
	delegation, found := k.GetDelegation(ctx, delAddr, val, denom)
	if !found {
		return sdk.Coins{}, stakingtypes.ErrNoDelegatorForAddress
	}
	asset := k.GetAssetByDenom(ctx, denom)

	// TODO check if there is a rewards rate change

	globalIndex := k.GlobalRewardIndex(ctx)

	// If index already == global index then there is nothing to claim
	if globalIndex == delegation.RewardIndex {
		return sdk.Coins{}, nil
	}

	claimWeight := delegation.Shares.MulInt(asset.TotalTokens).Mul(asset.RewardWeight)
	totalClaimable := (globalIndex.Sub(delegation.RewardIndex)).Mul(claimWeight)
	rewards := k.prorataAssetsInRewardPool(ctx, totalClaimable)

	delegation.RewardIndex = globalIndex
	k.SetDelegation(ctx, delAddr, val, denom, delegation)

	return rewards, nil
}

func (k Keeper) AddAssetsToRewardPool(ctx sdk.Context, from sdk.AccAddress, coins sdk.Coins) error {
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, from, types.RewardsPoolName, coins)
	if err != nil {
		return err
	}
	globalIndex := k.GlobalRewardIndex(ctx)
	totalRewards := sdk.ZeroDec()
	for _, c := range coins {
		totalRewards = totalRewards.Add(sdk.NewDecFromInt(c.Amount))
	}
	totalAssetWeight := k.totalAssetWeight(ctx)

	// We need some delegations before we can split rewards. Else rewards belong to no one
	if totalAssetWeight.IsZero() {
		return types.ErrZeroDelegations
	}

	globalIndex = globalIndex.Add(totalRewards.Quo(totalAssetWeight))
	k.SetGlobalRewardIndex(ctx, globalIndex)
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

func (k Keeper) prorataAssetsInRewardPool(ctx sdk.Context, totalClaimable sdk.Dec) sdk.Coins {
	rewardsPool := k.accountKeeper.GetModuleAddress(types.RewardsPoolName)
	coins := k.bankKeeper.GetAllBalances(ctx, rewardsPool)
	dCoins := sdk.NewDecCoinsFromCoins(coins...)
	totalCoins := sdk.NewDec(0)
	for _, c := range dCoins {
		totalCoins = totalCoins.Add(c.Amount)
	}
	weightedCoins := dCoins.MulDec(totalClaimable.Quo(totalCoins))
	rewards, _ := weightedCoins.TruncateDecimal()
	return rewards
}
