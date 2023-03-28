package tests_test

import (
	"testing"
	"time"

	test_helpers "github.com/terra-money/alliance/app"
	"github.com/terra-money/alliance/x/alliance"
	"github.com/terra-money/alliance/x/alliance/keeper"
	"github.com/terra-money/alliance/x/alliance/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	teststaking "github.com/cosmos/cosmos-sdk/x/staking/testutil"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slices"
)

var (
	AllianceDenom    = "alliance"
	AllianceDenomTwo = "alliance2"
)

func TestDelegationWithASingleAsset(t *testing.T) {
	app, ctx := createTestContext(t)
	genesisTime := ctx.BlockTime()
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			types.NewAllianceAsset(AllianceDenom, sdk.NewDec(2), sdk.ZeroDec(), sdk.NewDec(5), sdk.NewDec(0), genesisTime),
			types.NewAllianceAsset(AllianceDenomTwo, sdk.NewDec(10), sdk.NewDec(2), sdk.NewDec(12), sdk.NewDec(0), genesisTime),
		},
	})
	delegations := app.StakingKeeper.GetAllDelegations(ctx)
	require.Len(t, delegations, 1)

	// All the addresses needed
	delAddr, err := sdk.AccAddressFromBech32(delegations[0].DelegatorAddress)
	require.NoError(t, err)
	valAddr, err := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	require.NoError(t, err)
	moduleAddr := app.AccountKeeper.GetModuleAddress(types.ModuleName)
	val, err := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr)
	require.NoError(t, err)

	// Mint alliance tokens
	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(AllianceDenom, sdk.NewInt(2000_000))))
	require.NoError(t, err)
	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(AllianceDenomTwo, sdk.NewInt(2000_000))))
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr, sdk.NewCoins(sdk.NewCoin(AllianceDenom, sdk.NewInt(2000_000))))
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr, sdk.NewCoins(sdk.NewCoin(AllianceDenomTwo, sdk.NewInt(2000_000))))
	require.NoError(t, err)

	// Check current total staked tokens
	totalBonded := app.StakingKeeper.TotalBondedTokens(ctx)
	require.Equal(t, sdk.NewInt(1000_000), totalBonded)

	// Delegate
	_, err = app.AllianceKeeper.Delegate(ctx, delAddr, val, sdk.NewCoin(AllianceDenom, sdk.NewInt(1000_000)))
	require.NoError(t, err)

	// Manually trigger rebalancing
	app.AllianceKeeper.GetAssetByDenom(ctx, AllianceDenom)
	assets := app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)

	allianceBonded := app.StakingKeeper.GetDelegatorBonded(ctx, moduleAddr)
	// Total ALLIANCE tokens should be 2 * totalBonded
	require.Equal(t, totalBonded.Mul(sdk.NewInt(2)).String(), allianceBonded.String())

	// Check delegation in staking module
	delegations = app.StakingKeeper.GetAllDelegations(ctx)
	require.Len(t, delegations, 2)
	i := slices.IndexFunc(delegations, func(d stakingtypes.Delegation) bool {
		return d.DelegatorAddress == moduleAddr.String()
	})
	newDelegation := delegations[i]
	require.Equal(t, stakingtypes.Delegation{
		DelegatorAddress: moduleAddr.String(),
		ValidatorAddress: valAddr.String(),
		Shares:           sdk.NewDec(2),
	}, newDelegation)

	// Check delegation in alliance module
	allianceDelegation, found := app.AllianceKeeper.GetDelegation(ctx, delAddr, valAddr, AllianceDenom)
	require.True(t, found)
	require.Equal(t, types.Delegation{
		DelegatorAddress: delAddr.String(),
		ValidatorAddress: val.OperatorAddress,
		Denom:            AllianceDenom,
		Shares:           sdk.NewDec(1000_000),
		RewardHistory:    types.RewardHistories(nil),
	}, allianceDelegation)

	// Check asset
	asset, _ := app.AllianceKeeper.GetAssetByDenom(ctx, AllianceDenom)
	require.Equal(t, types.AllianceAsset{
		Denom:                AllianceDenom,
		RewardWeight:         sdk.NewDec(2),
		RewardWeightRange:    types.RewardWeightRange{Min: sdk.NewDec(0), Max: sdk.NewDec(5)},
		TakeRate:             sdk.NewDec(0),
		TotalTokens:          sdk.NewInt(1000_000),
		TotalValidatorShares: sdk.NewDec(1000_000),
		RewardStartTime:      genesisTime,
		RewardChangeRate:     sdk.OneDec(),
		RewardChangeInterval: 0,
	}, asset)

	// Delegate with same denom again
	_, err = app.AllianceKeeper.Delegate(ctx, delAddr, val, sdk.NewCoin(AllianceDenom, sdk.NewInt(1000_000)))
	require.NoError(t, err)

	// Manually trigger rebalancing
	assets = app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)

	// Check delegation in alliance module
	allianceDelegation, found = app.AllianceKeeper.GetDelegation(ctx, delAddr, valAddr, AllianceDenom)
	require.True(t, found)
	require.Equal(t, types.Delegation{
		DelegatorAddress: delAddr.String(),
		ValidatorAddress: val.OperatorAddress,
		Denom:            AllianceDenom,
		Shares:           sdk.NewDec(2000_000),
		RewardHistory:    types.RewardHistories(nil),
	}, allianceDelegation)

	// Check asset again
	asset, _ = app.AllianceKeeper.GetAssetByDenom(ctx, AllianceDenom)
	require.Equal(t, types.AllianceAsset{
		Denom:                AllianceDenom,
		RewardWeight:         sdk.NewDec(2),
		RewardWeightRange:    types.RewardWeightRange{Min: sdk.NewDec(0), Max: sdk.NewDec(5)},
		TakeRate:             sdk.NewDec(0),
		TotalTokens:          sdk.NewInt(2000_000),
		TotalValidatorShares: sdk.NewDec(2000_000),
		RewardStartTime:      genesisTime,
		RewardChangeRate:     sdk.OneDec(),
		RewardChangeInterval: 0,
	}, asset)

	// Check delegation in staking module total shares should not change
	delegations = app.StakingKeeper.GetAllDelegations(ctx)
	require.Len(t, delegations, 2)
	i = slices.IndexFunc(delegations, func(d stakingtypes.Delegation) bool {
		return d.DelegatorAddress == moduleAddr.String()
	})
	newDelegation = delegations[i]
	require.Equal(t, stakingtypes.Delegation{
		DelegatorAddress: moduleAddr.String(),
		ValidatorAddress: valAddr.String(),
		Shares:           sdk.NewDec(2),
	}, newDelegation)
}

