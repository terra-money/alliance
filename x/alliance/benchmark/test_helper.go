package benchmark

import (
	"fmt"
	"math/rand"
	"testing"

	test_helpers "alliance/app"
	"alliance/x/alliance/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/staking/teststaking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"time"
)

func SetupApp(t *testing.T, r *rand.Rand, numAssets int, numValidators int, numDelegators int) (app *test_helpers.App, ctx sdk.Context, assets []types.AllianceAsset, valAddrs []sdk.AccAddress, delAddrs []sdk.AccAddress) {
	app = test_helpers.Setup(t, false)
	ctx = app.BaseApp.NewContext(false, tmproto.Header{})
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime)
	for i := 0; i < numAssets; i += 1 {
		rewardWeight := simulation.RandomDecAmount(r, sdk.NewDec(1))
		takeRate := simulation.RandomDecAmount(r, sdk.NewDec(1))
		assets = append(assets, types.NewAllianceAsset(fmt.Sprintf("ASSET%d", i), rewardWeight, takeRate, startTime))
	}
	params := types.NewParams()
	params.RewardClaimInterval = time.Second * 5
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: params,
		Assets: assets,
	})

	// Accounts
	valAddrs = test_helpers.AddTestAddrsIncremental(app, ctx, numValidators, sdk.NewCoins())
	pks := test_helpers.CreateTestPubKeys(numValidators)

	for i := 0; i < numValidators; i += 1 {
		valAddr := sdk.ValAddress(valAddrs[i])
		_val := teststaking.NewValidator(t, valAddr, pks[i])
		_val.Commission = stakingtypes.Commission{
			CommissionRates: stakingtypes.CommissionRates{
				Rate:          sdk.NewDec(0),
				MaxRate:       sdk.NewDec(0),
				MaxChangeRate: sdk.NewDec(0),
			},
			UpdateTime: time.Now(),
		}
		_val.Status = stakingtypes.Bonded
		test_helpers.RegisterNewValidator(t, app, ctx, _val)
	}

	delAddrs = test_helpers.AddTestAddrsIncremental(app, ctx, numDelegators, sdk.NewCoins())
	return
}

func GenerateOperationSlots(operations ...int) func(r *rand.Rand) int {
	var slots []int
	for i := 0; i < len(operations); i += 1 {
		for o := 0; o < operations[i]; o += 1 {
			slots = append(slots, i)
		}
	}
	return func(r *rand.Rand) int {
		return slots[r.Intn(len(slots)-1)]
	}
}
