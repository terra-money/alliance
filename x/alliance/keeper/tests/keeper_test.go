package tests_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/terra-money/alliance/app"
)

func createTestContext(t *testing.T) (*app.App, sdk.Context) {
	app := app.Setup(t)
	ctx := app.NewContext(true)
	return app, ctx
}

func getOperator(val stakingtypes.ValidatorI) sdk.ValAddress {
	addr, _ := sdk.ValAddressFromBech32(val.GetOperator())
	return addr
}