func TestDelegationWithMultipleAssets(t *testing.T) {
	app, ctx := createTestContext(t)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			types.NewAllianceAsset(AllianceDenom, sdk.NewDec(2), sdk.ZeroDec(), sdk.NewDec(5), sdk.NewDec(0), ctx.BlockTime()),
			types.NewAllianceAsset(AllianceDenomTwo, sdk.NewDec(10), sdk.NewDec(2), sdk.NewDec(12), sdk.NewDec(0), ctx.BlockTime()),
		},
	})
	delegations := app.StakingKeeper.GetAllDelegations(ctx)
	require.Len(t, delegations, 1)

	// All the addresses needed
	delAddr, err := sdk.AccAddressFromBech32(delegations[0].DelegatorAddress)
	require.NoError(t, err)
	valAddr, err := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	require.NoError(t, err)
	moduleAddr := app.AccountKeeper.GetModuleAddress(types.ModuleName)
	val, err := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr)
	require.NoError(t, err)

	// Mint alliance tokens
	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(AllianceDenom, sdk.NewInt(2000_000))))
	require.NoError(t, err)
	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(AllianceDenomTwo, sdk.NewInt(2000_000))))
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr, sdk.NewCoins(sdk.NewCoin(AllianceDenom, sdk.NewInt(2000_000))))
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr, sdk.NewCoins(sdk.NewCoin(AllianceDenomTwo, sdk.NewInt(2000_000))))
	require.NoError(t, err)

	// Delegate
	_, err = app.AllianceKeeper.Delegate(ctx, delAddr, val, sdk.NewCoin(AllianceDenom, sdk.NewInt(2000_000)))
	require.NoError(t, err)
	// Delegate with another denom
	_, err = app.AllianceKeeper.Delegate(ctx, delAddr, val, sdk.NewCoin(AllianceDenomTwo, sdk.NewInt(1000_000)))
	require.NoError(t, err)

	// Manually trigger rebalancing
	app.AllianceKeeper.GetAssetByDenom(ctx, AllianceDenomTwo)
	assets := app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)

	// Check delegation in staking module
	delegations = app.StakingKeeper.GetAllDelegations(ctx)
	require.Len(t, delegations, 2)
	i := slices.IndexFunc(delegations, func(d stakingtypes.Delegation) bool {
		return d.DelegatorAddress == moduleAddr.String()
	})
	newDelegation := delegations[i]
	// 1 * 2 + 1 * 10 = 12
	require.Equal(t, stakingtypes.Delegation{
		DelegatorAddress: moduleAddr.String(),
		ValidatorAddress: valAddr.String(),
		Shares:           sdk.NewDec(12),
	}, newDelegation)

	// Check validator in x/staking
	val, err = app.AllianceKeeper.GetAllianceValidator(ctx, valAddr)
	require.NoError(t, err)
	require.Equal(t, sdk.NewDec(13), val.DelegatorShares)
}

