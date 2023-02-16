package alliance

import (
	"fmt"

	"github.com/terra-money/alliance/x/alliance/keeper"
	"github.com/terra-money/alliance/x/alliance/types"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// EndBlocker
func EndBlocker(ctx sdk.Context, k keeper.Keeper) []abci.ValidatorUpdate {
	defer telemetry.ModuleMeasureSince(types.ModuleName, ctx.BlockTime(), telemetry.MetricKeyEndBlocker)
	k.CompleteRedelegations(ctx)
	if err := k.CompleteUndelegations(ctx); err != nil {
		panic(fmt.Errorf("failed to complete undelegations from x/alliance module: %s", err))
	}

	assets := k.GetAllAssets(ctx)
	if _, err := k.DeductAssetsHook(ctx, assets); err != nil {
		panic(fmt.Errorf("failed to deduct take rate from alliance in x/alliance module: %s", err))
	}
	k.RewardWeightChangeHook(ctx, assets)
	if err := k.RebalanceHook(ctx, assets); err != nil {
		panic(fmt.Errorf("failed to rebalance assets in x/alliance module: %s", err))
	}
	return []abci.ValidatorUpdate{}
}
