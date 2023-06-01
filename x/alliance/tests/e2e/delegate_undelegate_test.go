package e2e

import (
	"testing"
	"time"

	"github.com/terra-money/alliance/x/alliance"

	"github.com/terra-money/alliance/x/alliance/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/terra-money/alliance/x/alliance/types"
)

var (
	allianceAsset1 = "asset1"
	allianceAsset2 = "asset2"
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
		RewardWeightRange:    types.RewardWeightRange{Min: sdk.ZeroDec(), Max: sdk.MustNewDecFromStr("0.1")},
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

	del0, found := app.AllianceKeeper.GetDelegation(ctx, dels[0], val0.GetOperator(), "test")
	require.True(t, found)
	tokens := types.GetDelegationTokens(del0, val0, asset)
	_, err = app.AllianceKeeper.Undelegate(ctx, dels[0], val0, tokens)
	require.NoError(t, err)

	_, found = app.AllianceKeeper.GetDelegation(ctx, dels[0], val0.GetOperator(), "test")
	require.False(t, found)

	val0, err = app.AllianceKeeper.GetAllianceValidator(ctx, vals[0])
	require.NoError(t, err)

	require.Equal(t, sdk.ZeroDec(), sdk.DecCoins(val0.ValidatorShares).AmountOf("test"))

	_, err = app.AllianceKeeper.Delegate(ctx, dels[0], val0, sdk.NewCoin("test", sdk.NewInt(33333)))
	require.NoError(t, err)

	_, stop := alliance.RunAllInvariants(ctx, app.AllianceKeeper)
	require.False(t, stop)
}

// TestDelegateThenTakeRateThenRedelegate
// This test makes sure that full redelegation after some take rate has been
// applied will not cause a division by zero error. Also ensure that dust delegations are not kept around
func TestDelegateThenTakeRateThenRedelegate(t *testing.T) {
	app, ctx, vals, dels := setupApp(t, 5, 2, sdk.NewCoins(sdk.NewCoin("test", sdk.NewInt(1000000000000000000))))
	err := app.AllianceKeeper.CreateAlliance(ctx, &types.MsgCreateAllianceProposal{
		Title:                "",
		Description:          "",
		Denom:                "test",
		RewardWeight:         sdk.MustNewDecFromStr("0.03"),
		RewardWeightRange:    types.RewardWeightRange{Min: sdk.ZeroDec(), Max: sdk.MustNewDecFromStr("0.1")},
		TakeRate:             sdk.MustNewDecFromStr("0.02"),
		RewardChangeRate:     sdk.MustNewDecFromStr("0.01"),
		RewardChangeInterval: time.Second * 60,
	})
	require.NoError(t, err)

	val0, err := app.AllianceKeeper.GetAllianceValidator(ctx, vals[0])
	require.NoError(t, err)
	val1, err := app.AllianceKeeper.GetAllianceValidator(ctx, vals[1])
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

	del0, found := app.AllianceKeeper.GetDelegation(ctx, dels[0], val0.GetOperator(), "test")
	require.True(t, found)
	tokens := types.GetDelegationTokens(del0, val0, asset)
	_, err = app.AllianceKeeper.Redelegate(ctx, dels[0], val0, val1, tokens)
	require.NoError(t, err)

	_, found = app.AllianceKeeper.GetDelegation(ctx, dels[0], val0.GetOperator(), "test")
	require.False(t, found)

	val0, err = app.AllianceKeeper.GetAllianceValidator(ctx, vals[0])
	require.NoError(t, err)

	require.Equal(t, sdk.ZeroDec(), sdk.DecCoins(val0.ValidatorShares).AmountOf("test"))

	_, err = app.AllianceKeeper.Delegate(ctx, dels[0], val0, sdk.NewCoin("test", sdk.NewInt(33333)))
	require.NoError(t, err)

	_, stop := alliance.RunAllInvariants(ctx, app.AllianceKeeper)
	require.False(t, stop)
}

