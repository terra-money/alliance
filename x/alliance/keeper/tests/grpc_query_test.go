package tests_test

import (
	"testing"
	"time"

	teststaking "github.com/cosmos/cosmos-sdk/x/staking/testutil"

	"cosmossdk.io/math"
	abcitypes "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	test_helpers "github.com/terra-money/alliance/app"
	"github.com/terra-money/alliance/x/alliance/keeper"
	"github.com/terra-money/alliance/x/alliance/types"
)

var ULunaAlliance = "uluna"

func TestQueryAlliances(t *testing.T) {
	// GIVEN: THE BLOCKCHAIN WITH ALLIANCES ON GENESIS
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			{
				Denom:             AllianceDenom,
				RewardWeight:      math.LegacyNewDec(2),
				RewardWeightRange: types.RewardWeightRange{Min: math.LegacyNewDec(0), Max: math.LegacyNewDec(5)},
				TakeRate:          math.LegacyNewDec(0),
				TotalTokens:       math.ZeroInt(),
			},
			{
				Denom:             AllianceDenomTwo,
				RewardWeight:      math.LegacyNewDec(10),
				RewardWeightRange: types.RewardWeightRange{Min: math.LegacyNewDec(2), Max: math.LegacyNewDec(12)},
				TakeRate:          math.LegacyMustNewDecFromStr("0.14159265359"),
				TotalTokens:       math.ZeroInt(),
			},
		},
	})
	queryServer := keeper.NewQueryServerImpl(app.AllianceKeeper)

	// WHEN: QUERYING THE ALLIANCES LIST
	alliances, err := queryServer.Alliances(ctx, &types.QueryAlliancesRequest{})

	// THEN: VALIDATE THAT BOTH ALLIANCES HAVE THE CORRECT MODEL WHEN QUERYING
	require.Nil(t, err)
	require.Equal(t, &types.QueryAlliancesResponse{
		Alliances: []types.AllianceAsset{
			{
				Denom:                "alliance",
				RewardWeight:         math.LegacyNewDec(2),
				RewardWeightRange:    types.RewardWeightRange{Min: math.LegacyNewDec(0), Max: math.LegacyNewDec(5)},
				TakeRate:             math.LegacyNewDec(0),
				TotalTokens:          math.ZeroInt(),
				TotalValidatorShares: math.LegacyNewDec(0),
				RewardChangeRate:     math.LegacyNewDec(0),
				RewardChangeInterval: 0,
			},
			{
				Denom:                "alliance2",
				RewardWeight:         math.LegacyNewDec(10),
				RewardWeightRange:    types.RewardWeightRange{Min: math.LegacyNewDec(2), Max: math.LegacyNewDec(12)},
				TakeRate:             math.LegacyMustNewDecFromStr("0.14159265359"),
				TotalTokens:          math.ZeroInt(),
				TotalValidatorShares: math.LegacyNewDec(0),
				RewardChangeRate:     math.LegacyNewDec(0),
				RewardChangeInterval: 0,
			},
		},
		Pagination: &query.PageResponse{
			NextKey: nil,
			Total:   2,
		},
	}, alliances)
}

func TestQueryAnUniqueAlliance(t *testing.T) {
	// GIVEN: THE BLOCKCHAIN WITH ALLIANCES ON GENESIS
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			{
				Denom:                AllianceDenom,
				RewardWeight:         math.LegacyNewDec(2),
				RewardWeightRange:    types.RewardWeightRange{Min: math.LegacyNewDec(0), Max: math.LegacyNewDec(5)},
				TakeRate:             math.LegacyNewDec(0),
				TotalTokens:          math.ZeroInt(),
				RewardChangeRate:     math.LegacyNewDec(0),
				RewardChangeInterval: 0,
			},
			{
				Denom:                AllianceDenomTwo,
				RewardWeight:         math.LegacyNewDec(10),
				RewardWeightRange:    types.RewardWeightRange{Min: math.LegacyNewDec(2), Max: math.LegacyNewDec(12)},
				TakeRate:             math.LegacyMustNewDecFromStr("0.14159265359"),
				TotalTokens:          math.ZeroInt(),
				RewardChangeRate:     math.LegacyNewDec(0),
				RewardChangeInterval: 0,
			},
		},
	})
	queryServer := keeper.NewQueryServerImpl(app.AllianceKeeper)

	// WHEN: QUERYING THE ALLIANCES LIST
	alliances, err := queryServer.Alliance(ctx, &types.QueryAllianceRequest{
		Denom: "alliance2",
	})

	// THEN: VALIDATE THAT BOTH ALLIANCES HAVE THE CORRECT MODEL WHEN QUERYING
	require.Nil(t, err)
	require.Equal(t, &types.QueryAllianceResponse{
		Alliance: &types.AllianceAsset{
			Denom:                "alliance2",
			RewardWeight:         math.LegacyNewDec(10),
			RewardWeightRange:    types.RewardWeightRange{Min: math.LegacyNewDec(2), Max: math.LegacyNewDec(12)},
			TakeRate:             math.LegacyMustNewDecFromStr("0.14159265359"),
			TotalTokens:          math.ZeroInt(),
			TotalValidatorShares: math.LegacyNewDec(0),
			RewardChangeRate:     math.LegacyNewDec(0),
			RewardChangeInterval: 0,
		},
	}, alliances)
}

