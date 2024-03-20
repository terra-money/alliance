package tests_test

import (
	"testing"
	"time"

	"cosmossdk.io/math"

	test_helpers "github.com/terra-money/alliance/app"
	"github.com/terra-money/alliance/x/alliance"
	"github.com/terra-money/alliance/x/alliance/keeper"
	"github.com/terra-money/alliance/x/alliance/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
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
			types.NewAllianceAsset(AllianceDenom, math.LegacyNewDec(2), math.LegacyZeroDec(), math.LegacyNewDec(5), math.LegacyNewDec(0), genesisTime),
			types.NewAllianceAsset(AllianceDenomTwo, math.LegacyNewDec(10), math.LegacyNewDec(2), math.LegacyNewDec(12), math.LegacyNewDec(0), genesisTime),
		},
	})
	delegations, err := app.StakingKeeper.GetAllDelegations(ctx)
	require.NoError(t, err)
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
	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(AllianceDenom, math.NewInt(2000_000))))
	require.NoError(t, err)
	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(AllianceDenomTwo, math.NewInt(2000_000))))
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr, sdk.NewCoins(sdk.NewCoin(AllianceDenom, math.NewInt(2000_000))))
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr, sdk.NewCoins(sdk.NewCoin(AllianceDenomTwo, math.NewInt(2000_000))))
	require.NoError(t, err)

	// Check current total staked tokens
	totalBonded, err := app.StakingKeeper.TotalBondedTokens(ctx)
	require.NoError(t, err)
	require.Equal(t, math.NewInt(1000_000), totalBonded)

	// Delegate
	_, err = app.AllianceKeeper.Delegate(ctx, delAddr, val, sdk.NewCoin(AllianceDenom, math.NewInt(1000_000)))
	require.NoError(t, err)

	// Manually trigger rebalancing
	app.AllianceKeeper.GetAssetByDenom(ctx, AllianceDenom)
	assets := app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)

	allianceBonded, err := app.StakingKeeper.GetDelegatorBonded(ctx, moduleAddr)
	require.NoError(t, err)
	// Total ALLIANCE tokens should be 2 * totalBonded
	require.Equal(t, totalBonded.Mul(math.NewInt(2)).String(), allianceBonded.String())

	// Check delegation in staking module
	delegations, err = app.StakingKeeper.GetAllDelegations(ctx)
	require.NoError(t, err)
	require.Len(t, delegations, 2)
	i := slices.IndexFunc(delegations, func(d stakingtypes.Delegation) bool {
		return d.DelegatorAddress == moduleAddr.String()
	})
	newDelegation := delegations[i]
	require.Equal(t, stakingtypes.Delegation{
		DelegatorAddress: moduleAddr.String(),
		ValidatorAddress: valAddr.String(),
		Shares:           math.LegacyNewDec(2),
	}, newDelegation)

	// Check delegation in alliance module
	allianceDelegation, found := app.AllianceKeeper.GetDelegation(ctx, delAddr, valAddr, AllianceDenom)
	require.True(t, found)
	require.Equal(t, types.Delegation{
		DelegatorAddress: delAddr.String(),
		ValidatorAddress: val.OperatorAddress,
		Denom:            AllianceDenom,
		Shares:           math.LegacyNewDec(1000_000),
		RewardHistory:    types.RewardHistories(nil),
	}, allianceDelegation)

	// Check asset
	asset, _ := app.AllianceKeeper.GetAssetByDenom(ctx, AllianceDenom)
	require.Equal(t, types.AllianceAsset{
		Denom:                AllianceDenom,
		RewardWeight:         math.LegacyNewDec(2),
		RewardWeightRange:    types.RewardWeightRange{Min: math.LegacyNewDec(0), Max: math.LegacyNewDec(5)},
		TakeRate:             math.LegacyNewDec(0),
		TotalTokens:          math.NewInt(1000_000),
		TotalValidatorShares: math.LegacyNewDec(1000_000),
		RewardStartTime:      genesisTime,
		RewardChangeRate:     math.LegacyOneDec(),
		RewardChangeInterval: 0,
	}, asset)

	// Delegate with same denom again
	_, err = app.AllianceKeeper.Delegate(ctx, delAddr, val, sdk.NewCoin(AllianceDenom, math.NewInt(1000_000)))
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
		Shares:           math.LegacyNewDec(2000_000),
		RewardHistory:    types.RewardHistories(nil),
	}, allianceDelegation)

	// Check asset again
	asset, _ = app.AllianceKeeper.GetAssetByDenom(ctx, AllianceDenom)
	require.Equal(t, types.AllianceAsset{
		Denom:                AllianceDenom,
		RewardWeight:         math.LegacyNewDec(2),
		RewardWeightRange:    types.RewardWeightRange{Min: math.LegacyNewDec(0), Max: math.LegacyNewDec(5)},
		TakeRate:             math.LegacyNewDec(0),
		TotalTokens:          math.NewInt(2000_000),
		TotalValidatorShares: math.LegacyNewDec(2000_000),
		RewardStartTime:      genesisTime,
		RewardChangeRate:     math.LegacyOneDec(),
		RewardChangeInterval: 0,
	}, asset)

	// Check delegation in staking module total shares should not change
	delegations, err = app.StakingKeeper.GetAllDelegations(ctx)
	require.NoError(t, err)
	require.Len(t, delegations, 2)
	i = slices.IndexFunc(delegations, func(d stakingtypes.Delegation) bool {
		return d.DelegatorAddress == moduleAddr.String()
	})
	newDelegation = delegations[i]
	require.Equal(t, stakingtypes.Delegation{
		DelegatorAddress: moduleAddr.String(),
		ValidatorAddress: valAddr.String(),
		Shares:           math.LegacyNewDec(2),
	}, newDelegation)
}

