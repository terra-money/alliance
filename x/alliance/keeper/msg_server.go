package keeper

import (
	"context"
	"github.com/terra-money/alliance/x/alliance/types"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgServer struct {
	Keeper
}

var _ types.MsgServer = MsgServer{}

func (m MsgServer) Delegate(ctx context.Context, msg *types.MsgDelegate) (*types.MsgDelegateResponse, error) {
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}

	delAddr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil, err
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	valAddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	validator, err := m.Keeper.GetAllianceValidator(sdkCtx, valAddr)
	if err != nil {
		return nil, err
	}

	newShares, err := m.Keeper.Delegate(sdkCtx, delAddr, validator, msg.Amount)
	if err != nil {
		return nil, err
	}

	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDelegate,
			sdk.NewAttribute(types.AttributeKeyValidator, msg.ValidatorAddress),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyNewShares, newShares.String()),
		),
	})

	return &types.MsgDelegateResponse{}, nil
}

func (m MsgServer) Redelegate(ctx context.Context, msg *types.MsgRedelegate) (*types.MsgRedelegateResponse, error) {
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	delAddr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	srcValAddr, err := sdk.ValAddressFromBech32(msg.ValidatorSrcAddress)
	if err != nil {
		return nil, err
	}
	dstValAddr, err := sdk.ValAddressFromBech32(msg.ValidatorDstAddress)
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

	completionTime, err := m.Keeper.Redelegate(sdkCtx, delAddr, srcValidator, dstValidator, msg.Amount)
	if err != nil {
		return nil, err
	}

	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRedelegate,
			sdk.NewAttribute(types.AttributeKeySrcValidator, msg.ValidatorSrcAddress),
			sdk.NewAttribute(types.AttributeKeyDstValidator, msg.ValidatorDstAddress),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyCompletionTime, completionTime.Format(time.RFC3339)),
		),
	})
	return &types.MsgRedelegateResponse{}, nil
}

func (m MsgServer) Undelegate(ctx context.Context, msg *types.MsgUndelegate) (*types.MsgUndelegateResponse, error) {
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	delAddr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil, err
	}
	valAddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	validator, err := m.Keeper.GetAllianceValidator(sdkCtx, valAddr)
	if err != nil {
		return nil, err
	}

	completionTime, err := m.Keeper.Undelegate(sdkCtx, delAddr, validator, msg.Amount)
	if err != nil {
		return nil, err
	}

	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUndelegate,
			sdk.NewAttribute(types.AttributeKeyValidator, msg.ValidatorAddress),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyCompletionTime, completionTime.Format(time.RFC3339)),
		),
	})
	return &types.MsgUndelegateResponse{}, nil
}

func (m MsgServer) ClaimDelegationRewards(ctx context.Context, msg *types.MsgClaimDelegationRewards) (*types.MsgClaimDelegationRewardsResponse, error) {
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}

	delAddr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil, err
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	valAddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	validator, err := m.Keeper.GetAllianceValidator(sdkCtx, valAddr)
	if err != nil {
		return nil, err
	}

	coins, err := m.Keeper.ClaimDelegationRewards(sdkCtx, delAddr, validator, msg.Denom)

	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeClaimDelegationRewards,
			sdk.NewAttribute(types.AttributeKeyValidator, msg.ValidatorAddress),
			sdk.NewAttribute(sdk.AttributeKeyAmount, coins.String()),
		),
	})
	return &types.MsgClaimDelegationRewardsResponse{}, err
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &MsgServer{Keeper: keeper}
}
