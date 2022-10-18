package keeper

import (
	"alliance/x/alliance/types"
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type msgServer struct {
	Keeper
}

var _ types.MsgServer = msgServer{}

func (m msgServer) Delegate(ctx context.Context, delegate *types.MsgDelegate) (*types.MsgDelegateResponse, error) {
	delAddr, err := sdk.AccAddressFromBech32(delegate.DelegatorAddress)
	if err != nil {
		return nil, err
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	valAddr, err := sdk.ValAddressFromBech32(delegate.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	validator, ok := m.Keeper.stakingKeeper.GetValidator(sdkCtx, valAddr)
	if !ok {
		return nil, stakingtypes.ErrNoValidatorFound
	}

	_, err = m.Keeper.Delegate(sdkCtx, delAddr, validator, delegate.Amount)
	if err != nil {
		return nil, err
	}
	return &types.MsgDelegateResponse{}, nil
}

func (m msgServer) Redelegate(ctx context.Context, redelegate *types.MsgRedelegate) (*types.MsgRedelegateResponse, error) {
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

	srcValidator, ok := m.Keeper.stakingKeeper.GetValidator(sdkCtx, srcValAddr)
	if !ok {
		return nil, stakingtypes.ErrNoValidatorFound
	}
	dstValidator, ok := m.Keeper.stakingKeeper.GetValidator(sdkCtx, dstValAddr)
	if !ok {
		return nil, stakingtypes.ErrNoValidatorFound
	}

	_, err = m.Keeper.Redelegate(sdkCtx, delAddr, srcValidator, dstValidator, redelegate.Amount)
	if err != nil {
		return nil, err
	}
	return &types.MsgRedelegateResponse{}, nil
}

func (m msgServer) Undelegate(ctx context.Context, undelegate *types.MsgUndelegate) (*types.MsgUndelegateResponse, error) {
	//TODO implement me
	panic("implement me")
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}
