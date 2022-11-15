package alliance

import (
	"time"

	"github.com/terra-money/alliance/x/alliance/keeper"
	"github.com/terra-money/alliance/x/alliance/types"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// EndBlocker
func EndBlocker(ctx sdk.Context, k keeper.Keeper) []abci.ValidatorUpdate {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)
	k.CompleteRedelegations(ctx)
	if err := k.CompleteUndelegations(ctx); err != nil {
		panic(err)
	}

	assets := k.GetAllAssets(ctx)
	if _, err := k.DeductAssetsHook(ctx, assets); err != nil {
		panic(err)
	}
	k.RewardWeightDecayHook(ctx, assets)
	if err := k.RebalanceHook(ctx, assets); err != nil {
		panic(err)
	}
	return []abci.ValidatorUpdate{}
}
