package tests_test

import (
	"cosmossdk.io/math"
	"testing"
	"time"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	test_helpers "github.com/terra-money/alliance/app"
	"github.com/terra-money/alliance/x/alliance"
	"github.com/terra-money/alliance/x/alliance/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	teststaking "github.com/cosmos/cosmos-sdk/x/staking/testutil"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
)

func TestRebalancingAfterRewardsRateChange(t *testing.T) {
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			types.NewAllianceAsset(AllianceDenom, math.LegacyNewDec(2), math.LegacyZeroDec(), math.LegacyNewDec(4), math.LegacyZeroDec(), startTime),
			types.NewAllianceAsset(AllianceDenomTwo, math.LegacyNewDec(10), math.LegacyNewDec(2), math.LegacyNewDec(12), math.LegacyZeroDec(), startTime),
		},
	})

	// Set tax and rewards to be zero for easier calculation
	distParams, err := app.DistrKeeper.Params.Get(ctx)
	require.NoError(t, err)
	distParams.CommunityTax = math.LegacyZeroDec()
	err = app.DistrKeeper.Params.Set(ctx, distParams)
	require.NoError(t, err)

	// Accounts
	addrs := test_helpers.AddTestAddrsIncremental(app, ctx, 4, sdk.NewCoins(
		sdk.NewCoin(AllianceDenom, math.NewInt(1000_000)),
		sdk.NewCoin(AllianceDenomTwo, math.NewInt(1000_000)),
	))
	pks := test_helpers.CreateTestPubKeys(2)
	powerReduction := app.StakingKeeper.PowerReduction(ctx)

	valAddr1 := sdk.ValAddress(addrs[0])
	_val1 := teststaking.NewValidator(t, valAddr1, pks[0])
	_val1.Commission = stakingtypes.Commission{
		CommissionRates: stakingtypes.CommissionRates{
			Rate:          math.LegacyNewDec(0),
			MaxRate:       math.LegacyNewDec(0),
			MaxChangeRate: math.LegacyNewDec(0),
		},
		UpdateTime: time.Now(),
	}
	_val1.Status = stakingtypes.Bonded
	test_helpers.RegisterNewValidator(t, app, ctx, _val1)
	val1, err := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr1)
	require.NoError(t, err)

	user1 := addrs[2]

	// Start by delegating
	_, err = app.AllianceKeeper.Delegate(ctx, user1, val1, sdk.NewCoin(AllianceDenom, math.NewInt(1000_000)))
	require.NoError(t, err)

	assets := app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)
	totalBonded, err := app.StakingKeeper.TotalBondedTokens(ctx)
	require.NoError(t, err)
	require.Equal(t, math.NewInt(3_000_000), totalBonded)

	// Expecting voting power for the alliance module
	val, err := app.StakingKeeper.GetValidator(ctx, valAddr1)
	require.NoError(t, err)
	require.Equal(t, int64(2), val.ConsensusPower(powerReduction))

	// Update but did not change reward weight
	err = app.AllianceKeeper.UpdateAllianceAsset(ctx, types.AllianceAsset{
		Denom:                AllianceDenom,
		RewardWeight:         math.LegacyNewDec(2),
		RewardWeightRange:    types.RewardWeightRange{Min: math.LegacyNewDec(0), Max: math.LegacyNewDec(5)},
		TakeRate:             math.LegacyNewDec(10),
		RewardChangeRate:     math.LegacyNewDec(0),
		RewardChangeInterval: 0,
	})
	require.NoError(t, err)

	// Expect no snapshots to be created
	iter, err := app.AllianceKeeper.IterateWeightChangeSnapshot(ctx, AllianceDenom, getOperator(val), 0)
	require.False(t, iter.Valid())

	err = app.AllianceKeeper.UpdateAllianceAsset(ctx, types.AllianceAsset{
		Denom:                AllianceDenom,
		RewardWeight:         math.LegacyNewDec(20),
		RewardWeightRange:    types.RewardWeightRange{Min: math.LegacyNewDec(5), Max: math.LegacyNewDec(25)},
		TakeRate:             math.LegacyNewDec(0),
		RewardChangeRate:     math.LegacyNewDec(0),
		RewardChangeInterval: 0,
	})
	require.NoError(t, err)

	// Expect a snapshot to be created
	iter, err = app.AllianceKeeper.IterateWeightChangeSnapshot(ctx, AllianceDenom, getOperator(val), 0)
	require.True(t, iter.Valid())

	assets = app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)
	totalBonded, err = app.StakingKeeper.TotalBondedTokens(ctx)
	require.NoError(t, err)
	require.Equal(t, math.NewInt(21_000_000), totalBonded)

	// Expecting voting power to increase
	val, err = app.StakingKeeper.GetValidator(ctx, valAddr1)
	require.NoError(t, err)
	require.Equal(t, int64(20), val.ConsensusPower(powerReduction))

	err = app.AllianceKeeper.UpdateAllianceAsset(ctx, types.AllianceAsset{
		Denom:                AllianceDenom,
		RewardWeight:         math.LegacyNewDec(1),
		RewardWeightRange:    types.RewardWeightRange{Min: math.LegacyNewDec(0), Max: math.LegacyNewDec(5)},
		TakeRate:             math.LegacyNewDec(0),
		RewardChangeRate:     math.LegacyNewDec(0),
		RewardChangeInterval: 0,
	})
	require.NoError(t, err)

	assets = app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)
	totalBonded, err = app.StakingKeeper.TotalBondedTokens(ctx)
	require.NoError(t, err)
	require.Equal(t, math.NewInt(2_000_000), totalBonded)

	// Expecting voting power to decrease
	val, err = app.StakingKeeper.GetValidator(ctx, valAddr1)
	require.NoError(t, err)
	require.Equal(t, int64(1), val.ConsensusPower(powerReduction))

	_, stop := alliance.RunAllInvariants(ctx, app.AllianceKeeper)
	require.False(t, stop)
}

