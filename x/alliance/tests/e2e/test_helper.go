package e2e

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	teststaking "github.com/cosmos/cosmos-sdk/x/staking/testutil"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	sdkmath "cosmossdk.io/math"
	test_helpers "github.com/terra-money/alliance/app"
)

func setupApp(t *testing.T, numValidators int, numDelegators int, initBalance sdk.Coins) (app *test_helpers.App, ctx sdk.Context, valAddrs []sdk.ValAddress, delAddrs []sdk.AccAddress) {
	app = test_helpers.Setup(t)
	ctx = app.BaseApp.NewContext(false)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime)

	// Accounts
	valAccAddrs := test_helpers.AddTestAddrsIncremental(app, ctx, numValidators, sdk.NewCoins())
	pks := test_helpers.CreateTestPubKeys(numValidators)

	for i := 0; i < numValidators; i++ {
		valAddr := sdk.ValAddress(valAccAddrs[i])
		valAddrs = append(valAddrs, valAddr)
		_val := teststaking.NewValidator(t, valAddr, pks[i])
		_val.Commission = stakingtypes.Commission{
			CommissionRates: stakingtypes.CommissionRates{
				Rate:          sdkmath.LegacyNewDec(0),
				MaxRate:       sdkmath.LegacyNewDec(0),
				MaxChangeRate: sdkmath.LegacyNewDec(0),
			},
			UpdateTime: time.Now(),
		}
		_val.Status = stakingtypes.Bonded
		test_helpers.RegisterNewValidator(t, app, ctx, _val)
	}

	delAddrs = test_helpers.AddTestAddrsIncremental(app, ctx, numDelegators, initBalance)
	return
}
