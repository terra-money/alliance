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
	if _, err := k.DeductAssetsHook(ctx); err != nil {
		panic(err)
	}
	if err := k.RewardWeightDecayHook(ctx); err != nil {
		panic(err)
	}
	if err := k.RebalanceHook(ctx); err != nil {
		panic(err)
	}
	return []abci.ValidatorUpdate{}
}