func TestRebalancingWithUnbondedValidator(t *testing.T) {
	var err error
	app, ctx := createTestContext(t)
	bondDenom, err := app.StakingKeeper.BondDenom(ctx)
	require.NoError(t, err)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			{
				Denom:        AllianceDenom,
				RewardWeight: math.LegacyMustNewDecFromStr("0.1"),
				TakeRate:     math.LegacyNewDec(0),
				TotalTokens:  math.ZeroInt(),
			},
			{
				Denom:        AllianceDenomTwo,
				RewardWeight: math.LegacyMustNewDecFromStr("0.5"),
				TakeRate:     math.LegacyNewDec(0),
				TotalTokens:  math.ZeroInt(),
			},
		},
	})

	// Set tax and rewards to be zero for easier calculation
	distParams, err := app.DistrKeeper.Params.Get(ctx)
	require.NoError(t, err)
	distParams.CommunityTax = math.LegacyZeroDec()
	err = app.DistrKeeper.Params.Set(ctx, distParams)
	require.NoError(t, err)

	// Accounts
	addrs := test_helpers.AddTestAddrsIncremental(app, ctx, 5, sdk.NewCoins(
		sdk.NewCoin(bondDenom, math.NewInt(10_000_000)),
		sdk.NewCoin(AllianceDenom, math.NewInt(50_000_000)),
		sdk.NewCoin(AllianceDenomTwo, math.NewInt(50_000_000)),
	))
	pks := test_helpers.CreateTestPubKeys(2)

	// Increase the stake on genesis validator
	delegations, err := app.StakingKeeper.GetAllDelegations(ctx)
	require.NoError(t, err)
	require.Len(t, delegations, 1)
	valAddr0, err := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	require.NoError(t, err)
	val0, _ := app.StakingKeeper.GetValidator(ctx, valAddr0)
	_, err = app.StakingKeeper.Delegate(ctx, addrs[4], math.NewInt(9_000_000), stakingtypes.Unbonded, val0, true)
	require.NoError(t, err)

	// Creating two validators: 1 with 0% commission, 1 with 100% commission
	valAddr1 := sdk.ValAddress(addrs[0])
	_val1 := teststaking.NewValidator(t, valAddr1, pks[0])
	_val1.Commission = stakingtypes.Commission{
		CommissionRates: stakingtypes.CommissionRates{
			Rate:          math.LegacyNewDec(0),
			MaxRate:       math.LegacyNewDec(0),
			MaxChangeRate: math.LegacyNewDec(0),
		},
		UpdateTime: time.Now(),
	}
	_val1.Description.Moniker = "val1" //nolint:goconst
	test_helpers.RegisterNewValidator(t, app, ctx, _val1)
	val1, err := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr1)
	require.NoError(t, err)

	valAddr2 := sdk.ValAddress(addrs[1])
	_val2 := teststaking.NewValidator(t, valAddr2, pks[1])
	_val2.Commission = stakingtypes.Commission{
		CommissionRates: stakingtypes.CommissionRates{
			Rate:          math.LegacyNewDec(1),
			MaxRate:       math.LegacyNewDec(1),
			MaxChangeRate: math.LegacyNewDec(0),
		},
		UpdateTime: time.Now(),
	}
	_val2.Description.Moniker = "val2" //nolint:goconst
	test_helpers.RegisterNewValidator(t, app, ctx, _val2)
	val2, err := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr2)
	require.NoError(t, err)

	user1 := addrs[2]
	user2 := addrs[3]

	// Users add delegations
	_, err = app.AllianceKeeper.Delegate(ctx, user1, val1, sdk.NewCoin(AllianceDenom, math.NewInt(20_000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user1, val2, sdk.NewCoin(AllianceDenom, math.NewInt(10_000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user2, val1, sdk.NewCoin(AllianceDenom, math.NewInt(10_000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user2, val2, sdk.NewCoin(AllianceDenom, math.NewInt(10_000_000)))
	require.NoError(t, err)

	_, err = app.AllianceKeeper.Delegate(ctx, user1, val1, sdk.NewCoin(AllianceDenomTwo, math.NewInt(10_000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user1, val2, sdk.NewCoin(AllianceDenomTwo, math.NewInt(10_000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user2, val1, sdk.NewCoin(AllianceDenomTwo, math.NewInt(10_000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user2, val2, sdk.NewCoin(AllianceDenomTwo, math.NewInt(10_000_000)))
	require.NoError(t, err)
	require.NoError(t, err)

	assets := app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)
	_, err = app.StakingKeeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	totalBonded, err := app.StakingKeeper.TotalBondedTokens(ctx)
	require.NoError(t, err)
	require.Equal(t, math.NewInt(16_000_000), totalBonded)

	val1, _ = app.AllianceKeeper.GetAllianceValidator(ctx, valAddr1)
	val2, _ = app.AllianceKeeper.GetAllianceValidator(ctx, valAddr2)
	require.Greater(t, val1.Tokens.Int64(), val2.Tokens.Int64())

	// Set max validators to be 2 to trigger unbonding
	params, err := app.StakingKeeper.GetParams(ctx)
	require.NoError(t, err)
	params.MaxValidators = 2
	err = app.StakingKeeper.SetParams(ctx, params)
	require.NoError(t, err)
	_, err = app.StakingKeeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	totalBonded, err = app.StakingKeeper.TotalBondedTokens(ctx)
	require.NoError(t, err)
	require.Equal(t, math.NewInt(13_100_000), totalBonded)

	vals, err := app.StakingKeeper.GetBondedValidatorsByPower(ctx)
	require.NoError(t, err)
	require.Equal(t, 2, len(vals))
	require.Equal(t, "val1", vals[1].GetMoniker())

	assets = app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)
	totalBonded, err = app.StakingKeeper.TotalBondedTokens(ctx)
	require.NoError(t, err)
	require.Equal(t, math.NewInt(16_000_000), totalBonded)

	_, err = app.AllianceKeeper.GetAllianceValidator(ctx, valAddr2)
	require.NoError(t, err)

	// Set max validators to be 3 to trigger rebonding
	params, err = app.StakingKeeper.GetParams(ctx)
	require.NoError(t, err)
	params.MaxValidators = 3
	err = app.StakingKeeper.SetParams(ctx, params)
	require.NoError(t, err)
	_, err = app.StakingKeeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	totalBonded, err = app.StakingKeeper.TotalBondedTokens(ctx)
	require.NoError(t, err)
	require.Equal(t, math.NewInt(18_900_000), totalBonded)

	assets = app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)
	totalBonded, err = app.StakingKeeper.TotalBondedTokens(ctx)
	require.NoError(t, err)
	require.Equal(t, math.NewInt(16_000_000), totalBonded)

	_, stop := alliance.RunAllInvariants(ctx, app.AllianceKeeper)
	require.False(t, stop)
}

func TestRebalancingWithJailedValidator(t *testing.T) {
	var err error
	app, ctx := createTestContext(t)
	bondDenom, err := app.StakingKeeper.BondDenom(ctx)
	require.NoError(t, err)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			{
				Denom:        AllianceDenom,
				RewardWeight: math.LegacyMustNewDecFromStr("0.1"),
				TakeRate:     math.LegacyNewDec(0),
				TotalTokens:  math.ZeroInt(),
			},
			{
				Denom:        AllianceDenomTwo,
				RewardWeight: math.LegacyMustNewDecFromStr("0.5"),
				TakeRate:     math.LegacyNewDec(0),
				TotalTokens:  math.ZeroInt(),
			},
		},
	})

	// Set tax and rewards to be zero for easier calculation
	distParams, err := app.DistrKeeper.Params.Get(ctx)
	require.NoError(t, err)
	distParams.CommunityTax = math.LegacyZeroDec()
	err = app.DistrKeeper.Params.Set(ctx, distParams)
	require.NoError(t, err)

	// Accounts
	addrs := test_helpers.AddTestAddrsIncremental(app, ctx, 5, sdk.NewCoins(
		sdk.NewCoin(bondDenom, math.NewInt(10_000_000)),
		sdk.NewCoin(AllianceDenom, math.NewInt(50_000_000)),
		sdk.NewCoin(AllianceDenomTwo, math.NewInt(50_000_000)),
	))
	pks := test_helpers.CreateTestPubKeys(2)

	// Increase the stake on genesis validator
	delegations, err := app.StakingKeeper.GetAllDelegations(ctx)
	require.NoError(t, err)
	require.Len(t, delegations, 1)
	valAddr0, err := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	require.NoError(t, err)
	val0, _ := app.StakingKeeper.GetValidator(ctx, valAddr0)
	_, err = app.StakingKeeper.Delegate(ctx, addrs[4], math.NewInt(9_000_000), stakingtypes.Unbonded, val0, true)
	require.NoError(t, err)

	// Creating two validators: 1 with 0% commission, 1 with 100% commission
	valAddr1 := sdk.ValAddress(addrs[0])
	_val1 := teststaking.NewValidator(t, valAddr1, pks[0])
	_val1.Commission = stakingtypes.Commission{
		CommissionRates: stakingtypes.CommissionRates{
			Rate:          math.LegacyNewDec(0),
			MaxRate:       math.LegacyNewDec(0),
			MaxChangeRate: math.LegacyNewDec(0),
		},
		UpdateTime: time.Now(),
	}
	_val1.Description.Moniker = "val1"
	test_helpers.RegisterNewValidator(t, app, ctx, _val1)
	_, err = app.StakingKeeper.Delegate(ctx, addrs[0], math.NewInt(1_000_000), stakingtypes.Unbonded, _val1, true)
	require.NoError(t, err)
	val1, err := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr1)
	require.NoError(t, err)

	valAddr2 := sdk.ValAddress(addrs[1])
	_val2 := teststaking.NewValidator(t, valAddr2, pks[1])
	_val2.Commission = stakingtypes.Commission{
		CommissionRates: stakingtypes.CommissionRates{
			Rate:          math.LegacyNewDec(1),
			MaxRate:       math.LegacyNewDec(1),
			MaxChangeRate: math.LegacyNewDec(0),
		},
		UpdateTime: time.Now(),
	}
	_val2.Description.Moniker = "val2"
	test_helpers.RegisterNewValidator(t, app, ctx, _val2)
	_, err = app.StakingKeeper.Delegate(ctx, addrs[1], math.NewInt(1_000_000), stakingtypes.Unbonded, _val2, true)
	require.NoError(t, err)
	val2, err := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr2)
	require.NoError(t, err)

	user1 := addrs[2]
	user2 := addrs[3]

	// Users add delegations
	_, err = app.AllianceKeeper.Delegate(ctx, user1, val1, sdk.NewCoin(AllianceDenom, math.NewInt(20_000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user1, val2, sdk.NewCoin(AllianceDenom, math.NewInt(10_000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user2, val1, sdk.NewCoin(AllianceDenom, math.NewInt(10_000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user2, val2, sdk.NewCoin(AllianceDenom, math.NewInt(10_000_000)))
	require.NoError(t, err)

	_, err = app.AllianceKeeper.Delegate(ctx, user1, val1, sdk.NewCoin(AllianceDenomTwo, math.NewInt(10_000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user1, val2, sdk.NewCoin(AllianceDenomTwo, math.NewInt(10_000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user2, val1, sdk.NewCoin(AllianceDenomTwo, math.NewInt(10_000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user2, val2, sdk.NewCoin(AllianceDenomTwo, math.NewInt(10_000_000)))
	require.NoError(t, err)
	require.NoError(t, err)

	_, err = app.StakingKeeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	assets := app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)
	totalBonded, err := app.StakingKeeper.TotalBondedTokens(ctx)
	require.NoError(t, err)
	// 12 * 1.6 = 19.2
	require.Equal(t, math.NewInt(19_200_000), totalBonded)

	val1, _ = app.AllianceKeeper.GetAllianceValidator(ctx, valAddr1)
	val2, _ = app.AllianceKeeper.GetAllianceValidator(ctx, valAddr2)
	require.Greater(t, val1.Tokens.Int64(), val2.Tokens.Int64())

	// Jail validator
	cons2, _ := val2.GetConsAddr()
	app.SlashingKeeper.Jail(ctx, cons2)
	_, err = app.StakingKeeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	totalBonded, err = app.StakingKeeper.TotalBondedTokens(ctx)
	require.NoError(t, err)
	require.Equal(t, math.NewInt(14_720_000), totalBonded)

	vals, err := app.StakingKeeper.GetBondedValidatorsByPower(ctx)
	require.NoError(t, err)
	require.Equal(t, 2, len(vals))
	require.Equal(t, "val1", vals[1].GetMoniker())

	assets = app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)
	totalBonded, err = app.StakingKeeper.TotalBondedTokens(ctx)
	require.NoError(t, err)
	// 11 * 1.6 = 17.6
	require.Equal(t, math.NewInt(17_600_000), totalBonded)

	_, err = app.AllianceKeeper.GetAllianceValidator(ctx, valAddr2)
	require.NoError(t, err)

	// Unjail validator
	err = app.SlashingKeeper.Unjail(ctx, valAddr2)
	require.NoError(t, err)
	_, err = app.StakingKeeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	totalBonded, err = app.StakingKeeper.TotalBondedTokens(ctx)
	require.NoError(t, err)
	require.Equal(t, math.NewInt(22_080_000), totalBonded)

	assets = app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)
	totalBonded, err = app.StakingKeeper.TotalBondedTokens(ctx)
	require.NoError(t, err)
	require.Equal(t, math.NewInt(19_200_000), totalBonded)

	_, stop := alliance.RunAllInvariants(ctx, app.AllianceKeeper)
	require.False(t, stop)
}

func TestRebalancingWithDelayedRewardsStartTime(t *testing.T) {
	var err error
	app, ctx := createTestContext(t)
	bondDenom, err := app.StakingKeeper.BondDenom(ctx)
	require.NoError(t, err)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)

	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			types.NewAllianceAsset(AllianceDenom, math.LegacyMustNewDecFromStr("0.5"), math.LegacyZeroDec(), math.LegacyOneDec(), math.LegacyMustNewDecFromStr("0.1"), startTime.Add(time.Hour*24)),
			types.NewAllianceAsset(AllianceDenomTwo, math.LegacyMustNewDecFromStr("0.2"), math.LegacyZeroDec(), math.LegacyOneDec(), math.LegacyMustNewDecFromStr("0.1"), startTime.Add(time.Hour*24*2)),
		},
	})

	// Set tax and rewards to be zero for easier calculation
	distParams, err := app.DistrKeeper.Params.Get(ctx)
	require.NoError(t, err)
	distParams.CommunityTax = math.LegacyZeroDec()
	err = app.DistrKeeper.Params.Set(ctx, distParams)
	require.NoError(t, err)

	// Accounts
	addrs := test_helpers.AddTestAddrsIncremental(app, ctx, 5, sdk.NewCoins(
		sdk.NewCoin(bondDenom, math.NewInt(10_000_000)),
		sdk.NewCoin(AllianceDenom, math.NewInt(50_000_000)),
		sdk.NewCoin(AllianceDenomTwo, math.NewInt(50_000_000)),
	))
	pks := test_helpers.CreateTestPubKeys(2)

	// Increase the stake on genesis validator
	delegations, err := app.StakingKeeper.GetAllDelegations(ctx)
	require.NoError(t, err)
	require.Len(t, delegations, 1)
	valAddr0, err := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	require.NoError(t, err)
	val0, _ := app.StakingKeeper.GetValidator(ctx, valAddr0)
	_, err = app.StakingKeeper.Delegate(ctx, addrs[4], math.NewInt(9_000_000), stakingtypes.Unbonded, val0, true)
	require.NoError(t, err)

	// Creating two validators: 1 with 0% commission, 1 with 100% commission
	valAddr1 := sdk.ValAddress(addrs[0])
	_val1 := teststaking.NewValidator(t, valAddr1, pks[0])
	_val1.Commission = stakingtypes.Commission{
		CommissionRates: stakingtypes.CommissionRates{
			Rate:          math.LegacyNewDec(0),
			MaxRate:       math.LegacyNewDec(0),
			MaxChangeRate: math.LegacyNewDec(0),
		},
		UpdateTime: time.Now(),
	}
	_val1.Description.Moniker = "val1"
	test_helpers.RegisterNewValidator(t, app, ctx, _val1)
	_, err = app.StakingKeeper.Delegate(ctx, addrs[0], math.NewInt(1_000_000), stakingtypes.Unbonded, _val1, true)
	require.NoError(t, err)
	val1, err := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr1)
	require.NoError(t, err)

	valAddr2 := sdk.ValAddress(addrs[1])
	_val2 := teststaking.NewValidator(t, valAddr2, pks[1])
	_val2.Commission = stakingtypes.Commission{
		CommissionRates: stakingtypes.CommissionRates{
			Rate:          math.LegacyNewDec(1),
			MaxRate:       math.LegacyNewDec(1),
			MaxChangeRate: math.LegacyNewDec(0),
		},
		UpdateTime: time.Now(),
	}
	_val2.Description.Moniker = "val2"
	test_helpers.RegisterNewValidator(t, app, ctx, _val2)
	_, err = app.StakingKeeper.Delegate(ctx, addrs[1], math.NewInt(1_000_000), stakingtypes.Unbonded, _val2, true)
	require.NoError(t, err)
	val2, err := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr2)
	require.NoError(t, err)

	user1 := addrs[2]
	user2 := addrs[3]

	// Users add delegations
	_, err = app.AllianceKeeper.Delegate(ctx, user1, val1, sdk.NewCoin(AllianceDenom, math.NewInt(20_000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user1, val2, sdk.NewCoin(AllianceDenom, math.NewInt(10_000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user2, val1, sdk.NewCoin(AllianceDenom, math.NewInt(10_000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user2, val2, sdk.NewCoin(AllianceDenom, math.NewInt(10_000_000)))
	require.NoError(t, err)

	_, err = app.AllianceKeeper.Delegate(ctx, user1, val1, sdk.NewCoin(AllianceDenomTwo, math.NewInt(10_000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user1, val2, sdk.NewCoin(AllianceDenomTwo, math.NewInt(10_000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user2, val1, sdk.NewCoin(AllianceDenomTwo, math.NewInt(10_000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user2, val2, sdk.NewCoin(AllianceDenomTwo, math.NewInt(10_000_000)))
	require.NoError(t, err)
	_, err = app.StakingKeeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)

	// Expect that rewards rates are not updated due to ctx being before rewards start time
	require.NoError(t, err)
	assets := app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)
	totalBonded, err := app.StakingKeeper.TotalBondedTokens(ctx)
	require.NoError(t, err)
	require.Equal(t, math.NewInt(12_000_000), totalBonded)

	// Expect that rewards rates are updated only for alliance 1
	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(time.Hour * 24))
	assets = app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)
	totalBonded, err = app.StakingKeeper.TotalBondedTokens(ctx)
	require.NoError(t, err)
	// 12 * 1.5 = 18
	require.Equal(t, math.NewInt(18_000_000), totalBonded)

	// Expect that rewards rates are updated all alliances
	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(time.Hour * 48))
	assets = app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)
	totalBonded, err = app.StakingKeeper.TotalBondedTokens(ctx)
	require.NoError(t, err)
	// 12 * 1.7 = 18
	require.Equal(t, math.NewInt(20_400_000), totalBonded)

	_, stop := alliance.RunAllInvariants(ctx, app.AllianceKeeper)
	require.False(t, stop)
}

func TestConsumingRebalancingEvent(t *testing.T) {
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)

	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			types.NewAllianceAsset(AllianceDenom, math.LegacyMustNewDecFromStr("0.5"), math.LegacyZeroDec(), math.LegacyOneDec(), math.LegacyMustNewDecFromStr("0.1"), startTime.Add(time.Hour*24)),
			types.NewAllianceAsset(AllianceDenomTwo, math.LegacyMustNewDecFromStr("0.2"), math.LegacyZeroDec(), math.LegacyOneDec(), math.LegacyMustNewDecFromStr("0.1"), startTime.Add(time.Hour*24*2)),
		},
	})

	app.AllianceKeeper.QueueAssetRebalanceEvent(ctx)
	store := app.AllianceKeeper.StoreService().OpenKVStore(ctx)
	key := types.AssetRebalanceQueueKey
	b, err := store.Get(key)
	require.NoError(t, err)
	require.NotNil(t, b)

	require.True(t, app.AllianceKeeper.ConsumeAssetRebalanceEvent(ctx))
	b, err = store.Get(key)
	require.NoError(t, err)
	require.Nil(t, b)

	require.False(t, app.AllianceKeeper.ConsumeAssetRebalanceEvent(ctx))
}

func TestRewardRangeWithChangeRateOverTime(t *testing.T) {
	app, ctx := createTestContext(t)
	decayInterval := time.Hour * 24 * 30
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			{
				Denom:                AllianceDenom,
				RewardWeight:         math.LegacyMustNewDecFromStr("0.075"),
				RewardWeightRange:    types.RewardWeightRange{Min: math.LegacyMustNewDecFromStr("0.05"), Max: math.LegacyMustNewDecFromStr("0.10")},
				TakeRate:             math.LegacyNewDec(0),
				TotalTokens:          math.ZeroInt(),
				RewardChangeRate:     math.LegacyMustNewDecFromStr("1.5"),
				RewardChangeInterval: decayInterval,
			},
			{
				Denom:                AllianceDenomTwo,
				RewardWeight:         math.LegacyMustNewDecFromStr("0.075"),
				RewardWeightRange:    types.RewardWeightRange{Min: math.LegacyMustNewDecFromStr("0.05"), Max: math.LegacyMustNewDecFromStr("0.10")},
				TakeRate:             math.LegacyNewDec(0),
				TotalTokens:          math.ZeroInt(),
				RewardChangeRate:     math.LegacyMustNewDecFromStr("0.5"),
				RewardChangeInterval: decayInterval,
			},
		},
	})
	asset, _ := app.AllianceKeeper.GetAssetByDenom(ctx, AllianceDenom)
	assets := app.AllianceKeeper.GetAllAssets(ctx)
	// Running the decay hook now should do nothing
	err := app.AllianceKeeper.RewardWeightChangeHook(ctx, assets)
	require.NoError(t, err)

	// Move block time to after change interval + one year
	ctx = ctx.WithBlockTime(asset.RewardStartTime.Add(decayInterval * 2))

	// Running the decay hook should update reward weight
	err = app.AllianceKeeper.RewardWeightChangeHook(ctx, assets)
	require.NoError(t, err)

	updatedAsset, _ := app.AllianceKeeper.GetAssetByDenom(ctx, AllianceDenom)
	require.Equal(t, types.AllianceAsset{
		Denom:                AllianceDenom,
		RewardWeight:         math.LegacyMustNewDecFromStr("0.10"),
		RewardWeightRange:    types.RewardWeightRange{Min: math.LegacyMustNewDecFromStr("0.05"), Max: math.LegacyMustNewDecFromStr("0.10")},
		TakeRate:             asset.TakeRate,
		TotalTokens:          asset.TotalTokens,
		TotalValidatorShares: asset.TotalValidatorShares,
		RewardStartTime:      asset.RewardStartTime,
		RewardChangeRate:     asset.RewardChangeRate,
		RewardChangeInterval: asset.RewardChangeInterval,
		LastRewardChangeTime: ctx.BlockTime(),
	}, updatedAsset)

	updatedAsset, _ = app.AllianceKeeper.GetAssetByDenom(ctx, AllianceDenomTwo)
	asset, _ = app.AllianceKeeper.GetAssetByDenom(ctx, AllianceDenomTwo)
	require.Equal(t, types.AllianceAsset{
		Denom:                AllianceDenomTwo,
		RewardWeight:         math.LegacyMustNewDecFromStr("0.05"),
		RewardWeightRange:    types.RewardWeightRange{Min: math.LegacyMustNewDecFromStr("0.05"), Max: math.LegacyMustNewDecFromStr("0.10")},
		TakeRate:             asset.TakeRate,
		TotalTokens:          asset.TotalTokens,
		TotalValidatorShares: asset.TotalValidatorShares,
		RewardStartTime:      asset.RewardStartTime,
		RewardChangeRate:     asset.RewardChangeRate,
		RewardChangeInterval: asset.RewardChangeInterval,
		LastRewardChangeTime: ctx.BlockTime(),
	}, updatedAsset)
}

func TestRewardWeightDecay(t *testing.T) {
	var err error
	app, ctx := createTestContext(t)
	bondDenom, err := app.StakingKeeper.BondDenom(ctx)
	require.NoError(t, err)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{},
	})

	// Accounts
	addrs := test_helpers.AddTestAddrsIncremental(app, ctx, 5, sdk.NewCoins(
		sdk.NewCoin(bondDenom, math.NewInt(10_000_000)),
		sdk.NewCoin(AllianceDenom, math.NewInt(50_000_000)),
		sdk.NewCoin(AllianceDenomTwo, math.NewInt(50_000_000)),
	))

	// Increase the stake on genesis validator
	delegations, err := app.StakingKeeper.GetAllDelegations(ctx)
	require.NoError(t, err)
	require.Len(t, delegations, 1)
	valAddr0, err := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	require.NoError(t, err)
	val0, _ := app.StakingKeeper.GetValidator(ctx, valAddr0)
	_, err = app.StakingKeeper.Delegate(ctx, addrs[4], math.NewInt(9_000_000), stakingtypes.Unbonded, val0, true)
	require.NoError(t, err)

	// Pass a proposal to add a new asset with a decay rate
	decayInterval := time.Hour * 24 * 30
	err = app.AllianceKeeper.CreateAlliance(ctx, &types.MsgCreateAllianceProposal{
		Title:                "",
		Description:          "",
		Denom:                AllianceDenom,
		RewardWeight:         math.LegacyNewDec(1),
		RewardWeightRange:    types.RewardWeightRange{Min: math.LegacyNewDec(0), Max: math.LegacyNewDec(5)},
		TakeRate:             math.LegacyZeroDec(),
		RewardChangeRate:     math.LegacyMustNewDecFromStr("0.5"),
		RewardChangeInterval: decayInterval,
	})
	require.NoError(t, err)

	asset, _ := app.AllianceKeeper.GetAssetByDenom(ctx, AllianceDenom)

	assets := app.AllianceKeeper.GetAllAssets(ctx)
	// Running the decay hook now should do nothing
	err = app.AllianceKeeper.RewardWeightChangeHook(ctx, assets)
	require.NoError(t, err)

	// Move block time to after change interval + one year
	ctx = ctx.WithBlockTime(asset.RewardStartTime.Add(decayInterval))

	// Running the decay hook should update reward weight
	err = app.AllianceKeeper.RewardWeightChangeHook(ctx, assets)
	require.NoError(t, err)
	updatedAsset, _ := app.AllianceKeeper.GetAssetByDenom(ctx, AllianceDenom)
	require.Equal(t, types.AllianceAsset{
		Denom:                AllianceDenom,
		RewardWeight:         math.LegacyMustNewDecFromStr("0.5"),
		RewardWeightRange:    types.RewardWeightRange{Min: math.LegacyNewDec(0), Max: math.LegacyNewDec(5)},
		TakeRate:             asset.TakeRate,
		TotalTokens:          asset.TotalTokens,
		TotalValidatorShares: asset.TotalValidatorShares,
		RewardStartTime:      asset.RewardStartTime,
		RewardChangeRate:     asset.RewardChangeRate,
		RewardChangeInterval: asset.RewardChangeInterval,
		LastRewardChangeTime: ctx.BlockTime(),
	}, updatedAsset)

	// There should be a rebalancing event stored
	require.True(t, app.AllianceKeeper.ConsumeAssetRebalanceEvent(ctx))

	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(time.Hour * 10))
	// Updating the alliance asset through proposal should queue another decay event
	err = app.AllianceKeeper.UpdateAlliance(ctx, &types.MsgUpdateAllianceProposal{
		Title:                "",
		Description:          "",
		Denom:                AllianceDenom,
		RewardWeight:         math.LegacyMustNewDecFromStr("0.5"),
		RewardWeightRange:    types.RewardWeightRange{Min: math.LegacyNewDec(0), Max: math.LegacyNewDec(5)},
		TakeRate:             math.LegacyZeroDec(),
		RewardChangeRate:     math.LegacyOneDec(),
		RewardChangeInterval: 0,
	})
	require.NoError(t, err)

	// Updating alliance asset again with a non-zero decay
	err = app.AllianceKeeper.UpdateAlliance(ctx, &types.MsgUpdateAllianceProposal{
		Title:                "",
		Description:          "",
		Denom:                AllianceDenom,
		RewardWeight:         math.LegacyMustNewDecFromStr("0.5"),
		RewardWeightRange:    types.RewardWeightRange{Min: math.LegacyNewDec(0), Max: math.LegacyNewDec(5)},
		TakeRate:             math.LegacyZeroDec(),
		RewardChangeRate:     math.LegacyMustNewDecFromStr("0.1"),
		RewardChangeInterval: decayInterval,
	})
	require.NoError(t, err)

	// Add a new asset with an initial 0 decay
	err = app.AllianceKeeper.CreateAlliance(ctx, &types.MsgCreateAllianceProposal{
		Title:                "",
		Description:          "",
		Denom:                AllianceDenomTwo,
		RewardWeight:         math.LegacyNewDec(1),
		RewardWeightRange:    types.RewardWeightRange{Min: math.LegacyNewDec(0), Max: math.LegacyNewDec(5)},
		TakeRate:             math.LegacyZeroDec(),
		RewardChangeRate:     math.LegacyOneDec(),
		RewardChangeInterval: decayInterval,
	})
	require.NoError(t, err)

	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(time.Hour))
	// Updating alliance asset again with a non-zero decay
	err = app.AllianceKeeper.UpdateAlliance(ctx, &types.MsgUpdateAllianceProposal{
		Title:                "",
		Description:          "",
		Denom:                AllianceDenomTwo,
		RewardWeight:         math.LegacyMustNewDecFromStr("0.5"),
		RewardWeightRange:    types.RewardWeightRange{Min: math.LegacyNewDec(0), Max: math.LegacyNewDec(5)},
		TakeRate:             math.LegacyZeroDec(),
		RewardChangeRate:     math.LegacyMustNewDecFromStr("0.1"),
		RewardChangeInterval: decayInterval,
	})
	require.NoError(t, err)
}

func TestRewardWeightDecayOverTime(t *testing.T) {
	var err error
	app, ctx := createTestContext(t)
	bondDenom, err := app.StakingKeeper.BondDenom(ctx)
	require.NoError(t, err)
	startTime := time.Now().UTC()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{},
	})

	// Accounts
	addrs := test_helpers.AddTestAddrsIncremental(app, ctx, 5, sdk.NewCoins(
		sdk.NewCoin(bondDenom, math.NewInt(1_000_000_000_000)),
		sdk.NewCoin(AllianceDenom, math.NewInt(5_000_000)),
		sdk.NewCoin(AllianceDenomTwo, math.NewInt(5_000_000)),
	))

	// Increase the stake on genesis validator
	delegations, err := app.StakingKeeper.GetAllDelegations(ctx)
	require.NoError(t, err)
	require.Len(t, delegations, 1)
	valAddr0, err := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	require.NoError(t, err)
	_val0, _ := app.StakingKeeper.GetValidator(ctx, valAddr0)
	_, err = app.StakingKeeper.Delegate(ctx, addrs[4], math.NewInt(9_000_000), stakingtypes.Unbonded, _val0, true)
	require.NoError(t, err)

	val0, _ := app.AllianceKeeper.GetAllianceValidator(ctx, getOperator(_val0))
	require.NoError(t, err)

	// Pass a proposal to add a new asset with a decay rate
	decayInterval := time.Minute
	decayRate := math.LegacyMustNewDecFromStr("0.99998")
	err = app.AllianceKeeper.CreateAlliance(ctx, &types.MsgCreateAllianceProposal{
		Title:                "",
		Description:          "",
		Denom:                AllianceDenom,
		RewardWeight:         math.LegacyNewDec(1),
		RewardWeightRange:    types.RewardWeightRange{Min: math.LegacyNewDec(0), Max: math.LegacyNewDec(5)},
		TakeRate:             math.LegacyZeroDec(),
		RewardChangeRate:     decayRate,
		RewardChangeInterval: decayInterval,
	})
	require.NoError(t, err)

	// Delegate to validator
	_, err = app.AllianceKeeper.Delegate(ctx, addrs[1], val0, sdk.NewCoin(AllianceDenom, math.NewInt(5_000_000)))
	require.NoError(t, err)
	//
	assets := app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceHook(ctx, assets)
	require.NoError(t, err)

	asset, _ := app.AllianceKeeper.GetAssetByDenom(ctx, AllianceDenom)

	// Simulate the chain running for 10 days with a defined block time
	blockTime := time.Second * 20
	for i := time.Duration(0); i <= time.Hour*24*10; i += blockTime {
		ctx = ctx.WithBlockTime(ctx.BlockTime().Add(blockTime)).WithBlockHeight(ctx.BlockHeight() + 1)
		assets = app.AllianceKeeper.GetAllAssets(ctx)
		// Running the decay hook should update reward weight
		err = app.AllianceKeeper.RewardWeightChangeHook(ctx, assets)
		require.NoError(t, err)
	}

	// time passed minus reward delay time (rewards and decay only start after the delay)
	totalDecayTime := (time.Hour * 24 * 10) - app.AllianceKeeper.RewardDelayTime(ctx)
	intervals := uint64(totalDecayTime / decayInterval)
	asset, _ = app.AllianceKeeper.GetAssetByDenom(ctx, AllianceDenom)
	require.Equal(t, startTime.Add(app.AllianceKeeper.RewardDelayTime(ctx)).Add(decayInterval*time.Duration(intervals)), asset.LastRewardChangeTime)
	require.True(t, decayRate.Power(intervals).Sub(asset.RewardWeight).LT(math.LegacyMustNewDecFromStr("0.0000000001")))
}

