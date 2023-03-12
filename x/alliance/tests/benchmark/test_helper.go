package benchmark

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/simulation"
	teststaking "github.com/cosmos/cosmos-sdk/x/staking/testutil"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	test_helpers "github.com/terra-money/alliance/app"
	"github.com/terra-money/alliance/x/alliance/types"
)

func SetupApp(t *testing.T, r *rand.Rand, numAssets int, numValidators int, numDelegators int) (app *test_helpers.App, ctx sdk.Context, assets []types.AllianceAsset, valAddrs []sdk.AccAddress, delAddrs []sdk.AccAddress) {
	app = test_helpers.Setup(t, false)
	ctx = app.BaseApp.NewContext(false, tmproto.Header{})
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime)
	for i := 0; i < numAssets; i++ {
		rewardWeight := simulation.RandomDecAmount(r, sdk.NewDec(1))
		takeRate := simulation.RandomDecAmount(r, sdk.MustNewDecFromStr("0.0001"))
		asset := types.NewAllianceAsset(fmt.Sprintf("ASSET%d", i), rewardWeight, sdk.ZeroDec(), sdk.NewDec(5), takeRate, startTime)
		asset.RewardChangeRate = sdk.OneDec().Sub(simulation.RandomDecAmount(r, sdk.MustNewDecFromStr("0.00001")))
		asset.RewardChangeInterval = time.Minute * 5
		assets = append(assets, asset)
	}
	params := types.NewParams()
	params.TakeRateClaimInterval = time.Minute * 5
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: params,
		Assets: assets,
	})

	// Accounts
	valAddrs = test_helpers.AddTestAddrsIncremental(app, ctx, numValidators, sdk.NewCoins())
	pks := test_helpers.CreateTestPubKeys(numValidators)

	for i := 0; i < numValidators; i++ {
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
	return app, ctx, assets, valAddrs, delAddrs
}

func GenerateOperationSlots(operations ...int) func(r *rand.Rand) int {
	var slots []int
	for i := 0; i < len(operations); i++ {
		for o := 0; o < operations[i]; o++ {
			slots = append(slots, i)
		}
	}
	return func(r *rand.Rand) int {
		return slots[r.Intn(len(slots)-1)]
	}
}