// Tests delegating a small amount that triggers a re-balancing event that adds < 1 utoken to a validator.
// Re-balancing event should ignore small delegations < 1 utoken since it rounds down to 0.
func TestDelegatingASmallAmount(t *testing.T) {
	app, ctx, vals, dels := setupApp(t, 2, 3, sdk.NewCoins(
		sdk.NewCoin(allianceAsset1, sdk.NewInt(1000000000000000000)),
		sdk.NewCoin(allianceAsset2, sdk.NewInt(1000000000000000000)),
	))
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	params := types.DefaultParams()
	params.LastTakeRateClaimTime = startTime
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: params,
		Assets: []types.AllianceAsset{
			types.NewAllianceAsset(allianceAsset1, sdk.NewDec(2), sdk.NewDec(0), sdk.NewDec(5), sdk.NewDec(0), ctx.BlockTime()),
			types.NewAllianceAsset(allianceAsset2, sdk.NewDec(10), sdk.NewDec(2), sdk.NewDec(12), sdk.MustNewDecFromStr("0.1"), ctx.BlockTime()),
		},
	})
	queryServer := keeper.NewQueryServerImpl(app.AllianceKeeper)

	// Set tax and rewards to be zero for easier calculation
	distParams := app.DistrKeeper.GetParams(ctx)
	distParams.CommunityTax = sdk.ZeroDec()
	err := app.DistrKeeper.SetParams(ctx, distParams)
	require.NoError(t, err)

	user1 := dels[0]
	user2 := dels[1]

	val1, err := app.AllianceKeeper.GetAllianceValidator(ctx, vals[0])
	require.NoError(t, err)
	val2, err := app.AllianceKeeper.GetAllianceValidator(ctx, vals[1])
	require.NoError(t, err)

	// Delegate token with non-zero take_rate
	_, err = app.AllianceKeeper.Delegate(ctx, user1, val1, sdk.NewCoin(allianceAsset2, sdk.NewInt(100)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user2, val2, sdk.NewCoin(allianceAsset2, sdk.NewInt(1000_000_000)))
	require.NoError(t, err)

	assets := app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)

	_, err = app.AllianceKeeper.DeductAssetsHook(ctx, assets)
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
	_, found := app.AllianceKeeper.GetDelegation(ctx, user1, val1.GetOperator(), allianceAsset2)
	require.False(t, found)

	// Delegate again
	_, err = app.AllianceKeeper.Delegate(ctx, user1, val1, sdk.NewCoin(allianceAsset2, sdk.NewInt(500_000_000)))
	require.NoError(t, err)

	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(time.Minute * 1)).WithBlockHeight(2)

	_, err = app.AllianceKeeper.Delegate(ctx, user1, val1, sdk.NewCoin(allianceAsset2, sdk.NewInt(400_000_000)))
	require.NoError(t, err)

	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(time.Minute * 5)).WithBlockHeight(3)
	coins, err := app.AllianceKeeper.DeductAssetsHook(ctx, assets)
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
	_, found = app.AllianceKeeper.GetDelegation(ctx, user1, val1.GetOperator(), allianceAsset2)
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

	unbondings, err := app.AllianceKeeper.GetAllUnbondings(ctx, allianceAsset2, user1)
	require.NoError(t, err)
	require.True(t, len(unbondings) == 2)
}

// This test replicates this issue where there are large amounts of tokens delegated,
// calculating token balances for a small delegation is rounded wrongly
// E.g. When user delegated 200 tokens, there was an issue such that it showed 199 tokens instead
func TestDelegateAndUndelegateWithSmallAmounts(t *testing.T) {
	app, ctx, vals, dels := setupApp(t, 5, 2, sdk.NewCoins(
		sdk.NewCoin(allianceAsset1, sdk.NewInt(2000_000_000_000_000_000)),
		sdk.NewCoin(allianceAsset2, sdk.NewInt(2000_000_000_000_000_000)),
	))
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	params := types.DefaultParams()
	params.LastTakeRateClaimTime = startTime
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: params,
		Assets: []types.AllianceAsset{
			types.NewAllianceAsset(allianceAsset1, sdk.NewDec(2), sdk.NewDec(0), sdk.NewDec(5), sdk.NewDec(0), ctx.BlockTime()),
			types.NewAllianceAsset(allianceAsset2, sdk.NewDec(10), sdk.NewDec(2), sdk.NewDec(12), sdk.MustNewDecFromStr("0.1"), ctx.BlockTime()),
		},
	})
	queryServer := keeper.NewQueryServerImpl(app.AllianceKeeper)

	// Set tax and rewards to be zero for easier calculation
	distParams := app.DistrKeeper.GetParams(ctx)
	distParams.CommunityTax = sdk.ZeroDec()

	err := app.DistrKeeper.SetParams(ctx, distParams)
	require.NoError(t, err)

	val1, err := app.AllianceKeeper.GetAllianceValidator(ctx, vals[0])
	require.NoError(t, err)
	val2, err := app.AllianceKeeper.GetAllianceValidator(ctx, vals[1])
	require.NoError(t, err)

	user1 := dels[0]
	user2 := dels[1]

	// Delegate token with non-zero take_rate
	_, err = app.AllianceKeeper.Delegate(ctx, user1, val1, sdk.NewCoin(allianceAsset2, sdk.NewInt(200)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user2, val2, sdk.NewCoin(allianceAsset2, sdk.NewInt(1000_000_000_000_000_000)))
	require.NoError(t, err)

	assets := app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)

	ctx = ctx.WithBlockTime(startTime.Add(time.Minute * 6)).WithBlockHeight(2)

	res, err := queryServer.AllianceDelegation(ctx, &types.QueryAllianceDelegationRequest{
		DelegatorAddr: user1.String(),
		ValidatorAddr: val1.GetOperator().String(),
		Denom:         allianceAsset2,
		Pagination:    nil,
	})
	require.NoError(t, err)
	del := res.GetDelegation()

	// Undelegate token with initial amount should fail
	_, err = app.AllianceKeeper.Undelegate(ctx, user1, val1, sdk.NewCoin(allianceAsset2, sdk.NewInt(1000_000_000_000_000_000)))
	require.Error(t, err)
	require.Equal(t, del.Balance.Amount, sdk.NewInt(200))

	// Undelegate token with more than current amount still pass
	_, err = app.AllianceKeeper.Undelegate(ctx, user1, val1, del.Balance)
	require.NoError(t, err)
}