func TestDelegationWithUnknownAssets(t *testing.T) {
	app, ctx := createTestContext(t)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			types.NewAllianceAsset(AllianceDenom, sdk.NewDec(2), sdk.ZeroDec(), sdk.NewDec(5), sdk.NewDec(0), ctx.BlockTime()),
			types.NewAllianceAsset(AllianceDenomTwo, sdk.NewDec(10), sdk.NewDec(2), sdk.NewDec(12), sdk.NewDec(0), ctx.BlockTime()),
		},
	})
	delegations := app.StakingKeeper.GetAllDelegations(ctx)
	require.Len(t, delegations, 1)

	// All the addresses needed
	delAddr, err := sdk.AccAddressFromBech32(delegations[0].DelegatorAddress)
	require.NoError(t, err)
	valAddr, err := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	require.NoError(t, err)
	val, err := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr)
	require.NoError(t, err)

	// Mint alliance tokens
	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin("UNKNOWN", sdk.NewInt(2000_000))))
	require.NoError(t, err)

	_, err = app.AllianceKeeper.Delegate(ctx, delAddr, val, sdk.NewCoin("UNKNOWN", sdk.NewInt(2000_000)))
	require.Error(t, err)
}

func TestSuccessfulRedelegation(t *testing.T) {
	app, ctx := createTestContext(t)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			types.NewAllianceAsset(AllianceDenom, sdk.NewDec(2), sdk.ZeroDec(), sdk.NewDec(5), sdk.NewDec(0), ctx.BlockTime()),
			types.NewAllianceAsset(AllianceDenomTwo, sdk.NewDec(10), sdk.NewDec(2), sdk.NewDec(12), sdk.NewDec(0), ctx.BlockTime()),
		},
	})

	// Get all the addresses needed for the test
	moduleAddr := app.AccountKeeper.GetModuleAddress(types.ModuleName)
	delegations := app.StakingKeeper.GetAllDelegations(ctx)
	valAddr1, err := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	require.NoError(t, err)
	val1, err := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr1)
	require.NoError(t, err)
	addrs := test_helpers.AddTestAddrsIncremental(app, ctx, 3, sdk.NewCoins(
		sdk.NewCoin(AllianceDenom, sdk.NewInt(1000_000)),
		sdk.NewCoin(AllianceDenomTwo, sdk.NewInt(1000_000)),
	))
	valAddr2 := sdk.ValAddress(addrs[0])
	_val2 := teststaking.NewValidator(t, valAddr2, test_helpers.CreateTestPubKeys(1)[0])
	test_helpers.RegisterNewValidator(t, app, ctx, _val2)
	val2, err := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr2)
	require.NoError(t, err)
	delAddr1 := addrs[1]
	delAddr2 := addrs[2]

	// Mint alliance tokens
	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(AllianceDenom, sdk.NewInt(2000_000))))
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr1, sdk.NewCoins(sdk.NewCoin(AllianceDenom, sdk.NewInt(1000_000))))
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr2, sdk.NewCoins(sdk.NewCoin(AllianceDenom, sdk.NewInt(1000_000))))
	require.NoError(t, err)

	// First delegate to validator 1
	_, err = app.AllianceKeeper.Delegate(ctx, delAddr1, val1, sdk.NewCoin(AllianceDenom, sdk.NewInt(500_000)))
	require.NoError(t, err)

	// Then redelegate to validator2
	_, err = app.AllianceKeeper.Redelegate(ctx, delAddr1, val1, val2, sdk.NewCoin(AllianceDenom, sdk.NewInt(500_000)))
	require.NoError(t, err)

	_, stop := alliance.RunAllInvariants(ctx, app.AllianceKeeper)
	require.False(t, stop)

	assets := app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)

	// Check delegation share amount
	delegations = app.StakingKeeper.GetAllDelegations(ctx)
	for _, d := range delegations {
		if d.DelegatorAddress == moduleAddr.String() {
			require.Equal(t, stakingtypes.Delegation{
				DelegatorAddress: moduleAddr.String(),
				ValidatorAddress: val2.OperatorAddress,
				Shares:           sdk.NewDec(2_000_000),
			}, d)
		}
	}

	// Check total bonded amount
	require.Equal(t, sdk.NewInt(3_000_000), app.StakingKeeper.TotalBondedTokens(ctx))

	// Check if there is a re-delegation event stored
	iter := app.AllianceKeeper.IterateRedelegationsByDelegator(ctx, delAddr1)
	defer iter.Close()
	require.True(t, iter.Valid())
	for ; iter.Valid(); iter.Next() {
		var redelegation types.Redelegation
		app.AppCodec().MustUnmarshal(iter.Value(), &redelegation)
		require.Equal(t, types.Redelegation{
			DelegatorAddress:    delAddr1.String(),
			SrcValidatorAddress: valAddr1.String(),
			DstValidatorAddress: valAddr2.String(),
			Balance:             sdk.NewCoin(AllianceDenom, sdk.NewInt(500_000)),
		}, redelegation)
	}

	// Check if the delegation objects are correct
	_, found := app.AllianceKeeper.GetDelegation(ctx, delAddr1, valAddr1, AllianceDenom)
	require.False(t, found)
	dstDelegation, found := app.AllianceKeeper.GetDelegation(ctx, delAddr1, valAddr2, AllianceDenom)
	require.Equal(t, types.Delegation{
		DelegatorAddress: delAddr1.String(),
		ValidatorAddress: val2.GetOperator().String(),
		Denom:            AllianceDenom,
		Shares:           sdk.NewDec(500_000),
		RewardHistory:    types.RewardHistories(nil),
	}, dstDelegation)
	require.True(t, found)

	// Check if index by src validator was saved
	iter = app.AllianceKeeper.IterateRedelegationsBySrcValidator(ctx, valAddr1)
	require.True(t, iter.Valid())

	// User then delegates to validator2
	_, err = app.AllianceKeeper.Delegate(ctx, delAddr2, val2, sdk.NewCoin(AllianceDenom, sdk.NewInt(500_000)))
	require.NoError(t, err)

	assets = app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)

	// Then redelegate to validator1 correctly
	_, err = app.AllianceKeeper.Redelegate(ctx, delAddr2, val2, val1, sdk.NewCoin(AllianceDenom, sdk.NewInt(500_000)))
	// Should pass since we removed the re-delegate attempt on x/staking that prevents this
	require.NoError(t, err)

	assets = app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)

	// Immediately calling complete re-delegation should do nothing
	deleted := app.AllianceKeeper.CompleteRedelegations(ctx)
	require.Equal(t, 0, deleted)

	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(app.StakingKeeper.UnbondingTime(ctx)).Add(time.Minute))
	// Calling after re-delegation has matured will delete it from the store
	deleted = app.AllianceKeeper.CompleteRedelegations(ctx)
	require.Equal(t, 2, deleted)

	// There shouldn't be any more delegations in the store
	iter = app.AllianceKeeper.IterateRedelegationsByDelegator(ctx, delAddr1)
	require.False(t, iter.Valid())
	iter = app.AllianceKeeper.IterateRedelegationsByDelegator(ctx, delAddr2)
	require.False(t, iter.Valid())

	// Calling again should not process anymore redelegations
	deleted = app.AllianceKeeper.CompleteRedelegations(ctx)
	require.Equal(t, 0, deleted)
}