func TestQueryAnUniqueIBCAlliance(t *testing.T) {
	// GIVEN: THE BLOCKCHAIN WITH ALLIANCES ON GENESIS
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			{
				Denom:                "ibc/" + AllianceDenomTwo,
				RewardWeight:         math.LegacyNewDec(10),
				RewardWeightRange:    types.RewardWeightRange{Min: math.LegacyNewDec(2), Max: math.LegacyNewDec(12)},
				TakeRate:             math.LegacyMustNewDecFromStr("0.14159265359"),
				TotalTokens:          math.ZeroInt(),
				RewardChangeRate:     math.LegacyNewDec(0),
				RewardChangeInterval: 0,
			},
		},
	})
	queryServer := keeper.NewQueryServerImpl(app.AllianceKeeper)

	// WHEN: QUERYING THE ALLIANCES LIST
	alliances, err := queryServer.Alliance(ctx, &types.QueryAllianceRequest{
		Denom: "ibc%2Falliance2",
	})
	alliancesIbcEndpoint, err2 := queryServer.IBCAlliance(ctx, &types.QueryIBCAllianceRequest{ //nolint:staticcheck // SA1019: types.types.QueryIBCAllianceRequest is deprecated
		Hash: "alliance2",
	})

	// THEN: VALIDATE THAT BOTH ALLIANCES HAVE THE CORRECT MODEL WHEN QUERYING
	require.Nil(t, err)
	require.Equal(t, &types.QueryAllianceResponse{
		Alliance: &types.AllianceAsset{
			Denom:                "ibc/alliance2",
			RewardWeight:         math.LegacyNewDec(10),
			RewardWeightRange:    types.RewardWeightRange{Min: math.LegacyNewDec(2), Max: math.LegacyNewDec(12)},
			TakeRate:             math.LegacyMustNewDecFromStr("0.14159265359"),
			TotalTokens:          math.ZeroInt(),
			TotalValidatorShares: math.LegacyNewDec(0),
			RewardChangeRate:     math.LegacyNewDec(0),
			RewardChangeInterval: 0,
		},
	}, alliances)

	require.Nil(t, err2)
	require.Equal(t, &types.QueryAllianceResponse{
		Alliance: &types.AllianceAsset{
			Denom:                "ibc/alliance2",
			RewardWeight:         math.LegacyNewDec(10),
			RewardWeightRange:    types.RewardWeightRange{Min: math.LegacyNewDec(2), Max: math.LegacyNewDec(12)},
			TakeRate:             math.LegacyMustNewDecFromStr("0.14159265359"),
			TotalTokens:          math.ZeroInt(),
			TotalValidatorShares: math.LegacyNewDec(0),
			RewardChangeRate:     math.LegacyNewDec(0),
			RewardChangeInterval: 0,
		},
	}, alliancesIbcEndpoint)
}

func TestQueryAllianceNotFound(t *testing.T) {
	// GIVEN: THE BLOCKCHAIN
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	queryServer := keeper.NewQueryServerImpl(app.AllianceKeeper)

	// WHEN: QUERYING THE ALLIANCE
	_, err := queryServer.Alliance(ctx, &types.QueryAllianceRequest{
		Denom: "alliance2",
	})

	// THEN: VALIDATE THE ERROR
	require.Equal(t, err.Error(), "alliance asset is not whitelisted")
}

func TestQueryAllAlliances(t *testing.T) {
	// GIVEN: THE BLOCKCHAIN
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	queryServer := keeper.NewQueryServerImpl(app.AllianceKeeper)

	// WHEN: QUERYING THE ALLIANCE
	res, err := queryServer.Alliances(ctx, &types.QueryAlliancesRequest{})

	// THEN: VALIDATE THE ERROR
	require.Nil(t, err)
	require.Equal(t, len(res.Alliances), 0)
	require.Equal(t, res.Pagination, &query.PageResponse{
		NextKey: nil,
		Total:   0,
	})
}

func TestQueryParams(t *testing.T) {
	// GIVEN: THE BLOCKCHAIN WITH AN ALLIANCE ON GENESIS
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			{
				Denom:                AllianceDenom,
				RewardWeight:         math.LegacyNewDec(2),
				TakeRate:             math.LegacyNewDec(0),
				TotalTokens:          math.ZeroInt(),
				RewardChangeRate:     math.LegacyNewDec(0),
				RewardChangeInterval: 0,
			},
		},
	})
	queryServer := keeper.NewQueryServerImpl(app.AllianceKeeper)

	// WHEN: QUERYING THE PARAMS...
	queyParams, err := queryServer.Params(ctx, &types.QueryParamsRequest{})

	// THEN: VALIDATE THAT NO ERRORS HAVE BEEN PRODUCED AND OUTPUT IS AS WE EXPECT
	require.Nil(t, err)

	require.Equal(t, queyParams.Params.RewardDelayTime, time.Hour*24*7)
	require.Equal(t, queyParams.Params.TakeRateClaimInterval, time.Minute*5)

	// there is no way to match the exact time when the module is being instantiated
	// but we know that this time should be older than actually the time when this
	// following two lines are executed
	require.NotNil(t, queyParams.Params.LastTakeRateClaimTime)
	require.LessOrEqual(t, queyParams.Params.LastTakeRateClaimTime, time.Now())
}

