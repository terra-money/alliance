package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"alliance/x/alliance/types"
)

func TestQueryAllianceDelegation(t *testing.T) {
	// GIVEN: THE BLOCKCHAIN WITH ALLIANCES ON GENESIS
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
		},
	})
	delegations := app.StakingKeeper.GetAllDelegations(ctx)
	delAddr, _ := sdk.AccAddressFromBech32(delegations[0].DelegatorAddress)
	valAddr, _ := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	val, _ := app.StakingKeeper.GetValidator(ctx, valAddr)
	app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(2000_000))))
	app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr, sdk.NewCoins(sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(2000_000))))

	// WHEN: DELEGATING AND QUERYING ...
	delegationTxRes, txErr := app.AllianceKeeper.Delegate(ctx, delAddr, val, sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(1000_000)))
	queryDelegation, queryErr := app.AllianceKeeper.AllianceDelegation(ctx, &types.QueryAllianceDelegationRequest{
		DelegatorAddr: delAddr.String(),
		ValidatorAddr: val.OperatorAddress,
		Denom:         ALLIANCE_TOKEN_DENOM,
	})

	// THEN: VALIDATE THAT NO ERRORS HAVE BEEN PRODUCED AND BOTH OUTPUTS ARE AS WE EXPECT
	require.Nil(t, txErr)
	require.Nil(t, queryErr)
	require.Equal(t, &types.QueryAllianceDelegationResponse{
		Delegation: types.DelegationResponse{
			Delegation: types.Delegation{
				DelegatorAddress: delAddr.String(),
				ValidatorAddress: val.OperatorAddress,
				Denom:            ALLIANCE_TOKEN_DENOM,
				Shares:           sdk.NewDec(1000_000),
				RewardIndices:    nil,
			},
			Balance: sdk.Coin{
				Denom:  ALLIANCE_TOKEN_DENOM,
				Amount: sdk.NewInt(1000_000),
			},
		},
	}, queryDelegation)
	require.Equal(t, &types.Delegation{
		DelegatorAddress: delAddr.String(),
		ValidatorAddress: val.OperatorAddress,
		Denom:            ALLIANCE_TOKEN_DENOM,
		Shares:           sdk.NewDec(1000_000),
		RewardIndices:    []types.RewardIndex{},
	}, delegationTxRes)
}

func TestQueryAllianceDelegationNotFound(t *testing.T) {
	// GIVEN: THE BLOCKCHAIN
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	delegations := app.StakingKeeper.GetAllDelegations(ctx)
	delAddr, _ := sdk.AccAddressFromBech32(delegations[0].DelegatorAddress)
	valAddr, _ := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	val, _ := app.StakingKeeper.GetValidator(ctx, valAddr)

	// WHEN: QUERYING ...
	_, err := app.AllianceKeeper.AllianceDelegation(ctx, &types.QueryAllianceDelegationRequest{
		DelegatorAddr: delAddr.String(),
		ValidatorAddr: val.OperatorAddress,
		Denom:         ALLIANCE_TOKEN_DENOM,
	})

	// THEN: VALIDATE THAT NO ERRORS HAVE BEEN PRODUCED AND BOTH OUTPUTS ARE AS WE EXPECT
	require.Equal(t, err, status.Error(codes.NotFound, "AllianceAsset not found by denom alliance"))
}

func TestQueryAllianceValidatorNotFound(t *testing.T) {
	// GIVEN: THE BLOCKCHAIN
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	delegations := app.StakingKeeper.GetAllDelegations(ctx)
	delAddr, _ := sdk.AccAddressFromBech32(delegations[0].DelegatorAddress)

	// WHEN: QUERYING ...
	_, err := app.AllianceKeeper.AllianceDelegation(ctx, &types.QueryAllianceDelegationRequest{
		DelegatorAddr: delAddr.String(),
		ValidatorAddr: "cosmosvaloper19lss6zgdh5vvcpjhfftdghrpsw7a4434elpwpu",
		Denom:         ALLIANCE_TOKEN_DENOM,
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
				Denom:        ALLIANCE_TOKEN_DENOM,
				RewardWeight: sdk.NewDec(2),
				TakeRate:     sdk.NewDec(0),
				TotalTokens:  sdk.ZeroInt(),
			},
		},
	})
	delegations := app.StakingKeeper.GetAllDelegations(ctx)
	delAddr, _ := sdk.AccAddressFromBech32(delegations[0].DelegatorAddress)
	valAddr, _ := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	val, _ := app.StakingKeeper.GetValidator(ctx, valAddr)
	app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(2000_000))))
	app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr, sdk.NewCoins(sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(2000_000))))

	// WHEN: DELEGATING AND QUERYING ...
	delegationTxRes, txErr := app.AllianceKeeper.Delegate(ctx, delAddr, val, sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(1000_000)))
	queryDelegation, queryErr := app.AllianceKeeper.AlliancesDelegationByValidator(ctx, &types.QueryAlliancesDelegationByValidatorRequest{
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
					DelegatorAddress: delAddr.String(),
					ValidatorAddress: val.OperatorAddress,
					Denom:            ALLIANCE_TOKEN_DENOM,
					Shares:           sdk.NewDec(1000_000),
					RewardIndices:    nil,
				},
				Balance: sdk.Coin{
					Denom:  ALLIANCE_TOKEN_DENOM,
					Amount: sdk.NewInt(1000_000),
				},
			},
		},
		Pagination: &query.PageResponse{
			NextKey: nil,
			Total:   1,
		},
	}, queryDelegation)
	require.Equal(t, &types.Delegation{
		DelegatorAddress: delAddr.String(),
		ValidatorAddress: val.OperatorAddress,
		Denom:            ALLIANCE_TOKEN_DENOM,
		Shares:           sdk.NewDec(1000_000),
		RewardIndices:    []types.RewardIndex{},
	}, delegationTxRes)
}

