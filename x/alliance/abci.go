package alliance

import (
	"time"

	"alliance/x/alliance/keeper"
	"alliance/x/alliance/types"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// BeginBlocker
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)
}

// EndBlocker
func EndBlocker(ctx sdk.Context, k keeper.Keeper) []abci.ValidatorUpdate {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)
	k.CompleteRedelegations(ctx)
	k.CompleteUndelegations(ctx)
	k.ClaimAssetsWithTakeRateRateLimited(ctx)
	return []abci.ValidatorUpdate{}
}
