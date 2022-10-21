package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"

	"alliance/x/alliance/types"
)

func TestQueryAlliances(t *testing.T) {
	// GIVEN: THE BLOCKCHAIN WITH ALLIANCES ON GENESIS
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
			{
				Denom:        ALLIANCE_2_TOKEN_DENOM,
				RewardWeight: sdk.NewDec(10),
				TakeRate:     sdk.MustNewDecFromStr("0.14159265359"),
				TotalTokens:  sdk.ZeroInt(),
			},
		},
	})

	// WHEN: QUERYING THE ALLIANCES LIST
	alliances, err := app.AllianceKeeper.Alliances(ctx, &types.QueryAlliancesRequest{})

	// THEN: VALIDATE THAT BOTH ALLIANCES HAVE THE CORRECT MODEL WHEN QUERYING
	require.Nil(t, err)
	require.Equal(t, &types.QueryAlliancesResponse{
		Alliances: []types.AllianceAsset{
			{
				Denom:                "alliance",
				RewardWeight:         sdk.NewDec(2),
				TakeRate:             sdk.NewDec(0),
				TotalTokens:          sdk.ZeroInt(),
				TotalValidatorShares: sdk.NewDec(0),
			},
			{
				Denom:                "alliance2",
				RewardWeight:         sdk.NewDec(10),
				TakeRate:             sdk.MustNewDecFromStr("0.14159265359"),
				TotalTokens:          sdk.ZeroInt(),
				TotalValidatorShares: sdk.NewDec(0),
			},
		},
		Pagination: &query.PageResponse{
			NextKey: nil,
			Total:   2,
		},
	}, alliances)
}

func TestQueryAnUniqueAlliance(t *testing.T) {
	// GIVEN: THE BLOCKCHAIN WITH ALLIANCES ON GENESIS
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
			{
				Denom:        ALLIANCE_2_TOKEN_DENOM,
				RewardWeight: sdk.NewDec(10),
				TakeRate:     sdk.MustNewDecFromStr("0.14159265359"),
				TotalTokens:  sdk.ZeroInt(),
			},
		},
	})

	// WHEN: QUERYING THE ALLIANCES LIST
	alliances, err := app.AllianceKeeper.Alliance(ctx, &types.QueryAllianceRequest{
		Denom: "alliance2",
	})

	// THEN: VALIDATE THAT BOTH ALLIANCES HAVE THE CORRECT MODEL WHEN QUERYING
	require.Nil(t, err)
	require.Equal(t, &types.QueryAllianceResponse{
		Alliance: &types.AllianceAsset{
			Denom:                "alliance2",
			RewardWeight:         sdk.NewDec(10),
			TakeRate:             sdk.MustNewDecFromStr("0.14159265359"),
			TotalTokens:          sdk.ZeroInt(),
			TotalValidatorShares: sdk.NewDec(0),
		},
	}, alliances)
}

func TestQueryAllianceNotFound(t *testing.T) {
	// GIVEN: THE BLOCKCHAIN
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)

	// WHEN: QUERYING THE ALLIANCE
	_, err := app.AllianceKeeper.Alliance(ctx, &types.QueryAllianceRequest{
		Denom: "alliance2",
	})

	// THEN: VALIDATE THE ERROR
	require.Equal(t, err.Error(), "alliance asset is not whitelisted")
}

func TestQueryAllAlliances(t *testing.T) {
	// GIVEN: THE BLOCKCHAIN
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)

	// WHEN: QUERYING THE ALLIANCE
	res, err := app.AllianceKeeper.Alliances(ctx, &types.QueryAlliancesRequest{})

	// THEN: VALIDATE THE ERROR
	require.Nil(t, err)
	require.Equal(t, len(res.Alliances), 0)
	require.Equal(t, res.Pagination, &query.PageResponse{
		NextKey: nil,
		Total:   0,
	})
}
