package keeper_test

import (
	"alliance/app"
	"alliance/x/alliance/keeper"
	"alliance/x/alliance/types"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func createTestContext(t *testing.T) (*app.App, sdk.Context) {
	app := app.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	app.AllianceKeeper = keeper.NewKeeper(app.AppCodec(), app.GetKey(types.StoreKey), app.GetSubspace(types.ModuleName))
	return app, ctx
}
