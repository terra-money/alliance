package keeper_test

import (
	"alliance/x/alliance/keeper"
	"alliance/x/alliance/types"
	"fmt"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreateAlliance(t *testing.T) {
	// GIVEN
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx.WithBlockTime(startTime).WithBlockHeight(1)
	msgServer := keeper.NewMsgServerImpl(app.AllianceKeeper)
	rewardDuration := app.AllianceKeeper.RewardDelayTime(ctx)

	// WHEN
	createRes, createErr := msgServer.CreateAlliance(ctx, &types.MsgCreateAlliance{
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Alliance: types.NewAllianceAssetMsg{
			Denom:        "uluna",
			RewardWeight: sdk.OneDec(),
			TakeRate:     sdk.OneDec(),
		},
	})
	alliancesRes, alliancesErr := app.AllianceKeeper.Alliances(ctx, &types.QueryAlliancesRequest{})

	// THEN
	require.Nil(t, createErr)
	require.Equal(t, createRes, &types.MsgCreateAllianceResponse{})
	require.Nil(t, alliancesErr)
	require.Equal(t, alliancesRes, &types.QueryAlliancesResponse{
		Alliances: []types.AllianceAsset{
			{
				Denom:                "uluna",
				RewardWeight:         sdk.NewDec(1),
				TakeRate:             sdk.NewDec(1),
				TotalTokens:          sdk.ZeroInt(),
				TotalValidatorShares: sdk.NewDec(0),
				RewardStartTime:      ctx.BlockTime().Add(rewardDuration),
			},
		},
		Pagination: &query.PageResponse{
			NextKey: nil,
			Total:   1,
		},
	})
}

func TestCreateAllianceWithNoDenom(t *testing.T) {
	// GIVEN
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx.WithBlockTime(startTime).WithBlockHeight(1)
	msgServer := keeper.NewMsgServerImpl(app.AllianceKeeper)

	// WHEN
	_, err := msgServer.CreateAlliance(ctx, &types.MsgCreateAlliance{
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Alliance:  types.NewAllianceAssetMsg{},
	})

	// THEN
	require.Equal(t, err, status.Errorf(codes.InvalidArgument, "Alliance denom must have a value"))
}

func TestCreateAllianceWithNoRewardWeight(t *testing.T) {
	// GIVEN
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx.WithBlockTime(startTime).WithBlockHeight(1)
	msgServer := keeper.NewMsgServerImpl(app.AllianceKeeper)

	// WHEN
	_, err := msgServer.CreateAlliance(ctx, &types.MsgCreateAlliance{
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Alliance: types.NewAllianceAssetMsg{
			Denom: "uluna",
		},
	})

	// THEN
	require.Equal(t, err, status.Errorf(codes.InvalidArgument, "Alliance rewardWeight must be a positive number"))
}

func TestCreateAllianceWithNoTakeRate(t *testing.T) {
	// GIVEN
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx.WithBlockTime(startTime).WithBlockHeight(1)
	msgServer := keeper.NewMsgServerImpl(app.AllianceKeeper)

	// WHEN
	_, err := msgServer.CreateAlliance(ctx, &types.MsgCreateAlliance{
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Alliance: types.NewAllianceAssetMsg{
			Denom:        "uluna",
			RewardWeight: sdk.OneDec(),
		},
	})

	// THEN
	require.Equal(t, err, status.Errorf(codes.InvalidArgument, "Alliance takeRate must be a positive number"))
}

func TestCreateAllianceWithWrongAuthority(t *testing.T) {
	// GIVEN
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx.WithBlockTime(startTime).WithBlockHeight(1)
	msgServer := keeper.NewMsgServerImpl(app.AllianceKeeper)

	// WHEN
	_, err := msgServer.CreateAlliance(ctx, &types.MsgCreateAlliance{
		Authority: "cosmosvaloper19lss6zgdh5vvcpjhfftdghrpsw7a4434elpwpu",
		Alliance: types.NewAllianceAssetMsg{
			Denom:        "uluna",
			RewardWeight: sdk.OneDec(),
			TakeRate:     sdk.OneDec(),
		},
	})

	// THEN
	require.Equal(t,
		err.Error(),
		fmt.Sprintf(
			"expected %s got cosmosvaloper19lss6zgdh5vvcpjhfftdghrpsw7a4434elpwpu: expected gov account as only signer for proposal message",
			authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		),
	)
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
				RewardWeight:         sdk.NewDec(2),
				TakeRate:             sdk.OneDec(),
				TotalTokens:          sdk.ZeroInt(),
				TotalValidatorShares: sdk.NewDec(0),
			},
		},
	})
	msgServer := keeper.NewMsgServerImpl(app.AllianceKeeper)

	// WHEN
	createRes, updateErr := msgServer.UpdateAlliance(ctx, &types.MsgUpdateAlliance{
		Authority:    authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Denom:        "uluna",
		RewardWeight: sdk.NewDec(6),
		TakeRate:     sdk.NewDec(7),
	})
	alliancesRes, alliancesErr := app.AllianceKeeper.Alliances(ctx, &types.QueryAlliancesRequest{})

	// THEN
	require.Nil(t, updateErr)
	require.Equal(t, createRes, &types.MsgUpdateAllianceResponse{})
	require.Nil(t, alliancesErr)
	require.Equal(t, alliancesRes, &types.QueryAlliancesResponse{
		Alliances: []types.AllianceAsset{
			{
				Denom:                "uluna",
				RewardWeight:         sdk.NewDec(6),
				TakeRate:             sdk.NewDec(7),
				TotalTokens:          sdk.ZeroInt(),
				TotalValidatorShares: sdk.NewDec(0),
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
				RewardWeight: sdk.NewDec(2),
				TakeRate:     sdk.OneDec(),
				TotalTokens:  sdk.ZeroInt(),
			},
		},
	})
	msgServer := keeper.NewMsgServerImpl(app.AllianceKeeper)

	// WHEN
	createRes, updateErr := msgServer.DeleteAlliance(ctx, &types.MsgDeleteAlliance{
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Denom:     "uluna",
	})
	alliancesRes, alliancesErr := app.AllianceKeeper.Alliances(ctx, &types.QueryAlliancesRequest{})

	// THEN
	require.Nil(t, updateErr)
	require.Equal(t, createRes, &types.MsgDeleteAllianceResponse{})
	require.Nil(t, alliancesErr)
	require.Equal(t, alliancesRes, &types.QueryAlliancesResponse{
		Alliances: nil,
		Pagination: &query.PageResponse{
			NextKey: nil,
			Total:   0,
		},
	})
}
