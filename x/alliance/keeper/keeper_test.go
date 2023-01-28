package keeper_test

import (
	"testing"

	"github.com/terra-money/alliance/app"

	sdk "github.com/cosmos/cosmos-sdk/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func createTestContext(t *testing.T) (*app.App, sdk.Context) {
	app := app.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	return app, ctx
}