func TestClaimTakeRate(t *testing.T) {
	app, ctx := createTestContext(t)
	startTime := time.Now().UTC()
	ctx = ctx.WithBlockTime(startTime)
	ctx = ctx.WithBlockHeight(1)
	takeRateInterval := time.Minute * 5
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.Params{
			RewardDelayTime:       time.Minute * 60,
			TakeRateClaimInterval: takeRateInterval,
			LastTakeRateClaimTime: startTime,
		},
		Assets: []types.AllianceAsset{
			types.NewAllianceAsset(AllianceDenom, math.LegacyNewDec(2), math.LegacyZeroDec(), math.LegacyNewDec(5), math.LegacyMustNewDecFromStr("0.5"), startTime),
			types.NewAllianceAsset(AllianceDenomTwo, math.LegacyNewDec(10), math.LegacyNewDec(2), math.LegacyNewDec(12), math.LegacyNewDec(0), startTime),
		},
	})

	// Accounts
	feeCollectorAddr := app.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
	delegations, err := app.StakingKeeper.GetAllDelegations(ctx)
	require.NoError(t, err)
	valAddr1, err := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	require.NoError(t, err)
	val1, err := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr1)
	require.NoError(t, err)
	addrs := test_helpers.AddTestAddrsIncremental(app, ctx, 1, sdk.NewCoins(
		sdk.NewCoin(AllianceDenom, math.NewInt(1000_000_000)),
		sdk.NewCoin(AllianceDenomTwo, math.NewInt(1000_000_000)),
	))
	user1 := addrs[0]

	_, err = app.AllianceKeeper.Delegate(ctx, user1, val1, sdk.NewCoin(AllianceDenom, math.NewInt(1000_000_000)))
	require.NoError(t, err)

	_, err = app.AllianceKeeper.Delegate(ctx, user1, val1, sdk.NewCoin(AllianceDenomTwo, math.NewInt(1000_000_000)))
	require.NoError(t, err)

	assets := app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)
	totalBonded, err := app.StakingKeeper.TotalBondedTokens(ctx)
	require.NoError(t, err)
	// Check total bonded amount
	require.Equal(t, math.NewInt(13_000_000), totalBonded)

	// Calling it immediately will not update anything
	coins, err := app.AllianceKeeper.DeductAssetsHook(ctx, assets)
	require.Nil(t, coins)
	require.Nil(t, err)

	// Advance block time
	timePassed := time.Minute*5 + time.Second
	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(timePassed))
	ctx = ctx.WithBlockHeight(2)
	coinsClaimed, _ := app.AllianceKeeper.DeductAssetsHook(ctx, assets)
	coins = app.BankKeeper.GetAllBalances(ctx, feeCollectorAddr)
	require.Equal(t, coinsClaimed, coins)

	expectedAmount := math.LegacyMustNewDecFromStr("0.5").Mul(math.LegacyNewDec(timePassed.Nanoseconds() / takeRateInterval.Nanoseconds())).MulInt(math.NewInt(1000_000_000))
	require.Equal(t, expectedAmount.TruncateInt(), coins.AmountOf(AllianceDenom))

	lastUpdate := app.AllianceKeeper.LastRewardClaimTime(ctx)
	require.Equal(t, startTime.Add(takeRateInterval), lastUpdate)

	asset, found := app.AllianceKeeper.GetAssetByDenom(ctx, AllianceDenom)
	require.True(t, found)
	require.Equal(t, math.LegacyNewDec(2), asset.RewardWeight)

	// At the next begin block, tokens will be distributed from the fee pool
	cons, _ := val1.GetConsAddr()
	app.DistrKeeper.AllocateTokens(ctx, 1, []abcitypes.VoteInfo{
		{
			Validator: abcitypes.Validator{
				Address: cons,
				Power:   1,
			},
			BlockIdFlag: 1,
		},
	})

	outstandingRewards, err := app.DistrKeeper.GetValidatorOutstandingRewards(ctx, valAddr1)
	require.NoError(t, err)
	rewards := outstandingRewards.Rewards
	feePool, err := app.DistrKeeper.FeePool.Get(ctx)
	require.NoError(t, err)
	community := feePool.CommunityPool
	// This is case, validator 1 has 0% commission
	accumulatedCommission, err := app.DistrKeeper.GetValidatorAccumulatedCommission(ctx, valAddr1)
	require.NoError(t, err)
	commission := accumulatedCommission.Commission
	require.Equal(t, sdk.DecCoins(nil), commission)
	// And rewards + community pool should add up to total coins claimed
	require.Equal(t,
		math.LegacyNewDecFromInt(coinsClaimed.AmountOf(AllianceDenom)),
		rewards.AmountOf(AllianceDenom).Add(community.AmountOf(AllianceDenom)),
	)
}

