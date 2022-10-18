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
			RewardDelayTime:     time.Duration(1000000),
			GlobalRewardIndices: types.RewardIndices{},
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
	globalIndices := app.AllianceKeeper.GlobalRewardIndices(ctx)
	require.Equal(t, types.RewardIndices{
		types.RewardIndex{
			Denom: "stake",
			Index: sdk.NewDec(1),
		},
	}, globalIndices)

	// New delegation from user 2
	_, err = app.AllianceKeeper.Delegate(ctx, user2, val1, sdk.NewCoin(ALLIANCE_2_TOKEN_DENOM, sdk.NewInt(1000_000)))
	require.NoError(t, err)

	// Transfer to reward pool
	err = app.AllianceKeeper.AddAssetsToRewardPool(ctx, mintPoolAddr, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(2000_000))))
	require.NoError(t, err)

	globalIndices = app.AllianceKeeper.GlobalRewardIndices(ctx)
	require.Equal(t, types.RewardIndices{
		types.RewardIndex{
			Denom: "stake",
			Index: sdk.NewDec(14).Quo(sdk.NewDec(12)),
		},
	}, globalIndices)

	// Transfer another token to reward pool
	err = app.AllianceKeeper.AddAssetsToRewardPool(ctx, mintPoolAddr, sdk.NewCoins(sdk.NewCoin("stake2", sdk.NewInt(4000_000))))
	require.NoError(t, err)

	// Expect global index to be updated
	// 14/12 + 4/12 = 18/12
	globalIndices = app.AllianceKeeper.GlobalRewardIndices(ctx)
	require.Equal(t, types.RewardIndices{
		types.RewardIndex{
			Denom: "stake",
			Index: sdk.NewDec(14).Quo(sdk.NewDec(12)),
		},
		types.RewardIndex{
			Denom: "stake2",
			Index: sdk.NewDec(4).Quo(sdk.NewDec(12)),
		},
	}, globalIndices)
}

func TestClaimRewards(t *testing.T) {
	app, ctx := createTestContext(t)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.Params{
			RewardDelayTime:     time.Duration(1000000),
			GlobalRewardIndices: types.RewardIndices{},
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
	mintPoolAddr := app.AccountKeeper.GetModuleAddress(minttypes.ModuleName)
	rewardsPoolAddr := app.AccountKeeper.GetModuleAddress(types.RewardsPoolName)
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

	// New delegation from user 1
	_, err = app.AllianceKeeper.Delegate(ctx, user1, val1, sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(1000_000)))
	require.NoError(t, err)

	// Transfer to reward pool
	err = app.AllianceKeeper.AddAssetsToRewardPool(ctx, mintPoolAddr, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(2000_000))))
	require.NoError(t, err)

	// New delegation from user 2
	_, err = app.AllianceKeeper.Delegate(ctx, user2, val1, sdk.NewCoin(ALLIANCE_2_TOKEN_DENOM, sdk.NewInt(1000_000)))
	require.NoError(t, err)

	// Transfer to reward pool
	err = app.AllianceKeeper.AddAssetsToRewardPool(ctx, mintPoolAddr, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(2000_000))))
	require.NoError(t, err)

	// Transfer another token to reward pool
	err = app.AllianceKeeper.AddAssetsToRewardPool(ctx, mintPoolAddr, sdk.NewCoins(sdk.NewCoin("stake2", sdk.NewInt(4000_000))))
	require.NoError(t, err)

	// before claiming, there should be tokens in rewards pool
	coins := app.BankKeeper.GetAllBalances(ctx, rewardsPoolAddr)

	// User 1 claims rewards
	// User 1 has 1 STAKE (2 Power)
	// Added 2 stake rewards (fully belonging to user 1)
	// User 2 has 1 STAKE (10 Power)
	// Added 2 stake rewards (user1: 2/12 * 2, user2: 10/12 * 2)
	// Added 4 stake2 rewards (user1: 2/12 * 4, user2: 10/12 * 4)
	coins, err = app.AllianceKeeper.ClaimDelegationRewards(ctx, user1, val1, ALLIANCE_TOKEN_DENOM)
	require.NoError(t, err)
	require.Equal(t, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(2_333_333)), sdk.NewCoin("stake2", sdk.NewInt(666_666))), coins)

	// User 2 claims rewards but doesn't the right denom
	_, err = app.AllianceKeeper.ClaimDelegationRewards(ctx, user2, val1, ALLIANCE_TOKEN_DENOM)
	require.Error(t, err)

	// User 2 claims rewards
	coins, err = app.AllianceKeeper.ClaimDelegationRewards(ctx, user2, val1, ALLIANCE_2_TOKEN_DENOM)
	require.NoError(t, err)
	require.Equal(t, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(1_666_666)), sdk.NewCoin("stake2", sdk.NewInt(3_333_333))), coins)

	// After claiming, there should be nothing left in rewards pool
	// Some rounding left
	coins = app.BankKeeper.GetAllBalances(ctx, rewardsPoolAddr)
	require.Equal(t, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(1)), sdk.NewCoin("stake2", sdk.NewInt(1))), coins)

	// Global indices
	indices := app.AllianceKeeper.GlobalRewardIndices(ctx)

	// Check that all delegations have updated local indices
	delegation, found := app.AllianceKeeper.GetDelegation(ctx, user1, val1, ALLIANCE_TOKEN_DENOM)
	require.True(t, found)
	require.Equal(t, indices, types.NewRewardIndices(delegation.RewardIndices))

	delegation, found = app.AllianceKeeper.GetDelegation(ctx, user2, val1, ALLIANCE_2_TOKEN_DENOM)
	require.True(t, found)
	require.Equal(t, indices, types.NewRewardIndices(delegation.RewardIndices))
}
