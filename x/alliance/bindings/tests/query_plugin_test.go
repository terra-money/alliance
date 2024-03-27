package bindings_test

import (
	"encoding/json"
	"testing"

	"cosmossdk.io/math"

	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/terra-money/alliance/app"
	"github.com/terra-money/alliance/x/alliance/bindings"
	bindingtypes "github.com/terra-money/alliance/x/alliance/bindings/types"
	"github.com/terra-money/alliance/x/alliance/types"
)

func createTestContext(t *testing.T) (*app.App, sdk.Context) {
	app := app.Setup(t)
	ctx := app.BaseApp.NewContext(false)
	return app, ctx
}

var AllianceDenom = "alliance"

func TestAssetQuery(t *testing.T) {
	app, ctx := createTestContext(t)
	genesisTime := ctx.BlockTime()
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			types.NewAllianceAsset(AllianceDenom, math.LegacyNewDec(2), math.LegacyZeroDec(), math.LegacyNewDec(5), math.LegacyNewDec(0), genesisTime),
		},
	})

	querierPlugin := bindings.NewAllianceQueryPlugin(&app.AllianceKeeper)
	querier := bindings.CustomQuerier(querierPlugin)

	assetQuery := bindingtypes.AllianceQuery{
		Alliance: &bindingtypes.Alliance{
			Denom: AllianceDenom,
		},
	}
	qBz, err := json.Marshal(assetQuery)
	require.NoError(t, err)
	rBz, err := querier(ctx, qBz)
	require.NoError(t, err)

	var assetResponse bindingtypes.AllianceResponse
	err = json.Unmarshal(rBz, &assetResponse)
	require.NoError(t, err)

	require.Equal(t, bindingtypes.AllianceResponse{
		Denom:                AllianceDenom,
		RewardWeight:         math.LegacyMustNewDecFromStr("2").String(),
		TakeRate:             math.LegacyMustNewDecFromStr("0").String(),
		TotalTokens:          "0",
		TotalValidatorShares: math.LegacyMustNewDecFromStr("0").String(),
		RewardStartTime:      uint64(genesisTime.Nanosecond()),
		RewardChangeRate:     math.LegacyMustNewDecFromStr("1").String(),
		LastRewardChangeTime: 0,
		RewardWeightRange: bindingtypes.RewardWeightRange{
			Min: math.LegacyMustNewDecFromStr("0").String(),
			Max: math.LegacyMustNewDecFromStr("5").String(),
		},
		IsInitialized: false,
	}, assetResponse)
}

func TestDelegationQuery(t *testing.T) {
	app, ctx := createTestContext(t)
	genesisTime := ctx.BlockTime()
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			types.NewAllianceAsset(AllianceDenom, math.LegacyNewDec(2), math.LegacyZeroDec(), math.LegacyNewDec(5), math.LegacyNewDec(0), genesisTime),
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
	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(AllianceDenom, math.NewInt(2000_000))))
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr, sdk.NewCoins(sdk.NewCoin(AllianceDenom, math.NewInt(2000_000))))
	require.NoError(t, err)

	// Check current total staked tokens
	totalBonded, err := app.StakingKeeper.TotalBondedTokens(ctx)
	require.NoError(t, err)
	require.Equal(t, math.NewInt(1000_000), totalBonded)

	// Delegate
	_, err = app.AllianceKeeper.Delegate(ctx, delAddr, val, sdk.NewCoin(AllianceDenom, math.NewInt(1000_000)))
	require.NoError(t, err)

	querierPlugin := bindings.NewAllianceQueryPlugin(&app.AllianceKeeper)
	querier := bindings.CustomQuerier(querierPlugin)

	delegationQuery := bindingtypes.AllianceQuery{
		Delegation: &bindingtypes.Delegation{
			Delegator: delAddr.String(),
			Validator: val.GetOperator(),
			Denom:     AllianceDenom,
		},
	}
	qBz, err := json.Marshal(delegationQuery)
	require.NoError(t, err)
	rBz, err := querier(ctx, qBz)
	require.NoError(t, err)

	var delegationResponse bindingtypes.DelegationResponse
	err = json.Unmarshal(rBz, &delegationResponse)
	require.NoError(t, err)

	require.Equal(t, bindingtypes.DelegationResponse{
		Delegator: delAddr.String(),
		Validator: val.GetOperator(),
		Denom:     AllianceDenom,
		Amount:    "1000000",
	}, delegationResponse)
}

func TestDelegationRewardsQuery(t *testing.T) {
	app, ctx := createTestContext(t)
	genesisTime := ctx.BlockTime()
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			types.NewAllianceAsset(AllianceDenom, math.LegacyNewDec(2), math.LegacyZeroDec(), math.LegacyNewDec(5), math.LegacyNewDec(0), genesisTime),
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
	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(AllianceDenom, math.NewInt(2000_000))))
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr, sdk.NewCoins(sdk.NewCoin(AllianceDenom, math.NewInt(2000_000))))
	require.NoError(t, err)

	// Check current total staked tokens
	totalBonded, err := app.StakingKeeper.TotalBondedTokens(ctx)
	require.NoError(t, err)
	require.Equal(t, math.NewInt(1000_000), totalBonded)

	// Delegate
	_, err = app.AllianceKeeper.Delegate(ctx, delAddr, val, sdk.NewCoin(AllianceDenom, math.NewInt(1000_000)))
	require.NoError(t, err)

	assets := app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)

	// Transfer to reward pool
	mintPoolAddr := app.AccountKeeper.GetModuleAddress(minttypes.ModuleName)
	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(4000_000))))
	require.NoError(t, err)
	err = app.AllianceKeeper.AddAssetsToRewardPool(ctx, mintPoolAddr, val, sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(2000_000))))
	require.NoError(t, err)

	querierPlugin := bindings.NewAllianceQueryPlugin(&app.AllianceKeeper)
	querier := bindings.CustomQuerier(querierPlugin)

	delegationQuery := bindingtypes.AllianceQuery{
		DelegationRewards: &bindingtypes.DelegationRewards{
			Delegator: delAddr.String(),
			Validator: val.GetOperator(),
			Denom:     AllianceDenom,
		},
	}
	qBz, err := json.Marshal(delegationQuery)
	require.NoError(t, err)
	rBz, err := querier(ctx, qBz)
	require.NoError(t, err)

	var response bindingtypes.DelegationRewardsResponse
	err = json.Unmarshal(rBz, &response)
	require.NoError(t, err)

	require.Equal(t, bindingtypes.DelegationRewardsResponse{
		Rewards: []sdk.Coin{
			{
				Denom:  "stake",
				Amount: math.NewInt(2000000),
			},
		},
	}, response)
}

func TestCustomQuerier(t *testing.T) {
	app, ctx := createTestContext(t)
	genesisTime := ctx.BlockTime()
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			types.NewAllianceAsset(AllianceDenom, math.LegacyNewDec(2), math.LegacyZeroDec(), math.LegacyNewDec(5), math.LegacyNewDec(0), genesisTime),
		},
	})

	querierPlugin := bindings.NewAllianceQueryPlugin(&app.AllianceKeeper)
	querier := bindings.CustomQuerier(querierPlugin)

	queryBytes := []byte("{\"random\": \"query\"}")
	_, err := querier(ctx, queryBytes)
	require.Error(t, err)
}
