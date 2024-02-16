package tests_test

import (
	"cosmossdk.io/math"
	"testing"
	"time"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/terra-money/alliance/x/alliance/keeper"
	"github.com/terra-money/alliance/x/alliance/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
)

func TestCreateAlliance(t *testing.T) {
	// GIVEN
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx.WithBlockTime(startTime).WithBlockHeight(1)
	queryServer := keeper.NewQueryServerImpl(app.AllianceKeeper)
	rewardDuration := app.AllianceKeeper.RewardDelayTime(ctx)

	// WHEN
	createErr := app.AllianceKeeper.CreateAlliance(ctx, &types.MsgCreateAllianceProposal{
		Title:                "",
		Description:          "",
		Denom:                "uluna",
		RewardWeight:         math.LegacyOneDec(),
		RewardWeightRange:    types.RewardWeightRange{Min: math.LegacyNewDec(0), Max: math.LegacyNewDec(5)},
		RewardChangeRate:     math.LegacyOneDec(),
		RewardChangeInterval: 0,
		TakeRate:             math.LegacyMustNewDecFromStr("0.5"),
	})
	alliancesRes, alliancesErr := queryServer.Alliances(ctx, &types.QueryAlliancesRequest{})

	// THEN
	require.NoError(t, createErr)
	require.NoError(t, alliancesErr)
	require.Equal(t, alliancesRes, &types.QueryAlliancesResponse{
		Alliances: []types.AllianceAsset{
			{
				Denom:                "uluna",
				RewardWeight:         math.LegacyNewDec(1),
				RewardWeightRange:    types.RewardWeightRange{Min: math.LegacyNewDec(0), Max: math.LegacyNewDec(5)},
				TakeRate:             math.LegacyMustNewDecFromStr("0.5"),
				TotalTokens:          math.ZeroInt(),
				TotalValidatorShares: math.LegacyNewDec(0),
				RewardStartTime:      ctx.BlockTime().Add(rewardDuration),
				RewardChangeRate:     math.LegacyOneDec(),
				RewardChangeInterval: 0,
				LastRewardChangeTime: ctx.BlockTime().Add(rewardDuration),
			},
		},
		Pagination: &query.PageResponse{
			NextKey: nil,
			Total:   1,
		},
	})
}

func TestCreateAllianceFailWithDuplicatedDenom(t *testing.T) {
	// GIVEN
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx.WithBlockTime(startTime).WithBlockHeight(1)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			types.NewAllianceAsset("uluna", math.LegacyNewDec(1), math.LegacyZeroDec(), math.LegacyNewDec(2), math.LegacyNewDec(0), startTime),
		},
	})

	// WHEN
	createErr := app.AllianceKeeper.CreateAlliance(ctx, &types.MsgCreateAllianceProposal{
		Title:        "",
		Description:  "",
		Denom:        "uluna",
		RewardWeight: math.LegacyOneDec(),
		TakeRate:     math.LegacyOneDec(),
	})

	// THEN
	require.Error(t, createErr)
}

func TestUpdateAlliance(t *testing.T) {
	// GIVEN
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx.WithBlockTime(startTime).WithBlockHeight(1)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			{
				Denom:                "uluna",
				RewardWeight:         math.LegacyNewDec(2),
				RewardWeightRange:    types.RewardWeightRange{Min: math.LegacyNewDec(0), Max: math.LegacyNewDec(10)},
				TakeRate:             math.LegacyOneDec(),
				TotalTokens:          math.ZeroInt(),
				TotalValidatorShares: math.LegacyNewDec(0),
			},
		},
	})
	queryServer := keeper.NewQueryServerImpl(app.AllianceKeeper)

	// WHEN
	updateErr := app.AllianceKeeper.UpdateAlliance(ctx, &types.MsgUpdateAllianceProposal{
		Title:        "",
		Description:  "",
		Denom:        "uluna",
		RewardWeight: math.LegacyNewDec(11),
		RewardWeightRange: types.RewardWeightRange{
			Min: math.LegacyNewDec(0),
			Max: math.LegacyNewDec(11),
		},
		TakeRate:             math.LegacyNewDec(0),
		RewardChangeInterval: 0,
		RewardChangeRate:     math.LegacyOneDec(),
	})
	alliancesRes, alliancesErr := queryServer.Alliances(ctx, &types.QueryAlliancesRequest{})

	// THEN
	require.NoError(t, updateErr)
	require.NoError(t, alliancesErr)
	require.Equal(t, alliancesRes, &types.QueryAlliancesResponse{
		Alliances: []types.AllianceAsset{
			{
				Denom:                "uluna",
				RewardWeight:         math.LegacyNewDec(11),
				RewardWeightRange:    types.RewardWeightRange{Min: math.LegacyNewDec(0), Max: math.LegacyNewDec(11)},
				TakeRate:             math.LegacyNewDec(0),
				TotalTokens:          math.ZeroInt(),
				TotalValidatorShares: math.LegacyNewDec(0),
				RewardChangeRate:     math.LegacyOneDec(),
				RewardChangeInterval: 0,
			},
		},
		Pagination: &query.PageResponse{
			NextKey: nil,
			Total:   1,
		},
	})
}

