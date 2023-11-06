package bindings_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"github.com/terra-money/alliance/app"
	"github.com/terra-money/alliance/x/alliance/bindings"
	bindingtypes "github.com/terra-money/alliance/x/alliance/bindings/types"
	"github.com/terra-money/alliance/x/alliance/types"
)

func createTestContext(t *testing.T) (*app.App, sdk.Context) {
	app := app.Setup(t)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	return app, ctx
}

var AllianceDenom = "alliance"

func TestAssetQuery(t *testing.T) {
	app, ctx := createTestContext(t)
	genesisTime := ctx.BlockTime()
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			types.NewAllianceAsset(AllianceDenom, sdk.NewDec(2), sdk.ZeroDec(), sdk.NewDec(5), sdk.NewDec(0), genesisTime),
		},
	})

	querierPlugin := bindings.NewAllianceQueryPlugin(app.AllianceKeeper)
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
		RewardWeight:         sdk.MustNewDecFromStr("2").String(),
		TakeRate:             sdk.MustNewDecFromStr("0").String(),
		TotalTokens:          "0",
		TotalValidatorShares: sdk.MustNewDecFromStr("0").String(),
		RewardStartTime:      uint64(genesisTime.Nanosecond()),
		RewardChangeRate:     sdk.MustNewDecFromStr("1").String(),
		LastRewardChangeTime: 0,
		RewardWeightRange: bindingtypes.RewardWeightRange{
			Min: sdk.MustNewDecFromStr("0").String(),
			Max: sdk.MustNewDecFromStr("5").String(),
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
			types.NewAllianceAsset(AllianceDenom, sdk.NewDec(2), sdk.ZeroDec(), sdk.NewDec(5), sdk.NewDec(0), genesisTime),
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
	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(AllianceDenom, sdk.NewInt(2000_000))))
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr, sdk.NewCoins(sdk.NewCoin(AllianceDenom, sdk.NewInt(2000_000))))
	require.NoError(t, err)

	// Check current total staked tokens
	totalBonded := app.StakingKeeper.TotalBondedTokens(ctx)
	require.Equal(t, sdk.NewInt(1000_000), totalBonded)

	// Delegate
	_, err = app.AllianceKeeper.Delegate(ctx, delAddr, val, sdk.NewCoin(AllianceDenom, sdk.NewInt(1000_000)))
	require.NoError(t, err)

	querierPlugin := bindings.NewAllianceQueryPlugin(app.AllianceKeeper)
	querier := bindings.CustomQuerier(querierPlugin)

	delegationQuery := bindingtypes.AllianceQuery{
		Delegation: &bindingtypes.Delegation{
			Delegator: delAddr.String(),
			Validator: val.GetOperator().String(),
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
		Validator: val.GetOperator().String(),
		Denom:     AllianceDenom,
		Amount: bindingtypes.Coin{
			Denom:  AllianceDenom,
			Amount: "1000000",
		},
	}, delegationResponse)
}

func TestDelegationRewardsQuery(t *testing.T) {
	app, ctx := createTestContext(t)
	genesisTime := ctx.BlockTime()
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			types.NewAllianceAsset(AllianceDenom, sdk.NewDec(2), sdk.ZeroDec(), sdk.NewDec(5), sdk.NewDec(0), genesisTime),
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
	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(AllianceDenom, sdk.NewInt(2000_000))))
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr, sdk.NewCoins(sdk.NewCoin(AllianceDenom, sdk.NewInt(2000_000))))
	require.NoError(t, err)

	// Check current total staked tokens
	totalBonded := app.StakingKeeper.TotalBondedTokens(ctx)
	require.Equal(t, sdk.NewInt(1000_000), totalBonded)

	// Delegate
	_, err = app.AllianceKeeper.Delegate(ctx, delAddr, val, sdk.NewCoin(AllianceDenom, sdk.NewInt(1000_000)))
	require.NoError(t, err)

	assets := app.AllianceKeeper.GetAllAssets(ctx)
	err = app.AllianceKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)

	// Transfer to reward pool
	mintPoolAddr := app.AccountKeeper.GetModuleAddress(minttypes.ModuleName)
	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(4000_000))))
	require.NoError(t, err)
	err = app.AllianceKeeper.AddAssetsToRewardPool(ctx, mintPoolAddr, val, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(2000_000))))
	require.NoError(t, err)

	querierPlugin := bindings.NewAllianceQueryPlugin(app.AllianceKeeper)
	querier := bindings.CustomQuerier(querierPlugin)

	delegationQuery := bindingtypes.AllianceQuery{
		DelegationRewards: &bindingtypes.DelegationRewards{
			Delegator: delAddr.String(),
			Validator: val.GetOperator().String(),
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
		Rewards: []bindingtypes.Coin{
			{
				Denom:  "stake",
				Amount: "2000000",
			},
		},
	}, response)
}
