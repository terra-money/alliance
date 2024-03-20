package tests_test

import (
	"testing"
	"time"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	test_helpers "github.com/terra-money/alliance/app"
	"github.com/terra-money/alliance/x/alliance/types"
)

func TestUnbondingMethods(t *testing.T) {
	// Setup the context with an alliance asset
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	allianceAsset := types.AllianceAsset{
		Denom:                AllianceDenom,
		RewardWeight:         math.LegacyNewDec(1),
		TakeRate:             math.LegacyNewDec(0),
		TotalTokens:          math.ZeroInt(),
		TotalValidatorShares: math.LegacyZeroDec(),
		RewardStartTime:      startTime,
		RewardChangeRate:     math.LegacyNewDec(0),
		RewardChangeInterval: time.Duration(0),
		LastRewardChangeTime: startTime,
		RewardWeightRange:    types.RewardWeightRange{Min: math.LegacyNewDec(0), Max: math.LegacyNewDec(1)},
		IsInitialized:        true,
	}
	app.AllianceKeeper.SetAsset(ctx, allianceAsset)

	// Query staking module unbonding time to assert later on
	unbondingTime := app.StakingKeeper.UnbondingTime(ctx)

	// Get the native delegations to have a validator address where to delegate
	delegations := app.StakingKeeper.GetAllDelegations(ctx)
	valAddr, err := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	require.NoError(t, err)

	// Get a delegator address with funds
	delAddrs := test_helpers.AddTestAddrsIncremental(app, ctx, 2, sdk.NewCoins(sdk.NewCoin(AllianceDenom, math.NewInt(1000_000_000))))
	delAddr := delAddrs[0]
	delAddr1 := delAddrs[1]

	// Get an alliance validator
	val, err := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr)
	require.NoError(t, err)

	// Delegate the alliance asset with both accounts
	res, err := app.AllianceKeeper.Delegate(ctx, delAddr, val, sdk.NewCoin(AllianceDenom, math.NewInt(1000_000_000)))
	require.Nil(t, err)
	require.Equal(t, math.LegacyNewDec(1000_000_000), *res)
	res2, err := app.AllianceKeeper.Delegate(ctx, delAddr1, val, sdk.NewCoin(AllianceDenom, math.NewInt(1000_000_000)))
	require.Nil(t, err)
	require.Equal(t, math.LegacyNewDec(1000_000_000), *res2)

	// Undelegate the alliance assets with both accounts
	undelRes, err := app.AllianceKeeper.Undelegate(ctx, delAddr, val, sdk.NewCoin(AllianceDenom, math.NewInt(1000_000_000)))
	require.Nil(t, err)
	require.Equal(t, ctx.BlockHeader().Time.Add(unbondingTime), *undelRes)
	undelRes2, err := app.AllianceKeeper.Undelegate(ctx, delAddr1, val, sdk.NewCoin(AllianceDenom, math.NewInt(1000_000_000)))
	require.Nil(t, err)
	require.Equal(t, ctx.BlockHeader().Time.Add(unbondingTime), *undelRes2)

	// Validate that both user delegations executed the unbonding process
	unbondings, err := app.AllianceKeeper.GetUnbondingsByDelegator(ctx, delAddr)
	require.NoError(t, err)
	require.Equal(t,
		[]types.UnbondingDelegation{{
			ValidatorAddress: valAddr.String(),
			Amount:           math.NewInt(1000_000_000),
			CompletionTime:   ctx.BlockHeader().Time.Add(unbondingTime),
			Denom:            AllianceDenom,
		}},
		unbondings,
	)

	// Validate that both user delegations executed the unbonding process
	unbondings, err = app.AllianceKeeper.GetUnbondingsByDenomAndDelegator(ctx, AllianceDenom, delAddr1)
	require.NoError(t, err)
	require.Equal(t,
		[]types.UnbondingDelegation{{
			ValidatorAddress: valAddr.String(),
			Amount:           math.NewInt(1000_000_000),
			CompletionTime:   ctx.BlockHeader().Time.Add(unbondingTime),
			Denom:            AllianceDenom,
		}},
		unbondings,
	)

	unbondings, err = app.AllianceKeeper.GetUnbondings(ctx, AllianceDenom, delAddr1, valAddr)
	require.NoError(t, err)
	require.Equal(t,
		[]types.UnbondingDelegation{{
			ValidatorAddress: valAddr.String(),
			Amount:           math.NewInt(1000_000_000),
			CompletionTime:   ctx.BlockHeader().Time.Add(unbondingTime),
			Denom:            AllianceDenom,
		}},
		unbondings,
	)
}

