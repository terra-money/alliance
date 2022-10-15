package keeper_test

import (
	"alliance/x/alliance/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slices"
	"testing"
	"time"
)

var ALLIANCE_TOKEN_DENOM = "alliance"
var ALLIANCE_2_TOKEN_DENOM = "alliance2"

func TestDelegation(t *testing.T) {
	app, ctx := createTestContext(t)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.Params{
			RewardDelayTime: time.Duration(1000000),
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
	delegations := app.StakingKeeper.GetAllDelegations(ctx)
	require.Len(t, delegations, 1)

	// All the addresses needed
	delAddr, err := sdk.AccAddressFromBech32(delegations[0].DelegatorAddress)
	require.NoError(t, err)
	valAddr, err := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	require.NoError(t, err)
	val, _ := app.StakingKeeper.GetValidator(ctx, valAddr)
	moduleAddr := app.AccountKeeper.GetModuleAddress(types.ModuleName)

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
	}, allianceDelegation)

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
	}, allianceDelegation)

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
