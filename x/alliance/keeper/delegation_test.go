package keeper_test

import (
	test_helpers "alliance/app"
	"alliance/x/alliance/types"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/staking/teststaking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slices"
)

var ALLIANCE_TOKEN_DENOM = "alliance"
var ALLIANCE_2_TOKEN_DENOM = "alliance2"

func TestDelegation(t *testing.T) {
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
	delegations := app.StakingKeeper.GetAllDelegations(ctx)
	require.Len(t, delegations, 1)

	// All the addresses needed
	delAddr, err := sdk.AccAddressFromBech32(delegations[0].DelegatorAddress)
	require.NoError(t, err)
	valAddr, err := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	require.NoError(t, err)
	val, _ := app.StakingKeeper.GetValidator(ctx, valAddr)
	moduleAddr := app.AccountKeeper.GetModuleAddress(types.ModuleName)

	// Mint alliance tokens
	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(2000_000))))
	require.NoError(t, err)
	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(ALLIANCE_2_TOKEN_DENOM, sdk.NewInt(2000_000))))
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr, sdk.NewCoins(sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(2000_000))))
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr, sdk.NewCoins(sdk.NewCoin(ALLIANCE_2_TOKEN_DENOM, sdk.NewInt(2000_000))))
	require.NoError(t, err)

	// Delegate
	_, err = app.AllianceKeeper.Delegate(ctx, delAddr, val, sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(1000_000)))
	require.NoError(t, err)

	// Check delegation in staking module
	delegations = app.StakingKeeper.GetAllDelegations(ctx)
	require.Len(t, delegations, 2)
	i := slices.IndexFunc(delegations, func(d stakingtypes.Delegation) bool {
		return d.DelegatorAddress == moduleAddr.String()
	})
	newDelegation := delegations[i]
	require.Equal(t, stakingtypes.Delegation{
		DelegatorAddress: moduleAddr.String(),
		ValidatorAddress: valAddr.String(),
		Shares:           sdk.NewDec(2),
	}, newDelegation)

	// Check delegation in alliance module
	allianceDelegation, found := app.AllianceKeeper.GetDelegation(ctx, delAddr, val, ALLIANCE_TOKEN_DENOM)
	require.True(t, found)
	require.Equal(t, types.Delegation{
		DelegatorAddress: delAddr.String(),
		ValidatorAddress: val.OperatorAddress,
		Denom:            ALLIANCE_TOKEN_DENOM,
		Shares:           sdk.NewDec(1000_000),
		RewardIndices:    types.RewardIndices(nil),
	}, allianceDelegation)

	// Check asset
	asset, _ := app.AllianceKeeper.GetAssetByDenom(ctx, ALLIANCE_TOKEN_DENOM)
	require.Equal(t, types.AllianceAsset{
		Denom:                ALLIANCE_TOKEN_DENOM,
		RewardWeight:         sdk.NewDec(2),
		TakeRate:             sdk.NewDec(0),
		TotalTokens:          sdk.NewInt(1000_000),
		TotalValidatorShares: sdk.NewDec(1000_000),
	}, asset)

	// Delegate with same denom again
	_, err = app.AllianceKeeper.Delegate(ctx, delAddr, val, sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(1000_000)))
	require.NoError(t, err)
	// Check delegation in alliance module
	allianceDelegation, found = app.AllianceKeeper.GetDelegation(ctx, delAddr, val, ALLIANCE_TOKEN_DENOM)
	require.True(t, found)
	require.Equal(t, types.Delegation{
		DelegatorAddress: delAddr.String(),
		ValidatorAddress: val.OperatorAddress,
		Denom:            ALLIANCE_TOKEN_DENOM,
		Shares:           sdk.NewDec(2000_000),
		RewardIndices:    types.RewardIndices(nil),
	}, allianceDelegation)

	// Check asset again
	asset, _ = app.AllianceKeeper.GetAssetByDenom(ctx, ALLIANCE_TOKEN_DENOM)
	require.Equal(t, types.AllianceAsset{
		Denom:                ALLIANCE_TOKEN_DENOM,
		RewardWeight:         sdk.NewDec(2),
		TakeRate:             sdk.NewDec(0),
		TotalTokens:          sdk.NewInt(2000_000),
		TotalValidatorShares: sdk.NewDec(2000_000),
	}, asset)

	// Delegate with another denom
	_, err = app.AllianceKeeper.Delegate(ctx, delAddr, val, sdk.NewCoin(ALLIANCE_2_TOKEN_DENOM, sdk.NewInt(1000_000)))
	require.NoError(t, err)

	// Check delegation in staking module
	delegations = app.StakingKeeper.GetAllDelegations(ctx)
	require.Len(t, delegations, 2)
	i = slices.IndexFunc(delegations, func(d stakingtypes.Delegation) bool {
		return d.DelegatorAddress == moduleAddr.String()
	})
	newDelegation = delegations[i]
	require.Equal(t, stakingtypes.Delegation{
		DelegatorAddress: moduleAddr.String(),
		ValidatorAddress: valAddr.String(),
		Shares:           sdk.NewDec(14),
	}, newDelegation)
}

