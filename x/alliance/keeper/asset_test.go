package keeper_test

import (
	test_helpers "alliance/app"
	"alliance/x/alliance/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/teststaking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestUpdateRewardRates(t *testing.T) {
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx.WithBlockTime(startTime)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			types.NewAllianceAsset(ALLIANCE_TOKEN_DENOM, sdk.NewDec(2), sdk.ZeroDec()),
			types.NewAllianceAsset(ALLIANCE_2_TOKEN_DENOM, sdk.NewDec(10), sdk.ZeroDec()),
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
		sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(1000_000)),
		sdk.NewCoin(ALLIANCE_2_TOKEN_DENOM, sdk.NewInt(1000_000)),
	))
	pks := test_helpers.CreateTestPubKeys(2)
	powerReduction := app.StakingKeeper.PowerReduction(ctx)

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
	_val1.Status = stakingtypes.Bonded
	test_helpers.RegisterNewValidator(t, app, ctx, _val1)
	val1, err := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr1)
	require.NoError(t, err)

	user1 := addrs[2]

	// Start by delegating
	_, err = app.AllianceKeeper.Delegate(ctx, user1, val1, sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(1000_000)))
	require.NoError(t, err)

	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx)
	require.NoError(t, err)
	require.Equal(t, sdk.NewInt(3_000_000), app.StakingKeeper.TotalBondedTokens(ctx))

	// Expecting voting power for the alliance module
	val, found := app.StakingKeeper.GetValidator(ctx, valAddr1)
	require.True(t, found)
	require.Equal(t, int64(2), val.ConsensusPower(powerReduction))

	err = app.AllianceKeeper.UpdateAllianceAsset(ctx, types.AllianceAsset{
		Denom:        ALLIANCE_TOKEN_DENOM,
		RewardWeight: sdk.NewDec(20),
		TakeRate:     sdk.NewDec(0),
	})
	require.NoError(t, err)

	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx)
	require.NoError(t, err)
	require.Equal(t, sdk.NewInt(21_000_000), app.StakingKeeper.TotalBondedTokens(ctx))

	// Expecting voting power to increase
	val, found = app.StakingKeeper.GetValidator(ctx, valAddr1)
	require.True(t, found)
	require.Equal(t, int64(20), val.ConsensusPower(powerReduction))

	err = app.AllianceKeeper.UpdateAllianceAsset(ctx, types.AllianceAsset{
		Denom:        ALLIANCE_TOKEN_DENOM,
		RewardWeight: sdk.NewDec(1),
		TakeRate:     sdk.NewDec(0),
	})
	require.NoError(t, err)

	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx)
	require.NoError(t, err)
	require.Equal(t, sdk.NewInt(2_000_000), app.StakingKeeper.TotalBondedTokens(ctx))

	// Expecting voting power to decrease
	val, found = app.StakingKeeper.GetValidator(ctx, valAddr1)
	require.True(t, found)
	require.Equal(t, int64(1), val.ConsensusPower(powerReduction))
}

