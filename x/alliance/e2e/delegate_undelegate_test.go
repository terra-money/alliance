package e2e

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/terra-money/alliance/x/alliance/types"
)

// TestDelegateAndUndelegate
// This test makes sure that full undelegation after some take rate has been
// applied will not cause a division by zero error.
func TestDelegateThenTakeRateThenUndelegate(t *testing.T) {
	app, ctx, vals, dels := setupApp(t, 5, 2, sdk.NewCoins(sdk.NewCoin("test", sdk.NewInt(1000000000000000000))))
	err := app.AllianceKeeper.CreateAlliance(ctx, &types.MsgCreateAllianceProposal{
		Title:                "",
		Description:          "",
		Denom:                "test",
		RewardWeight:         sdk.MustNewDecFromStr("0.03"),
		TakeRate:             sdk.MustNewDecFromStr("0.02"),
		RewardChangeRate:     sdk.MustNewDecFromStr("0.01"),
		RewardChangeInterval: time.Second * 60,
	})
	require.NoError(t, err)

	val0, err := app.AllianceKeeper.GetAllianceValidator(ctx, vals[0])
	require.NoError(t, err)

	_, err = app.AllianceKeeper.Delegate(ctx, dels[0], val0, sdk.NewCoin("test", sdk.NewInt(100033333333333333)))
	require.NoError(t, err)

	val0, err = app.AllianceKeeper.GetAllianceValidator(ctx, vals[0])
	require.NoError(t, err)
	require.Equal(t, sdk.NewDec(100033333333333333), sdk.DecCoins(val0.TotalDelegatorShares).AmountOf("test"))

	lastClaim := ctx.BlockTime()
	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(time.Hour))
	assets := app.AllianceKeeper.GetAllAssets(ctx)
	_, err = app.AllianceKeeper.DeductAssetsWithTakeRate(ctx, lastClaim, assets)
	require.NoError(t, err)

	asset, _ := app.AllianceKeeper.GetAssetByDenom(ctx, "test")

	del0, found := app.AllianceKeeper.GetDelegation(ctx, dels[0], val0, "test")
	require.True(t, found)
	tokens := types.GetDelegationTokens(del0, val0, asset)
	_, err = app.AllianceKeeper.Undelegate(ctx, dels[0], val0, tokens)
	require.NoError(t, err)

	_, found = app.AllianceKeeper.GetDelegation(ctx, dels[0], val0, "test")
	require.False(t, found)

	val0, err = app.AllianceKeeper.GetAllianceValidator(ctx, vals[0])
	require.NoError(t, err)

	require.Equal(t, sdk.ZeroDec(), sdk.DecCoins(val0.ValidatorShares).AmountOf("test"))

	_, err = app.AllianceKeeper.Delegate(ctx, dels[0], val0, sdk.NewCoin("test", sdk.NewInt(33333)))
	require.NoError(t, err)
}