func TestClaimQueryReward(t *testing.T) {
	// GIVEN: THE BLOCKCHAIN WITH ACCOUNTS
	app, ctx := createTestContext(t)
	startTime := time.Now().UTC()
	ctx = ctx.WithBlockTime(startTime)
	ctx = ctx.WithBlockHeight(1)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.Params{
			RewardDelayTime:       time.Minute * 60,
			TakeRateClaimInterval: time.Minute * 5,
			LastTakeRateClaimTime: startTime,
		},
		Assets: []types.AllianceAsset{
			{
				Denom:                ULunaAlliance,
				RewardWeight:         math.LegacyNewDec(2),
				TakeRate:             math.LegacyMustNewDecFromStr("0.00005"),
				TotalTokens:          math.ZeroInt(),
				TotalValidatorShares: math.LegacyNewDec(0),
				RewardChangeRate:     math.LegacyNewDec(0),
				RewardChangeInterval: 0,
			},
		},
	})
	queryServer := keeper.NewQueryServerImpl(app.AllianceKeeper)
	feeCollectorAddr := app.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
	delegations, err := app.StakingKeeper.GetAllDelegations(ctx)
	require.NoError(t, err)
	valAddr, _ := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	val1, _ := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr)
	delAddr := test_helpers.AddTestAddrsIncremental(app, ctx, 1, sdk.NewCoins(sdk.NewCoin(ULunaAlliance, math.NewInt(1000_000_000))))[0]

	// WHEN: DELEGATING ...
	delRes, delErr := app.AllianceKeeper.Delegate(ctx, delAddr, val1, sdk.NewCoin(ULunaAlliance, math.NewInt(1000_000_000)))
	require.Nil(t, delErr)
	require.Equal(t, math.LegacyNewDec(1000000000), *delRes)
	assets := app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)

	// ...and advance block...
	timePassed := time.Minute*5 + time.Second
	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(timePassed))
	ctx = ctx.WithBlockHeight(2)
	_, err = app.AllianceKeeper.DeductAssetsHook(ctx, assets)
	require.NoError(t, err)

	app.BankKeeper.GetAllBalances(ctx, feeCollectorAddr)
	require.Equal(t, startTime.Add(time.Minute*5), app.AllianceKeeper.LastRewardClaimTime(ctx))
	app.AllianceKeeper.GetAssetByDenom(ctx, ULunaAlliance)

	// ... at the next begin block, tokens will be distributed from the fee pool...
	cons, _ := val1.GetConsAddr()
	err = app.DistrKeeper.AllocateTokens(ctx, 1, []abcitypes.VoteInfo{
		{
			Validator: abcitypes.Validator{
				Address: cons,
				Power:   1,
			},
		},
	})
	require.NoError(t, err)

	// THEN: Query the delegation rewards ...
	queryDelegation, queryErr := queryServer.AllianceDelegationRewards(ctx, &types.QueryAllianceDelegationRewardsRequest{
		DelegatorAddr: delAddr.String(),
		ValidatorAddr: valAddr.String(),
		Denom:         ULunaAlliance,
	})

	// ... validate that no error has been produced.
	require.Nil(t, queryErr)
	require.Equal(t, &types.QueryAllianceDelegationRewardsResponse{
		Rewards: []sdk.Coin{
			{
				Denom:  ULunaAlliance,
				Amount: math.NewInt(32666),
			},
		},
	}, queryDelegation)
}

