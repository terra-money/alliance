package keeper_test

import (
	"context"
	"testing"

	keepertest "alliance/testutil/keeper"
	"alliance/x/alliance/keeper"
	"alliance/x/alliance/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.AllianceKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}
