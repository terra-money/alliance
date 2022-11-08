package benchmark_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/simulation"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/stretchr/testify/require"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	test_helpers "github.com/terra-money/alliance/app"
	"github.com/terra-money/alliance/x/alliance/benchmark"
	"github.com/terra-money/alliance/x/alliance/types"
	"math/rand"
	"testing"
	"time"
)

var (
	SEED              = 1
	NUM_OF_BLOCKS     = 250
	BLOCKTIME_IN_S    = 5
	VOTE_RATE         = 0.9
	NUM_OF_VALIDATORS = 300
	NUM_OF_ASSETS     = 1280
	NUM_OF_DELEGATORS = 100

	OPERATIONS_PER_BLOCK = 50
	DELEGATION_RATE      = 10
	REDELEGATION_RATE    = 1
	UNDELEGATION_RATE    = 2
	REWARD_CLAIM_RATE    = 2
)

var createdDelegations = []types.Delegation{}

func TestRunBenchmarks(t *testing.T) {
	r := rand.New(rand.NewSource(1))
	app, ctx, assets, vals, dels := benchmark.SetupApp(t, r, NUM_OF_ASSETS, NUM_OF_VALIDATORS, NUM_OF_DELEGATORS)
	powerReduction := app.StakingKeeper.PowerReduction(ctx)
	operations := make(map[string]int)

	for b := 0; b < NUM_OF_BLOCKS; b += 1 {
		t.Logf("Block: %d\n Time: %s", ctx.BlockHeight(), ctx.BlockTime())
		ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1).WithBlockTime(ctx.BlockTime().Add(time.Second * time.Duration(BLOCKTIME_IN_S)))
		totalVotingPower := int64(0)
		var voteInfo []abcitypes.VoteInfo
		for i := 0; i < NUM_OF_VALIDATORS; i += 1 {
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
				SignedLastBlock: r.Float64() < VOTE_RATE,
			})
		}

		idx := simulation.RandIntBetween(r, 0, len(vals)-1)
		proposerAddr := sdk.ValAddress(vals[idx])
		proposer, err := app.AllianceKeeper.GetAllianceValidator(ctx, proposerAddr)
		require.NoError(t, err)
		proposerCons, _ := proposer.GetConsAddr()

		// Begin block
		app.DistrKeeper.AllocateTokens(ctx, totalVotingPower, totalVotingPower, proposerCons, voteInfo)

		// Delegator Actions
		operationFunc := benchmark.GenerateOperationSlots(DELEGATION_RATE, REDELEGATION_RATE, UNDELEGATION_RATE, REWARD_CLAIM_RATE)
		for o := 0; o < OPERATIONS_PER_BLOCK; o += 1 {
			switch operationFunc(r) {
			case 0:
				delegateOperation(ctx, app, r, assets, vals, dels)
				operations["delegate"] += 1
				break
			case 1:
				redelegateOperation(ctx, app, r, assets, vals, dels)
				operations["redelegate"] += 1
				break
			case 2:
				undelegateOperation(ctx, app, r)
				operations["undelegate"] += 1
				break
			case 3:
				claimRewardsOperation(ctx, app, r)
				operations["claim"] += 1
				break
			}
		}

		// Endblock
		app.AllianceKeeper.CompleteRedelegations(ctx)
		err = app.AllianceKeeper.CompleteUndelegations(ctx)
		if err != nil {
			panic(err)
		}
		_, err = app.AllianceKeeper.DeductAssetsHook(ctx)
		if err != nil {
			panic(err)
		}
		err = app.AllianceKeeper.RebalanceHook(ctx)
		if err != nil {
			panic(err)
		}
	}
	t.Logf("%v\n", operations)
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

	amountToDelegate := simulation.RandomAmount(r, sdk.NewInt(1000_000_000))
	if amountToDelegate.IsZero() {
		return
	}
	coins := sdk.NewCoin(asset.Denom, amountToDelegate)

	app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(coins))
	app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr, sdk.NewCoins(coins))

	val, _ := app.AllianceKeeper.GetAllianceValidator(ctx, valAddr)
	app.AllianceKeeper.Delegate(ctx, delAddr, val, coins)
	createdDelegations = append(createdDelegations, types.NewDelegation(ctx, delAddr, valAddr, asset.Denom, sdk.ZeroDec(), []types.RewardHistory{}))
}

func redelegateOperation(ctx sdk.Context, app *test_helpers.App, r *rand.Rand, assets []types.AllianceAsset, vals []sdk.AccAddress, dels []sdk.AccAddress) {
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

	if app.AllianceKeeper.IterateRedelegations(ctx, delAddr, srcValAddr, asset.Denom).Valid() {
		return
	}

	dstValAddr := sdk.ValAddress(vals[r.Intn(len(vals)-1)])
	dstValidator, _ := app.AllianceKeeper.GetAllianceValidator(ctx, dstValAddr)

	amountToRedelegate := simulation.RandomAmount(r, types.GetDelegationTokens(delegation, srcValidator, asset).Amount)
	if amountToRedelegate.IsZero() {
		return
	}
	app.AllianceKeeper.Redelegate(ctx, delAddr, srcValidator, dstValidator, sdk.NewCoin(delegation.Denom, amountToRedelegate))
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

	amountToUndelegate := simulation.RandomAmount(r, types.GetDelegationTokens(delegation, validator, asset).Amount)
	if amountToUndelegate.IsZero() {
		return
	}
	app.AllianceKeeper.Undelegate(ctx, delAddr, validator, sdk.NewCoin(asset.Denom, amountToUndelegate))
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

	app.AllianceKeeper.ClaimDelegationRewards(ctx, delAddr, validator, delegation.Denom)
}
