package keeper

import (
	"alliance/x/alliance/types"
	"context"
)

type msgServer struct {
	Keeper
}

var _ types.MsgServer = msgServer{}

func (m msgServer) Delegate(ctx context.Context, delegate *types.MsgDelegate) (*types.MsgDelegateResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (m msgServer) Redelegate(ctx context.Context, redelegate *types.MsgRedelegate) (*types.MsgRedelegateResponse, error) {
	//TODO implement me
	panic("implement me")
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
