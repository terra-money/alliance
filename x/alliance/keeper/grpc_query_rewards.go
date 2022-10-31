package keeper

import (
	"alliance/x/alliance/types"
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (k Keeper) AllianceDelegationRewards(context context.Context, request *types.QueryAllianceDelegationRewardsRequest) (*types.QueryAllianceDelegationRewardsResponse, error) {
	ctx := sdk.UnwrapSDKContext(context)
	delAddr, err := sdk.AccAddressFromBech32(request.DelegatorAddr)
	if err != nil {
		return nil, err
	}
	valAddr, err := sdk.ValAddressFromBech32(request.ValidatorAddr)
	if err != nil {
		return nil, err
	}
	_, found := k.GetAssetByDenom(ctx, request.Denom)
	if !found {
		return nil, types.ErrUnknownAsset
	}

	val, err := k.GetAllianceValidator(ctx, valAddr)
	if err != nil {
		return nil, err
	}

	_, found = k.GetDelegation(ctx, delAddr, val, request.Denom)
	if !found {
		return nil, stakingtypes.ErrNoDelegation
	}

	rewards, err := k.ClaimDelegationRewards(ctx, delAddr, val, request.Denom)
	if err != nil {
		return nil, err
	}
	return &types.QueryAllianceDelegationRewardsResponse{
		Rewards: rewards,
	}, nil
}
