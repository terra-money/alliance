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
			types.NewAllianceAsset(ALLIANCE_TOKEN_DENOM, sdk.NewDec(2), sdk.ZeroDec()),
			types.NewAllianceAsset(ALLIANCE_2_TOKEN_DENOM, sdk.NewDec(10), sdk.ZeroDec()),
		},
	})

	// Set tax and rewards to be zero for easier calculation
	distParams := app.DistrKeeper.GetParams(ctx)
	distParams.CommunityTax = sdk.ZeroDec()
	distParams.BaseProposerReward = sdk.ZeroDec()
	distParams.BonusProposerReward = sdk.ZeroDec()
	app.DistrKeeper.SetParams(ctx, distParams)

	// Accounts
	addrs := test_helpers.AddTestAddrsIncremental(app, ctx, 4, sdk.NewCoins(
		sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(1000_000)),
		sdk.NewCoin(ALLIANCE_2_TOKEN_DENOM, sdk.NewInt(1000_000)),
	))
	pks := test_helpers.CreateTestPubKeys(2)
	powerReduction := app.StakingKeeper.PowerReduction(ctx)

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
	_val1.Status = stakingtypes.Bonded
	test_helpers.RegisterNewValidator(t, app, ctx, _val1)
	val1, err := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr1)
	require.NoError(t, err)

	user1 := addrs[2]

	// Start by delegating
	_, err = app.AllianceKeeper.Delegate(ctx, user1, val1, sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(1000_000)))
	require.NoError(t, err)

	err = app.AllianceKeeper.RebalanceInternalStakeWeights(ctx)
	require.NoError(t, err)
	require.Equal(t, sdk.NewInt(3_000_000), app.StakingKeeper.TotalBondedTokens(ctx))

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

	err = app.AllianceKeeper.RebalanceInternalStakeWeights(ctx)
	require.NoError(t, err)
	require.Equal(t, sdk.NewInt(21_000_000), app.StakingKeeper.TotalBondedTokens(ctx))

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

	err = app.AllianceKeeper.RebalanceInternalStakeWeights(ctx)
	require.NoError(t, err)
	require.Equal(t, sdk.NewInt(2_000_000), app.StakingKeeper.TotalBondedTokens(ctx))

	// Expecting voting power to decrease
	val, found = app.StakingKeeper.GetValidator(ctx, valAddr1)
	require.True(t, found)
	require.Equal(t, int64(1), val.ConsensusPower(powerReduction))
}