func TestQueryAllianceDelegation(t *testing.T) {
	// GIVEN: THE BLOCKCHAIN WITH ALLIANCES ON GENESIS
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			{
				Denom:                AllianceDenom,
				RewardWeight:         math.LegacyNewDec(2),
				TakeRate:             math.LegacyNewDec(0),
				TotalTokens:          math.ZeroInt(),
				RewardChangeRate:     math.LegacyNewDec(0),
				RewardChangeInterval: 0,
			},
		},
	})
	queryServer := keeper.NewQueryServerImpl(app.AllianceKeeper)
	delegations, err := app.StakingKeeper.GetAllDelegations(ctx)
	require.NoError(t, err)
	delAddr, _ := sdk.AccAddressFromBech32(delegations[0].DelegatorAddress)
	valAddr, _ := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	val, _ := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr)
	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(AllianceDenom, math.NewInt(2000_000))))
	require.NoError(t, err)

	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr, sdk.NewCoins(sdk.NewCoin(AllianceDenom, math.NewInt(2000_000))))
	require.NoError(t, err)

	// WHEN: DELEGATING AND QUERYING ...
	delegationTxRes, txErr := app.AllianceKeeper.Delegate(ctx, delAddr, val, sdk.NewCoin(AllianceDenom, math.NewInt(1000_000)))
	queryDelegation, queryErr := queryServer.AllianceDelegation(ctx, &types.QueryAllianceDelegationRequest{
		DelegatorAddr: delAddr.String(),
		ValidatorAddr: val.OperatorAddress,
		Denom:         AllianceDenom,
	})

	// THEN: VALIDATE THAT NO ERRORS HAVE BEEN PRODUCED AND BOTH OUTPUTS ARE AS WE EXPECT
	require.Nil(t, txErr)
	require.Nil(t, queryErr)
	require.Equal(t, &types.QueryAllianceDelegationResponse{
		Delegation: types.DelegationResponse{
			Delegation: types.Delegation{
				DelegatorAddress:      delAddr.String(),
				ValidatorAddress:      val.OperatorAddress,
				Denom:                 AllianceDenom,
				Shares:                math.LegacyNewDec(1000_000),
				RewardHistory:         nil,
				LastRewardClaimHeight: uint64(ctx.BlockHeight()),
			},
			Balance: sdk.Coin{
				Denom:  AllianceDenom,
				Amount: math.NewInt(1000_000),
			},
		},
	}, queryDelegation)
	require.Equal(t, math.LegacyNewDec(1000000), *delegationTxRes)
}

func TestQueryAllianceDelegationNotFound(t *testing.T) {
	// GIVEN: THE BLOCKCHAIN
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	delegations, err := app.StakingKeeper.GetAllDelegations(ctx)
	require.NoError(t, err)
	delAddr, _ := sdk.AccAddressFromBech32(delegations[0].DelegatorAddress)
	valAddr, _ := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	val, _ := app.StakingKeeper.GetValidator(ctx, valAddr)
	queryServer := keeper.NewQueryServerImpl(app.AllianceKeeper)

	// WHEN: QUERYING ...
	_, err = queryServer.AllianceDelegation(ctx, &types.QueryAllianceDelegationRequest{
		DelegatorAddr: delAddr.String(),
		ValidatorAddr: val.OperatorAddress,
		Denom:         AllianceDenom,
	})

	// THEN: VALIDATE THAT NO ERRORS HAVE BEEN PRODUCED AND BOTH OUTPUTS ARE AS WE EXPECT
	require.Equal(t, err, status.Error(codes.NotFound, "AllianceAsset not found by denom alliance"))
}

func TestQueryAllianceValidatorNotFound(t *testing.T) {
	// GIVEN: THE BLOCKCHAIN
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	delegations, err := app.StakingKeeper.GetAllDelegations(ctx)
	require.NoError(t, err)
	delAddr, _ := sdk.AccAddressFromBech32(delegations[0].DelegatorAddress)
	queryServer := keeper.NewQueryServerImpl(app.AllianceKeeper)

	// WHEN: QUERYING ...
	_, err = queryServer.AllianceDelegation(ctx, &types.QueryAllianceDelegationRequest{
		DelegatorAddr: delAddr.String(),
		ValidatorAddr: "cosmosvaloper19lss6zgdh5vvcpjhfftdghrpsw7a4434elpwpu",
		Denom:         AllianceDenom,
	})

	// THEN: VALIDATE THAT NO ERRORS HAVE BEEN PRODUCED AND BOTH OUTPUTS ARE AS WE EXPECT
	require.Equal(t, err, status.Error(codes.NotFound, "Validator not found by address cosmosvaloper19lss6zgdh5vvcpjhfftdghrpsw7a4434elpwpu"))
}

func TestQueryAlliancesDelegationByValidator(t *testing.T) {
	// GIVEN: THE BLOCKCHAIN WITH ALLIANCES ON GENESIS
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			{
				Denom:                AllianceDenom,
				RewardWeight:         math.LegacyNewDec(2),
				TakeRate:             math.LegacyNewDec(0),
				TotalTokens:          math.ZeroInt(),
				RewardChangeRate:     math.LegacyNewDec(0),
				RewardChangeInterval: 0,
			},
		},
	})
	queryServer := keeper.NewQueryServerImpl(app.AllianceKeeper)
	delegations, err := app.StakingKeeper.GetAllDelegations(ctx)
	require.NoError(t, err)
	delAddr, _ := sdk.AccAddressFromBech32(delegations[0].DelegatorAddress)
	valAddr, _ := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	val, _ := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr)
	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(AllianceDenom, math.NewInt(2000_000))))
	require.NoError(t, err)

	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr, sdk.NewCoins(sdk.NewCoin(AllianceDenom, math.NewInt(2000_000))))
	require.NoError(t, err)

	// WHEN: DELEGATING AND QUERYING ...
	delegationTxRes, txErr := app.AllianceKeeper.Delegate(ctx, delAddr, val, sdk.NewCoin(AllianceDenom, math.NewInt(1000_000)))
	queryDelegation, queryErr := queryServer.AlliancesDelegationByValidator(ctx, &types.QueryAlliancesDelegationByValidatorRequest{
		DelegatorAddr: delAddr.String(),
		ValidatorAddr: val.OperatorAddress,
	})

	// THEN: VALIDATE THAT NO ERRORS HAVE BEEN PRODUCED AND BOTH OUTPUTS ARE AS WE EXPECT
	require.Nil(t, txErr)
	require.Nil(t, queryErr)
	require.Equal(t, &types.QueryAlliancesDelegationsResponse{
		Delegations: []types.DelegationResponse{
			{
				Delegation: types.Delegation{
					DelegatorAddress:      delAddr.String(),
					ValidatorAddress:      val.OperatorAddress,
					Denom:                 AllianceDenom,
					Shares:                math.LegacyNewDec(1000_000),
					RewardHistory:         nil,
					LastRewardClaimHeight: uint64(ctx.BlockHeight()),
				},
				Balance: sdk.Coin{
					Denom:  AllianceDenom,
					Amount: math.NewInt(1000_000),
				},
			},
		},
		Pagination: &query.PageResponse{
			NextKey: nil,
			Total:   1,
		},
	}, queryDelegation)
	require.Equal(t, math.LegacyNewDec(1000_000), *delegationTxRes)
}