func TestRedelegationFailsWithNoDelegations(t *testing.T) {
	app, ctx := createTestContext(t)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			types.NewAllianceAsset(AllianceDenom, sdk.NewDec(2), sdk.ZeroDec(), sdk.NewDec(5), sdk.NewDec(0), ctx.BlockTime()),
			types.NewAllianceAsset(AllianceDenomTwo, sdk.NewDec(10), sdk.NewDec(2), sdk.NewDec(12), sdk.NewDec(0), ctx.BlockTime()),
		},
	})

	// Get all the addresses needed for the test
	delegations := app.StakingKeeper.GetAllDelegations(ctx)
	valAddr1, err := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	require.NoError(t, err)
	val1, err := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr1)
	require.NoError(t, err)
	addrs := test_helpers.AddTestAddrsIncremental(app, ctx, 3, sdk.NewCoins(
		sdk.NewCoin(AllianceDenom, sdk.NewInt(1000_000)),
		sdk.NewCoin(AllianceDenomTwo, sdk.NewInt(1000_000)),
	))
	valAddr2 := sdk.ValAddress(addrs[0])
	_val2 := teststaking.NewValidator(t, valAddr2, test_helpers.CreateTestPubKeys(1)[0])
	test_helpers.RegisterNewValidator(t, app, ctx, _val2)
	val2, err := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr2)
	require.NoError(t, err)
	delAddr1 := addrs[1]
	delAddr2 := addrs[2]

	// Mint alliance tokens
	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(AllianceDenom, sdk.NewInt(2000_000))))
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr1, sdk.NewCoins(sdk.NewCoin(AllianceDenom, sdk.NewInt(1000_000))))
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr2, sdk.NewCoins(sdk.NewCoin(AllianceDenom, sdk.NewInt(1000_000))))
	require.NoError(t, err)

	// User tries to re-delegate without having an initial delegation fails
	_, err = app.AllianceKeeper.Redelegate(ctx, delAddr2, val2, val1, sdk.NewCoin(AllianceDenom, sdk.NewInt(500_000)))
	require.Error(t, err)
}