func TestDelegationWithMultipleAssets(t *testing.T) {
	app, ctx := createTestContext(t)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			types.NewAllianceAsset(AllianceDenom, math.LegacyNewDec(2), math.LegacyZeroDec(), math.LegacyNewDec(5), math.LegacyNewDec(0), ctx.BlockTime()),
			types.NewAllianceAsset(AllianceDenomTwo, math.LegacyNewDec(10), math.LegacyNewDec(2), math.LegacyNewDec(12), math.LegacyNewDec(0), ctx.BlockTime()),
		},
	})
	delegations, err := app.StakingKeeper.GetAllDelegations(ctx)
	require.NoError(t, err)
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
	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(AllianceDenom, math.NewInt(2000_000))))
	require.NoError(t, err)
	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(AllianceDenomTwo, math.NewInt(2000_000))))
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr, sdk.NewCoins(sdk.NewCoin(AllianceDenom, math.NewInt(2000_000))))
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr, sdk.NewCoins(sdk.NewCoin(AllianceDenomTwo, math.NewInt(2000_000))))
	require.NoError(t, err)

	// Delegate
	_, err = app.AllianceKeeper.Delegate(ctx, delAddr, val, sdk.NewCoin(AllianceDenom, math.NewInt(2000_000)))
	require.NoError(t, err)
	// Delegate with another denom
	_, err = app.AllianceKeeper.Delegate(ctx, delAddr, val, sdk.NewCoin(AllianceDenomTwo, math.NewInt(1000_000)))
	require.NoError(t, err)

	// Manually trigger rebalancing
	app.AllianceKeeper.GetAssetByDenom(ctx, AllianceDenomTwo)
	assets := app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)

	// Check delegation in staking module
	delegations, err = app.StakingKeeper.GetAllDelegations(ctx)
	require.NoError(t, err)
	require.Len(t, delegations, 2)
	i := slices.IndexFunc(delegations, func(d stakingtypes.Delegation) bool {
		return d.DelegatorAddress == moduleAddr.String()
	})
	newDelegation := delegations[i]
	// 1 * 2 + 1 * 10 = 12
	require.Equal(t, stakingtypes.Delegation{
		DelegatorAddress: moduleAddr.String(),
		ValidatorAddress: valAddr.String(),
		Shares:           math.LegacyNewDec(12),
	}, newDelegation)

	// Check validator in x/staking
	val, err = app.AllianceKeeper.GetAllianceValidator(ctx, valAddr)
	require.NoError(t, err)
	require.Equal(t, math.LegacyNewDec(13), val.DelegatorShares)
}

func TestDelegationWithUnknownAssets(t *testing.T) {
	app, ctx := createTestContext(t)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			types.NewAllianceAsset(AllianceDenom, math.LegacyNewDec(2), math.LegacyZeroDec(), math.LegacyNewDec(5), math.LegacyNewDec(0), ctx.BlockTime()),
			types.NewAllianceAsset(AllianceDenomTwo, math.LegacyNewDec(10), math.LegacyNewDec(2), math.LegacyNewDec(12), math.LegacyNewDec(0), ctx.BlockTime()),
		},
	})
	delegations, err := app.StakingKeeper.GetAllDelegations(ctx)
	require.NoError(t, err)
	require.Len(t, delegations, 1)

	// All the addresses needed
	delAddr, err := sdk.AccAddressFromBech32(delegations[0].DelegatorAddress)
	require.NoError(t, err)
	valAddr, err := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	require.NoError(t, err)
	val, err := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr)
	require.NoError(t, err)

	// Mint alliance tokens
	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin("UNKNOWN", math.NewInt(2000_000))))
	require.NoError(t, err)

	_, err = app.AllianceKeeper.Delegate(ctx, delAddr, val, sdk.NewCoin("UNKNOWN", math.NewInt(2000_000)))
	require.Error(t, err)
}

