package tests_test

import (
	"testing"

	"github.com/terra-money/alliance/app"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func createTestContext(t *testing.T) (*app.App, sdk.Context) {
	app := app.Setup(t)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	return app, ctx
}