func TestRedelegationFailsWithTransitiveDelegation(t *testing.T) {
	app, ctx := createTestContext(t)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			types.NewAllianceAsset(AllianceDenom, sdk.NewDec(2), sdk.ZeroDec(), sdk.NewDec(5), sdk.NewDec(0), ctx.BlockTime()),
			types.NewAllianceAsset(AllianceDenomTwo, sdk.NewDec(10), sdk.NewDec(2), sdk.NewDec(12), sdk.NewDec(0), ctx.BlockTime()),
		},
	})

	// Get all the addresses needed for the test
	delegations := app.StakingKeeper.GetAllDelegations(ctx)
	valAddr1, err := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	require.NoError(t, err)
	val1, err := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr1)
	require.NoError(t, err)
	addrs := test_helpers.AddTestAddrsIncremental(app, ctx, 3, sdk.NewCoins(
		sdk.NewCoin(AllianceDenom, sdk.NewInt(1000_000)),
		sdk.NewCoin(AllianceDenomTwo, sdk.NewInt(1000_000)),
	))
	valAddr2 := sdk.ValAddress(addrs[0])
	_val2 := teststaking.NewValidator(t, valAddr2, test_helpers.CreateTestPubKeys(1)[0])
	test_helpers.RegisterNewValidator(t, app, ctx, _val2)
	val2, err := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr2)
	require.NoError(t, err)
	delAddr1 := addrs[1]
	delAddr2 := addrs[2]

	// Mint alliance tokens
	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(AllianceDenom, sdk.NewInt(2000_000))))
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr1, sdk.NewCoins(sdk.NewCoin(AllianceDenom, sdk.NewInt(1000_000))))
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr2, sdk.NewCoins(sdk.NewCoin(AllianceDenom, sdk.NewInt(1000_000))))
	require.NoError(t, err)

	// First delegate to validator 1
	_, err = app.AllianceKeeper.Delegate(ctx, delAddr1, val1, sdk.NewCoin(AllianceDenom, sdk.NewInt(500_000)))
	require.NoError(t, err)

	// Then redelegate to validator2
	_, err = app.AllianceKeeper.Redelegate(ctx, delAddr1, val1, val2, sdk.NewCoin(AllianceDenom, sdk.NewInt(500_000)))
	require.NoError(t, err)

	assets := app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)

	// Should fail when re-delegating back to validator1
	// Same user who re-delegated to from 1 -> 2 cannot re-re-delegate from 2 -> X
	_, err = app.AllianceKeeper.Redelegate(ctx, delAddr1, val2, val1, sdk.NewCoin(AllianceDenom, sdk.NewInt(500_000)))
	require.Error(t, err)
}

func TestRedelegationFailsWithGreaterAmount(t *testing.T) {
	app, ctx := createTestContext(t)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			types.NewAllianceAsset(AllianceDenom, sdk.NewDec(2), sdk.ZeroDec(), sdk.NewDec(5), sdk.NewDec(0), ctx.BlockTime()),
			types.NewAllianceAsset(AllianceDenomTwo, sdk.NewDec(10), sdk.NewDec(2), sdk.NewDec(12), sdk.NewDec(0), ctx.BlockTime()),
		},
	})

	// Get all the addresses needed for the test
	delegations := app.StakingKeeper.GetAllDelegations(ctx)
	valAddr1, err := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	require.NoError(t, err)
	val1, err := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr1)
	require.NoError(t, err)
	addrs := test_helpers.AddTestAddrsIncremental(app, ctx, 3, sdk.NewCoins(
		sdk.NewCoin(AllianceDenom, sdk.NewInt(1000_000)),
		sdk.NewCoin(AllianceDenomTwo, sdk.NewInt(1000_000)),
	))
	valAddr2 := sdk.ValAddress(addrs[0])
	_val2 := teststaking.NewValidator(t, valAddr2, test_helpers.CreateTestPubKeys(1)[0])
	test_helpers.RegisterNewValidator(t, app, ctx, _val2)
	val2, err := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr2)
	require.NoError(t, err)
	delAddr1 := addrs[1]
	delAddr2 := addrs[2]

	// Mint alliance tokens
	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(AllianceDenom, sdk.NewInt(2000_000))))
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr1, sdk.NewCoins(sdk.NewCoin(AllianceDenom, sdk.NewInt(1000_000))))
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr2, sdk.NewCoins(sdk.NewCoin(AllianceDenom, sdk.NewInt(1000_000))))
	require.NoError(t, err)

	// First delegate to validator 1
	_, err = app.AllianceKeeper.Delegate(ctx, delAddr1, val1, sdk.NewCoin(AllianceDenom, sdk.NewInt(500_000)))
	require.NoError(t, err)

	// User then delegates to validator2
	_, err = app.AllianceKeeper.Delegate(ctx, delAddr2, val2, sdk.NewCoin(AllianceDenom, sdk.NewInt(500_000)))
	require.NoError(t, err)

	assets := app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)

	// Then redelegate to validator1 with more than what was delegated but fails
	_, err = app.AllianceKeeper.Redelegate(ctx, delAddr2, val2, val1, sdk.NewCoin(AllianceDenom, sdk.NewInt(1000_000)))
	require.Error(t, err)
}

