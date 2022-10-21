package keeper_test

import (
	"alliance/x/alliance/types"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/stretchr/testify/require"
)

func TestQueryRewards(t *testing.T) {
	// GIVEN: THE BLOCKCHAIN WITH ALLIANCES ON GENESIS
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			{
				Denom:        ALLIANCE_TOKEN_DENOM,
				RewardWeight: sdk.NewDec(1),
				TakeRate:     sdk.NewDec(1),
				TotalTokens:  sdk.ZeroInt(),
			},
		},
	})
	delegations := app.StakingKeeper.GetAllDelegations(ctx)
	require.Len(t, delegations, 1)
	delAddr, _ := sdk.AccAddressFromBech32(delegations[0].DelegatorAddress)
	valAddr, _ := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	val, _ := app.StakingKeeper.GetValidator(ctx, valAddr)
	app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin("alliance", sdk.NewInt(2_000_000_000_000))))
	app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr, sdk.NewCoins(sdk.NewCoin("alliance", sdk.NewInt(2_000_000_000_000))))

	// WHEN: DELEGATING AND QUERYING ...
	delegationTxRes, txErr := app.AllianceKeeper.Delegate(ctx, delAddr, val, sdk.NewCoin("alliance", sdk.NewInt(1_900_000_000_000)))
	queryDelegation, queryErr := app.AllianceKeeper.AllianceDelegationRewards(ctx, &types.QueryAllianceDelegationRewardsRequest{
		DelegatorAddr: delAddr.String(),
		ValidatorAddr: valAddr.String(),
		Denom:         "alliance",
	})

	// THEN: VALIDATE THAT NO ERRORS HAVE BEEN PRODUCED AND BOTH OUTPUTS ARE AS WE EXPECT
	require.Nil(t, txErr)
	require.Nil(t, queryErr)
	require.Equal(t, &types.QueryAllianceDelegationRewardsResponse{}, queryDelegation)
	require.Equal(t, &types.Delegation{
		DelegatorAddress: delAddr.String(),
		ValidatorAddress: val.OperatorAddress,
		Denom:            "alliance",
		Shares:           sdk.NewDec(1_900_000_000_000),
		RewardIndices:    []types.RewardIndex{},
	}, delegationTxRes)
}