func TestClaimTakeRateToZero(t *testing.T) {
	app, ctx := createTestContext(t)
	startTime := time.Now().UTC()
	ctx = ctx.WithBlockTime(startTime)
	ctx = ctx.WithBlockHeight(1)
	takeRateInterval := time.Minute * 5
	asset := types.NewAllianceAsset(AllianceDenom, math.LegacyNewDec(2), math.LegacyZeroDec(), math.LegacyNewDec(5), math.LegacyMustNewDecFromStr("0.8"), startTime)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.Params{
			RewardDelayTime:       time.Minute * 60,
			TakeRateClaimInterval: takeRateInterval,
			LastTakeRateClaimTime: startTime,
		},
		Assets: []types.AllianceAsset{
			asset,
		},
	})

	// Accounts
	delegations, err := app.StakingKeeper.GetAllDelegations(ctx)
	require.NoError(t, err)
	valAddr1, err := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	require.NoError(t, err)
	val1, err := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr1)
	require.NoError(t, err)
	addrs := test_helpers.AddTestAddrsIncremental(app, ctx, 1, sdk.NewCoins(
		sdk.NewCoin(AllianceDenom, math.NewInt(1000_000_000)),
	))
	user1 := addrs[0]

	_, err = app.AllianceKeeper.Delegate(ctx, user1, val1, sdk.NewCoin(AllianceDenom, math.NewInt(1000_000_000)))
	require.NoError(t, err)

	assets := app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)

	timePassed := time.Minute * 5
	// Advance block time
	for i := 0; i < 1000; i++ {
		ctx = ctx.WithBlockTime(ctx.BlockTime().Add(timePassed))
		ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
		_, err = app.AllianceKeeper.DeductAssetsHook(ctx, assets)
		require.NoError(t, err)
	}

	asset, _ = app.AllianceKeeper.GetAssetByDenom(ctx, AllianceDenom)
	require.True(t, asset.TotalTokens.GTE(math.OneInt()))
}

