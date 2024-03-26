package benchmark_test

import (
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/terra-money/alliance/x/alliance/tests/benchmark"

	sdkmath "cosmossdk.io/math"
	abcitypes "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/simulation"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/stretchr/testify/require"

	test_helpers "github.com/terra-money/alliance/app"
	"github.com/terra-money/alliance/x/alliance"
	"github.com/terra-money/alliance/x/alliance/types"
)

var (
	SEED               = int64(1)
	NumOfBlocks        = 200
	BlocktimeInSeconds = 5
	VoteRate           = 0.8
	NumOfValidators    = 160
	NumOfAssets        = 20
	NumOfDelegators    = 10

	OperationsPerBlock = 30
	DelegationRate     = 10
	RedelegationRate   = 2
	UndelegationRate   = 2
	RewardClaimRate    = 2
)

var createdDelegations = []types.Delegation{}

func TestRunBenchmarks(t *testing.T) {
	r := rand.New(rand.NewSource(SEED))
	app, ctx, assets, vals, dels := benchmark.SetupApp(t, r, NumOfAssets, NumOfValidators, NumOfDelegators)
	powerReduction := sdkmath.OneInt()
	operations := make(map[string]int)

	for b := 0; b < NumOfBlocks; b++ {
		t.Logf("Block: %d\n Time: %s", ctx.BlockHeight(), ctx.BlockTime())
		ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1).WithBlockTime(ctx.BlockTime().Add(time.Second * time.Duration(BlocktimeInSeconds)))
		totalVotingPower := int64(0)
		var voteInfo []abcitypes.VoteInfo
		for i := 0; i < NumOfValidators; i++ {
			valAddr := sdk.ValAddress(vals[i])
			val, err := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr)
			require.NoError(t, err)
			cons, _ := val.GetConsAddr()
			votingPower := val.ConsensusPower(powerReduction)
			totalVotingPower += votingPower
			voteInfo = append(voteInfo, abcitypes.VoteInfo{
				Validator: abcitypes.Validator{
					Address: cons,
					Power:   votingPower,
				},
			})
		}

		// Begin block
		app.DistrKeeper.AllocateTokens(ctx, totalVotingPower, voteInfo)

		// Delegator Actions
		operationFunc := benchmark.GenerateOperationSlots(DelegationRate, RedelegationRate, UndelegationRate, RewardClaimRate)
		for o := 0; o < OperationsPerBlock; o++ {
			switch operationFunc(r) {
			case 0:
				delegateOperation(ctx, app, r, assets, vals, dels)
				operations["delegate"]++
			case 1:
				redelegateOperation(ctx, app, r, vals)
				operations["redelegate"]++
			case 2:
				undelegateOperation(ctx, app, r)
				operations["undelegate"]++
			case 3:
				claimRewardsOperation(ctx, app, r)
				operations["claim"]++
			}
		}

		// Endblock
		assets := app.AllianceKeeper.GetAllAssets(ctx)
		app.AllianceKeeper.CompleteRedelegations(ctx)
		err := app.AllianceKeeper.CompleteUnbondings(ctx)
		if err != nil {
			panic(err)
		}
		_, err = app.AllianceKeeper.DeductAssetsHook(ctx, assets)
		if err != nil {
			panic(err)
		}
		err = app.AllianceKeeper.RewardWeightChangeHook(ctx, assets)
		if err != nil {
			panic(err)
		}
		err = app.AllianceKeeper.RebalanceHook(ctx, assets)
		if err != nil {
			panic(err)
		}
		res, stop := alliance.RunAllInvariants(ctx, app.AllianceKeeper)
		if stop {
			panic(res)
		}
	}
	t.Logf("%v\n", operations)

	state := app.AllianceKeeper.ExportGenesis(ctx)
	file, _ := os.Create("./benchmark_genesis.json")
	defer file.Close()
	file.Write(app.AppCodec().MustMarshalJSON(state)) //nolint:errcheck
}

func delegateOperation(ctx sdk.Context, app *test_helpers.App, r *rand.Rand, assets []types.AllianceAsset, vals []sdk.AccAddress, dels []sdk.AccAddress) {
	var asset types.AllianceAsset
	if len(assets) == 0 {
		return
	}
	if len(assets) == 1 {
		asset = assets[0]
	} else {
		asset = assets[r.Intn(len(assets)-1)]
	}
	valAddr := sdk.ValAddress(vals[r.Intn(len(vals)-1)])
	delAddr := dels[r.Intn(len(dels)-1)]

	amountToDelegate := simulation.RandomAmount(r, sdkmath.NewInt(1000_000_000))
	if amountToDelegate.IsZero() {
		return
	}
	coins := sdk.NewCoin(asset.Denom, amountToDelegate)

	app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(coins))                             //nolint:errcheck
	app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr, sdk.NewCoins(coins)) //nolint:errcheck

	val, _ := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr)
	app.AllianceKeeper.Delegate(ctx, delAddr, val, coins) //nolint:errcheck
	createdDelegations = append(createdDelegations, types.NewDelegation(ctx, delAddr, valAddr, asset.Denom, sdkmath.LegacyZeroDec(), []types.RewardHistory{}))
}