func TestSuccessfulRedelegation(t *testing.T) {
	app, ctx := createTestContext(t)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			types.NewAllianceAsset(AllianceDenom, math.LegacyNewDec(2), math.LegacyZeroDec(), math.LegacyNewDec(5), math.LegacyNewDec(0), ctx.BlockTime()),
			types.NewAllianceAsset(AllianceDenomTwo, math.LegacyNewDec(10), math.LegacyNewDec(2), math.LegacyNewDec(12), math.LegacyNewDec(0), ctx.BlockTime()),
		},
	})
	queryServer := keeper.NewQueryServerImpl(app.AllianceKeeper)

	// Get all the addresses needed for the test
	moduleAddr := app.AccountKeeper.GetModuleAddress(types.ModuleName)
	delegations, err := app.StakingKeeper.GetAllDelegations(ctx)
	require.NoError(t, err)
	valAddr1, err := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	require.NoError(t, err)
	val1, err := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr1)
	require.NoError(t, err)
	addrs := test_helpers.AddTestAddrsIncremental(app, ctx, 3, sdk.NewCoins(
		sdk.NewCoin(AllianceDenom, math.NewInt(1000_000)),
		sdk.NewCoin(AllianceDenomTwo, math.NewInt(1000_000)),
	))
	valAddr2 := sdk.ValAddress(addrs[0])
	_val2 := teststaking.NewValidator(t, valAddr2, test_helpers.CreateTestPubKeys(1)[0])
	test_helpers.RegisterNewValidator(t, app, ctx, _val2)
	val2, err := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr2)
	require.NoError(t, err)
	delAddr1 := addrs[1]
	delAddr2 := addrs[2]

	// Mint alliance tokens
	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(AllianceDenom, math.NewInt(2000_000))))
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr1, sdk.NewCoins(sdk.NewCoin(AllianceDenom, math.NewInt(1000_000))))
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr2, sdk.NewCoins(sdk.NewCoin(AllianceDenom, math.NewInt(1000_000))))
	require.NoError(t, err)

	// First delegate to validator 1
	_, err = app.AllianceKeeper.Delegate(ctx, delAddr1, val1, sdk.NewCoin(AllianceDenom, math.NewInt(500_000)))
	require.NoError(t, err)

	// Then redelegate to validator2
	_, err = app.AllianceKeeper.Redelegate(ctx, delAddr1, val1, val2, sdk.NewCoin(AllianceDenom, math.NewInt(500_000)))
	require.NoError(t, err)

	_, stop := alliance.RunAllInvariants(ctx, app.AllianceKeeper)
	require.False(t, stop)

	assets := app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)

	// Check delegation share amount
	delegations, err = app.StakingKeeper.GetAllDelegations(ctx)
	require.NoError(t, err)
	for _, d := range delegations {
		if d.DelegatorAddress == moduleAddr.String() {
			require.Equal(t, stakingtypes.Delegation{
				DelegatorAddress: moduleAddr.String(),
				ValidatorAddress: val2.OperatorAddress,
				Shares:           math.LegacyNewDec(2_000_000),
			}, d)
		}
	}

	totalBonded, err := app.StakingKeeper.TotalBondedTokens(ctx)
	require.NoError(t, err)
	// Check total bonded amount
	require.Equal(t, math.NewInt(3_000_000), totalBonded)

	// Query redelegations by the delegator address
	redelegationsByDelegator, err := queryServer.AllianceRedelegationsByDelegator(ctx, &types.QueryAllianceRedelegationsByDelegatorRequest{
		DelegatorAddr: delAddr1.String(),
	})
	require.NoError(t, err)
	require.Equal(t, &types.QueryAllianceRedelegationsByDelegatorResponse{
		Redelegations: []types.RedelegationEntry{
			{
				DelegatorAddress:    delAddr1.String(),
				SrcValidatorAddress: valAddr1.String(),
				DstValidatorAddress: valAddr2.String(),
				Balance:             sdk.NewCoin(AllianceDenom, math.NewInt(500_000)),
				CompletionTime:      time.Date(1, time.January, 22, 0, 0, 0, 0, time.UTC),
			},
		},
		Pagination: &query.PageResponse{
			Total: 1,
		},
	}, redelegationsByDelegator)

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
			Balance:             sdk.NewCoin(AllianceDenom, math.NewInt(500_000)),
		}, redelegation)
	}

	// Check if the delegation objects are correct
	_, found := app.AllianceKeeper.GetDelegation(ctx, delAddr1, valAddr1, AllianceDenom)
	require.False(t, found)
	dstDelegation, found := app.AllianceKeeper.GetDelegation(ctx, delAddr1, valAddr2, AllianceDenom)
	require.Equal(t, types.Delegation{
		DelegatorAddress: delAddr1.String(),
		ValidatorAddress: val2.GetOperator(),
		Denom:            AllianceDenom,
		Shares:           math.LegacyNewDec(500_000),
		RewardHistory:    types.RewardHistories(nil),
	}, dstDelegation)
	require.True(t, found)

	// Check if index by src validator was saved
	iter = app.AllianceKeeper.IterateRedelegationsBySrcValidator(ctx, valAddr1)
	require.True(t, iter.Valid())

	// User then delegates to validator2
	_, err = app.AllianceKeeper.Delegate(ctx, delAddr2, val2, sdk.NewCoin(AllianceDenom, math.NewInt(500_000)))
	require.NoError(t, err)

	assets = app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)

	// Then redelegate to validator1 correctly
	_, err = app.AllianceKeeper.Redelegate(ctx, delAddr2, val2, val1, sdk.NewCoin(AllianceDenom, math.NewInt(500_000)))
	// Should pass since we removed the re-delegate attempt on x/staking that prevents this
	require.NoError(t, err)

	assets = app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)

	// Query all redelegations
	redelegationsRes, err := queryServer.AllianceRedelegations(ctx,
		&types.QueryAllianceRedelegationsRequest{
			Denom:         AllianceDenom,
			DelegatorAddr: delAddr1.String(),
			Pagination:    nil,
		})
	require.NoError(t, err)
	require.Equal(t, &types.QueryAllianceRedelegationsResponse{
		Redelegations: []types.RedelegationEntry{
			{
				DelegatorAddress:    delAddr1.String(),
				SrcValidatorAddress: valAddr1.String(),
				DstValidatorAddress: valAddr2.String(),
				Balance:             sdk.NewCoin(AllianceDenom, math.NewInt(500_000)),
				CompletionTime:      time.Date(1, time.January, 22, 0, 0, 0, 0, time.UTC),
			},
		},
		Pagination: &query.PageResponse{
			Total: 1,
		},
	}, redelegationsRes)

	unbondingPeriod, err := app.StakingKeeper.UnbondingTime(ctx)
	require.NoError(t, err)
	require.Equal(t, types.RedelegationEntry{
		DelegatorAddress:    delAddr1.String(),
		SrcValidatorAddress: valAddr1.String(),
		DstValidatorAddress: valAddr2.String(),
		Balance:             sdk.NewCoin(AllianceDenom, math.NewInt(500_000)),
		CompletionTime:      ctx.BlockTime().Add(unbondingPeriod),
	}, redelegationsRes.Redelegations[0])

	// Immediately calling complete re-delegation should do nothing
	deleted := app.AllianceKeeper.CompleteRedelegations(ctx)
	require.Equal(t, 0, deleted)

	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(unbondingPeriod).Add(time.Minute))
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
			types.NewAllianceAsset(AllianceDenom, math.LegacyNewDec(2), math.LegacyZeroDec(), math.LegacyNewDec(5), math.LegacyNewDec(0), ctx.BlockTime()),
			types.NewAllianceAsset(AllianceDenomTwo, math.LegacyNewDec(10), math.LegacyNewDec(2), math.LegacyNewDec(12), math.LegacyNewDec(0), ctx.BlockTime()),
		},
	})

	// Get all the addresses needed for the test
	delegations, err := app.StakingKeeper.GetAllDelegations(ctx)
	require.NoError(t, err)
	valAddr1, err := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	require.NoError(t, err)
	val1, err := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr1)
	require.NoError(t, err)
	addrs := test_helpers.AddTestAddrsIncremental(app, ctx, 3, sdk.NewCoins(
		sdk.NewCoin(AllianceDenom, math.NewInt(1000_000)),
		sdk.NewCoin(AllianceDenomTwo, math.NewInt(1000_000)),
	))
	valAddr2 := sdk.ValAddress(addrs[0])
	_val2 := teststaking.NewValidator(t, valAddr2, test_helpers.CreateTestPubKeys(1)[0])
	test_helpers.RegisterNewValidator(t, app, ctx, _val2)
	val2, err := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr2)
	require.NoError(t, err)
	delAddr1 := addrs[1]
	delAddr2 := addrs[2]

	// Mint alliance tokens
	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(AllianceDenom, math.NewInt(2000_000))))
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr1, sdk.NewCoins(sdk.NewCoin(AllianceDenom, math.NewInt(1000_000))))
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr2, sdk.NewCoins(sdk.NewCoin(AllianceDenom, math.NewInt(1000_000))))
	require.NoError(t, err)

	// User tries to re-delegate without having an initial delegation fails
	_, err = app.AllianceKeeper.Redelegate(ctx, delAddr2, val2, val1, sdk.NewCoin(AllianceDenom, math.NewInt(500_000)))
	require.Error(t, err)
}

