package keeper_test

import (
	test_helpers "alliance/app"
	"alliance/x/alliance/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestRewardPoolAndGlobalIndex(t *testing.T) {
	app, ctx := createTestContext(t)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.Params{
			RewardDelayTime: time.Duration(1000000),
			GlobalIndex:     sdk.NewDec(0),
		},
		Assets: []types.AllianceAsset{
			{
				Denom:        ALLIANCE_TOKEN_DENOM,
				RewardWeight: sdk.NewDec(2),
				TakeRate:     sdk.NewDec(0),
				TotalShares:  sdk.NewDec(0),
				TotalTokens:  sdk.ZeroInt(),
			},
			{
				Denom:        ALLIANCE_2_TOKEN_DENOM,
				RewardWeight: sdk.NewDec(10),
				TakeRate:     sdk.NewDec(0),
				TotalShares:  sdk.NewDec(0),
				TotalTokens:  sdk.ZeroInt(),
			},
		},
	})

	// Accounts
	rewardsPoolAddr := app.AccountKeeper.GetModuleAddress(types.RewardsPoolName)
	mintPoolAddr := app.AccountKeeper.GetModuleAddress(minttypes.ModuleName)
	delegations := app.StakingKeeper.GetAllDelegations(ctx)
	valAddr1, err := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	require.NoError(t, err)
	val1, found := app.StakingKeeper.GetValidator(ctx, valAddr1)
	require.True(t, found)
	addrs := test_helpers.AddTestAddrsIncremental(app, ctx, 2, sdk.NewCoins(
		sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(1000_000)),
		sdk.NewCoin(ALLIANCE_2_TOKEN_DENOM, sdk.NewInt(1000_000)),
	))
	user1 := addrs[0]
	user2 := addrs[1]

	// Mint tokens
	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(4000_000))))
	require.NoError(t, err)
	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin("stake2", sdk.NewInt(4000_000))))
	require.NoError(t, err)
	coin := app.BankKeeper.GetBalance(ctx, mintPoolAddr, "stake")
	require.Equal(t, sdk.NewCoin("stake", sdk.NewInt(4000_000)), coin)

	// Transfer to reward pool without delegations will fail
	err = app.AllianceKeeper.AddAssetsToRewardPool(ctx, mintPoolAddr, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(2000_000))))
	require.Error(t, err)

	_, err = app.AllianceKeeper.Delegate(ctx, user1, val1, sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(1000_000)))
	require.NoError(t, err)

	// Transfer to reward pool
	err = app.AllianceKeeper.AddAssetsToRewardPool(ctx, mintPoolAddr, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(2000_000))))
	require.NoError(t, err)

	// Expect rewards pool to have something
	balance := app.BankKeeper.GetBalance(ctx, rewardsPoolAddr, "stake")
	require.Equal(t, sdk.NewCoin("stake", sdk.NewInt(2000_000)), balance)

	// Expect global index to be updated
	globalIndex := app.AllianceKeeper.GlobalRewardIndex(ctx)
	require.Equal(t, sdk.NewDec(1), globalIndex)

	// New delegation from user 2
	_, err = app.AllianceKeeper.Delegate(ctx, user2, val1, sdk.NewCoin(ALLIANCE_2_TOKEN_DENOM, sdk.NewInt(1000_000)))
	require.NoError(t, err)

	// Transfer to reward pool
	err = app.AllianceKeeper.AddAssetsToRewardPool(ctx, mintPoolAddr, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(2000_000))))
	require.NoError(t, err)

	// Expect global index to be updated
	// 1 + 2 / (2 + 10) = 14/12
	globalIndex = app.AllianceKeeper.GlobalRewardIndex(ctx)
	require.Equal(t, sdk.NewDec(14).Quo(sdk.NewDec(12)), globalIndex)

	// Transfer another token to reward pool
	err = app.AllianceKeeper.AddAssetsToRewardPool(ctx, mintPoolAddr, sdk.NewCoins(sdk.NewCoin("stake2", sdk.NewInt(4000_000))))
	require.NoError(t, err)

	// Expect global index to be updated
	// 14/12 + 4/12 = 18/12
	globalIndex = app.AllianceKeeper.GlobalRewardIndex(ctx)
	require.Equal(t, sdk.NewDec(18).Quo(sdk.NewDec(12)), globalIndex)
}