func TestSlashingEvent(t *testing.T) {
	var err error
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			{
				Denom:        ALLIANCE_TOKEN_DENOM,
				RewardWeight: sdk.NewDec(2),
				TakeRate:     sdk.NewDec(0),
				TotalTokens:  sdk.ZeroInt(),
			},
			{
				Denom:        ALLIANCE_2_TOKEN_DENOM,
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
		sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(20_000_000)),
		sdk.NewCoin(ALLIANCE_2_TOKEN_DENOM, sdk.NewInt(20_000_000)),
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
	_, err = app.AllianceKeeper.Delegate(ctx, user1, val1, sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(10_000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user1, val2, sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(10_000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user2, val1, sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(10_000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user2, val2, sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(10_000_000)))
	require.NoError(t, err)

	_, err = app.AllianceKeeper.Delegate(ctx, user1, val1, sdk.NewCoin(ALLIANCE_2_TOKEN_DENOM, sdk.NewInt(10_000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user1, val2, sdk.NewCoin(ALLIANCE_2_TOKEN_DENOM, sdk.NewInt(10_000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user2, val1, sdk.NewCoin(ALLIANCE_2_TOKEN_DENOM, sdk.NewInt(10_000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user2, val2, sdk.NewCoin(ALLIANCE_2_TOKEN_DENOM, sdk.NewInt(10_000_000)))
	require.NoError(t, err)
	require.NoError(t, err)

	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx)
	require.NoError(t, err)
	require.Equal(t, sdk.NewInt(13_000_000), app.StakingKeeper.TotalBondedTokens(ctx))

	val1, _ = app.AllianceKeeper.GetAllianceValidator(ctx, valAddr1)
	valPower1 := val1.GetConsensusPower(app.StakingKeeper.PowerReduction(ctx))
	valConAddr1, err := val1.GetConsAddr()

	// Tokens should remain the same before slashing
	asset1, _ := app.AllianceKeeper.GetAssetByDenom(ctx, ALLIANCE_TOKEN_DENOM)
	tokens := val1.TotalTokensWithAsset(asset1)
	require.Equal(t, sdk.NewInt(20_000_000), tokens)
	asset2, _ := app.AllianceKeeper.GetAssetByDenom(ctx, ALLIANCE_2_TOKEN_DENOM)
	tokens = val1.TotalTokensWithAsset(asset2)
	require.Equal(t, sdk.NewInt(20_000_000), tokens)

	app.SlashingKeeper.Slash(ctx, valConAddr1, app.SlashingKeeper.SlashFractionDoubleSign(ctx), valPower1, 1)
	// Slashing will first reduce tokens from validator
	require.NotEqual(t, sdk.NewInt(13_000_000), app.StakingKeeper.TotalBondedTokens(ctx))

	// After rebalancing, it should recover the tokens
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx)
	require.NoError(t, err)
	require.Equal(t, sdk.NewInt(13_000_000), app.StakingKeeper.TotalBondedTokens(ctx))

	// Expect that total tokens with validator 1 are reduced
	val1, _ = app.AllianceKeeper.GetAllianceValidator(ctx, valAddr1)
	asset1, _ = app.AllianceKeeper.GetAssetByDenom(ctx, ALLIANCE_TOKEN_DENOM)
	tokens = val1.TotalTokensWithAsset(asset1)
	require.Greater(t, sdk.NewInt(20_000_000).Int64(), tokens.Int64())
	asset2, _ = app.AllianceKeeper.GetAssetByDenom(ctx, ALLIANCE_2_TOKEN_DENOM)
	tokens = val1.TotalTokensWithAsset(asset2)
	require.Greater(t, sdk.NewInt(20_000_000).Int64(), tokens.Int64())

	// Expect that total tokens with validator 2 increased (redistributed from slashing)
	val2, _ = app.AllianceKeeper.GetAllianceValidator(ctx, valAddr2)
	asset1, _ = app.AllianceKeeper.GetAssetByDenom(ctx, ALLIANCE_TOKEN_DENOM)
	tokens = val2.TotalTokensWithAsset(asset1)
	require.Less(t, sdk.NewInt(20_000_000).Int64(), tokens.Int64())
	asset2, _ = app.AllianceKeeper.GetAssetByDenom(ctx, ALLIANCE_2_TOKEN_DENOM)
	tokens = val2.TotalTokensWithAsset(asset2)
	require.Less(t, sdk.NewInt(20_000_000).Int64(), tokens.Int64())

	// Expect that consensus power for val1 dropped
	newValPower1 := val1.GetConsensusPower(app.StakingKeeper.PowerReduction(ctx))
	require.Less(t, newValPower1, valPower1)
}

func TestUnbondedValidator(t *testing.T) {
	var err error
	app, ctx := createTestContext(t)
	bondDenom := app.StakingKeeper.BondDenom(ctx)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			{
				Denom:        ALLIANCE_TOKEN_DENOM,
				RewardWeight: sdk.MustNewDecFromStr("0.1"),
				TakeRate:     sdk.NewDec(0),
				TotalTokens:  sdk.ZeroInt(),
			},
			{
				Denom:        ALLIANCE_2_TOKEN_DENOM,
				RewardWeight: sdk.MustNewDecFromStr("0.5"),
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
	addrs := test_helpers.AddTestAddrsIncremental(app, ctx, 5, sdk.NewCoins(
		sdk.NewCoin(bondDenom, sdk.NewInt(10_000_000)),
		sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(50_000_000)),
		sdk.NewCoin(ALLIANCE_2_TOKEN_DENOM, sdk.NewInt(50_000_000)),
	))
	pks := test_helpers.CreateTestPubKeys(2)

	// Increase the stake on genesis validator
	delegations := app.StakingKeeper.GetAllDelegations(ctx)
	require.Len(t, delegations, 1)
	valAddr0, err := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	require.NoError(t, err)
	val0, _ := app.StakingKeeper.GetValidator(ctx, valAddr0)
	_, err = app.StakingKeeper.Delegate(ctx, addrs[4], sdk.NewInt(9_000_000), stakingtypes.Unbonded, val0, true)
	require.NoError(t, err)

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
	_val1.Description.Moniker = "val1"
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
	_val2.Description.Moniker = "val2"
	test_helpers.RegisterNewValidator(t, app, ctx, _val2)
	val2, err := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr2)
	require.NoError(t, err)

	user1 := addrs[2]
	user2 := addrs[3]

	// Users add delegations
	_, err = app.AllianceKeeper.Delegate(ctx, user1, val1, sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(20_000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user1, val2, sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(10_000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user2, val1, sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(10_000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user2, val2, sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(10_000_000)))
	require.NoError(t, err)

	_, err = app.AllianceKeeper.Delegate(ctx, user1, val1, sdk.NewCoin(ALLIANCE_2_TOKEN_DENOM, sdk.NewInt(10_000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user1, val2, sdk.NewCoin(ALLIANCE_2_TOKEN_DENOM, sdk.NewInt(10_000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user2, val1, sdk.NewCoin(ALLIANCE_2_TOKEN_DENOM, sdk.NewInt(10_000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user2, val2, sdk.NewCoin(ALLIANCE_2_TOKEN_DENOM, sdk.NewInt(10_000_000)))
	require.NoError(t, err)
	require.NoError(t, err)

	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx)
	require.NoError(t, err)
	_, err = app.StakingKeeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	require.Equal(t, sdk.NewInt(16_000_000), app.StakingKeeper.TotalBondedTokens(ctx))

	val1, _ = app.AllianceKeeper.GetAllianceValidator(ctx, valAddr1)
	val2, _ = app.AllianceKeeper.GetAllianceValidator(ctx, valAddr2)
	require.Greater(t, val1.Tokens.Int64(), val2.Tokens.Int64())

	// Set max validators to be 2 to trigger unbonding
	params := app.StakingKeeper.GetParams(ctx)
	params.MaxValidators = 2
	app.StakingKeeper.SetParams(ctx, params)
	_, err = app.StakingKeeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	require.Equal(t, sdk.NewInt(13_100_000), app.StakingKeeper.TotalBondedTokens(ctx))

	vals := app.StakingKeeper.GetBondedValidatorsByPower(ctx)
	require.Equal(t, 2, len(vals))
	require.Equal(t, "val1", vals[1].GetMoniker())

	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx)
	require.NoError(t, err)
	require.Equal(t, sdk.NewInt(16_000_000).String(), app.StakingKeeper.TotalBondedTokens(ctx).String())

	val2, _ = app.AllianceKeeper.GetAllianceValidator(ctx, valAddr2)

	// Set max validators to be 3 to trigger rebonding
	params = app.StakingKeeper.GetParams(ctx)
	params.MaxValidators = 3
	app.StakingKeeper.SetParams(ctx, params)
	_, err = app.StakingKeeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	require.Equal(t, sdk.NewInt(18_900_000), app.StakingKeeper.TotalBondedTokens(ctx))

	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx)
	require.NoError(t, err)
	require.Equal(t, sdk.NewInt(16_000_000).String(), app.StakingKeeper.TotalBondedTokens(ctx).String())
}

func TestJailedValidator(t *testing.T) {
	var err error
	app, ctx := createTestContext(t)
	bondDenom := app.StakingKeeper.BondDenom(ctx)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			{
				Denom:        ALLIANCE_TOKEN_DENOM,
				RewardWeight: sdk.MustNewDecFromStr("0.1"),
				TakeRate:     sdk.NewDec(0),
				TotalTokens:  sdk.ZeroInt(),
			},
			{
				Denom:        ALLIANCE_2_TOKEN_DENOM,
				RewardWeight: sdk.MustNewDecFromStr("0.5"),
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
	addrs := test_helpers.AddTestAddrsIncremental(app, ctx, 5, sdk.NewCoins(
		sdk.NewCoin(bondDenom, sdk.NewInt(10_000_000)),
		sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(50_000_000)),
		sdk.NewCoin(ALLIANCE_2_TOKEN_DENOM, sdk.NewInt(50_000_000)),
	))
	pks := test_helpers.CreateTestPubKeys(2)

	// Increase the stake on genesis validator
	delegations := app.StakingKeeper.GetAllDelegations(ctx)
	require.Len(t, delegations, 1)
	valAddr0, err := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	require.NoError(t, err)
	val0, _ := app.StakingKeeper.GetValidator(ctx, valAddr0)
	_, err = app.StakingKeeper.Delegate(ctx, addrs[4], sdk.NewInt(9_000_000), stakingtypes.Unbonded, val0, true)
	require.NoError(t, err)

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
	_val1.Description.Moniker = "val1"
	test_helpers.RegisterNewValidator(t, app, ctx, _val1)
	_, err = app.StakingKeeper.Delegate(ctx, addrs[0], sdk.NewInt(1_000_000), stakingtypes.Unbonded, _val1, true)
	require.NoError(t, err)
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
	_val2.Description.Moniker = "val2"
	test_helpers.RegisterNewValidator(t, app, ctx, _val2)
	_, err = app.StakingKeeper.Delegate(ctx, addrs[1], sdk.NewInt(1_000_000), stakingtypes.Unbonded, _val2, true)
	require.NoError(t, err)
	val2, err := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr2)
	require.NoError(t, err)

	user1 := addrs[2]
	user2 := addrs[3]

	// Users add delegations
	_, err = app.AllianceKeeper.Delegate(ctx, user1, val1, sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(20_000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user1, val2, sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(10_000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user2, val1, sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(10_000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user2, val2, sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(10_000_000)))
	require.NoError(t, err)

	_, err = app.AllianceKeeper.Delegate(ctx, user1, val1, sdk.NewCoin(ALLIANCE_2_TOKEN_DENOM, sdk.NewInt(10_000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user1, val2, sdk.NewCoin(ALLIANCE_2_TOKEN_DENOM, sdk.NewInt(10_000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user2, val1, sdk.NewCoin(ALLIANCE_2_TOKEN_DENOM, sdk.NewInt(10_000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user2, val2, sdk.NewCoin(ALLIANCE_2_TOKEN_DENOM, sdk.NewInt(10_000_000)))
	require.NoError(t, err)
	require.NoError(t, err)

	_, err = app.StakingKeeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx)
	require.NoError(t, err)
	// 12 * 1.6 = 19.2
	require.Equal(t, sdk.NewInt(19_200_000), app.StakingKeeper.TotalBondedTokens(ctx))

	val1, _ = app.AllianceKeeper.GetAllianceValidator(ctx, valAddr1)
	val2, _ = app.AllianceKeeper.GetAllianceValidator(ctx, valAddr2)
	require.Greater(t, val1.Tokens.Int64(), val2.Tokens.Int64())

	// Jail validator
	cons2, _ := val2.GetConsAddr()
	app.SlashingKeeper.Jail(ctx, cons2)
	_, err = app.StakingKeeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	require.Equal(t, sdk.NewInt(14_720_000), app.StakingKeeper.TotalBondedTokens(ctx))

	vals := app.StakingKeeper.GetBondedValidatorsByPower(ctx)
	require.Equal(t, 2, len(vals))
	require.Equal(t, "val1", vals[1].GetMoniker())

	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx)
	require.NoError(t, err)
	// 11 * 1.6 = 17.6
	require.Equal(t, sdk.NewInt(17_600_000).String(), app.StakingKeeper.TotalBondedTokens(ctx).String())

	val2, _ = app.AllianceKeeper.GetAllianceValidator(ctx, valAddr2)

	// Unjail validator
	err = app.SlashingKeeper.Unjail(ctx, valAddr2)
	require.NoError(t, err)
	_, err = app.StakingKeeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	require.Equal(t, sdk.NewInt(22_080_000), app.StakingKeeper.TotalBondedTokens(ctx))

	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx)
	require.NoError(t, err)
	require.Equal(t, sdk.NewInt(19_200_000).String(), app.StakingKeeper.TotalBondedTokens(ctx).String())
}