func TestRedelegationFailsWithTransitiveDelegation(t *testing.T) {
	app, ctx := createTestContext(t)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			types.NewAllianceAsset(AllianceDenom, math.LegacyNewDec(2), math.LegacyZeroDec(), math.LegacyNewDec(5), math.LegacyNewDec(0), ctx.BlockTime()),
			types.NewAllianceAsset(AllianceDenomTwo, math.LegacyNewDec(10), math.LegacyNewDec(2), math.LegacyNewDec(12), math.LegacyNewDec(0), ctx.BlockTime()),
		},
	})

	// Get all the addresses needed for the test
	delegations, err := app.StakingKeeper.GetAllDelegations(ctx)
	require.NoError(t, err)
	valAddr1, err := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	require.NoError(t, err)
	val1, err := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr1)
	require.NoError(t, err)
	addrs := test_helpers.AddTestAddrsIncremental(app, ctx, 3, sdk.NewCoins(
		sdk.NewCoin(AllianceDenom, math.NewInt(1000_000)),
		sdk.NewCoin(AllianceDenomTwo, math.NewInt(1000_000)),
	))
	valAddr2 := sdk.ValAddress(addrs[0])
	_val2 := teststaking.NewValidator(t, valAddr2, test_helpers.CreateTestPubKeys(1)[0])
	test_helpers.RegisterNewValidator(t, app, ctx, _val2)
	val2, err := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr2)
	require.NoError(t, err)
	delAddr1 := addrs[1]
	delAddr2 := addrs[2]

	// Mint alliance tokens
	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(AllianceDenom, math.NewInt(2000_000))))
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr1, sdk.NewCoins(sdk.NewCoin(AllianceDenom, math.NewInt(1000_000))))
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr2, sdk.NewCoins(sdk.NewCoin(AllianceDenom, math.NewInt(1000_000))))
	require.NoError(t, err)

	// First delegate to validator 1
	_, err = app.AllianceKeeper.Delegate(ctx, delAddr1, val1, sdk.NewCoin(AllianceDenom, math.NewInt(500_000)))
	require.NoError(t, err)

	// Then redelegate to validator2
	_, err = app.AllianceKeeper.Redelegate(ctx, delAddr1, val1, val2, sdk.NewCoin(AllianceDenom, math.NewInt(500_000)))
	require.NoError(t, err)

	assets := app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)

	// Should fail when re-delegating back to validator1
	// Same user who re-delegated to from 1 -> 2 cannot re-re-delegate from 2 -> X
	_, err = app.AllianceKeeper.Redelegate(ctx, delAddr1, val2, val1, sdk.NewCoin(AllianceDenom, math.NewInt(500_000)))
	require.Error(t, err)
}