func TestQueryAlliancesDelegationByValidatorNotFound(t *testing.T) {
	// GIVEN: THE BLOCKCHAIN
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	delegations, err := app.StakingKeeper.GetAllDelegations(ctx)
	require.NoError(t, err)
	delAddr, _ := sdk.AccAddressFromBech32(delegations[0].DelegatorAddress)
	queryServer := keeper.NewQueryServerImpl(app.AllianceKeeper)

	// WHEN: QUERYING ...
	_, err = queryServer.AlliancesDelegationByValidator(ctx, &types.QueryAlliancesDelegationByValidatorRequest{
		DelegatorAddr: delAddr.String(),
		ValidatorAddr: "cosmosvaloper19lss6zgdh5vvcpjhfftdghrpsw7a4434elpwpu",
	})

	// THEN: VALIDATE THAT NO ERRORS HAVE BEEN PRODUCED AND BOTH OUTPUTS ARE AS WE EXPECT
	require.Equal(t, err, status.Error(codes.NotFound, "Validator not found by address cosmosvaloper19lss6zgdh5vvcpjhfftdghrpsw7a4434elpwpu"))
}

func TestQueryAlliancesAlliancesDelegation(t *testing.T) {
	// GIVEN: THE BLOCKCHAIN WITH ALLIANCES ON GENESIS
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			{
				Denom:                AllianceDenom,
				RewardWeight:         math.LegacyNewDec(2),
				TakeRate:             math.LegacyNewDec(0),
				TotalTokens:          math.ZeroInt(),
				TotalValidatorShares: math.LegacyNewDec(0),
				RewardChangeRate:     math.LegacyNewDec(0),
				RewardChangeInterval: 0,
			},
			{
				Denom:                AllianceDenomTwo,
				RewardWeight:         math.LegacyNewDec(10),
				TakeRate:             math.LegacyMustNewDecFromStr("0.14159265359"),
				TotalTokens:          math.ZeroInt(),
				TotalValidatorShares: math.LegacyNewDec(0),
				RewardChangeRate:     math.LegacyNewDec(0),
				RewardChangeInterval: 0,
			},
		},
	})
	queryServer := keeper.NewQueryServerImpl(app.AllianceKeeper)
	delegations, err := app.StakingKeeper.GetAllDelegations(ctx)
	require.NoError(t, err)
	delAddr, _ := sdk.AccAddressFromBech32(delegations[0].DelegatorAddress)
	valAddr, _ := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	val, _ := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr)
	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(AllianceDenom, math.NewInt(2000_000))))
	require.NoError(t, err)

	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr, sdk.NewCoins(sdk.NewCoin(AllianceDenom, math.NewInt(2000_000))))
	require.NoError(t, err)

	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(AllianceDenomTwo, math.NewInt(2000_000))))
	require.NoError(t, err)

	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr, sdk.NewCoins(sdk.NewCoin(AllianceDenomTwo, math.NewInt(2000_000))))
	require.NoError(t, err)

	// WHEN: DELEGATING AND QUERYING ...
	delegationTxRes, txErr := app.AllianceKeeper.Delegate(ctx, delAddr, val, sdk.NewCoin(AllianceDenom, math.NewInt(1000_000)))
	delegation2TxRes, tx2Err := app.AllianceKeeper.Delegate(ctx, delAddr, val, sdk.NewCoin(AllianceDenomTwo, math.NewInt(1000_000)))
	queryDelegation, queryErr := queryServer.AlliancesDelegation(ctx, &types.QueryAlliancesDelegationsRequest{
		DelegatorAddr: delAddr.String(),
	})

	// THEN: VALIDATE THAT NO ERRORS HAVE BEEN PRODUCED AND BOTH OUTPUTS ARE AS WE EXPECT
	require.Nil(t, txErr)
	require.Nil(t, tx2Err)
	require.Nil(t, queryErr)
	require.Equal(t, &types.QueryAlliancesDelegationsResponse{
		Delegations: []types.DelegationResponse{
			{
				Delegation: types.Delegation{
					DelegatorAddress:      delAddr.String(),
					ValidatorAddress:      val.OperatorAddress,
					Denom:                 AllianceDenom,
					Shares:                math.LegacyNewDec(1000_000),
					RewardHistory:         nil,
					LastRewardClaimHeight: uint64(ctx.BlockHeight()),
				},
				Balance: sdk.Coin{
					Denom:  AllianceDenom,
					Amount: math.NewInt(1000_000),
				},
			},
			{
				Delegation: types.Delegation{
					DelegatorAddress:      delAddr.String(),
					ValidatorAddress:      val.OperatorAddress,
					Denom:                 AllianceDenomTwo,
					Shares:                math.LegacyNewDec(1000_000),
					RewardHistory:         nil,
					LastRewardClaimHeight: uint64(ctx.BlockHeight()),
				},
				Balance: sdk.Coin{
					Denom:  AllianceDenomTwo,
					Amount: math.NewInt(1000_000),
				},
			},
		},
		Pagination: &query.PageResponse{
			NextKey: nil,
			Total:   2,
		},
	}, queryDelegation)
	require.Equal(t, math.LegacyNewDec(1000_000), *delegationTxRes)
	require.Equal(t, math.LegacyNewDec(1000_000), *delegation2TxRes)
}