func TestDeleteAlliance(t *testing.T) {
	// GIVEN
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx.WithBlockTime(startTime).WithBlockHeight(1)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			{
				Denom:        "uluna",
				RewardWeight: math.LegacyNewDec(2),
				TakeRate:     math.LegacyOneDec(),
				TotalTokens:  math.ZeroInt(),
			},
		},
	})
	queryServer := keeper.NewQueryServerImpl(app.AllianceKeeper)

	// WHEN
	deleteErr := app.AllianceKeeper.DeleteAlliance(ctx, &types.MsgDeleteAllianceProposal{
		Denom: "uluna",
	})
	alliancesRes, alliancesErr := queryServer.Alliances(ctx, &types.QueryAlliancesRequest{})

	// THEN
	require.Nil(t, deleteErr)
	require.Nil(t, alliancesErr)
	require.Equal(t, alliancesRes, &types.QueryAlliancesResponse{
		Alliances: nil,
		Pagination: &query.PageResponse{
			NextKey: nil,
			Total:   0,
		},
	})
}

func TestUpdateParams(t *testing.T) {
	// GIVEN
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx.WithBlockTime(startTime).WithBlockHeight(1)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			types.NewAllianceAsset("uluna", math.LegacyNewDec(1), math.LegacyZeroDec(), math.LegacyNewDec(2), math.LegacyNewDec(0), startTime),
		},
	})
	timeNow := time.Now().UTC()
	govAddr := authtypes.NewModuleAddress(govtypes.ModuleName).String()

	// WHEN
	msgServer := keeper.MsgServer{Keeper: app.AllianceKeeper}
	_, err := msgServer.UpdateParams(sdk.WrapSDKContext(ctx), &types.MsgUpdateParams{
		Authority: govAddr,
		Params: types.Params{
			RewardDelayTime:       100,
			TakeRateClaimInterval: 100,
			LastTakeRateClaimTime: timeNow,
		},
	})
	require.NoError(t, err)

	// THEN
	params := app.AllianceKeeper.GetParams(ctx)
	require.Equal(t, time.Duration(100), params.RewardDelayTime)
	require.Equal(t, time.Duration(100), params.TakeRateClaimInterval)
	require.Equal(t, timeNow, params.LastTakeRateClaimTime)
}

func TestUnauthorizedUpdateParams(t *testing.T) {
	// GIVEN
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx.WithBlockTime(startTime).WithBlockHeight(1)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			types.NewAllianceAsset("uluna", math.LegacyNewDec(1), math.LegacyZeroDec(), math.LegacyNewDec(2), math.LegacyNewDec(0), startTime),
		},
	})
	timeNow := time.Now().UTC()

	// WHEN
	msgServer := keeper.MsgServer{Keeper: app.AllianceKeeper}
	_, err := msgServer.UpdateParams(sdk.WrapSDKContext(ctx), &types.MsgUpdateParams{
		Authority: sdk.MustBech32ifyAddressBytes(sdk.GetConfig().GetBech32AccountAddrPrefix(), []byte("random")),
		Params: types.Params{
			RewardDelayTime:       100,
			TakeRateClaimInterval: 100,
			LastTakeRateClaimTime: timeNow,
		},
	})

	// THEN
	require.NotNil(t, err)
}
