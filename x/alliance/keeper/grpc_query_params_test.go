package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"alliance/x/alliance/types"
)

func TestQueryParams(t *testing.T) {
	// GIVEN: THE BLOCKCHAIN WITH AN ALLIANCE ON GENESIS
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			{
				Denom:        ALLIANCE_TOKEN_DENOM,
				RewardWeight: sdk.NewDec(2),
				TakeRate:     sdk.NewDec(0),
				TotalTokens:  sdk.ZeroInt(),
			},
		},
	})
	// WHEN: QUERYING THE PARAMS...
	queyParams, err := app.AllianceKeeper.Params(ctx, &types.QueryParamsRequest{})

	// THEN: VALIDATE THAT NO ERRORS HAVE BEEN PRODUCED AND OUTPUT IS AS WE EXPECT
	require.Nil(t, err)

	require.Equal(t, queyParams.Params.RewardDelayTime, time.Hour)
	require.Equal(t, queyParams.Params.RewardClaimInterval, time.Minute*5)

	// there is no way to match the exact time when the module is being instantiated
	// but we know that this time should be older than actually the time when this
	// following two lines are executed
	require.NotNil(t, queyParams.Params.LastRewardClaimTime)
	require.LessOrEqual(t, queyParams.Params.LastRewardClaimTime, time.Now())
}