func redelegateOperation(ctx sdk.Context, app *test_helpers.App, r *rand.Rand, vals []sdk.AccAddress) {
	var delegation types.Delegation
	if len(createdDelegations) == 0 {
		return
	}
	if len(createdDelegations) == 1 {
		delegation = createdDelegations[0]
	} else {
		delegation = createdDelegations[r.Intn(len(createdDelegations)-1)]
	}

	delAddr, _ := sdk.AccAddressFromBech32(delegation.DelegatorAddress)
	srcValAddr, _ := sdk.ValAddressFromBech32(delegation.ValidatorAddress)
	srcValidator, _ := app.AllianceKeeper.GetAllianceValidator(ctx, srcValAddr)
	asset, _ := app.AllianceKeeper.GetAssetByDenom(ctx, delegation.Denom)

	if app.AllianceKeeper.HasRedelegation(ctx, delAddr, srcValAddr, asset.Denom) {
		return
	}

	dstValAddr := getRandomValAddress(r, vals, srcValAddr)
	dstValidator, _ := app.AllianceKeeper.GetAllianceValidator(ctx, dstValAddr)

	srcValidatorAddress, err := sdk.ValAddressFromBech32(delegation.ValidatorAddress)
	if err != nil {
		panic(err)
	}

	delegation, found := app.AllianceKeeper.GetDelegation(ctx, delAddr, srcValidatorAddress, asset.Denom)
	if !found {
		return
	}
	amountToRedelegate := simulation.RandomAmount(r, types.GetDelegationTokens(delegation, srcValidator, asset).Amount)
	if amountToRedelegate.LTE(sdkmath.OneInt()) {
		return
	}

	_, err = app.AllianceKeeper.Redelegate(ctx, delAddr, srcValidator, dstValidator, sdk.NewCoin(delegation.Denom, amountToRedelegate))
	if err != nil {
		panic(err)
	}
}

func getRandomValAddress(r *rand.Rand, vals []sdk.AccAddress, srcValAddr sdk.ValAddress) sdk.ValAddress {
	var dstValAddr sdk.ValAddress

	for {
		// Get a random destination validator address
		dstValAddr = sdk.ValAddress(vals[r.Intn(len(vals)-1)])

		// Break the loop if the destination validator address is different from the source validator address
		if !dstValAddr.Equals(srcValAddr) {
			break
		}
	}

	return dstValAddr
}

func undelegateOperation(ctx sdk.Context, app *test_helpers.App, r *rand.Rand) {
	if len(createdDelegations) == 0 {
		return
	}
	var delegation types.Delegation

	if len(createdDelegations) == 1 {
		delegation = createdDelegations[0]
	} else {
		delegation = createdDelegations[r.Intn(len(createdDelegations)-1)]
	}

	delAddr, _ := sdk.AccAddressFromBech32(delegation.DelegatorAddress)
	valAddr, _ := sdk.ValAddressFromBech32(delegation.ValidatorAddress)
	validator, _ := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr)
	asset, _ := app.AllianceKeeper.GetAssetByDenom(ctx, delegation.Denom)

	validatorAddress, err := sdk.ValAddressFromBech32(delegation.ValidatorAddress)
	if err != nil {
		panic(err)
	}

	delegation, found := app.AllianceKeeper.GetDelegation(ctx, delAddr, validatorAddress, asset.Denom)
	if !found {
		return
	}
	amountToUndelegate := simulation.RandomAmount(r, types.GetDelegationTokens(delegation, validator, asset).Amount)
	if amountToUndelegate.IsZero() {
		return
	}

	_, err = app.AllianceKeeper.Undelegate(ctx, delAddr, validator, sdk.NewCoin(asset.Denom, amountToUndelegate))
	if err != nil {
		panic(err)
	}
}

func claimRewardsOperation(ctx sdk.Context, app *test_helpers.App, r *rand.Rand) {
	var delegation types.Delegation
	if len(createdDelegations) == 0 {
		return
	}
	if len(createdDelegations) == 1 {
		delegation = createdDelegations[0]
	} else {
		delegation = createdDelegations[r.Intn(len(createdDelegations)-1)]
	}
	delAddr, _ := sdk.AccAddressFromBech32(delegation.DelegatorAddress)
	valAddr, _ := sdk.ValAddressFromBech32(delegation.ValidatorAddress)
	validator, _ := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr)

	validatorAddress, err := sdk.ValAddressFromBech32(delegation.ValidatorAddress)
	if err != nil {
		panic(err)
	}

	delegation, found := app.AllianceKeeper.GetDelegation(ctx, delAddr, validatorAddress, delegation.Denom)
	if !found {
		return
	}

	_, err = app.AllianceKeeper.ClaimDelegationRewards(ctx, delAddr, validator, delegation.Denom)
	if err != nil {
		panic(err)
	}
}