func TestSuccessfulUndelegation(t *testing.T) {
	app, ctx := createTestContext(t)
	ctx = ctx.WithBlockTime(time.Now())
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			types.NewAllianceAsset(AllianceDenom, sdk.NewDec(2), sdk.ZeroDec(), sdk.NewDec(5), sdk.NewDec(0), ctx.BlockTime()),
			types.NewAllianceAsset(AllianceDenomTwo, sdk.NewDec(10), sdk.NewDec(2), sdk.NewDec(12), sdk.NewDec(0), ctx.BlockTime()),
		},
	})
	delegations := app.StakingKeeper.GetAllDelegations(ctx)
	require.Len(t, delegations, 1)
	unbondingTime := app.StakingKeeper.UnbondingTime(ctx)

	// All the addresses needed
	delAddr, err := sdk.AccAddressFromBech32(delegations[0].DelegatorAddress)
	require.NoError(t, err)
	valAddr, err := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	require.NoError(t, err)
	val, _ := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr)
	moduleAddr := app.AccountKeeper.GetModuleAddress(types.ModuleName)

	// Mint alliance tokens
	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(AllianceDenom, sdk.NewInt(2000_000))))
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr, sdk.NewCoins(sdk.NewCoin(AllianceDenom, sdk.NewInt(2000_000))))
	require.NoError(t, err)

	// Delegate to a validator
	_, err = app.AllianceKeeper.Delegate(ctx, delAddr, val, sdk.NewCoin(AllianceDenom, sdk.NewInt(1000_000)))
	require.NoError(t, err)

	assets := app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)

	// Check total bonded amount
	require.Equal(t, sdk.NewInt(3_000_000), app.StakingKeeper.TotalBondedTokens(ctx))

	// Check that balance dropped
	coin := app.BankKeeper.GetBalance(ctx, delAddr, AllianceDenom)
	require.Equal(t, sdk.NewCoin(AllianceDenom, sdk.NewInt(1000_000)), coin)

	// Check that staked balance increased
	d, _ := app.StakingKeeper.GetDelegation(ctx, moduleAddr, valAddr)
	require.Equal(t, stakingtypes.Delegation{
		DelegatorAddress: moduleAddr.String(),
		ValidatorAddress: valAddr.String(),
		Shares:           sdk.NewDec(2),
	}, d)

	// Immediately undelegate from the validator
	_, err = app.AllianceKeeper.Undelegate(ctx, delAddr, val, sdk.NewCoin(AllianceDenom, sdk.NewInt(250_000)))
	require.NoError(t, err)

	_, err = app.AllianceKeeper.Undelegate(ctx, delAddr, val, sdk.NewCoin(AllianceDenom, sdk.NewInt(250_000)))
	require.NoError(t, err)

	// Check if undelegations were stored correctly
	iter := app.AllianceKeeper.IterateUndelegationsByCompletionTime(ctx, ctx.BlockTime().Add(unbondingTime).Add(time.Second))
	require.True(t, iter.Valid())
	var queuedUndelegations types.QueuedUndelegation
	b := iter.Value()
	app.AppCodec().MustUnmarshal(b, &queuedUndelegations)
	require.Equal(t, types.QueuedUndelegation{Entries: []*types.Undelegation{
		{
			DelegatorAddress: delAddr.String(),
			ValidatorAddress: val.GetOperator().String(),
			Balance:          sdk.NewCoin(AllianceDenom, sdk.NewInt(250_000)),
		},
		{
			DelegatorAddress: delAddr.String(),
			ValidatorAddress: val.GetOperator().String(),
			Balance:          sdk.NewCoin(AllianceDenom, sdk.NewInt(250_000)),
		},
	}}, queuedUndelegations)

	assets = app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)

	// Check total bonded amount
	require.Equal(t, sdk.NewInt(3_000_000), app.StakingKeeper.TotalBondedTokens(ctx))

	// Check that staked balance stays the same
	d, _ = app.StakingKeeper.GetDelegation(ctx, moduleAddr, valAddr)
	require.Equal(t, stakingtypes.Delegation{
		DelegatorAddress: moduleAddr.String(),
		ValidatorAddress: valAddr.String(),
		Shares:           sdk.NewDec(2),
	}, d)

	// Immediately try to complete undelegation
	err = app.AllianceKeeper.CompleteUndelegations(ctx)
	require.NoError(t, err)

	// Check that balance stayed the same
	coin = app.BankKeeper.GetBalance(ctx, delAddr, AllianceDenom)
	require.Equal(t, sdk.NewCoin(AllianceDenom, sdk.NewInt(1000_000)), coin)

	// Advance time to after unbonding period
	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(app.StakingKeeper.UnbondingTime(ctx)).Add(time.Minute))

	err = app.AllianceKeeper.CompleteUndelegations(ctx)
	require.NoError(t, err)

	// Check that balance increased
	coin = app.BankKeeper.GetBalance(ctx, delAddr, AllianceDenom)
	require.Equal(t, sdk.NewCoin(AllianceDenom, sdk.NewInt(1500_000)), coin)

	// Completing again should not process anymore undelegations
	err = app.AllianceKeeper.CompleteUndelegations(ctx)
	require.NoError(t, err)
}

