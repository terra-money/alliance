package keeper_test

import (
	test_helpers "alliance/app"
	"alliance/x/alliance/types"
	"testing"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"
	abcitypes "github.com/tendermint/tendermint/abci/types"
)

var ULUNA_ALLIANCE = "uluna"

func TestClaimQueryReward(t *testing.T) {
	// GIVEN: THE BLOCKCHAIN WITH ACCOUNTS
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
				Denom:                ULUNA_ALLIANCE,
				RewardWeight:         sdk.NewDec(2),
				TakeRate:             sdk.MustNewDecFromStr("0.5"),
				TotalTokens:          sdk.ZeroInt(),
				TotalValidatorShares: sdk.NewDec(0),
				TotalStakeTokens:     sdk.ZeroInt(),
			},
		},
	})
	feeCollectorAddr := app.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
	delegations := app.StakingKeeper.GetAllDelegations(ctx)
	valAddr, _ := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	val1, _ := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr)
	delAddr := test_helpers.AddTestAddrsIncremental(app, ctx, 1, sdk.NewCoins(sdk.NewCoin(ULUNA_ALLIANCE, sdk.NewInt(1000_000_000))))[0]

	// WHEN: DELEGATING ...
	delRes, delErr := app.AllianceKeeper.Delegate(ctx, delAddr, val1, sdk.NewCoin(ULUNA_ALLIANCE, sdk.NewInt(1000_000_000)))
	require.Nil(t, delErr)
	require.Equal(t, delRes, &types.Delegation{
		DelegatorAddress:      delAddr.String(),
		ValidatorAddress:      valAddr.String(),
		Denom:                 "uluna",
		Shares:                sdk.NewDec(1000_000_000),
		RewardHistory:         []types.RewardHistory{},
		LastRewardClaimHeight: uint64(ctx.BlockHeight()),
	})
	err := app.AllianceKeeper.RebalanceBondTokenWeights(ctx)
	require.NoError(t, err)

	// ...and advance block...
	timePassed := time.Minute*5 + time.Second
	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(timePassed))
	ctx = ctx.WithBlockHeight(2)
	app.AllianceKeeper.DeductAssetsHook(ctx)
	app.BankKeeper.GetAllBalances(ctx, feeCollectorAddr)
	sdk.MustNewDecFromStr("0.5").Mul(sdk.NewDec(timePassed.Nanoseconds()).Quo(sdk.NewDec(31_557_000_000_000_000))).MulInt(sdk.NewInt(1000_000_000))
	app.AllianceKeeper.LastRewardClaimTime(ctx)
	app.AllianceKeeper.GetAssetByDenom(ctx, ULUNA_ALLIANCE)
	sdk.MustNewDecFromStr("2").Mul(sdk.OneDec().Add(sdk.MustNewDecFromStr("0.5").Mul(sdk.NewDec(timePassed.Nanoseconds()).Quo(sdk.NewDec(31_557_000_000_000_000)))))

	// ... at the next begin block, tokens will be distributed from the fee pool...
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

	// THEN: Query the delegation rewards ...
	queryDelegation, queryErr := app.AllianceKeeper.AllianceDelegationRewards(ctx, &types.QueryAllianceDelegationRewardsRequest{
		DelegatorAddr: delAddr.String(),
		ValidatorAddr: valAddr.String(),
		Denom:         ULUNA_ALLIANCE,
	})

	// ... validate that no error has been produced.
	require.Nil(t, queryErr)
	require.Equal(t, &types.QueryAllianceDelegationRewardsResponse{
		Rewards: []sdk.Coin{
			{
				Denom:  ULUNA_ALLIANCE,
				Amount: math.NewInt(3115),
			},
		},
	}, queryDelegation)
}
