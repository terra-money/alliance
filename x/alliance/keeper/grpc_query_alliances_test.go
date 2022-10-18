package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"

	"alliance/testutil/keeper"
	"alliance/testutil/nullify"
	"alliance/x/alliance/types"
)

func TestQueryAlliances(t *testing.T) {
	k, ctx := keeper.AllianceKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	asset := keeper.CreateNewAllianceAsset(&k, ctx, 2)

	k.SetAsset(ctx, asset)

	for _, tc := range []struct {
		desc     string
		request  *types.QueryAlliancesRequest
		response *types.QueryAlliancesResponse
		err      error
	}{
		{
			desc:    "First",
			request: &types.QueryAlliancesRequest{},
			response: &types.QueryAlliancesResponse{
				Alliances: []types.AllianceAsset{asset},
				Pagination: &query.PageResponse{
					Total: 1,
				},
			},
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := k.Alliances(wctx, tc.request)

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

func TestQueryAlliance(t *testing.T) {
	k, ctx := keeper.AllianceKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	asset := keeper.CreateNewAllianceAsset(&k, ctx, 1)

	k.SetAsset(ctx, asset)

	for _, tc := range []struct {
		desc     string
		request  *types.QueryAllianceRequest
		response *types.QueryAllianceResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryAllianceRequest{
				Denom: "uluna",
			},
			response: &types.QueryAllianceResponse{
				Alliance: &asset,
			},
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := k.Alliance(wctx, tc.request)

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
