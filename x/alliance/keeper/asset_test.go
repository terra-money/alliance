package keeper_test

import (
	test_helpers "alliance/app"
	"alliance/x/alliance/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/teststaking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestUpdateRewardRates(t *testing.T) {
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx.WithBlockTime(startTime)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			types.NewAsset(ALLIANCE_TOKEN_DENOM, sdk.NewDec(2), sdk.ZeroDec()),
			types.NewAsset(ALLIANCE_2_TOKEN_DENOM, sdk.NewDec(10), sdk.ZeroDec()),
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
	powerReduction := app.StakingKeeper.PowerReduction(ctx)

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
	val1.Status = stakingtypes.Bonded
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
	//user2 := addrs[3]

	// Start by delegating
	_, err = app.AllianceKeeper.Delegate(ctx, user1, val1, sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(1000_000)))
	require.NoError(t, err)

	// Expecting voting power for the alliance module
	val, found := app.StakingKeeper.GetValidator(ctx, valAddr1)
	require.True(t, found)
	require.Equal(t, int64(2), val.ConsensusPower(powerReduction))

	err = app.AllianceKeeper.UpdateAllianceAsset(ctx, types.AllianceAsset{
		Denom:        ALLIANCE_TOKEN_DENOM,
		RewardWeight: sdk.NewDec(20),
		TakeRate:     sdk.NewDec(0),
	})
	require.NoError(t, err)

	// Expecting voting power to increase
	val, found = app.StakingKeeper.GetValidator(ctx, valAddr1)
	require.True(t, found)
	require.Equal(t, int64(20), val.ConsensusPower(powerReduction))

	err = app.AllianceKeeper.UpdateAllianceAsset(ctx, types.AllianceAsset{
		Denom:        ALLIANCE_TOKEN_DENOM,
		RewardWeight: sdk.NewDec(1),
		TakeRate:     sdk.NewDec(0),
	})
	require.NoError(t, err)

	// Expecting voting power to decrease
	val, found = app.StakingKeeper.GetValidator(ctx, valAddr1)
	require.True(t, found)
	require.Equal(t, int64(1), val.ConsensusPower(powerReduction))
}
