package keeper_test

import (
	test_helpers "alliance/app"
	"alliance/x/alliance/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/staking/teststaking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	"testing"
	"time"
)

func TestRewardPoolAndGlobalIndex(t *testing.T) {
	app, ctx := createTestContext(t)
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
	err = app.AllianceKeeper.AddAssetsToRewardPool(ctx, mintPoolAddr, valAddr1, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(2000_000))))
	require.Error(t, err)

	_, err = app.AllianceKeeper.Delegate(ctx, user1, val1, sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(1000_000)))
	require.NoError(t, err)

	// Transfer to reward pool
	err = app.AllianceKeeper.AddAssetsToRewardPool(ctx, mintPoolAddr, valAddr1, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(2000_000))))
	require.NoError(t, err)

	// Expect rewards pool to have something
	balance := app.BankKeeper.GetBalance(ctx, rewardsPoolAddr, "stake")
	require.Equal(t, sdk.NewCoin("stake", sdk.NewInt(2000_000)), balance)

	// Expect validator global index to be updated
	aVal := app.AllianceKeeper.GetOrCreateValidator(ctx, valAddr1)
	globalIndices := types.NewRewardIndices(aVal.RewardIndices)
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
	err = app.AllianceKeeper.AddAssetsToRewardPool(ctx, mintPoolAddr, valAddr1, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(2000_000))))
	require.NoError(t, err)

	aVal = app.AllianceKeeper.GetOrCreateValidator(ctx, valAddr1)
	globalIndices = types.NewRewardIndices(aVal.RewardIndices)
	require.Equal(t, types.RewardIndices{
		types.RewardIndex{
			Denom: "stake",
			Index: sdk.NewDec(14).Quo(sdk.NewDec(12)),
		},
	}, globalIndices)

	// Transfer another token to reward pool
	err = app.AllianceKeeper.AddAssetsToRewardPool(ctx, mintPoolAddr, valAddr1, sdk.NewCoins(sdk.NewCoin("stake2", sdk.NewInt(4000_000))))
	require.NoError(t, err)

	// Expect global index to be updated
	// 14/12 + 4/12 = 18/12
	aVal = app.AllianceKeeper.GetOrCreateValidator(ctx, valAddr1)
	globalIndices = types.NewRewardIndices(aVal.RewardIndices)
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
	err = app.AllianceKeeper.AddAssetsToRewardPool(ctx, mintPoolAddr, valAddr1, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(2000_000))))
	require.NoError(t, err)

	// New delegation from user 2
	_, err = app.AllianceKeeper.Delegate(ctx, user2, val1, sdk.NewCoin(ALLIANCE_2_TOKEN_DENOM, sdk.NewInt(1000_000)))
	require.NoError(t, err)

	// Transfer to reward pool
	err = app.AllianceKeeper.AddAssetsToRewardPool(ctx, mintPoolAddr, valAddr1, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(2000_000))))
	require.NoError(t, err)

	aVal1 := app.AllianceKeeper.GetOrCreateValidator(ctx, valAddr1)
	asset, _ := app.AllianceKeeper.GetAssetByDenom(ctx, ALLIANCE_TOKEN_DENOM)
	require.Equal(t,
		sdk.NewInt(1000_000),
		aVal1.TotalTokensWithAsset(asset),
	)
	asset, _ = app.AllianceKeeper.GetAssetByDenom(ctx, ALLIANCE_2_TOKEN_DENOM)
	require.Equal(t,
		sdk.NewInt(1000_000),
		aVal1.TotalTokensWithAsset(asset),
	)

	// Transfer another token to reward pool
	err = app.AllianceKeeper.AddAssetsToRewardPool(ctx, mintPoolAddr, valAddr1, sdk.NewCoins(sdk.NewCoin("stake2", sdk.NewInt(4000_000))))
	require.NoError(t, err)

	// Make sure reward indices are right
	aVal1 = app.AllianceKeeper.GetOrCreateValidator(ctx, valAddr1)
	require.Equal(t,
		types.NewRewardIndices([]types.RewardIndex{
			{
				Denom: "stake",
				Index: sdk.MustNewDecFromStr("1.166666666666666667"),
			},
			{
				Denom: "stake2",
				Index: sdk.MustNewDecFromStr("0.333333333333333333"),
			},
		}),
		types.NewRewardIndices(aVal1.RewardIndices),
	)

	// before claiming, there should be tokens in rewards pool
	coins := app.BankKeeper.GetAllBalances(ctx, rewardsPoolAddr)
	require.Equal(t,
		sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(4000_000)), sdk.NewCoin("stake2", sdk.NewInt(4000_000))),
		coins,
	)

	// User 1 claims rewards
	// User 1 has 1 STAKE (2 Power)
	// Added 2 stake rewards (fully belonging to user 1)
	// User 2 has 1 STAKE (10 Power)
	// Added 2 stake rewards (user1: 2/12 * 2, user2: 10/12 * 2)
	// Added 4 stake2 rewards (user1: 2/12 * 4, user2: 10/12 * 4)
	coins, err = app.AllianceKeeper.ClaimDelegationRewards(ctx, user1, val1, ALLIANCE_TOKEN_DENOM)
	require.NoError(t, err)
	require.Equal(t, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(2_333_333)), sdk.NewCoin("stake2", sdk.NewInt(666_666))), coins)

	// User 2 claims rewards but doesn't use the right denom
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
	aVal := app.AllianceKeeper.GetOrCreateValidator(ctx, valAddr1)
	indices := types.NewRewardIndices(aVal.RewardIndices)

	// Check that all delegations have updated local indices
	delegation, found := app.AllianceKeeper.GetDelegation(ctx, user1, val1, ALLIANCE_TOKEN_DENOM)
	require.True(t, found)
	require.Equal(t, indices, types.NewRewardIndices(delegation.RewardIndices))

	delegation, found = app.AllianceKeeper.GetDelegation(ctx, user2, val1, ALLIANCE_2_TOKEN_DENOM)
	require.True(t, found)
	require.Equal(t, indices, types.NewRewardIndices(delegation.RewardIndices))
}