// This test replicates un-delegating slightly more (1 utoken more) than the balance of token
// Due to truncation of shares, un-delegation's validation might allow more tokens to be removed than there exists in
// the delegation.
func TestUnDelegatingSlightlyMoreCoin(t *testing.T) {
	app, ctx, vals, dels := setupApp(t, 5, 2, sdk.NewCoins(
		sdk.NewCoin(allianceAsset1, sdk.NewInt(1000000000000000000)),
		sdk.NewCoin(allianceAsset2, sdk.NewInt(1000000000000000000)),
	))
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	params := types.DefaultParams()
	params.LastTakeRateClaimTime = startTime
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: params,
		Assets: []types.AllianceAsset{
			types.NewAllianceAsset(allianceAsset1, sdk.NewDec(2), sdk.NewDec(0), sdk.NewDec(5), sdk.NewDec(0), ctx.BlockTime()),
			types.NewAllianceAsset(allianceAsset2, sdk.NewDec(10), sdk.NewDec(2), sdk.NewDec(12), sdk.MustNewDecFromStr("0.1"), ctx.BlockTime()),
		},
	})
	queryServer := keeper.NewQueryServerImpl(app.AllianceKeeper)

	// Set tax and rewards to be zero for easier calculation
	distParams := app.DistrKeeper.GetParams(ctx)
	distParams.CommunityTax = sdk.ZeroDec()

	err := app.DistrKeeper.SetParams(ctx, distParams)
	require.NoError(t, err)

	val1, err := app.AllianceKeeper.GetAllianceValidator(ctx, vals[0])
	require.NoError(t, err)
	val2, err := app.AllianceKeeper.GetAllianceValidator(ctx, vals[1])
	require.NoError(t, err)

	user1 := dels[0]
	user2 := dels[1]

	// Delegate token with non-zero take_rate
	_, err = app.AllianceKeeper.Delegate(ctx, user1, val1, sdk.NewCoin(allianceAsset2, sdk.NewInt(5000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user2, val2, sdk.NewCoin(allianceAsset2, sdk.NewInt(1000_000_000)))
	require.NoError(t, err)

	assets := app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)

	ctx = ctx.WithBlockTime(startTime.Add(time.Minute * 6)).WithBlockHeight(2)
	res, err := queryServer.AllianceDelegation(ctx, &types.QueryAllianceDelegationRequest{
		DelegatorAddr: user1.String(),
		ValidatorAddr: val1.GetOperator().String(),
		Denom:         allianceAsset2,
		Pagination:    nil,
	})
	require.NoError(t, err)
	del := res.GetDelegation()

	// Undelegate token with initial amount should fail
	_, err = app.AllianceKeeper.Undelegate(ctx, user1, val1, sdk.NewCoin(allianceAsset2, sdk.NewInt(1000_000_000)))
	require.Error(t, err)

	// Undelegate token with more than current amount should fail
	_, err = app.AllianceKeeper.Undelegate(ctx, user1, val1, sdk.NewCoin(allianceAsset2, del.Balance.Amount.AddRaw(1)))
	require.Error(t, err)
}