func TestQueryAllDelegations(t *testing.T) {
	// GIVEN: THE BLOCKCHAIN WITH ALLIANCES ON GENESIS
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			{
				Denom:                AllianceDenom,
				RewardWeight:         math.LegacyNewDec(2),
				TakeRate:             math.LegacyNewDec(0),
				TotalTokens:          math.ZeroInt(),
				TotalValidatorShares: math.LegacyNewDec(0),
				RewardChangeRate:     math.LegacyNewDec(0),
				RewardChangeInterval: 0,
			},
			{
				Denom:                AllianceDenomTwo,
				RewardWeight:         math.LegacyNewDec(10),
				TakeRate:             math.LegacyMustNewDecFromStr("0.14159265359"),
				TotalTokens:          math.ZeroInt(),
				TotalValidatorShares: math.LegacyNewDec(0),
				RewardChangeRate:     math.LegacyNewDec(0),
				RewardChangeInterval: 0,
			},
		},
	})
	queryServer := keeper.NewQueryServerImpl(app.AllianceKeeper)
	delegations, err := app.StakingKeeper.GetAllDelegations(ctx)
	require.NoError(t, err)
	delAddr, _ := sdk.AccAddressFromBech32(delegations[0].DelegatorAddress)
	valAddr, _ := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	val, _ := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr)
	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(AllianceDenom, math.NewInt(2000_000))))
	require.NoError(t, err)

	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr, sdk.NewCoins(sdk.NewCoin(AllianceDenom, math.NewInt(2000_000))))
	require.NoError(t, err)

	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(AllianceDenomTwo, math.NewInt(2000_000))))
	require.NoError(t, err)

	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr, sdk.NewCoins(sdk.NewCoin(AllianceDenomTwo, math.NewInt(2000_000))))
	require.NoError(t, err)

	// WHEN: DELEGATING AND QUERYING ...
	_, txErr := app.AllianceKeeper.Delegate(ctx, delAddr, val, sdk.NewCoin(AllianceDenom, math.NewInt(1000_000)))
	require.NoError(t, txErr)
	_, tx2Err := app.AllianceKeeper.Delegate(ctx, delAddr, val, sdk.NewCoin(AllianceDenomTwo, math.NewInt(1000_000)))
	require.NoError(t, tx2Err)
	queryDelegations, queryErr := queryServer.AllAlliancesDelegations(ctx, &types.QueryAllAlliancesDelegationsRequest{
		Pagination: &query.PageRequest{
			Key:        nil,
			Offset:     0,
			Limit:      1,
			CountTotal: false,
			Reverse:    false,
		},
	})
	require.NoError(t, queryErr)
	require.Equal(t, 1, len(queryDelegations.Delegations))

	require.Equal(t, types.DelegationResponse{
		Delegation: types.Delegation{
			DelegatorAddress:      delAddr.String(),
			ValidatorAddress:      val.OperatorAddress,
			Denom:                 AllianceDenom,
			Shares:                math.LegacyNewDec(1000_000),
			RewardHistory:         nil,
			LastRewardClaimHeight: uint64(ctx.BlockHeight()),
		},
		Balance: sdk.Coin{
			Denom:  AllianceDenom,
			Amount: math.NewInt(1000_000),
		},
	}, queryDelegations.Delegations[0])

	queryDelegations, queryErr = queryServer.AllAlliancesDelegations(ctx, &types.QueryAllAlliancesDelegationsRequest{
		Pagination: &query.PageRequest{
			Key:        queryDelegations.Pagination.NextKey,
			Offset:     0,
			Limit:      1,
			CountTotal: false,
			Reverse:    false,
		},
	})
	require.NoError(t, queryErr)
	require.Equal(t, types.DelegationResponse{
		Delegation: types.Delegation{
			DelegatorAddress:      delAddr.String(),
			ValidatorAddress:      val.OperatorAddress,
			Denom:                 AllianceDenomTwo,
			Shares:                math.LegacyNewDec(1000_000),
			RewardHistory:         nil,
			LastRewardClaimHeight: uint64(ctx.BlockHeight()),
		},
		Balance: sdk.Coin{
			Denom:  AllianceDenomTwo,
			Amount: math.NewInt(1000_000),
		},
	}, queryDelegations.Delegations[0])
}