func TestUnbondingMethodsLargeNumbers(t *testing.T) {
	// Setup the context with an alliance asset
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	allianceAsset := types.AllianceAsset{
		Denom:                AllianceDenom,
		RewardWeight:         math.LegacyMustNewDecFromStr("0.025"),
		TakeRate:             math.LegacyNewDec(0),
		TotalTokens:          math.ZeroInt(),
		TotalValidatorShares: math.LegacyZeroDec(),
		RewardStartTime:      startTime,
		RewardChangeRate:     math.LegacyNewDec(0),
		RewardChangeInterval: time.Duration(0),
		LastRewardChangeTime: startTime,
		RewardWeightRange:    types.RewardWeightRange{Min: math.LegacyNewDec(0), Max: math.LegacyNewDec(1)},
		IsInitialized:        true,
	}
	app.AllianceKeeper.SetAsset(ctx, allianceAsset)

	// Query staking module unbonding time to assert later on
	unbondingTime := app.StakingKeeper.UnbondingTime(ctx)

	// Get the native delegations to have a validator address where to delegate
	delegations := app.StakingKeeper.GetAllDelegations(ctx)
	valAddr, err := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	require.NoError(t, err)

	// Get a delegator address with funds
	delAddrs := test_helpers.AddTestAddrsIncremental(app, ctx, 2, sdk.NewCoins(sdk.NewCoin(AllianceDenom, math.NewInt(1000_000_000_000_000))))
	delAddr := delAddrs[0]
	delAddr1 := delAddrs[1]

	// Get an alliance validator
	val, err := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr)
	require.NoError(t, err)

	// Delegate the alliance asset with both accounts
	res, err := app.AllianceKeeper.Delegate(ctx, delAddr, val, sdk.NewCoin(AllianceDenom, math.NewInt(830697941465481)))
	require.Nil(t, err)
	require.Equal(t, math.LegacyNewDec(830697941465481), *res)
	res2, err := app.AllianceKeeper.Delegate(ctx, delAddr1, val, sdk.NewCoin(AllianceDenom, math.NewInt(975933204219431)))
	require.Nil(t, err)
	require.Equal(t, math.LegacyNewDec(975933204219431), *res2)

	// Undelegate the alliance assets with both accounts
	undelRes, err := app.AllianceKeeper.Undelegate(ctx, delAddr, val, sdk.NewCoin(AllianceDenom, math.NewInt(564360383558874)))
	require.Nil(t, err)
	require.Equal(t, ctx.BlockHeader().Time.Add(unbondingTime), *undelRes)
	undelRes2, err := app.AllianceKeeper.Undelegate(ctx, delAddr1, val, sdk.NewCoin(AllianceDenom, math.NewInt(384108763572096)))
	require.Nil(t, err)
	require.Equal(t, ctx.BlockHeader().Time.Add(unbondingTime), *undelRes2)

	// Validate that both user delegations executed the unbonding process
	unbondings, err := app.AllianceKeeper.GetUnbondingsByDelegator(ctx, delAddr)
	require.NoError(t, err)
	require.Equal(t,
		[]types.UnbondingDelegation{{
			ValidatorAddress: valAddr.String(),
			Amount:           math.NewInt(564360383558874),
			CompletionTime:   ctx.BlockHeader().Time.Add(unbondingTime),
			Denom:            AllianceDenom,
		}},
		unbondings,
	)

	// Validate that both user delegations executed the unbonding process
	unbondings, err = app.AllianceKeeper.GetUnbondingsByDenomAndDelegator(ctx, AllianceDenom, delAddr1)
	require.NoError(t, err)
	require.Equal(t,
		[]types.UnbondingDelegation{{
			ValidatorAddress: valAddr.String(),
			Amount:           math.NewInt(384108763572096),
			CompletionTime:   ctx.BlockHeader().Time.Add(unbondingTime),
			Denom:            AllianceDenom,
		}},
		unbondings,
	)

	unbondings, err = app.AllianceKeeper.GetUnbondings(ctx, AllianceDenom, delAddr1, valAddr)
	require.NoError(t, err)
	require.Equal(t,
		[]types.UnbondingDelegation{{
			ValidatorAddress: valAddr.String(),
			Amount:           math.NewInt(384108763572096),
			CompletionTime:   ctx.BlockHeader().Time.Add(unbondingTime),
			Denom:            AllianceDenom,
		}},
		unbondings,
	)
}