// This test replicates re-delegating slightly more (1 utoken more) than the balance of token
// Due to truncation of shares, re-delegation's validation might allow more tokens to be removed than there exists in
// the delegation.
func TestReDelegatingSlightlyMoreCoin(t *testing.T) {
	app, ctx, vals, dels := setupApp(t, 5, 2, sdk.NewCoins(
		sdk.NewCoin(allianceAsset1, sdk.NewInt(1000000000000000000)),
		sdk.NewCoin(allianceAsset2, sdk.NewInt(1000000000000000000)),
	))
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	params := types.DefaultParams()
	params.LastTakeRateClaimTime = startTime
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: params,
		Assets: []types.AllianceAsset{
			types.NewAllianceAsset(allianceAsset1, sdk.NewDec(2), sdk.NewDec(0), sdk.NewDec(5), sdk.NewDec(0), ctx.BlockTime()),
			types.NewAllianceAsset(allianceAsset2, sdk.NewDec(10), sdk.NewDec(2), sdk.NewDec(12), sdk.MustNewDecFromStr("0.1"), ctx.BlockTime()),
		},
	})
	queryServer := keeper.NewQueryServerImpl(app.AllianceKeeper)

	// Set tax and rewards to be zero for easier calculation
	distParams := app.DistrKeeper.GetParams(ctx)
	distParams.CommunityTax = sdk.ZeroDec()

	err := app.DistrKeeper.SetParams(ctx, distParams)
	require.NoError(t, err)

	val1, err := app.AllianceKeeper.GetAllianceValidator(ctx, vals[0])
	require.NoError(t, err)
	val2, err := app.AllianceKeeper.GetAllianceValidator(ctx, vals[1])
	require.NoError(t, err)

	user1 := dels[0]
	user2 := dels[1]

	// Delegate token with non-zero take_rate
	_, err = app.AllianceKeeper.Delegate(ctx, user1, val1, sdk.NewCoin(allianceAsset2, sdk.NewInt(5000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user2, val2, sdk.NewCoin(allianceAsset2, sdk.NewInt(1000_000_000)))
	require.NoError(t, err)

	assets := app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)

	ctx = ctx.WithBlockTime(startTime.Add(time.Minute * 6)).WithBlockHeight(2)
	res, err := queryServer.AllianceDelegation(ctx, &types.QueryAllianceDelegationRequest{
		DelegatorAddr: user1.String(),
		ValidatorAddr: val1.GetOperator().String(),
		Denom:         allianceAsset2,
		Pagination:    nil,
	})
	require.NoError(t, err)
	del := res.GetDelegation()

	// Undelegate token with initial amount should fail
	_, err = app.AllianceKeeper.Undelegate(ctx, user1, val1, sdk.NewCoin(allianceAsset2, sdk.NewInt(1000_000_000)))
	require.Error(t, err)

	// Undelegate token with more than current amount should fail
	_, err = app.AllianceKeeper.Redelegate(ctx, user1, val1, val2, sdk.NewCoin(allianceAsset2, del.Balance.Amount.AddRaw(1)))
	require.Error(t, err)
}

