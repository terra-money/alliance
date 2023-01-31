package e2e

import (
	"github.com/terra-money/alliance/x/alliance/keeper"
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

// Tests delegating a small amount that triggers a re-balancing event that adds < 1 utoken to a validator.
// Re-balancing event should ignore small delegations < 1 utoken since it rounds down to 0.
func TestDelegatingASmallAmount(t *testing.T) {
	allianceAsset1 := "asset1"
	allianceAsset2 := "asset2"

	app, ctx, vals, dels := setupApp(t, 5, 2, sdk.NewCoins(
		sdk.NewCoin(allianceAsset1, sdk.NewInt(1000000000000000000)),
		sdk.NewCoin(allianceAsset2, sdk.NewInt(1000000000000000000)),
	))
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			types.NewAllianceAsset(allianceAsset1, sdk.NewDec(2), sdk.NewDec(0), ctx.BlockTime()),
			types.NewAllianceAsset(allianceAsset2, sdk.NewDec(10), sdk.MustNewDecFromStr("0.1"), ctx.BlockTime()),
		},
	})
	queryServer := keeper.NewQueryServerImpl(app.AllianceKeeper)

	// Set tax and rewards to be zero for easier calculation
	distParams := app.DistrKeeper.GetParams(ctx)
	distParams.CommunityTax = sdk.ZeroDec()
	distParams.BaseProposerReward = sdk.ZeroDec()
	distParams.BonusProposerReward = sdk.ZeroDec()
	app.DistrKeeper.SetParams(ctx, distParams)

	user1 := dels[0]
	user2 := dels[1]

	val1, err := app.AllianceKeeper.GetAllianceValidator(ctx, vals[0])
	val2, err := app.AllianceKeeper.GetAllianceValidator(ctx, vals[1])

	// Delegate token with non-zero take_rate
	_, err = app.AllianceKeeper.Delegate(ctx, user1, val1, sdk.NewCoin(allianceAsset2, sdk.NewInt(100)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user2, val2, sdk.NewCoin(allianceAsset2, sdk.NewInt(1000_000_000)))
	require.NoError(t, err)

	assets := app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)

	coins, err := app.AllianceKeeper.DeductAssetsHook(ctx, assets)
	require.NoError(t, err)

	res, err := queryServer.AllianceDelegation(ctx, &types.QueryAllianceDelegationRequest{
		DelegatorAddr: user1.String(),
		ValidatorAddr: val1.GetOperator().String(),
		Denom:         allianceAsset2,
		Pagination:    nil,
	})
	require.NoError(t, err)

	del := res.GetDelegation()
	_, err = app.AllianceKeeper.Undelegate(ctx, user1, val1, sdk.NewCoin(allianceAsset2, sdk.NewInt(1000_000_000)))
	require.Error(t, err)

	// Undelegate token with current amount should pass
	_, err = app.AllianceKeeper.Undelegate(ctx, user1, val1, del.Balance)
	require.NoError(t, err)

	// User should have everything withdrawn
	_, found := app.AllianceKeeper.GetDelegation(ctx, user1, val1, allianceAsset2)
	require.False(t, found)

	// Delegate again
	_, err = app.AllianceKeeper.Delegate(ctx, user1, val1, sdk.NewCoin(allianceAsset2, sdk.NewInt(500_000_000)))
	require.NoError(t, err)

	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(time.Minute * 1)).WithBlockHeight(2)

	_, err = app.AllianceKeeper.Delegate(ctx, user1, val1, sdk.NewCoin(allianceAsset2, sdk.NewInt(400_000_000)))
	require.NoError(t, err)

	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(time.Minute * 5)).WithBlockHeight(3)
	coins, err = app.AllianceKeeper.DeductAssetsHook(ctx, assets)
	require.NoError(t, err)
	require.False(t, coins.IsZero())

	res, err = queryServer.AllianceDelegation(ctx, &types.QueryAllianceDelegationRequest{
		DelegatorAddr: user1.String(),
		ValidatorAddr: val1.GetOperator().String(),
		Denom:         allianceAsset2,
		Pagination:    nil,
	})
	require.NoError(t, err)
	del = res.GetDelegation()
	require.True(t, del.GetBalance().Amount.LT(sdk.NewInt(900_000_000)), "%s should be less than %s", del.GetBalance().Amount, sdk.NewInt(1000_000_000))

	// Undelegate token with current amount should pass
	_, err = app.AllianceKeeper.Undelegate(ctx, user1, val1, sdk.NewCoin(allianceAsset2, del.Balance.Amount))
	require.NoError(t, err)

	// User should have everything withdrawn
	_, found = app.AllianceKeeper.GetDelegation(ctx, user1, val1, allianceAsset2)
	require.False(t, found)

	res, err = queryServer.AllianceDelegation(ctx, &types.QueryAllianceDelegationRequest{
		DelegatorAddr: user1.String(),
		ValidatorAddr: val1.GetOperator().String(),
		Denom:         allianceAsset2,
		Pagination:    nil,
	})
	require.NoError(t, err)
	del = res.GetDelegation()
	require.True(t, del.Balance.Amount.IsZero())
}