func TestClaimTakeRateForNewlyAddedAssets(t *testing.T) {
	app, ctx := createTestContext(t)
	startTime := time.Now().UTC()
	ctx = ctx.WithBlockTime(startTime)
	ctx = ctx.WithBlockHeight(1)
	takeRateInterval := time.Minute * 5
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.Params{
			RewardDelayTime:       time.Minute * 60,
			TakeRateClaimInterval: takeRateInterval,
			LastTakeRateClaimTime: startTime,
		},
		Assets: []types.AllianceAsset{
			types.NewAllianceAsset(AllianceDenom, math.LegacyNewDec(2), math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec(), startTime),
		},
	})

	// Accounts
	// feeCollectorAddr := app.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
	delegations, err := app.StakingKeeper.GetAllDelegations(ctx)
	require.NoError(t, err)
	valAddr1, err := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	require.NoError(t, err)
	val1, err := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr1)
	require.NoError(t, err)
	addrs := test_helpers.AddTestAddrsIncremental(app, ctx, 1, sdk.NewCoins(
		sdk.NewCoin(AllianceDenom, math.NewInt(1000_000_000)),
		sdk.NewCoin(AllianceDenomTwo, math.NewInt(1000_000_000)),
	))
	user1 := addrs[0]

	_, err = app.AllianceKeeper.Delegate(ctx, user1, val1, sdk.NewCoin(AllianceDenom, math.NewInt(1000_000_000)))
	require.NoError(t, err)

	assets := app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)

	// Calling it immediately will not update anything
	coins, err := app.AllianceKeeper.DeductAssetsHook(ctx, assets)
	require.Nil(t, coins)
	require.Nil(t, err)

	// Advance block time
	blockTime := ctx.BlockTime().Add(time.Minute*5 + time.Second)
	ctx = ctx.WithBlockTime(blockTime)
	ctx = ctx.WithBlockHeight(2)
	_, err = app.AllianceKeeper.DeductAssetsHook(ctx, assets)
	require.NoError(t, err)

	// Last take rate claim time should be updated even though nothing has been taxed
	lastClaimTime := app.AllianceKeeper.LastRewardClaimTime(ctx)
	require.Equal(t, blockTime, lastClaimTime)

	err = app.AllianceKeeper.CreateAlliance(ctx, &types.MsgCreateAllianceProposal{
		Title:                "New alliance",
		Description:          "",
		Denom:                AllianceDenomTwo,
		RewardWeight:         math.LegacyNewDec(1),
		TakeRate:             math.LegacyMustNewDecFromStr("0.1"),
		RewardChangeRate:     math.LegacyOneDec(),
		RewardChangeInterval: 0,
		RewardWeightRange: types.RewardWeightRange{
			Min: math.LegacyZeroDec(),
			Max: math.LegacyOneDec(),
		},
	})
	require.NoError(t, err)
	tax, err := app.AllianceKeeper.DeductAssetsHook(ctx, assets)
	require.NoError(t, err)

	_, err = app.AllianceKeeper.Delegate(ctx, user1, val1, sdk.NewCoin(AllianceDenomTwo, math.NewInt(1000_000_000)))
	require.NoError(t, err)
	require.Len(t, tax, 0)

	assets = app.AllianceKeeper.GetAllAssets(ctx)

	// Advance block time but not yet reward delay time
	blockTime = ctx.BlockTime().Add(time.Minute*5 + time.Second)
	ctx = ctx.WithBlockTime(blockTime)
	ctx = ctx.WithBlockHeight(3)
	tax, err = app.AllianceKeeper.DeductAssetsHook(ctx, assets)
	require.NoError(t, err)
	require.Len(t, tax, 0)

	// Advance block time after reward delay time
	blockTime = ctx.BlockTime().Add(time.Minute * 60)
	ctx = ctx.WithBlockTime(blockTime)
	ctx = ctx.WithBlockHeight(4)
	tax, err = app.AllianceKeeper.DeductAssetsHook(ctx, assets)
	require.NoError(t, err)
	require.Len(t, tax, 1)
}