func TestDustValidatorSharesAfterUndelegationError(t *testing.T) {
	app, ctx, vals, addrs := setupApp(t, 5, 2, sdk.NewCoins(
		sdk.NewCoin(allianceAsset1, sdk.NewInt(1000000000000000000)),
		sdk.NewCoin(allianceAsset2, sdk.NewInt(1000000000000000000)),
	))
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			types.NewAllianceAsset(allianceAsset1, sdk.NewDec(2), sdk.NewDec(0), sdk.NewDec(5), sdk.NewDec(0), ctx.BlockTime()),
			types.NewAllianceAsset(allianceAsset2, sdk.MustNewDecFromStr("10"), sdk.NewDec(5), sdk.NewDec(0), sdk.MustNewDecFromStr("0.1"), ctx.BlockTime()),
		},
	})
	queryServer := keeper.NewQueryServerImpl(app.AllianceKeeper)

	// Set tax and rewards to be zero for easier calculation
	distParams := app.DistrKeeper.GetParams(ctx)
	distParams.CommunityTax = sdk.ZeroDec()

	err := app.DistrKeeper.SetParams(ctx, distParams)
	require.NoError(t, err)

	val1, err := app.AllianceKeeper.GetAllianceValidator(ctx, vals[0])
	require.NoError(t, err)
	val2, err := app.AllianceKeeper.GetAllianceValidator(ctx, vals[1])
	require.NoError(t, err)

	user1 := addrs[0]
	user2 := addrs[1]

	// Delegate token with non-zero take_rate
	_, err = app.AllianceKeeper.Delegate(ctx, user1, val1, sdk.NewCoin(allianceAsset2, sdk.NewInt(1000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user2, val2, sdk.NewCoin(allianceAsset2, sdk.NewInt(0)))
	require.NoError(t, err)

	assets := app.AllianceKeeper.GetAllAssets(ctx)

	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)

	require.NoError(t, err)

	ctx = ctx.WithBlockTime(startTime.Add(time.Minute * 6)).WithBlockHeight(2)
	_, err = app.AllianceKeeper.DeductAssetsHook(ctx, assets)
	require.NoError(t, err)

	res, err := queryServer.AllianceDelegation(ctx, &types.QueryAllianceDelegationRequest{
		DelegatorAddr: user1.String(),
		ValidatorAddr: val1.GetOperator().String(),
		Denom:         allianceAsset2,
		Pagination:    nil,
	})
	require.NoError(t, err)
	del := res.GetDelegation()

	_, err = app.AllianceKeeper.Undelegate(ctx, user1, val1, sdk.NewCoin(allianceAsset2, del.Balance.Amount))
	require.NoError(t, err)

	assets = app.AllianceKeeper.GetAllAssets(ctx)

	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)

	require.NoError(t, err)

	_, err = app.AllianceKeeper.Delegate(ctx, user1, val1, sdk.NewCoin(allianceAsset2, sdk.NewInt(200)))
	require.NoError(t, err)
}

func TestDustValidatorSharesAfterRedelegationError(t *testing.T) {
	app, ctx, vals, addrs := setupApp(t, 5, 2, sdk.NewCoins(
		sdk.NewCoin(allianceAsset1, sdk.NewInt(1000000000000000000)),
		sdk.NewCoin(allianceAsset2, sdk.NewInt(1000000000000000000)),
	))
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			types.NewAllianceAsset(allianceAsset1, sdk.NewDec(2), sdk.NewDec(0), sdk.NewDec(5), sdk.NewDec(0), ctx.BlockTime()),
			types.NewAllianceAsset(allianceAsset2, sdk.MustNewDecFromStr("10"), sdk.NewDec(5), sdk.NewDec(0), sdk.MustNewDecFromStr("0.1"), ctx.BlockTime()),
		},
	})
	queryServer := keeper.NewQueryServerImpl(app.AllianceKeeper)

	// Set tax and rewards to be zero for easier calculation
	distParams := app.DistrKeeper.GetParams(ctx)
	distParams.CommunityTax = sdk.ZeroDec()

	err := app.DistrKeeper.SetParams(ctx, distParams)
	require.NoError(t, err)

	val1, err := app.AllianceKeeper.GetAllianceValidator(ctx, vals[0])
	require.NoError(t, err)
	val2, err := app.AllianceKeeper.GetAllianceValidator(ctx, vals[1])
	require.NoError(t, err)

	user1 := addrs[0]
	user2 := addrs[1]

	// Delegate token with non-zero take_rate
	_, err = app.AllianceKeeper.Delegate(ctx, user1, val1, sdk.NewCoin(allianceAsset2, sdk.NewInt(1000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user2, val2, sdk.NewCoin(allianceAsset2, sdk.NewInt(0)))
	require.NoError(t, err)

	assets := app.AllianceKeeper.GetAllAssets(ctx)

	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)

	require.NoError(t, err)

	ctx = ctx.WithBlockTime(startTime.Add(time.Minute * 6)).WithBlockHeight(2)
	_, err = app.AllianceKeeper.DeductAssetsHook(ctx, assets)
	require.NoError(t, err)

	res, err := queryServer.AllianceDelegation(ctx, &types.QueryAllianceDelegationRequest{
		DelegatorAddr: user1.String(),
		ValidatorAddr: val1.GetOperator().String(),
		Denom:         allianceAsset2,
		Pagination:    nil,
	})
	require.NoError(t, err)
	del := res.GetDelegation()

	_, err = app.AllianceKeeper.Redelegate(ctx, user1, val1, val2, sdk.NewCoin(allianceAsset2, del.Balance.Amount))
	require.NoError(t, err)

	assets = app.AllianceKeeper.GetAllAssets(ctx)

	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)

	require.NoError(t, err)

	_, err = app.AllianceKeeper.Delegate(ctx, user1, val1, sdk.NewCoin(allianceAsset2, sdk.NewInt(200)))
	require.NoError(t, err)
}
