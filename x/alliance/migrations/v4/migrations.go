package v4

import (
	"math"

	cmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	alliancekeeper "github.com/terra-money/alliance/x/alliance/keeper"
	"github.com/terra-money/alliance/x/alliance/types"
)

func Migrate(k alliancekeeper.Keeper) func(ctx sdk.Context) error {
	return func(ctx sdk.Context) error {
		err := migrateAssetsWithDefaultRewardWeightRange(ctx, k)
		if err != nil {
			return err
		}
		return nil
	}
}

func migrateAssetsWithDefaultRewardWeightRange(ctx sdk.Context, k alliancekeeper.Keeper) error {
	assets := k.GetAllAssets(ctx)
	for _, asset := range assets {
		asset.RewardWeightRange = types.RewardWeightRange{
			Min: cmath.LegacyZeroDec(),
			Max: cmath.LegacyNewDec(math.MaxInt),
		}
		if asset.RewardsStarted(ctx.BlockTime()) {
			asset.IsInitialized = true
		}
		if err := k.SetAsset(ctx, *asset); err != nil {
			return err
		}
	}
	return nil
}
