package tests_test

import (
	"testing"
	"time"

	test_helpers "github.com/terra-money/alliance/app"
	"github.com/terra-money/alliance/x/alliance"
	"github.com/terra-money/alliance/x/alliance/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/teststaking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
)

func TestSlashingEvent(t *testing.T) {
	var err error
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			{
				Denom:        AllianceDenom,
				RewardWeight: sdk.NewDec(2),
				TakeRate:     sdk.NewDec(0),
				TotalTokens:  sdk.ZeroInt(),
			},
			{
				Denom:        AllianceDenomTwo,
				RewardWeight: sdk.NewDec(10),
				TakeRate:     sdk.NewDec(0),
				TotalTokens:  sdk.ZeroInt(),
			},
		},
	})

	// Set tax and rewards to be zero for easier calculation
	distParams := app.DistrKeeper.GetParams(ctx)
	distParams.CommunityTax = sdk.ZeroDec()
	distParams.BaseProposerReward = sdk.ZeroDec()
	distParams.BonusProposerReward = sdk.ZeroDec()
	app.DistrKeeper.SetParams(ctx, distParams)

	// Accounts
	addrs := test_helpers.AddTestAddrsIncremental(app, ctx, 4, sdk.NewCoins(
		sdk.NewCoin(AllianceDenom, sdk.NewInt(20_000_000)),
		sdk.NewCoin(AllianceDenomTwo, sdk.NewInt(20_000_000)),
	))
	pks := test_helpers.CreateTestPubKeys(2)

	// Creating two validators: 1 with 0% commission, 1 with 100% commission
	valAddr1 := sdk.ValAddress(addrs[0])
	_val1 := teststaking.NewValidator(t, valAddr1, pks[0])
	_val1.Commission = stakingtypes.Commission{
		CommissionRates: stakingtypes.CommissionRates{
			Rate:          sdk.NewDec(0),
			MaxRate:       sdk.NewDec(0),
			MaxChangeRate: sdk.NewDec(0),
		},
		UpdateTime: time.Now(),
	}
	test_helpers.RegisterNewValidator(t, app, ctx, _val1)
	val1, err := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr1)
	require.NoError(t, err)

	valAddr2 := sdk.ValAddress(addrs[1])
	_val2 := teststaking.NewValidator(t, valAddr2, pks[1])
	_val2.Commission = stakingtypes.Commission{
		CommissionRates: stakingtypes.CommissionRates{
			Rate:          sdk.NewDec(1),
			MaxRate:       sdk.NewDec(1),
			MaxChangeRate: sdk.NewDec(0),
		},
		UpdateTime: time.Now(),
	}
	test_helpers.RegisterNewValidator(t, app, ctx, _val2)
	val2, err := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr2)
	require.NoError(t, err)

	user1 := addrs[2]
	user2 := addrs[3]

	// Users add delegations
	_, err = app.AllianceKeeper.Delegate(ctx, user1, val1, sdk.NewCoin(AllianceDenom, sdk.NewInt(10_000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user1, val2, sdk.NewCoin(AllianceDenom, sdk.NewInt(10_000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user2, val1, sdk.NewCoin(AllianceDenom, sdk.NewInt(10_000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user2, val2, sdk.NewCoin(AllianceDenom, sdk.NewInt(10_000_000)))
	require.NoError(t, err)

	_, err = app.AllianceKeeper.Delegate(ctx, user1, val1, sdk.NewCoin(AllianceDenomTwo, sdk.NewInt(10_000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user1, val2, sdk.NewCoin(AllianceDenomTwo, sdk.NewInt(10_000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user2, val1, sdk.NewCoin(AllianceDenomTwo, sdk.NewInt(10_000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user2, val2, sdk.NewCoin(AllianceDenomTwo, sdk.NewInt(10_000_000)))
	require.NoError(t, err)

	assets := app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)
	require.Equal(t, sdk.NewInt(13_000_000), app.StakingKeeper.TotalBondedTokens(ctx))

	val1, _ = app.AllianceKeeper.GetAllianceValidator(ctx, valAddr1)
	valPower1 := val1.GetConsensusPower(app.StakingKeeper.PowerReduction(ctx))
	valConAddr1, _ := val1.GetConsAddr()

	// Tokens should remain the same before slashing
	asset1, _ := app.AllianceKeeper.GetAssetByDenom(ctx, AllianceDenom)
	tokens := val1.TotalTokensWithAsset(asset1).TruncateInt()
	require.Equal(t, sdk.NewInt(20_000_000), tokens)
	asset2, _ := app.AllianceKeeper.GetAssetByDenom(ctx, AllianceDenomTwo)
	tokens = val1.TotalTokensWithAsset(asset2).TruncateInt()
	require.Equal(t, sdk.NewInt(20_000_000), tokens)

	app.SlashingKeeper.Slash(ctx, valConAddr1, app.SlashingKeeper.SlashFractionDoubleSign(ctx), valPower1, 1)
	// Slashing will first reduce tokens from validator
	require.NotEqual(t, sdk.NewInt(13_000_000), app.StakingKeeper.TotalBondedTokens(ctx))

	// After rebalancing, it should recover the tokens
	assets = app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)
	require.Equal(t, sdk.NewInt(12_999_999), app.StakingKeeper.TotalBondedTokens(ctx))

	// Expect that total tokens with validator 1 are reduced
	val1, _ = app.AllianceKeeper.GetAllianceValidator(ctx, valAddr1)
	asset1, _ = app.AllianceKeeper.GetAssetByDenom(ctx, AllianceDenom)
	tokens = val1.TotalTokensWithAsset(asset1).TruncateInt()
	require.Greater(t, sdk.NewInt(20_000_000).Int64(), tokens.Int64())
	asset2, _ = app.AllianceKeeper.GetAssetByDenom(ctx, AllianceDenomTwo)
	tokens = val1.TotalTokensWithAsset(asset2).TruncateInt()
	require.Greater(t, sdk.NewInt(20_000_000).Int64(), tokens.Int64())

	// Expect that total tokens with validator 2 increased (redistributed from slashing)
	val2, _ = app.AllianceKeeper.GetAllianceValidator(ctx, valAddr2)
	asset1, _ = app.AllianceKeeper.GetAssetByDenom(ctx, AllianceDenom)
	tokens = val2.TotalTokensWithAsset(asset1).TruncateInt()
	require.Less(t, sdk.NewInt(20_000_000).Int64(), tokens.Int64())
	asset2, _ = app.AllianceKeeper.GetAssetByDenom(ctx, AllianceDenomTwo)
	tokens = val2.TotalTokensWithAsset(asset2).TruncateInt()
	require.Less(t, sdk.NewInt(20_000_000).Int64(), tokens.Int64())

	// Expect that consensus power for val1 dropped
	newValPower1 := val1.GetConsensusPower(app.StakingKeeper.PowerReduction(ctx))
	require.Less(t, newValPower1, valPower1)

	_, stop := alliance.RunAllInvariants(ctx, app.AllianceKeeper)
	require.False(t, stop)
}

func TestSlashingAfterRedelegation(t *testing.T) {
	var err error
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			{
				Denom:        AllianceDenom,
				RewardWeight: sdk.NewDec(2),
				TakeRate:     sdk.NewDec(0),
				TotalTokens:  sdk.ZeroInt(),
			},
			{
				Denom:        AllianceDenomTwo,
				RewardWeight: sdk.NewDec(10),
				TakeRate:     sdk.NewDec(0),
				TotalTokens:  sdk.ZeroInt(),
			},
		},
	})

	// Set tax and rewards to be zero for easier calculation
	distParams := app.DistrKeeper.GetParams(ctx)
	distParams.CommunityTax = sdk.ZeroDec()
	distParams.BaseProposerReward = sdk.ZeroDec()
	distParams.BonusProposerReward = sdk.ZeroDec()
	app.DistrKeeper.SetParams(ctx, distParams)

	// Accounts
	addrs := test_helpers.AddTestAddrsIncremental(app, ctx, 4, sdk.NewCoins(
		sdk.NewCoin(AllianceDenom, sdk.NewInt(20_000_000)),
		sdk.NewCoin(AllianceDenomTwo, sdk.NewInt(20_000_000)),
	))
	pks := test_helpers.CreateTestPubKeys(2)

	// Creating two validators: 1 with 0% commission, 1 with 100% commission
	valAddr1 := sdk.ValAddress(addrs[0])
	_val1 := teststaking.NewValidator(t, valAddr1, pks[0])
	_val1.Commission = stakingtypes.Commission{
		CommissionRates: stakingtypes.CommissionRates{
			Rate:          sdk.NewDec(0),
			MaxRate:       sdk.NewDec(0),
			MaxChangeRate: sdk.NewDec(0),
		},
		UpdateTime: time.Now(),
	}
	test_helpers.RegisterNewValidator(t, app, ctx, _val1)
	val1, err := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr1)
	require.NoError(t, err)

	valAddr2 := sdk.ValAddress(addrs[1])
	_val2 := teststaking.NewValidator(t, valAddr2, pks[1])
	_val2.Commission = stakingtypes.Commission{
		CommissionRates: stakingtypes.CommissionRates{
			Rate:          sdk.NewDec(1),
			MaxRate:       sdk.NewDec(1),
			MaxChangeRate: sdk.NewDec(0),
		},
		UpdateTime: time.Now(),
	}
	test_helpers.RegisterNewValidator(t, app, ctx, _val2)
	val2, err := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr2)
	require.NoError(t, err)

	user1 := addrs[2]
	user2 := addrs[3]

	// Users add delegations
	_, err = app.AllianceKeeper.Delegate(ctx, user1, val1, sdk.NewCoin(AllianceDenom, sdk.NewInt(10_000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user1, val2, sdk.NewCoin(AllianceDenom, sdk.NewInt(10_000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user2, val2, sdk.NewCoin(AllianceDenom, sdk.NewInt(10_000_000)))
	require.NoError(t, err)

	_, err = app.AllianceKeeper.Delegate(ctx, user1, val1, sdk.NewCoin(AllianceDenomTwo, sdk.NewInt(10_000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user1, val2, sdk.NewCoin(AllianceDenomTwo, sdk.NewInt(10_000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user2, val1, sdk.NewCoin(AllianceDenomTwo, sdk.NewInt(10_000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user2, val2, sdk.NewCoin(AllianceDenomTwo, sdk.NewInt(10_000_000)))
	require.NoError(t, err)

	assets := app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)
	require.Equal(t, sdk.NewInt(12_999_999), app.StakingKeeper.TotalBondedTokens(ctx))

	_, err = app.AllianceKeeper.Redelegate(ctx, user1, val1, val2, sdk.NewCoin(AllianceDenom, sdk.NewInt(10_000_000)))
	require.NoError(t, err)

	assets = app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)
	require.Equal(t, sdk.NewInt(13_000_000), app.StakingKeeper.TotalBondedTokens(ctx))

	// Expect that delegation has increased
	delegation, _ := app.AllianceKeeper.GetDelegation(ctx, user1, valAddr2, AllianceDenom)
	asset, _ := app.AllianceKeeper.GetAssetByDenom(ctx, AllianceDenom)
	tokens := types.GetDelegationTokens(delegation, val2, asset)
	require.Equal(t, sdk.NewInt(20_000_000), tokens.Amount)

	// Now we slash val 1
	val1, _ = app.AllianceKeeper.GetAllianceValidator(ctx, valAddr1)
	valPower1 := val1.GetConsensusPower(app.StakingKeeper.PowerReduction(ctx))
	valConAddr1, _ := val1.GetConsAddr()
	slashFraction := app.SlashingKeeper.SlashFractionDoubleSign(ctx)
	app.SlashingKeeper.Slash(ctx, valConAddr1, slashFraction, valPower1, 1)

	// Expect that delegation decreased
	delegation, _ = app.AllianceKeeper.GetDelegation(ctx, user1, valAddr2, AllianceDenom)
	asset, _ = app.AllianceKeeper.GetAssetByDenom(ctx, AllianceDenom)
	tokens = types.GetDelegationTokens(delegation, val2, asset)
	require.Greater(t, sdk.NewInt(20_000_000).Int64(), tokens.Amount.Int64())

	// Move time to after redelegation completes
	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(app.StakingKeeper.UnbondingTime(ctx)).Add(time.Second))

	// Now we slash val 1
	_, err = app.AllianceKeeper.GetAllianceValidator(ctx, valAddr1)
	require.NoError(t, err)
	app.SlashingKeeper.Slash(ctx, valConAddr1, slashFraction, valPower1, 1)

	// Expect that delegation stayed the same
	delegation, _ = app.AllianceKeeper.GetDelegation(ctx, user1, valAddr2, AllianceDenom)
	asset, _ = app.AllianceKeeper.GetAssetByDenom(ctx, AllianceDenom)
	require.Equal(t, tokens.Amount.Int64(), types.GetDelegationTokens(delegation, val2, asset).Amount.Int64())

	_, stop := alliance.RunAllInvariants(ctx, app.AllianceKeeper)
	require.False(t, stop)
}

func TestSlashingAfterUndelegation(t *testing.T) {
	var err error
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			{
				Denom:        AllianceDenom,
				RewardWeight: sdk.NewDec(2),
				TakeRate:     sdk.NewDec(0),
				TotalTokens:  sdk.ZeroInt(),
			},
			{
				Denom:        AllianceDenomTwo,
				RewardWeight: sdk.NewDec(10),
				TakeRate:     sdk.NewDec(0),
				TotalTokens:  sdk.ZeroInt(),
			},
		},
	})

	// Set tax and rewards to be zero for easier calculation
	distParams := app.DistrKeeper.GetParams(ctx)
	distParams.CommunityTax = sdk.ZeroDec()
	distParams.BaseProposerReward = sdk.ZeroDec()
	distParams.BonusProposerReward = sdk.ZeroDec()
	app.DistrKeeper.SetParams(ctx, distParams)

	// Accounts
	addrs := test_helpers.AddTestAddrsIncremental(app, ctx, 4, sdk.NewCoins(
		sdk.NewCoin(AllianceDenom, sdk.NewInt(20_000_000)),
		sdk.NewCoin(AllianceDenomTwo, sdk.NewInt(20_000_000)),
	))
	pks := test_helpers.CreateTestPubKeys(2)

	// Creating two validators: 1 with 0% commission, 1 with 100% commission
	valAddr1 := sdk.ValAddress(addrs[0])
	_val1 := teststaking.NewValidator(t, valAddr1, pks[0])
	_val1.Commission = stakingtypes.Commission{
		CommissionRates: stakingtypes.CommissionRates{
			Rate:          sdk.NewDec(0),
			MaxRate:       sdk.NewDec(0),
			MaxChangeRate: sdk.NewDec(0),
		},
		UpdateTime: time.Now(),
	}
	test_helpers.RegisterNewValidator(t, app, ctx, _val1)
	val1, err := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr1)
	require.NoError(t, err)

	valAddr2 := sdk.ValAddress(addrs[1])
	_val2 := teststaking.NewValidator(t, valAddr2, pks[1])
	_val2.Commission = stakingtypes.Commission{
		CommissionRates: stakingtypes.CommissionRates{
			Rate:          sdk.NewDec(1),
			MaxRate:       sdk.NewDec(1),
			MaxChangeRate: sdk.NewDec(0),
		},
		UpdateTime: time.Now(),
	}
	test_helpers.RegisterNewValidator(t, app, ctx, _val2)
	val2, err := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr2)
	require.NoError(t, err)

	user1 := addrs[2]
	user2 := addrs[3]

	// Users add delegations
	_, err = app.AllianceKeeper.Delegate(ctx, user1, val1, sdk.NewCoin(AllianceDenom, sdk.NewInt(10_000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user1, val2, sdk.NewCoin(AllianceDenom, sdk.NewInt(10_000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user2, val2, sdk.NewCoin(AllianceDenom, sdk.NewInt(10_000_000)))
	require.NoError(t, err)

	_, err = app.AllianceKeeper.Delegate(ctx, user1, val1, sdk.NewCoin(AllianceDenomTwo, sdk.NewInt(10_000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user1, val2, sdk.NewCoin(AllianceDenomTwo, sdk.NewInt(10_000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user2, val1, sdk.NewCoin(AllianceDenomTwo, sdk.NewInt(10_000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user2, val2, sdk.NewCoin(AllianceDenomTwo, sdk.NewInt(10_000_000)))
	require.NoError(t, err)

	assets := app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)
	require.Equal(t, sdk.NewInt(12_999_999), app.StakingKeeper.TotalBondedTokens(ctx))

	_, err = app.AllianceKeeper.Undelegate(ctx, user1, val1, sdk.NewCoin(AllianceDenom, sdk.NewInt(10_000_000)))
	require.NoError(t, err)

	// Expect to have undelegation index saved
	undelegationIndexIter := app.AllianceKeeper.IterateUndelegationsBySrcValidator(ctx, valAddr1)
	require.True(t, undelegationIndexIter.Valid())

	assets = app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)
	require.Equal(t, sdk.NewInt(13_000_000), app.StakingKeeper.TotalBondedTokens(ctx))

	// Now we slash val 1
	val1, _ = app.AllianceKeeper.GetAllianceValidator(ctx, valAddr1)
	valPower1 := val1.GetConsensusPower(app.StakingKeeper.PowerReduction(ctx))
	valConAddr1, _ := val1.GetConsAddr()
	slashFraction := app.SlashingKeeper.SlashFractionDoubleSign(ctx)
	app.SlashingKeeper.Slash(ctx, valConAddr1, slashFraction, valPower1, 1)

	// Expect something to be slashed from undelegation entry
	undelegationsIter := app.AllianceKeeper.IterateUndelegationsByCompletionTime(ctx, ctx.BlockTime().Add(app.StakingKeeper.UnbondingTime(ctx)).Add(time.Second))
	require.True(t, undelegationsIter.Valid())
	var undelegations types.QueuedUndelegation
	app.AppCodec().MustUnmarshal(undelegationsIter.Value(), &undelegations)
	require.Equal(t, 1, len(undelegations.Entries))
	entry := undelegations.Entries[0]
	require.Greater(t, sdk.NewInt(10_000_000).Int64(), entry.Balance.Amount.Int64())

	// Move time to after undelegation completes
	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(app.StakingKeeper.UnbondingTime(ctx)).Add(time.Second))

	// Now we slash val 1
	_, err = app.AllianceKeeper.GetAllianceValidator(ctx, valAddr1)
	require.NoError(t, err)
	app.SlashingKeeper.Slash(ctx, valConAddr1, slashFraction, valPower1, 1)

	// Expect that delegation stayed the same
	undelegationsIter = app.AllianceKeeper.IterateUndelegationsByCompletionTime(ctx, ctx.BlockTime())
	require.True(t, undelegationsIter.Valid())
	var newUndelegations types.QueuedUndelegation
	app.AppCodec().MustUnmarshal(undelegationsIter.Value(), &newUndelegations)
	require.Equal(t, 1, len(newUndelegations.Entries))
	entry2 := newUndelegations.Entries[0]
	require.Equal(t, entry.Balance.Amount.Int64(), entry2.Balance.Amount.Int64())

	_, stop := alliance.RunAllInvariants(ctx, app.AllianceKeeper)
	require.False(t, stop)
}

func TestSlashingIncorrectAmount(t *testing.T) {
	// SETUP
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			{
				Denom:        AllianceDenom,
				RewardWeight: sdk.NewDec(2),
				TakeRate:     sdk.NewDec(0),
				TotalTokens:  sdk.ZeroInt(),
			},
		},
	})

	// Create and register the validator
	addrs := test_helpers.AddTestAddrsIncremental(app, ctx, 4, sdk.NewCoins(
		sdk.NewCoin(AllianceDenom, sdk.NewInt(20_000_000)),
	))
	pks := test_helpers.CreateTestPubKeys(2)
	valAddr1 := sdk.ValAddress(addrs[0])

	_val1 := teststaking.NewValidator(t, valAddr1, pks[0])
	_val1.Commission = stakingtypes.Commission{
		CommissionRates: stakingtypes.CommissionRates{
			Rate:          sdk.NewDec(0),
			MaxRate:       sdk.NewDec(0),
			MaxChangeRate: sdk.NewDec(0),
		},
		UpdateTime: time.Now(),
	}
	test_helpers.RegisterNewValidator(t, app, ctx, _val1)

	// Slash validator with incorrect amounts
	err := app.AllianceKeeper.SlashValidator(ctx, sdk.ValAddress(addrs[0]), sdk.NewDec(2))
	require.EqualErrorf(t, err, "slashed fraction must be greater than 0 and less than or equal to 1: 2.000000000000000000", "")

	err = app.AllianceKeeper.SlashValidator(ctx, sdk.ValAddress(addrs[0]), sdk.NewDec(-1))
	require.EqualErrorf(t, err, "slashed fraction must be greater than 0 and less than or equal to 1: -1.000000000000000000", "")
}
