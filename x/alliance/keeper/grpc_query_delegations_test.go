package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	keeper "alliance/testutil/keeper"
	"alliance/testutil/nullify"
	"alliance/x/alliance/types"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func TestQueryDelegations(t *testing.T) {
	k, ctx := keeper.AllianceKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	d := keeper.CreateNewDelegation(&k, ctx, 2)

	delAddr, _ := sdk.AccAddressFromBech32(d.DelegatorAddress)

	k.SetDelegation(ctx, delAddr, stakingtypes.Validator{OperatorAddress: d.ValidatorAddress}, d.Denom, d)

	for _, tc := range []struct {
		desc     string
		request  *types.QueryAllianceDelegationRequest
		response *types.QueryAllianceDelegationResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryAllianceDelegationRequest{
				DelegatorAddr: d.DelegatorAddress,
				ValidatorAddr: d.ValidatorAddress,
				Denom:         d.Denom,
			},
			response: &types.QueryAllianceDelegationResponse{
				Delegation: types.DelegationResponse{
					Delegation: d,
					Balance:    sdk.NewCoin("uluna", sdk.NewInt(10)),
				},
			},
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := k.AllianceDelegation(wctx, tc.request)

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
