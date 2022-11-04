package keeper

import (
	"alliance/x/alliance/types"
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgServer struct {
	Keeper
}

var _ types.MsgServer = MsgServer{}

func (m MsgServer) Delegate(ctx context.Context, delegate *types.MsgDelegate) (*types.MsgDelegateResponse, error) {
	err := delegate.ValidateBasic()
	if err != nil {
		return nil, err
	}

	delAddr, err := sdk.AccAddressFromBech32(delegate.DelegatorAddress)
	if err != nil {
		return nil, err
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	valAddr, err := sdk.ValAddressFromBech32(delegate.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	validator, err := m.Keeper.GetAllianceValidator(sdkCtx, valAddr)
	if err != nil {
		return nil, err
	}

	_, err = m.Keeper.Delegate(sdkCtx, delAddr, validator, delegate.Amount)
	if err != nil {
		return nil, err
	}
	return &types.MsgDelegateResponse{}, nil
}

func (m MsgServer) Redelegate(ctx context.Context, redelegate *types.MsgRedelegate) (*types.MsgRedelegateResponse, error) {
	err := redelegate.ValidateBasic()
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	delAddr, err := sdk.AccAddressFromBech32(redelegate.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	srcValAddr, err := sdk.ValAddressFromBech32(redelegate.ValidatorSrcAddress)
	if err != nil {
		return nil, err
	}
	dstValAddr, err := sdk.ValAddressFromBech32(redelegate.ValidatorDstAddress)
	if err != nil {
		return nil, err
	}

	srcValidator, err := m.Keeper.GetAllianceValidator(sdkCtx, srcValAddr)
	if err != nil {
		return nil, err
	}

	dstValidator, err := m.Keeper.GetAllianceValidator(sdkCtx, dstValAddr)
	if err != nil {
		return nil, err
	}

	_, err = m.Keeper.Redelegate(sdkCtx, delAddr, srcValidator, dstValidator, redelegate.Amount)
	if err != nil {
		return nil, err
	}
	return &types.MsgRedelegateResponse{}, nil
}

func (m MsgServer) Undelegate(ctx context.Context, undelegate *types.MsgUndelegate) (*types.MsgUndelegateResponse, error) {
	err := undelegate.ValidateBasic()
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	delAddr, err := sdk.AccAddressFromBech32(undelegate.DelegatorAddress)
	if err != nil {
		return nil, err
	}
	valAddr, err := sdk.ValAddressFromBech32(undelegate.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	validator, err := m.Keeper.GetAllianceValidator(sdkCtx, valAddr)
	if err != nil {
		return nil, err
	}

	err = m.Keeper.Undelegate(sdkCtx, delAddr, validator, undelegate.Amount)
	if err != nil {
		return nil, err
	}
	return &types.MsgUndelegateResponse{}, nil
}

func (m MsgServer) ClaimDelegationRewards(ctx context.Context, request *types.MsgClaimDelegationRewards) (*types.MsgClaimDelegationRewardsResponse, error) {
	err := request.ValidateBasic()
	if err != nil {
		return nil, err
	}

	delAddr, err := sdk.AccAddressFromBech32(request.DelegatorAddress)
	if err != nil {
		return nil, err
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	valAddr, err := sdk.ValAddressFromBech32(request.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	validator, err := m.Keeper.GetAllianceValidator(sdkCtx, valAddr)
	if err != nil {
		return nil, err
	}

	_, err = m.Keeper.ClaimDelegationRewards(sdkCtx, delAddr, validator, request.Denom)
	return &types.MsgClaimDelegationRewardsResponse{}, err
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &MsgServer{Keeper: keeper}
}
