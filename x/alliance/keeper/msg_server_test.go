package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	keepertest "alliance/testutil/keeper"
	"alliance/testutil/nullify"
	"alliance/x/alliance/keeper"
	"alliance/x/alliance/types"
)

func TestCreateAlliance(t *testing.T) {
	k, ctx := keepertest.AllianceKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgServer := keeper.NewMsgServerImpl(k)
	asset := keepertest.CreateNewAllianceAsset(&k, ctx, 2)

	for _, tc := range []struct {
		desc     string
		request  *types.MsgCreateAlliance
		response *types.MsgCreateAllianceResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.MsgCreateAlliance{
				Authority: "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",
				Alliance:  asset,
			},
			response: &types.MsgCreateAllianceResponse{},
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := msgServer.CreateAlliance(wctx, tc.request)

			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t,
					nullify.Fill(tc.response),
					nullify.Fill(response),
				)
			}
		})
	}
}

func TestUpdateAlliance(t *testing.T) {
	k, ctx := keepertest.AllianceKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgServer := keeper.NewMsgServerImpl(k)
	msgServer.CreateAlliance(wctx, &types.MsgCreateAlliance{
		Authority: "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",
		Alliance:  keepertest.CreateNewAllianceAsset(&k, ctx, 2),
	})

	for _, tc := range []struct {
		desc     string
		request  *types.MsgUpdateAlliance
		response *types.MsgUpdateAllianceResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.MsgUpdateAlliance{
				Authority:    "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",
				Denom:        "uluna",
				RewardWeight: sdk.NewDec(2),
			},
			response: &types.MsgUpdateAllianceResponse{},
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := msgServer.UpdateAlliance(wctx, tc.request)

			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t,
					nullify.Fill(tc.response),
					nullify.Fill(response),
				)
			}
		})
	}
}

func TestDeleteAlliance(t *testing.T) {
	k, ctx := keepertest.AllianceKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgServer := keeper.NewMsgServerImpl(k)
	msgServer.CreateAlliance(wctx, &types.MsgCreateAlliance{
		Authority: "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",
		Alliance:  keepertest.CreateNewAllianceAsset(&k, ctx, 2),
	})

	for _, tc := range []struct {
		desc     string
		request  *types.MsgDeleteAlliance
		response *types.MsgDeleteAllianceResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.MsgDeleteAlliance{
				Authority: "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",
				Denom:     "uluna",
			},
			response: &types.MsgDeleteAllianceResponse{},
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := msgServer.DeleteAlliance(wctx, tc.request)

			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t,
					nullify.Fill(tc.response),
					nullify.Fill(response),
				)
			}
		})
	}
}