func TestRedelegationFailsWithGreaterAmount(t *testing.T) {
	app, ctx := createTestContext(t)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			types.NewAllianceAsset(AllianceDenom, math.LegacyNewDec(2), math.LegacyZeroDec(), math.LegacyNewDec(5), math.LegacyNewDec(0), ctx.BlockTime()),
			types.NewAllianceAsset(AllianceDenomTwo, math.LegacyNewDec(10), math.LegacyNewDec(2), math.LegacyNewDec(12), math.LegacyNewDec(0), ctx.BlockTime()),
		},
	})

	// Get all the addresses needed for the test
	delegations, err := app.StakingKeeper.GetAllDelegations(ctx)
	require.NoError(t, err)
	valAddr1, err := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	require.NoError(t, err)
	val1, err := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr1)
	require.NoError(t, err)
	addrs := test_helpers.AddTestAddrsIncremental(app, ctx, 3, sdk.NewCoins(
		sdk.NewCoin(AllianceDenom, math.NewInt(1000_000)),
		sdk.NewCoin(AllianceDenomTwo, math.NewInt(1000_000)),
	))
	valAddr2 := sdk.ValAddress(addrs[0])
	_val2 := teststaking.NewValidator(t, valAddr2, test_helpers.CreateTestPubKeys(1)[0])
	test_helpers.RegisterNewValidator(t, app, ctx, _val2)
	val2, err := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr2)
	require.NoError(t, err)
	delAddr1 := addrs[1]
	delAddr2 := addrs[2]

	// Mint alliance tokens
	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(AllianceDenom, math.NewInt(2000_000))))
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr1, sdk.NewCoins(sdk.NewCoin(AllianceDenom, math.NewInt(1000_000))))
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr2, sdk.NewCoins(sdk.NewCoin(AllianceDenom, math.NewInt(1000_000))))
	require.NoError(t, err)

	// First delegate to validator 1
	_, err = app.AllianceKeeper.Delegate(ctx, delAddr1, val1, sdk.NewCoin(AllianceDenom, math.NewInt(500_000)))
	require.NoError(t, err)

	// User then delegates to validator2
	_, err = app.AllianceKeeper.Delegate(ctx, delAddr2, val2, sdk.NewCoin(AllianceDenom, math.NewInt(500_000)))
	require.NoError(t, err)

	assets := app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)

	// Then redelegate to validator1 with more than what was delegated but fails
	_, err = app.AllianceKeeper.Redelegate(ctx, delAddr2, val2, val1, sdk.NewCoin(AllianceDenom, math.NewInt(1000_000)))
	require.Error(t, err)
}