// TODO: test using unsupported denoms

func TestRedelegation(t *testing.T) {
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

	// Get all the addresses needed for the test
	moduleAddr := app.AccountKeeper.GetModuleAddress(types.ModuleName)
	delegations := app.StakingKeeper.GetAllDelegations(ctx)
	valAddr1, err := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	require.NoError(t, err)
	val1, found := app.StakingKeeper.GetValidator(ctx, valAddr1)
	require.True(t, found)
	addrs := test_helpers.AddTestAddrsIncremental(app, ctx, 3, sdk.NewCoins(
		sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(1000_000)),
		sdk.NewCoin(ALLIANCE_2_TOKEN_DENOM, sdk.NewInt(1000_000)),
	))
	valAddr2 := sdk.ValAddress(addrs[0])
	val2 := teststaking.NewValidator(t, valAddr2, test_helpers.CreateTestPubKeys(1)[0])
	test_helpers.RegisterNewValidator(t, app, ctx, val2)
	delAddr1 := addrs[1]
	delAddr2 := addrs[2]

	// Mint alliance tokens
	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(2000_000))))
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr1, sdk.NewCoins(sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(1000_000))))
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr2, sdk.NewCoins(sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(1000_000))))
	require.NoError(t, err)

	// First delegate to validator 1
	_, err = app.AllianceKeeper.Delegate(ctx, delAddr1, val1, sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(500_000)))
	require.NoError(t, err)

	// Then redelegate to validator2
	_, err = app.AllianceKeeper.Redelegate(ctx, delAddr1, val1, val2, sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(500_000)))
	require.NoError(t, err)

	delegations = app.StakingKeeper.GetAllDelegations(ctx)
	for _, d := range delegations {
		if d.DelegatorAddress == moduleAddr.String() {
			require.Equal(t, stakingtypes.Delegation{
				DelegatorAddress: moduleAddr.String(),
				ValidatorAddress: val2.OperatorAddress,
				Shares:           sdk.NewDec(1_000_000),
			}, d)
		}
	}

	// Check if there is a re-delegation event stored
	iter := app.AllianceKeeper.IterateRedelegationsByDelegator(ctx, delAddr1)
	require.True(t, iter.Valid())
	for ; iter.Valid(); iter.Next() {
		var redelegation types.Redelegation
		app.AppCodec().MustUnmarshal(iter.Value(), &redelegation)
		require.Equal(t, types.Redelegation{
			DelegatorAddress:    delAddr1.String(),
			SrcValidatorAddress: valAddr1.String(),
			DstValidatorAddress: valAddr2.String(),
			Balance:             sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(500_000)),
		}, redelegation)
	}

	// Check if the delegation objects are correct
	_, found = app.AllianceKeeper.GetDelegation(ctx, delAddr1, val1, ALLIANCE_TOKEN_DENOM)
	require.False(t, found)
	dstDelegation, found := app.AllianceKeeper.GetDelegation(ctx, delAddr1, val2, ALLIANCE_TOKEN_DENOM)
	require.Equal(t, types.Delegation{
		DelegatorAddress: delAddr1.String(),
		ValidatorAddress: val2.GetOperator().String(),
		Denom:            ALLIANCE_TOKEN_DENOM,
		Shares:           sdk.NewDec(500_000),
		RewardIndices:    types.RewardIndices(nil),
	}, dstDelegation)
	require.True(t, found)

	// Should fail when re-delegating back to validator1
	// Same user who re-delegated to from 1 -> 2 cannot re-re-delegate from 2 -> X
	_, err = app.AllianceKeeper.Redelegate(ctx, delAddr1, val2, val1, sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(500_000)))
	require.Error(t, err)

	// Another user tries to re-delegate without having an initial delegation but fails
	_, err = app.AllianceKeeper.Redelegate(ctx, delAddr2, val2, val1, sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(500_000)))
	require.Error(t, err)

	// User then delegates to validator2
	_, err = app.AllianceKeeper.Delegate(ctx, delAddr2, val2, sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(500_000)))
	require.NoError(t, err)

	// Then redelegate to validator1 with more than what was delegated but fails
	_, err = app.AllianceKeeper.Redelegate(ctx, delAddr2, val2, val1, sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(1000_000)))
	require.Error(t, err)

	// Then redelegate to validator1 correctly
	_, err = app.AllianceKeeper.Redelegate(ctx, delAddr2, val2, val1, sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(500_000)))
	// Should pass since we removed the re-delegate attempt on x/staking that prevents this
	require.NoError(t, err)

	// Immediately calling complete re-delegation should do nothing
	deleted := app.AllianceKeeper.CompleteRedelegations(ctx)
	require.Equal(t, 0, deleted)

	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(app.StakingKeeper.UnbondingTime(ctx)).Add(time.Minute))
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