func TestQueryValidator(t *testing.T) {
	// GIVEN: THE BLOCKCHAIN WITH ALLIANCES ON GENESIS
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			{
				Denom:                AllianceDenom,
				RewardWeight:         math.LegacyNewDec(2),
				TakeRate:             math.LegacyNewDec(0),
				TotalTokens:          math.ZeroInt(),
				TotalValidatorShares: math.LegacyNewDec(0),
				RewardChangeRate:     math.LegacyNewDec(0),
				RewardChangeInterval: 0,
			},
			{
				Denom:                AllianceDenomTwo,
				RewardWeight:         math.LegacyNewDec(10),
				TakeRate:             math.LegacyMustNewDecFromStr("0.14159265359"),
				TotalTokens:          math.ZeroInt(),
				TotalValidatorShares: math.LegacyNewDec(0),
				RewardChangeRate:     math.LegacyNewDec(0),
				RewardChangeInterval: 0,
			},
		},
	})
	queryServer := keeper.NewQueryServerImpl(app.AllianceKeeper)
	delegations, err := app.StakingKeeper.GetAllDelegations(ctx)
	require.NoError(t, err)
	delAddr, _ := sdk.AccAddressFromBech32(delegations[0].DelegatorAddress)
	valAddr, _ := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	val, _ := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr)
	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(AllianceDenom, math.NewInt(2000_000))))
	require.NoError(t, err)

	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr, sdk.NewCoins(sdk.NewCoin(AllianceDenom, math.NewInt(2000_000))))
	require.NoError(t, err)

	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(AllianceDenomTwo, math.NewInt(2000_000))))
	require.NoError(t, err)

	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr, sdk.NewCoins(sdk.NewCoin(AllianceDenomTwo, math.NewInt(2000_000))))
	require.NoError(t, err)

	// WHEN: DELEGATING AND QUERYING ...
	_, txErr := app.AllianceKeeper.Delegate(ctx, delAddr, val, sdk.NewCoin(AllianceDenom, math.NewInt(1000_000)))
	require.NoError(t, txErr)
	_, tx2Err := app.AllianceKeeper.Delegate(ctx, delAddr, val, sdk.NewCoin(AllianceDenomTwo, math.NewInt(1000_000)))
	require.NoError(t, tx2Err)

	queryVal, queryErr := queryServer.AllianceValidator(ctx, &types.QueryAllianceValidatorRequest{
		ValidatorAddr: getOperator(val).String(),
	})

	require.NoError(t, queryErr)
	require.Equal(t, &types.QueryAllianceValidatorResponse{
		ValidatorAddr: getOperator(val).String(),
		TotalDelegationShares: sdk.NewDecCoins(
			sdk.NewDecCoinFromDec(AllianceDenom, math.LegacyNewDec(1000000)),
			sdk.NewDecCoinFromDec(AllianceDenomTwo, math.LegacyNewDec(1000000)),
		),
		ValidatorShares: sdk.NewDecCoins(
			sdk.NewDecCoinFromDec(AllianceDenom, math.LegacyNewDec(1000000)),
			sdk.NewDecCoinFromDec(AllianceDenomTwo, math.LegacyNewDec(1000000)),
		),
		TotalStaked: sdk.NewDecCoins(
			sdk.NewDecCoinFromDec(AllianceDenom, math.LegacyNewDec(1000_000)),
			sdk.NewDecCoinFromDec(AllianceDenomTwo, math.LegacyNewDec(1000_000)),
		),
	}, queryVal)
}