func TestClaimRewardsWithMultipleValidators(t *testing.T) {
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx.WithBlockTime(startTime)
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

	// remove genesis validator delegations
	delegations := app.StakingKeeper.GetAllDelegations(ctx)
	require.Len(t, delegations, 1)
	err := app.StakingKeeper.RemoveDelegation(ctx, stakingtypes.Delegation{
		ValidatorAddress: delegations[0].ValidatorAddress,
		DelegatorAddress: delegations[0].DelegatorAddress,
	})
	require.NoError(t, err)

	// Set tax and rewards to be zero for easier calculation
	distParams := app.DistrKeeper.GetParams(ctx)
	distParams.CommunityTax = sdk.ZeroDec()
	distParams.BaseProposerReward = sdk.ZeroDec()
	distParams.BonusProposerReward = sdk.ZeroDec()
	app.DistrKeeper.SetParams(ctx, distParams)

	// Accounts
	//mintPoolAddr := app.AccountKeeper.GetModuleAddress(minttypes.ModuleName)
	//rewardsPoolAddr := app.AccountKeeper.GetModuleAddress(types.RewardsPoolName)
	addrs := test_helpers.AddTestAddrsIncremental(app, ctx, 4, sdk.NewCoins(
		sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(1000_000)),
		sdk.NewCoin(ALLIANCE_2_TOKEN_DENOM, sdk.NewInt(1000_000)),
	))
	pks := test_helpers.CreateTestPubKeys(2)

	// Creating two validators: 1 with 0% commission, 1 with 100% commission
	valAddr1 := sdk.ValAddress(addrs[0])
	val1 := teststaking.NewValidator(t, valAddr1, pks[0])
	val1.Commission = stakingtypes.Commission{
		CommissionRates: stakingtypes.CommissionRates{
			Rate:          sdk.NewDec(0),
			MaxRate:       sdk.NewDec(0),
			MaxChangeRate: sdk.NewDec(0),
		},
		UpdateTime: time.Now(),
	}
	test_helpers.RegisterNewValidator(t, app, ctx, val1)

	valAddr2 := sdk.ValAddress(addrs[1])
	val2 := teststaking.NewValidator(t, valAddr2, pks[1])
	val2.Commission = stakingtypes.Commission{
		CommissionRates: stakingtypes.CommissionRates{
			Rate:          sdk.NewDec(1),
			MaxRate:       sdk.NewDec(1),
			MaxChangeRate: sdk.NewDec(0),
		},
		UpdateTime: time.Now(),
	}
	test_helpers.RegisterNewValidator(t, app, ctx, val2)

	user1 := addrs[2]
	user2 := addrs[3]

	// Mint tokens
	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(4000_000))))
	require.NoError(t, err)
	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin("stake2", sdk.NewInt(4000_000))))
	require.NoError(t, err)

	// New delegation from user 1 to val 1
	_, err = app.AllianceKeeper.Delegate(ctx, user1, val1, sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(1000_000)))
	require.NoError(t, err)

	// New delegation from user 2 to val 2
	_, err = app.AllianceKeeper.Delegate(ctx, user2, val2, sdk.NewCoin(ALLIANCE_2_TOKEN_DENOM, sdk.NewInt(1000_000)))
	require.NoError(t, err)

	// Transfer to rewards to fee pool to be distributed
	app.BankKeeper.SendCoinsFromModuleToModule(ctx, minttypes.ModuleName, authtypes.FeeCollectorName, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(4000_000))))

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	// Distribute in the next begin block
	// At the next begin block, tokens will be distributed from the fee pool
	cons1, _ := val1.GetConsAddr()
	cons2, _ := val2.GetConsAddr()
	var votingPower int64 = 12
	app.DistrKeeper.AllocateTokens(ctx, votingPower, votingPower, cons1, []abcitypes.VoteInfo{
		{
			Validator: abcitypes.Validator{
				Address: cons1,
				Power:   2,
			},
			SignedLastBlock: true,
		},
		{
			Validator: abcitypes.Validator{
				Address: cons2,
				Power:   10,
			},
			SignedLastBlock: true,
		},
	})

	commission := app.DistrKeeper.GetValidatorAccumulatedCommission(ctx, val1.GetOperator()).Commission
	require.Equal(t, sdk.NewInt(0), commission.AmountOf("stake").TruncateInt())
	commission = app.DistrKeeper.GetValidatorAccumulatedCommission(ctx, val2.GetOperator()).Commission
	require.Equal(t, sdk.NewInt(3333333), commission.AmountOf("stake").TruncateInt())

	rewards := app.DistrKeeper.GetValidatorCurrentRewards(ctx, val1.GetOperator()).Rewards
	require.Equal(t, sdk.NewInt(666666), rewards.AmountOf("stake").TruncateInt())
	rewards = app.DistrKeeper.GetValidatorCurrentRewards(ctx, val2.GetOperator()).Rewards
	require.Equal(t, sdk.NewInt(0), rewards.AmountOf("stake").TruncateInt())

	// User 1 should be getting all the rewards from validator 1 since it has 0 commission
	coins, err := app.AllianceKeeper.ClaimDelegationRewards(ctx, user1, val1, ALLIANCE_TOKEN_DENOM)
	require.NoError(t, err)
	require.Equal(t, sdk.NewInt(666666), coins.AmountOf("stake"))

	// User 2 should be getting no rewards since validator 2 has 100% commission
	coins, err = app.AllianceKeeper.ClaimDelegationRewards(ctx, user2, val2, ALLIANCE_2_TOKEN_DENOM)
	require.NoError(t, err)
	require.Equal(t, sdk.NewInt(0), coins.AmountOf("stake"))
}

