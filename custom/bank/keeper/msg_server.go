package keeper

import (
	"context"

	"cosmossdk.io/core/address"
	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/hashicorp/go-metrics"
)

type msgServer struct {
	types.MsgServer

	keeper       bankkeeper.Keeper
	addressCodec address.Codec
}

var _ types.MsgServer = msgServer{}

func NewMsgServerImpl(keeper Keeper, addressCodec address.Codec) types.MsgServer {
	return &msgServer{
		MsgServer:    bankkeeper.NewMsgServerImpl(keeper),
		keeper:       keeper,
		addressCodec: addressCodec,
	}
}

func (k msgServer) Send(goCtx context.Context, msg *types.MsgSend) (*types.MsgSendResponse, error) {
	from, err := k.addressCodec.StringToBytes(msg.FromAddress)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", err)
	}
	to, err := k.addressCodec.StringToBytes(msg.ToAddress)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid to address: %s", err)
	}

	if !msg.Amount.IsValid() {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}

	if !msg.Amount.IsAllPositive() {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := k.keeper.IsSendEnabledCoins(ctx, msg.Amount...); err != nil {
		return nil, err
	}

	if k.keeper.BlockedAddr(to) {
		return nil, errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to receive funds", msg.ToAddress)
	}

	err = k.keeper.SendCoins(ctx, from, to, msg.Amount)
	if err != nil {
		return nil, err
	}

	defer func() {
		for _, a := range msg.Amount {
			if a.Amount.IsInt64() {
				telemetry.SetGaugeWithLabels(
					[]string{"tx", "msg", "send"},
					float32(a.Amount.Int64()),
					[]metrics.Label{telemetry.NewLabel("denom", a.Denom)},
				)
			}
		}
	}()

	return &types.MsgSendResponse{}, nil
}

func (k msgServer) MultiSend(goCtx context.Context, msg *types.MsgMultiSend) (*types.MsgMultiSendResponse, error) {
	if len(msg.Inputs) == 0 {
		return nil, types.ErrNoInputs
	}

	if len(msg.Inputs) != 1 {
		return nil, types.ErrMultipleSenders
	}

	if len(msg.Outputs) == 0 {
		return nil, types.ErrNoOutputs
	}

	if err := types.ValidateInputOutputs(msg.Inputs[0], msg.Outputs); err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// NOTE: totalIn == totalOut should already have been checked
	for _, in := range msg.Inputs {
		if err := k.keeper.IsSendEnabledCoins(ctx, in.Coins...); err != nil {
			return nil, err
		}
	}

	for _, out := range msg.Outputs {
		accAddr, err := k.addressCodec.StringToBytes(out.Address)
		if err != nil {
			return nil, err
		}

		if k.keeper.BlockedAddr(accAddr) {
			return nil, errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to receive funds", out.Address)
		}
	}

	err := k.keeper.InputOutputCoins(ctx, msg.Inputs[0], msg.Outputs)
	if err != nil {
		return nil, err
	}

	return &types.MsgMultiSendResponse{}, nil
}