func TestQueryAlliancesDelegationByValidatorNotFound(t *testing.T) {
	// GIVEN: THE BLOCKCHAIN
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	delegations := app.StakingKeeper.GetAllDelegations(ctx)
	delAddr, _ := sdk.AccAddressFromBech32(delegations[0].DelegatorAddress)

	// WHEN: QUERYING ...
	_, err := app.AllianceKeeper.AlliancesDelegationByValidator(ctx, &types.QueryAlliancesDelegationByValidatorRequest{
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
				Denom:        ALLIANCE_TOKEN_DENOM,
				RewardWeight: sdk.NewDec(2),
				TakeRate:     sdk.NewDec(0),
				TotalTokens:  sdk.ZeroInt(),
			},
			{
				Denom:        ALLIANCE_2_TOKEN_DENOM,
				RewardWeight: sdk.NewDec(10),
				TakeRate:     sdk.MustNewDecFromStr("0.14159265359"),
				TotalTokens:  sdk.ZeroInt(),
			},
		},
	})
	delegations := app.StakingKeeper.GetAllDelegations(ctx)
	delAddr, _ := sdk.AccAddressFromBech32(delegations[0].DelegatorAddress)
	valAddr, _ := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	val, _ := app.StakingKeeper.GetValidator(ctx, valAddr)
	app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(2000_000))))
	app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr, sdk.NewCoins(sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(2000_000))))
	app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(ALLIANCE_2_TOKEN_DENOM, sdk.NewInt(2000_000))))
	app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr, sdk.NewCoins(sdk.NewCoin(ALLIANCE_2_TOKEN_DENOM, sdk.NewInt(2000_000))))

	// WHEN: DELEGATING AND QUERYING ...
	delegationTxRes, txErr := app.AllianceKeeper.Delegate(ctx, delAddr, val, sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(1000_000)))
	delegation2TxRes, tx2Err := app.AllianceKeeper.Delegate(ctx, delAddr, val, sdk.NewCoin(ALLIANCE_2_TOKEN_DENOM, sdk.NewInt(1000_000)))
	queryDelegation, queryErr := app.AllianceKeeper.AlliancesDelegation(ctx, &types.QueryAlliancesDelegationsRequest{
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
					DelegatorAddress: delAddr.String(),
					ValidatorAddress: val.OperatorAddress,
					Denom:            ALLIANCE_TOKEN_DENOM,
					Shares:           sdk.NewDec(1000_000),
					RewardIndices:    nil,
				},
				Balance: sdk.Coin{
					Denom:  ALLIANCE_TOKEN_DENOM,
					Amount: sdk.NewInt(1000_000),
				},
			},
			{
				Delegation: types.Delegation{
					DelegatorAddress: delAddr.String(),
					ValidatorAddress: val.OperatorAddress,
					Denom:            ALLIANCE_2_TOKEN_DENOM,
					Shares:           sdk.NewDec(1000_000),
					RewardIndices:    nil,
				},
				Balance: sdk.Coin{
					Denom:  ALLIANCE_2_TOKEN_DENOM,
					Amount: sdk.NewInt(1000_000),
				},
			},
		},
		Pagination: &query.PageResponse{
			NextKey: nil,
			Total:   2,
		},
	}, queryDelegation)
	require.Equal(t, &types.Delegation{
		DelegatorAddress: delAddr.String(),
		ValidatorAddress: val.OperatorAddress,
		Denom:            ALLIANCE_TOKEN_DENOM,
		Shares:           sdk.NewDec(1000_000),
		RewardIndices:    []types.RewardIndex{},
	}, delegationTxRes)
	require.Equal(t, &types.Delegation{
		DelegatorAddress: delAddr.String(),
		ValidatorAddress: val.OperatorAddress,
		Denom:            ALLIANCE_2_TOKEN_DENOM,
		Shares:           sdk.NewDec(1000_000),
		RewardIndices:    nil,
	}, delegation2TxRes)
}