func TestUndelegation(t *testing.T) {
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
	delegations := app.StakingKeeper.GetAllDelegations(ctx)
	require.Len(t, delegations, 1)

	// All the addresses needed
	delAddr, err := sdk.AccAddressFromBech32(delegations[0].DelegatorAddress)
	require.NoError(t, err)
	valAddr, err := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	require.NoError(t, err)
	val, _ := app.StakingKeeper.GetValidator(ctx, valAddr)
	moduleAddr := app.AccountKeeper.GetModuleAddress(types.ModuleName)

	// Mint alliance tokens
	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(2000_000))))
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr, sdk.NewCoins(sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(2000_000))))
	require.NoError(t, err)

	// Undelegating without a delegation will fail
	err = app.AllianceKeeper.Undelegate(ctx, delAddr, val, sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(1000_000)))
	require.Error(t, err)

	// Delegate to a validator
	_, err = app.AllianceKeeper.Delegate(ctx, delAddr, val, sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(1000_000)))
	require.NoError(t, err)

	// Check that balance dropped
	coin := app.BankKeeper.GetBalance(ctx, delAddr, ALLIANCE_TOKEN_DENOM)
	require.Equal(t, sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(1000_000)), coin)

	// Check that staked balance increased
	d, _ := app.StakingKeeper.GetDelegation(ctx, moduleAddr, valAddr)
	require.Equal(t, stakingtypes.Delegation{
		DelegatorAddress: moduleAddr.String(),
		ValidatorAddress: valAddr.String(),
		Shares:           sdk.NewDec(2),
	}, d)

	// Immediately undelegate from the validator
	err = app.AllianceKeeper.Undelegate(ctx, delAddr, val, sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(500_000)))
	require.NoError(t, err)

	// Check that staked balance decreased
	d, _ = app.StakingKeeper.GetDelegation(ctx, moduleAddr, valAddr)
	require.Equal(t, stakingtypes.Delegation{
		DelegatorAddress: moduleAddr.String(),
		ValidatorAddress: valAddr.String(),
		Shares:           sdk.NewDec(1),
	}, d)

	// Immediately try to complete undelegation
	processed := app.AllianceKeeper.CompleteUndelegations(ctx)
	require.Equal(t, 0, processed)

	// Check that balance stayed the same
	coin = app.BankKeeper.GetBalance(ctx, delAddr, ALLIANCE_TOKEN_DENOM)
	require.Equal(t, sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(1000_000)), coin)

	// Advance time to after unbonding period
	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(app.StakingKeeper.UnbondingTime(ctx)).Add(time.Minute))

	processed = app.AllianceKeeper.CompleteUndelegations(ctx)
	require.Equal(t, 1, processed)

	// Check that balance increased
	coin = app.BankKeeper.GetBalance(ctx, delAddr, ALLIANCE_TOKEN_DENOM)
	require.Equal(t, sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(1500_000)), coin)

	// Completing again should not process anymore undelegations
	processed = app.AllianceKeeper.CompleteUndelegations(ctx)
	require.Equal(t, 0, processed)
}

func TestUndelegateAfterClaimingTakeRate(t *testing.T) {
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
				TakeRate:     sdk.MustNewDecFromStr("0.5"),
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
		sdk.NewCoin(ALLIANCE_TOKEN_DENOM, sdk.NewInt(1000_000_000)),
		sdk.NewCoin(ALLIANCE_2_TOKEN_DENOM, sdk.NewInt(1000_000_000)),
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

	// Delegate token with non-zero take_rate
	app.AllianceKeeper.Delegate(ctx, user1, val1, sdk.NewCoin(ALLIANCE_2_TOKEN_DENOM, sdk.NewInt(1000_000_000)))
	app.AllianceKeeper.Delegate(ctx, user2, val2, sdk.NewCoin(ALLIANCE_2_TOKEN_DENOM, sdk.NewInt(1000_000_000)))

	ctx = ctx.WithBlockTime(startTime.Add(time.Minute * 6)).WithBlockHeight(2)
	coins, err := app.AllianceKeeper.ClaimAssetsWithTakeRateRateLimited(ctx)
	require.NoError(t, err)
	require.False(t, coins.IsZero())

	res, err := app.AllianceKeeper.AllianceDelegation(ctx, &types.QueryAllianceDelegationRequest{
		DelegatorAddr: user1.String(),
		ValidatorAddr: val1.GetOperator().String(),
		Denom:         ALLIANCE_2_TOKEN_DENOM,
		Pagination:    nil,
	})
	require.NoError(t, err)
	del := res.GetDelegation()
	require.True(t, del.GetBalance().Amount.LT(sdk.NewInt(1000_000_000)), "%s should be less than %s", del.GetBalance().Amount, sdk.NewInt(1000_000_000))
	// Undelegate token with initial amount should fail
	err = app.AllianceKeeper.Undelegate(ctx, user1, val1, sdk.NewCoin(ALLIANCE_2_TOKEN_DENOM, sdk.NewInt(1000_000_000)))
	require.Error(t, err)

	// Undelegate token with current amount should pass
	err = app.AllianceKeeper.Undelegate(ctx, user1, val1, sdk.NewCoin(ALLIANCE_2_TOKEN_DENOM, del.Balance.Amount))
	require.NoError(t, err)

	// User should have everything withdrawn
	_, found := app.AllianceKeeper.GetDelegation(ctx, user1, val1, ALLIANCE_2_TOKEN_DENOM)
	require.False(t, found)
}