func TestQueryValidators(t *testing.T) {
	// GIVEN: THE BLOCKCHAIN WITH ALLIANCES ON GENESIS
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			{
				Denom:                AllianceDenom,
				RewardWeight:         math.LegacyNewDec(2),
				TakeRate:             math.LegacyNewDec(0),
				TotalTokens:          math.ZeroInt(),
				TotalValidatorShares: math.LegacyNewDec(0),
				RewardChangeRate:     math.LegacyNewDec(0),
				RewardChangeInterval: 0,
			},
			{
				Denom:                AllianceDenomTwo,
				RewardWeight:         math.LegacyNewDec(10),
				TakeRate:             math.LegacyMustNewDecFromStr("0.14159265359"),
				TotalTokens:          math.ZeroInt(),
				TotalValidatorShares: math.LegacyNewDec(0),
				RewardChangeRate:     math.LegacyNewDec(0),
				RewardChangeInterval: 0,
			},
		},
	})
	queryServer := keeper.NewQueryServerImpl(app.AllianceKeeper)
	delegations, err := app.StakingKeeper.GetAllDelegations(ctx)
	require.NoError(t, err)
	delAddr, _ := sdk.AccAddressFromBech32(delegations[0].DelegatorAddress)
	valAddr, _ := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	val, _ := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr)

	addrs := test_helpers.AddTestAddrsIncremental(app, ctx, 3, sdk.NewCoins(
		sdk.NewCoin(AllianceDenom, math.NewInt(1000_000)),
		sdk.NewCoin(AllianceDenomTwo, math.NewInt(1000_000)),
	))
	valAddr2 := sdk.ValAddress(addrs[0])
	_val2 := teststaking.NewValidator(t, valAddr2, test_helpers.CreateTestPubKeys(1)[0])
	test_helpers.RegisterNewValidator(t, app, ctx, _val2)
	val2, err := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr2)
	require.NoError(t, err)

	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(AllianceDenom, math.NewInt(2000_000))))
	require.NoError(t, err)

	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr, sdk.NewCoins(sdk.NewCoin(AllianceDenom, math.NewInt(2000_000))))
	require.NoError(t, err)

	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(AllianceDenomTwo, math.NewInt(2000_000))))
	require.NoError(t, err)

	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr, sdk.NewCoins(sdk.NewCoin(AllianceDenomTwo, math.NewInt(2000_000))))
	require.NoError(t, err)

	// WHEN: DELEGATING AND QUERYING ...
	_, txErr := app.AllianceKeeper.Delegate(ctx, delAddr, val, sdk.NewCoin(AllianceDenom, math.NewInt(1000_000)))
	require.NoError(t, txErr)
	_, tx2Err := app.AllianceKeeper.Delegate(ctx, delAddr, val2, sdk.NewCoin(AllianceDenomTwo, math.NewInt(1000_000)))
	require.NoError(t, tx2Err)

	queryVal, queryErr := queryServer.AllAllianceValidators(ctx, &types.QueryAllAllianceValidatorsRequest{
		Pagination: &query.PageRequest{
			Key:        nil,
			Offset:     0,
			Limit:      1,
			CountTotal: false,
			Reverse:    false,
		},
	})

	require.NoError(t, queryErr)
	// Order in which validators are returned is not deterministic since we randomly generate validator addresses
	if queryVal.Validators[0].ValidatorAddr == getOperator(val).String() {
		require.Equal(t, &types.QueryAllianceValidatorsResponse{
			Validators: []types.QueryAllianceValidatorResponse{
				{
					ValidatorAddr: getOperator(val).String(),
					TotalDelegationShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec(AllianceDenom, math.LegacyNewDec(1000000)),
					),
					ValidatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec(AllianceDenom, math.LegacyNewDec(1000000)),
					),
					TotalStaked: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec(AllianceDenom, math.LegacyNewDec(1000_000)),
					),
				},
			},
			Pagination: queryVal.Pagination,
		}, queryVal)
	} else {
		require.Equal(t, &types.QueryAllianceValidatorsResponse{
			Validators: []types.QueryAllianceValidatorResponse{
				{
					ValidatorAddr: val2.GetOperator(),
					TotalDelegationShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec(AllianceDenomTwo, math.LegacyNewDec(1000000)),
					),
					ValidatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec(AllianceDenomTwo, math.LegacyNewDec(1000000)),
					),
					TotalStaked: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec(AllianceDenomTwo, math.LegacyNewDec(1000_000)),
					),
				},
			},
			Pagination: queryVal.Pagination,
		}, queryVal)
	}

	queryVal2, queryErr := queryServer.AllAllianceValidators(ctx, &types.QueryAllAllianceValidatorsRequest{
		Pagination: &query.PageRequest{
			Key:        queryVal.Pagination.NextKey,
			Offset:     0,
			Limit:      1,
			CountTotal: false,
			Reverse:    false,
		},
	})

	require.NoError(t, queryErr)
	// Order in which validators are returned is not deterministic since we randomly generate validator addresses
	if queryVal2.Validators[0].ValidatorAddr == getOperator(val).String() {
		require.Equal(t, &types.QueryAllianceValidatorsResponse{
			Validators: []types.QueryAllianceValidatorResponse{
				{
					ValidatorAddr: getOperator(val).String(),
					TotalDelegationShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec(AllianceDenom, math.LegacyNewDec(1000000)),
					),
					ValidatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec(AllianceDenom, math.LegacyNewDec(1000000)),
					),
					TotalStaked: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec(AllianceDenom, math.LegacyNewDec(1000_000)),
					),
				},
			},
			Pagination: queryVal2.Pagination,
		}, queryVal2)
	} else {
		require.Equal(t, &types.QueryAllianceValidatorsResponse{
			Validators: []types.QueryAllianceValidatorResponse{
				{
					ValidatorAddr: val2.GetOperator(),
					TotalDelegationShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec(AllianceDenomTwo, math.LegacyNewDec(1000000)),
					),
					ValidatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec(AllianceDenomTwo, math.LegacyNewDec(1000000)),
					),
					TotalStaked: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec(AllianceDenomTwo, math.LegacyNewDec(1000_000)),
					),
				},
			},
			Pagination: queryVal2.Pagination,
		}, queryVal2)
	}
}