func TestSuccessfulUndelegation(t *testing.T) {
	app, ctx := createTestContext(t)
	ctx = ctx.WithBlockTime(time.Now())
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			types.NewAllianceAsset(AllianceDenom, math.LegacyNewDec(2), math.LegacyZeroDec(), math.LegacyNewDec(5), math.LegacyNewDec(0), ctx.BlockTime()),
			types.NewAllianceAsset(AllianceDenomTwo, math.LegacyNewDec(10), math.LegacyNewDec(2), math.LegacyNewDec(12), math.LegacyNewDec(0), ctx.BlockTime()),
		},
	})
	delegations, err := app.StakingKeeper.GetAllDelegations(ctx)
	require.NoError(t, err)
	require.Len(t, delegations, 1)
	unbondingTime, err := app.StakingKeeper.UnbondingTime(ctx)
	require.NoError(t, err)

	// All the addresses needed
	delAddr, err := sdk.AccAddressFromBech32(delegations[0].DelegatorAddress)
	require.NoError(t, err)
	valAddr, err := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	require.NoError(t, err)
	val, _ := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr)
	moduleAddr := app.AccountKeeper.GetModuleAddress(types.ModuleName)

	// Mint alliance tokens
	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(AllianceDenom, math.NewInt(2000_000))))
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr, sdk.NewCoins(sdk.NewCoin(AllianceDenom, math.NewInt(2000_000))))
	require.NoError(t, err)

	// Delegate to a validator
	_, err = app.AllianceKeeper.Delegate(ctx, delAddr, val, sdk.NewCoin(AllianceDenom, math.NewInt(1000_000)))
	require.NoError(t, err)

	assets := app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)

	totalBonded, err := app.StakingKeeper.TotalBondedTokens(ctx)
	require.NoError(t, err)
	// Check total bonded amount
	require.Equal(t, math.NewInt(3_000_000), totalBonded)

	// Check that balance dropped
	coin := app.BankKeeper.GetBalance(ctx, delAddr, AllianceDenom)
	require.Equal(t, sdk.NewCoin(AllianceDenom, math.NewInt(1000_000)), coin)

	// Check that staked balance increased
	d, _ := app.StakingKeeper.GetDelegation(ctx, moduleAddr, valAddr)
	require.Equal(t, stakingtypes.Delegation{
		DelegatorAddress: moduleAddr.String(),
		ValidatorAddress: valAddr.String(),
		Shares:           math.LegacyNewDec(2),
	}, d)

	// Immediately undelegate from the validator
	_, err = app.AllianceKeeper.Undelegate(ctx, delAddr, val, sdk.NewCoin(AllianceDenom, math.NewInt(250_000)))
	require.NoError(t, err)

	_, err = app.AllianceKeeper.Undelegate(ctx, delAddr, val, sdk.NewCoin(AllianceDenom, math.NewInt(250_000)))
	require.NoError(t, err)

	// Query unbondings directly from the entry point
	queryServer := keeper.NewQueryServerImpl(app.AllianceKeeper)
	res, err := queryServer.AllianceUnbondingsByDelegator(ctx, &types.QueryAllianceUnbondingsByDelegatorRequest{
		DelegatorAddr: delAddr.String(),
	})
	require.NoError(t, err)
	require.Equal(t, &types.QueryAllianceUnbondingsByDelegatorResponse{
		Unbondings: []types.UnbondingDelegation{
			{
				CompletionTime:   ctx.BlockTime().Add(unbondingTime),
				ValidatorAddress: valAddr.String(),
				Amount:           math.NewInt(250_000),
				Denom:            AllianceDenom,
			},
			{
				CompletionTime:   ctx.BlockTime().Add(unbondingTime),
				ValidatorAddress: valAddr.String(),
				Amount:           math.NewInt(250_000),
				Denom:            AllianceDenom,
			},
		},
	}, res)

	// Check if undelegations were stored correctly
	iter := app.AllianceKeeper.IterateUndelegationsByCompletionTime(ctx, ctx.BlockTime().Add(unbondingTime).Add(time.Second))
	require.True(t, iter.Valid())
	var queuedUndelegations types.QueuedUndelegation
	b := iter.Value()
	app.AppCodec().MustUnmarshal(b, &queuedUndelegations)
	require.Equal(t, types.QueuedUndelegation{Entries: []*types.Undelegation{
		{
			DelegatorAddress: delAddr.String(),
			ValidatorAddress: getOperator(val).String(),
			Balance:          sdk.NewCoin(AllianceDenom, math.NewInt(250_000)),
		},
		{
			DelegatorAddress: delAddr.String(),
			ValidatorAddress: getOperator(val).String(),
			Balance:          sdk.NewCoin(AllianceDenom, math.NewInt(250_000)),
		},
	}}, queuedUndelegations)

	assets = app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)

	totalBonded, err = app.StakingKeeper.TotalBondedTokens(ctx)
	require.NoError(t, err)
	// Check total bonded amount
	require.Equal(t, math.NewInt(3_000_000), totalBonded)

	// Check that staked balance stays the same
	d, _ = app.StakingKeeper.GetDelegation(ctx, moduleAddr, valAddr)
	require.Equal(t, stakingtypes.Delegation{
		DelegatorAddress: moduleAddr.String(),
		ValidatorAddress: valAddr.String(),
		Shares:           math.LegacyNewDec(2),
	}, d)

	// Immediately try to complete undelegation
	err = app.AllianceKeeper.CompleteUnbondings(ctx)
	require.NoError(t, err)

	// Check that balance stayed the same
	coin = app.BankKeeper.GetBalance(ctx, delAddr, AllianceDenom)
	require.Equal(t, sdk.NewCoin(AllianceDenom, math.NewInt(1000_000)), coin)

	// Advance time to after unbonding period
	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(unbondingTime).Add(time.Minute))

	err = app.AllianceKeeper.CompleteUnbondings(ctx)
	require.NoError(t, err)

	// Check that balance increased
	coin = app.BankKeeper.GetBalance(ctx, delAddr, AllianceDenom)
	require.Equal(t, sdk.NewCoin(AllianceDenom, math.NewInt(1500_000)), coin)

	// Completing again should not process anymore undelegations
	err = app.AllianceKeeper.CompleteUnbondings(ctx)
	require.NoError(t, err)
}