func TestRewardWeightRateChange(t *testing.T) {
	app, ctx := createTestContext(t)
	startTime := time.Now().UTC()
	ctx = ctx.WithBlockTime(startTime)
	ctx = ctx.WithBlockHeight(1)
	takeRateInterval := time.Minute * 5
	alliance := types.NewAllianceAsset(AllianceDenom, math.LegacyNewDec(2), math.LegacyZeroDec(), math.LegacyNewDec(5), math.LegacyZeroDec(), startTime)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.Params{
			RewardDelayTime:       time.Minute * 60,
			TakeRateClaimInterval: takeRateInterval,
			LastTakeRateClaimTime: startTime,
		},
		Assets: []types.AllianceAsset{
			alliance,
		},
	})

	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(time.Hour * 10))

	err := app.AllianceKeeper.UpdateAlliance(ctx, &types.MsgUpdateAllianceProposal{
		Title:                "Update",
		Description:          "",
		Denom:                alliance.Denom,
		RewardWeight:         alliance.RewardWeight,
		RewardWeightRange:    types.RewardWeightRange{Min: math.LegacyNewDec(0), Max: math.LegacyNewDec(5)},
		TakeRate:             alliance.TakeRate,
		RewardChangeRate:     math.LegacyMustNewDecFromStr("1.001"),
		RewardChangeInterval: time.Minute * 5,
	})
	require.NoError(t, err)

	alliance, _ = app.AllianceKeeper.GetAssetByDenom(ctx, alliance.Denom)
	require.Equal(t, alliance.LastRewardChangeTime, ctx.BlockTime())
}