func TestUndelegationWithoutDelegation(t *testing.T) {
	app, ctx := createTestContext(t)
	ctx = ctx.WithBlockTime(time.Now())
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			types.NewAllianceAsset(AllianceDenom, sdk.NewDec(2), sdk.ZeroDec(), sdk.NewDec(5), sdk.NewDec(0), ctx.BlockTime()),
			types.NewAllianceAsset(AllianceDenomTwo, sdk.NewDec(10), sdk.NewDec(2), sdk.NewDec(12), sdk.NewDec(0), ctx.BlockTime()),
		},
	})
	delegations := app.StakingKeeper.GetAllDelegations(ctx)
	require.Len(t, delegations, 1)

	// All the addresses needed
	delAddr, err := sdk.AccAddressFromBech32(delegations[0].DelegatorAddress)
	require.NoError(t, err)
	valAddr, err := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	require.NoError(t, err)
	val, _ := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr)

	// Mint alliance tokens
	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(AllianceDenom, sdk.NewInt(2000_000))))
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr, sdk.NewCoins(sdk.NewCoin(AllianceDenom, sdk.NewInt(2000_000))))
	require.NoError(t, err)

	// Undelegating without a delegation will fail
	_, err = app.AllianceKeeper.Undelegate(ctx, delAddr, val, sdk.NewCoin(AllianceDenom, sdk.NewInt(1000_000)))
	require.Error(t, err)
}

func TestUndelegateAfterClaimingTakeRate(t *testing.T) {
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	params := types.DefaultParams()
	params.LastTakeRateClaimTime = startTime
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: params,
		Assets: []types.AllianceAsset{
			types.NewAllianceAsset(AllianceDenom, sdk.NewDec(2), sdk.ZeroDec(), sdk.NewDec(5), sdk.NewDec(0), ctx.BlockTime()),
			types.NewAllianceAsset(AllianceDenomTwo, sdk.NewDec(10), sdk.NewDec(2), sdk.NewDec(12), sdk.MustNewDecFromStr("0.5"), ctx.BlockTime()),
		},
	})
	queryServer := keeper.NewQueryServerImpl(app.AllianceKeeper)

	// Set tax and rewards to be zero for easier calculation
	distParams := app.DistrKeeper.GetParams(ctx)
	distParams.CommunityTax = sdk.ZeroDec()

	err := app.DistrKeeper.SetParams(ctx, distParams)
	require.NoError(t, err)

	// Accounts

	// rewardsPoolAddr := app.AccountKeeper.GetModuleAddress(types.RewardsPoolName)
	addrs := test_helpers.AddTestAddrsIncremental(app, ctx, 4, sdk.NewCoins(
		sdk.NewCoin(AllianceDenom, sdk.NewInt(1000_000_000)),
		sdk.NewCoin(AllianceDenomTwo, sdk.NewInt(2000_000_000)),
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

	// Delegate token with non-zero take_rate
	_, err = app.AllianceKeeper.Delegate(ctx, user1, val1, sdk.NewCoin(AllianceDenomTwo, sdk.NewInt(1000_000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user2, val2, sdk.NewCoin(AllianceDenomTwo, sdk.NewInt(1000_000_000)))
	require.NoError(t, err)

	assets := app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)
	// Check total bonded amount
	require.Equal(t, sdk.NewInt(11_000_000), app.StakingKeeper.TotalBondedTokens(ctx))

	ctx = ctx.WithBlockTime(startTime.Add(time.Minute * 6)).WithBlockHeight(2)
	coins, err := app.AllianceKeeper.DeductAssetsHook(ctx, assets)
	require.NoError(t, err)
	require.False(t, coins.IsZero())

	res, err := queryServer.AllianceDelegation(ctx, &types.QueryAllianceDelegationRequest{
		DelegatorAddr: user1.String(),
		ValidatorAddr: val1.GetOperator().String(),
		Denom:         AllianceDenomTwo,
		Pagination:    nil,
	})
	require.NoError(t, err)
	del := res.GetDelegation()
	require.True(t, del.GetBalance().Amount.LT(sdk.NewInt(1000_000_000)), "%s should be less than %s", del.GetBalance().Amount, sdk.NewInt(1000_000_000))
	// Undelegate token with initial amount should fail
	_, err = app.AllianceKeeper.Undelegate(ctx, user1, val1, sdk.NewCoin(AllianceDenomTwo, sdk.NewInt(1000_000_000)))
	require.Error(t, err)

	// Undelegate token with current amount should pass
	_, err = app.AllianceKeeper.Undelegate(ctx, user1, val1, sdk.NewCoin(AllianceDenomTwo, del.Balance.Amount))
	require.NoError(t, err)

	// User should have everything withdrawn
	_, found := app.AllianceKeeper.GetDelegation(ctx, user1, valAddr1, AllianceDenomTwo)
	require.False(t, found)

	// Delegate again
	_, err = app.AllianceKeeper.Delegate(ctx, user1, val1, sdk.NewCoin(AllianceDenomTwo, sdk.NewInt(500_000_000)))
	require.NoError(t, err)

	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(time.Minute * 1)).WithBlockHeight(2)

	_, err = app.AllianceKeeper.Delegate(ctx, user1, val1, sdk.NewCoin(AllianceDenomTwo, sdk.NewInt(400_000_000)))
	require.NoError(t, err)

	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(time.Minute * 5)).WithBlockHeight(3)
	coins, err = app.AllianceKeeper.DeductAssetsHook(ctx, assets)
	require.NoError(t, err)
	require.False(t, coins.IsZero())

	res, err = queryServer.AllianceDelegation(ctx, &types.QueryAllianceDelegationRequest{
		DelegatorAddr: user1.String(),
		ValidatorAddr: val1.GetOperator().String(),
		Denom:         AllianceDenomTwo,
		Pagination:    nil,
	})
	require.NoError(t, err)
	del = res.GetDelegation()
	require.True(t, del.GetBalance().Amount.LT(sdk.NewInt(900_000_000)), "%s should be less than %s", del.GetBalance().Amount, sdk.NewInt(1000_000_000))

	// Undelegate token with current amount should pass
	_, err = app.AllianceKeeper.Undelegate(ctx, user1, val1, sdk.NewCoin(AllianceDenomTwo, del.Balance.Amount))
	require.NoError(t, err)

	// User should have everything withdrawn
	_, found = app.AllianceKeeper.GetDelegation(ctx, user1, valAddr1, AllianceDenomTwo)
	require.False(t, found)

	res, err = queryServer.AllianceDelegation(ctx, &types.QueryAllianceDelegationRequest{
		DelegatorAddr: user1.String(),
		ValidatorAddr: val1.GetOperator().String(),
		Denom:         AllianceDenomTwo,
		Pagination:    nil,
	})
	require.NoError(t, err)
	del = res.GetDelegation()
	require.True(t, del.Balance.Amount.IsZero())
}