func TestUndelegationWithoutDelegation(t *testing.T) {
	app, ctx := createTestContext(t)
	ctx = ctx.WithBlockTime(time.Now())
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			types.NewAllianceAsset(AllianceDenom, math.LegacyNewDec(2), math.LegacyZeroDec(), math.LegacyNewDec(5), math.LegacyNewDec(0), ctx.BlockTime()),
			types.NewAllianceAsset(AllianceDenomTwo, math.LegacyNewDec(10), math.LegacyNewDec(2), math.LegacyNewDec(12), math.LegacyNewDec(0), ctx.BlockTime()),
		},
	})
	delegations, err := app.StakingKeeper.GetAllDelegations(ctx)
	require.NoError(t, err)
	require.Len(t, delegations, 1)

	// All the addresses needed
	delAddr, err := sdk.AccAddressFromBech32(delegations[0].DelegatorAddress)
	require.NoError(t, err)
	valAddr, err := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	require.NoError(t, err)
	val, _ := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr)

	// Mint alliance tokens
	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(AllianceDenom, math.NewInt(2000_000))))
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr, sdk.NewCoins(sdk.NewCoin(AllianceDenom, math.NewInt(2000_000))))
	require.NoError(t, err)

	// Undelegating without a delegation will fail
	_, err = app.AllianceKeeper.Undelegate(ctx, delAddr, val, sdk.NewCoin(AllianceDenom, math.NewInt(1000_000)))
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
			types.NewAllianceAsset(AllianceDenom, math.LegacyNewDec(2), math.LegacyZeroDec(), math.LegacyNewDec(5), math.LegacyNewDec(0), ctx.BlockTime()),
			types.NewAllianceAsset(AllianceDenomTwo, math.LegacyNewDec(10), math.LegacyNewDec(2), math.LegacyNewDec(12), math.LegacyMustNewDecFromStr("0.5"), ctx.BlockTime()),
		},
	})
	queryServer := keeper.NewQueryServerImpl(app.AllianceKeeper)

	// Set tax and rewards to be zero for easier calculation
	distParams, err := app.DistrKeeper.Params.Get(ctx)
	require.NoError(t, err)
	distParams.CommunityTax = math.LegacyZeroDec()
	err = app.DistrKeeper.Params.Set(ctx, distParams)
	require.NoError(t, err)

	// Accounts

	// rewardsPoolAddr := app.AccountKeeper.GetModuleAddress(types.RewardsPoolName)
	addrs := test_helpers.AddTestAddrsIncremental(app, ctx, 4, sdk.NewCoins(
		sdk.NewCoin(AllianceDenom, math.NewInt(1000_000_000)),
		sdk.NewCoin(AllianceDenomTwo, math.NewInt(2000_000_000)),
	))
	pks := test_helpers.CreateTestPubKeys(2)

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
	test_helpers.RegisterNewValidator(t, app, ctx, _val2)
	val2, err := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr2)
	require.NoError(t, err)

	user1 := addrs[2]
	user2 := addrs[3]

	// Delegate token with non-zero take_rate
	_, err = app.AllianceKeeper.Delegate(ctx, user1, val1, sdk.NewCoin(AllianceDenomTwo, math.NewInt(1000_000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user2, val2, sdk.NewCoin(AllianceDenomTwo, math.NewInt(1000_000_000)))
	require.NoError(t, err)

	assets := app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)
	// Check total bonded amount
	totalBonded, err := app.StakingKeeper.TotalBondedTokens(ctx)
	require.NoError(t, err)
	require.Equal(t, math.NewInt(11_000_000), totalBonded)

	ctx = ctx.WithBlockTime(startTime.Add(time.Minute * 6)).WithBlockHeight(2)
	coins, err := app.AllianceKeeper.DeductAssetsHook(ctx, assets)
	require.NoError(t, err)
	require.False(t, coins.IsZero())

	res, err := queryServer.AllianceDelegation(ctx, &types.QueryAllianceDelegationRequest{
		DelegatorAddr: user1.String(),
		ValidatorAddr: val1.GetOperator(),
		Denom:         AllianceDenomTwo,
		Pagination:    nil,
	})
	require.NoError(t, err)
	del := res.GetDelegation()
	require.True(t, del.GetBalance().Amount.LT(math.NewInt(1000_000_000)), "%s should be less than %s", del.GetBalance().Amount, math.NewInt(1000_000_000))
	// Undelegate token with initial amount should fail
	_, err = app.AllianceKeeper.Undelegate(ctx, user1, val1, sdk.NewCoin(AllianceDenomTwo, math.NewInt(1000_000_000)))
	require.Error(t, err)

	// Undelegate token with current amount should pass
	_, err = app.AllianceKeeper.Undelegate(ctx, user1, val1, sdk.NewCoin(AllianceDenomTwo, del.Balance.Amount))
	require.NoError(t, err)

	// User should have everything withdrawn
	_, found := app.AllianceKeeper.GetDelegation(ctx, user1, valAddr1, AllianceDenomTwo)
	require.False(t, found)

	// Delegate again
	_, err = app.AllianceKeeper.Delegate(ctx, user1, val1, sdk.NewCoin(AllianceDenomTwo, math.NewInt(500_000_000)))
	require.NoError(t, err)

	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(time.Minute * 1)).WithBlockHeight(2)

	_, err = app.AllianceKeeper.Delegate(ctx, user1, val1, sdk.NewCoin(AllianceDenomTwo, math.NewInt(400_000_000)))
	require.NoError(t, err)

	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(time.Minute * 5)).WithBlockHeight(3)
	coins, err = app.AllianceKeeper.DeductAssetsHook(ctx, assets)
	require.NoError(t, err)
	require.False(t, coins.IsZero())

	res, err = queryServer.AllianceDelegation(ctx, &types.QueryAllianceDelegationRequest{
		DelegatorAddr: user1.String(),
		ValidatorAddr: val1.GetOperator(),
		Denom:         AllianceDenomTwo,
		Pagination:    nil,
	})
	require.NoError(t, err)
	del = res.GetDelegation()
	require.True(t, del.GetBalance().Amount.LT(math.NewInt(900_000_000)), "%s should be less than %s", del.GetBalance().Amount, math.NewInt(1000_000_000))

	// Undelegate token with current amount should pass
	_, err = app.AllianceKeeper.Undelegate(ctx, user1, val1, sdk.NewCoin(AllianceDenomTwo, del.Balance.Amount))
	require.NoError(t, err)

	// User should have everything withdrawn
	_, found = app.AllianceKeeper.GetDelegation(ctx, user1, valAddr1, AllianceDenomTwo)
	require.False(t, found)

	res, err = queryServer.AllianceDelegation(ctx, &types.QueryAllianceDelegationRequest{
		DelegatorAddr: user1.String(),
		ValidatorAddr: val1.GetOperator(),
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
			types.NewAllianceAsset(AllianceDenom, math.LegacyNewDec(2), math.LegacyZeroDec(), math.LegacyNewDec(5), math.LegacyNewDec(0), ctx.BlockTime()),
			types.NewAllianceAsset(AllianceDenomTwo, math.LegacyNewDec(10), math.LegacyNewDec(2), math.LegacyNewDec(12), math.LegacyMustNewDecFromStr("0.5"), ctx.BlockTime()),
		},
	})

	// Set tax and rewards to be zero for easier calculation
	distParams, err := app.DistrKeeper.Params.Get(ctx)
	require.NoError(t, err)
	distParams.CommunityTax = math.LegacyZeroDec()
	err = app.DistrKeeper.Params.Set(ctx, distParams)
	require.NoError(t, err)

	// Accounts

	// rewardsPoolAddr := app.AccountKeeper.GetModuleAddress(types.RewardsPoolName)
	bondDenom, err := app.StakingKeeper.BondDenom(ctx)
	require.NoError(t, err)
	addrs := test_helpers.AddTestAddrsIncremental(app, ctx, 4, sdk.NewCoins(
		sdk.NewCoin(bondDenom, math.NewInt(1000_000_000)),
		sdk.NewCoin(AllianceDenom, math.NewInt(1000_000_000)),
		sdk.NewCoin(AllianceDenomTwo, math.NewInt(2000_000_000)),
	))
	pks := test_helpers.CreateTestPubKeys(2)

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
	test_helpers.RegisterNewValidator(t, app, ctx, _val2)
	val2, err := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr2)
	require.NoError(t, err)

	user1 := addrs[2]
	user2 := addrs[3]

	// Stake some alliance tokens
	_, err = app.AllianceKeeper.Delegate(ctx, user1, val1, sdk.NewCoin(AllianceDenom, math.NewInt(2000_000)))
	require.NoError(t, err)
	_, err = app.AllianceKeeper.Delegate(ctx, user2, val2, sdk.NewCoin(AllianceDenomTwo, math.NewInt(2000_000)))
	require.NoError(t, err)
	assets := app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)

	// Check total bonded tokens
	totalBonded, err := app.StakingKeeper.TotalBondedTokens(ctx)
	require.NoError(t, err)
	require.Equal(t, math.NewInt(13_000_000), totalBonded)

	// Stake some native tokens
	_, err = app.StakingKeeper.Delegate(ctx, user2, math.NewInt(2000_000), stakingtypes.Unbonded, *val2.Validator, true)
	require.NoError(t, err)
	// Check total bonded tokens
	totalBonded, err = app.StakingKeeper.TotalBondedTokens(ctx)
	require.NoError(t, err)
	require.Equal(t, math.NewInt(15_000_000), totalBonded)

	assets = app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)
	// Check total bonded tokens
	totalBonded, err = app.StakingKeeper.TotalBondedTokens(ctx)
	require.NoError(t, err)
	require.Equal(t, math.NewInt(39_000_000), totalBonded)

	// Redelegate some native tokens
	_, err = app.StakingKeeper.BeginRedelegation(ctx, user2, valAddr2, valAddr1, math.LegacyNewDec(1))
	require.NoError(t, err)

	assets = app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)
	// Check total bonded tokens
	totalBonded, err = app.StakingKeeper.TotalBondedTokens(ctx)
	require.NoError(t, err)
	require.Equal(t, math.NewInt(39_000_000), totalBonded)

	// Undelegate some native tokens
	shares, _ := app.StakingKeeper.ValidateUnbondAmount(ctx, user2, valAddr2, math.NewInt(1000_000))
	_, _, err = app.StakingKeeper.Undelegate(ctx, user2, valAddr2, shares)
	require.NoError(t, err)

	assets = app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)
	// Check total bonded tokens
	totalBonded, err = app.StakingKeeper.TotalBondedTokens(ctx)
	require.NoError(t, err)
	require.Equal(t, math.NewInt(26_000_000), totalBonded)
}

func TestUndelegatingLargeNumbers(t *testing.T) {
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
	unbondingTime, err := app.StakingKeeper.UnbondingTime(ctx)
	require.NoError(t, err)

	// Get the native delegations to have a validator address where to delegate
	delegations, err := app.StakingKeeper.GetAllDelegations(ctx)
	require.NoError(t, err)
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
	unbondings, err := app.AllianceKeeper.GetUnbondingsByDenomAndDelegator(ctx, AllianceDenom, delAddr)
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