func TestClaimTakeRate(t *testing.T) {
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime)
	ctx = ctx.WithBlockHeight(1)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.Params{
			RewardDelayTime:     time.Minute * 60,
			RewardClaimInterval: time.Minute * 5,
			LastRewardClaimTime: startTime,
		},
		Assets: []types.AllianceAsset{
			{
				Denom:        ALLIANCE_TOKEN_DENOM,
				RewardWeight: sdk.NewDec(2),
				TakeRate:     sdk.MustNewDecFromStr("0.5"),
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

	// Accounts
	//mintPoolAddr := app.AccountKeeper.GetModuleAddress(minttypes.ModuleName)
	//rewardsPoolAddr := app.AccountKeeper.GetModuleAddress(types.RewardsPoolName)
	feeCollectorAddr := app.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
	delegations := app.StakingKeeper.GetAllDelegations(ctx)
	valAddr1, err := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	require.NoError(t, err)
	val1, found := app.StakingKeeper.GetValidator(ctx, valAddr1)
	require.True(t, found)
	addrs := test_helpers.AddTestAddrsIncremental(app, ctx, 1, sdk.NewCoins(
		sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(1000_000_000)),
		sdk.NewCoin(ALLIANCE_2_TOKEN_DENOM, sdk.NewInt(1000_000_000)),
	))
	user1 := addrs[0]

	app.AllianceKeeper.Delegate(ctx, user1, val1, sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(1000_000_000)))
	app.AllianceKeeper.Delegate(ctx, user1, val1, sdk.NewCoin(ALLIANCE_2_TOKEN_DENOM, sdk.NewInt(1000_000_000)))

	// Calling it immediately will not update anything
	coins, err := app.AllianceKeeper.ClaimAssetsWithTakeRateRateLimited(ctx)
	require.Nil(t, coins)
	require.Nil(t, err)

	// Advance block time
	timePassed := time.Minute*5 + time.Second
	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(timePassed))
	ctx = ctx.WithBlockHeight(2)
	coinsClaimed, err := app.AllianceKeeper.ClaimAssetsWithTakeRateRateLimited(ctx)
	coins = app.BankKeeper.GetAllBalances(ctx, feeCollectorAddr)
	require.Equal(t, coinsClaimed, coins)

	expectedAmount := sdk.MustNewDecFromStr("0.5").Mul(sdk.NewDec(timePassed.Nanoseconds()).Quo(sdk.NewDec(31_557_000_000_000_000))).MulInt(sdk.NewInt(1000_000_000))
	require.Equal(t, expectedAmount.TruncateInt(), coins.AmountOf(ALLIANCE_TOKEN_DENOM))

	lastUpdate := app.AllianceKeeper.LastRewardClaimTime(ctx)
	require.Equal(t, ctx.BlockTime(), lastUpdate)

	asset, found := app.AllianceKeeper.GetAssetByDenom(ctx, ALLIANCE_TOKEN_DENOM)
	require.True(t, found)
	expectedRewardRate := sdk.MustNewDecFromStr("2").Mul(sdk.OneDec().Add(sdk.MustNewDecFromStr("0.5").Mul(sdk.NewDec(timePassed.Nanoseconds()).Quo(sdk.NewDec(31_557_000_000_000_000)))))
	require.Equal(t, expectedRewardRate, asset.RewardWeight)

	// At the next begin block, tokens will be distributed from the fee pool
	cons, _ := val1.GetConsAddr()
	app.DistrKeeper.AllocateTokens(ctx, 1, 1, cons, []abcitypes.VoteInfo{
		{
			Validator: abcitypes.Validator{
				Address: cons,
				Power:   1,
			},
			SignedLastBlock: true,
		},
	})

	rewards := app.DistrKeeper.GetValidatorOutstandingRewards(ctx, valAddr1).Rewards
	community := app.DistrKeeper.GetFeePool(ctx).CommunityPool
	// This is case, validator 1 has 0% commission
	commission := app.DistrKeeper.GetValidatorAccumulatedCommission(ctx, valAddr1).Commission
	require.Equal(t, sdk.DecCoins(nil), commission)
	// And rewards + community pool should add up to total coins claimed
	require.Equal(t,
		sdk.NewDecFromInt(coinsClaimed.AmountOf(ALLIANCE_TOKEN_DENOM)),
		rewards.AmountOf(ALLIANCE_TOKEN_DENOM).Add(community.AmountOf(ALLIANCE_TOKEN_DENOM)),
	)
}