func TestDelegationWithNativeStakingChanges(t *testing.T) {
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			types.NewAllianceAsset(AllianceDenom, sdk.NewDec(2), sdk.ZeroDec(), sdk.NewDec(5), sdk.NewDec(0), ctx.BlockTime()),
			types.NewAllianceAsset(AllianceDenomTwo, sdk.NewDec(10), sdk.NewDec(2), sdk.NewDec(12), sdk.MustNewDecFromStr("0.5"), ctx.BlockTime()),
		},
	})

	// Set tax and rewards to be zero for easier calculation
	distParams := app.DistrKeeper.GetParams(ctx)
	distParams.CommunityTax = sdk.ZeroDec()

	err := app.DistrKeeper.SetParams(ctx, distParams)
	require.NoError(t, err)

	// Accounts

	// rewardsPoolAddr := app.AccountKeeper.GetModuleAddress(types.RewardsPoolName)
	bondDenom := app.StakingKeeper.BondDenom(ctx)
	addrs := test_helpers.AddTestAddrsIncremental(app, ctx, 4, sdk.NewCoins(
		sdk.NewCoin(bondDenom, sdk.NewInt(1000_000_000)),
		sdk.NewCoin(AllianceDenom, sdk.NewInt(1000_000_000)),
		sdk.NewCoin(AllianceDenomTwo, sdk.NewInt(2000_000_000)),
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

	// Stake some alliance tokens
	_, err = app.AllianceKeeper.Delegate(ctx, user1, val1, sdk.NewCoin(AllianceDenom, sdk.NewInt(2000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user2, val2, sdk.NewCoin(AllianceDenomTwo, sdk.NewInt(2000_000)))
	require.NoError(t, err)
	assets := app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)

	// Check total bonded tokens
	require.Equal(t, sdk.NewInt(13_000_000), app.StakingKeeper.TotalBondedTokens(ctx))

	// Stake some native tokens
	_, err = app.StakingKeeper.Delegate(ctx, user2, sdk.NewInt(2000_000), stakingtypes.Unbonded, *val2.Validator, true)
	require.NoError(t, err)
	// Check total bonded tokens
	require.Equal(t, sdk.NewInt(15_000_000), app.StakingKeeper.TotalBondedTokens(ctx))

	assets = app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)
	// Check total bonded tokens
	require.Equal(t, sdk.NewInt(39_000_000), app.StakingKeeper.TotalBondedTokens(ctx))

	// Redelegate some native tokens
	_, err = app.StakingKeeper.BeginRedelegation(ctx, user2, valAddr2, valAddr1, sdk.NewDec(1))
	require.NoError(t, err)

	assets = app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)
	// Check total bonded tokens
	require.Equal(t, sdk.NewInt(39_000_000), app.StakingKeeper.TotalBondedTokens(ctx))

	// Undelegate some native tokens
	shares, _ := app.StakingKeeper.ValidateUnbondAmount(ctx, user2, valAddr2, sdk.NewInt(1000_000))
	_, err = app.StakingKeeper.Undelegate(ctx, user2, valAddr2, shares)
	require.NoError(t, err)

	assets = app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)
	// Check total bonded tokens
	require.Equal(t, sdk.NewInt(26_000_000), app.StakingKeeper.TotalBondedTokens(ctx))
}
