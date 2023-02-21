package keeper

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	accountkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	"github.com/cosmos/cosmos-sdk/x/gov/types"
	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	custombankkeeper "github.com/terra-money/alliance/custom/bank/keeper"
	alliancekeeper "github.com/terra-money/alliance/x/alliance/keeper"
	alliancetypes "github.com/terra-money/alliance/x/alliance/types"
)

type Keeper struct {
	govkeeper.Keeper

	key  storetypes.StoreKey
	ak   alliancekeeper.Keeper
	sk   stakingkeeper.Keeper
	acck accountkeeper.AccountKeeper
	bk   custombankkeeper.Keeper
}

var _ = govkeeper.Keeper{}

func NewKeeper(
	cdc codec.BinaryCodec,
	key storetypes.StoreKey,
	paramSpace types.ParamSubspace,
	ak accountkeeper.AccountKeeper,
	bk custombankkeeper.Keeper,
	sk *stakingkeeper.Keeper,
	legacyRouter v1beta1.Router,
	router *baseapp.MsgServiceRouter,
	config types.Config,
) Keeper {
	keeper := Keeper{
		Keeper: govkeeper.NewKeeper(cdc, key, paramSpace, ak, bk, sk, legacyRouter, router, config),
		ak:     alliancekeeper.Keeper{},
		bk:     custombankkeeper.Keeper{},
		sk:     stakingkeeper.Keeper{},
		acck:   ak,
		key:    key,
	}
	return keeper
}

func (k *Keeper) RegisterKeepers(ak alliancekeeper.Keeper, bk custombankkeeper.Keeper, sk stakingkeeper.Keeper) {
	k.ak = ak
	k.bk = bk
	k.sk = sk
}

// deleteVote deletes a vote from a given proposalID and voter from the store
func (k Keeper) deleteVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress) {
	store := ctx.KVStore(k.key)
	store.Delete(types.VoteKey(proposalID, voterAddr))
}

func (k *Keeper) Tally(ctx sdk.Context, proposal v1.Proposal) (passes bool, burnDeposits bool, tallyResults v1.TallyResult) {
	results := make(map[v1.VoteOption]sdk.Dec)
	results[v1.OptionYes] = sdk.ZeroDec()
	results[v1.OptionAbstain] = sdk.ZeroDec()
	results[v1.OptionNo] = sdk.ZeroDec()
	results[v1.OptionNoWithVeto] = sdk.ZeroDec()

	totalVotingPower := sdk.ZeroDec()
	currValidators := make(map[string]v1.ValidatorGovInfo)

	// fetch all the bonded validators, insert them into currValidators
	k.sk.IterateBondedValidatorsByPower(ctx, func(index int64, validator stakingtypes.ValidatorI) (stop bool) {
		currValidators[validator.GetOperator().String()] = v1.NewValidatorGovInfo(
			validator.GetOperator(),
			validator.GetBondedTokens(),
			validator.GetDelegatorShares(),
			sdk.ZeroDec(),
			v1.WeightedVoteOptions{},
		)

		return false
	})

	k.IterateVotes(ctx, proposal.Id, func(vote v1.Vote) bool {
		// if validator, just record it in the map
		voter := sdk.MustAccAddressFromBech32(vote.Voter)

		valAddrStr := sdk.ValAddress(voter.Bytes()).String()
		if val, ok := currValidators[valAddrStr]; ok {
			val.Vote = vote.Options
			currValidators[valAddrStr] = val
		}

		// iterate over all delegations from voter, deduct from any delegated-to validators
		k.sk.IterateDelegations(ctx, voter, func(index int64, delegation stakingtypes.DelegationI) (stop bool) {
			valAddrStr := delegation.GetValidatorAddr().String()

			if val, ok := currValidators[valAddrStr]; ok {
				// There is no need to handle the special case that validator address equal to voter address.
				// Because voter's voting power will tally again even if there will be deduction of voter's voting power from validator.
				val.DelegatorDeductions = val.DelegatorDeductions.Add(delegation.GetShares())
				currValidators[valAddrStr] = val

				// delegation shares * bonded / total shares
				votingPower := delegation.GetShares().MulInt(val.BondedTokens).Quo(val.DelegatorShares)

				for _, option := range vote.Options {
					weight, _ := sdk.NewDecFromStr(option.Weight)
					subPower := votingPower.Mul(weight)
					results[option.Option] = results[option.Option].Add(subPower)
				}
				totalVotingPower = totalVotingPower.Add(votingPower)
			}

			return false
		})

		// iterate over all alliance asset delegations from voter, change option to abstain
		k.ak.IterateDelegations(ctx, func(delegation alliancetypes.Delegation) (stop bool) {
			valAddr := delegation.DelegatorAddress

			if val, ok := currValidators[valAddr]; ok {
				// delegation shares * bonded / total shares
				votingPower := delegation.Shares.MulInt(val.BondedTokens).Quo(val.DelegatorShares)

				for _, option := range vote.Options {
					weight, _ := sdk.NewDecFromStr(option.Weight)
					subPower := votingPower.Mul(weight)
					results[option.Option] = results[option.Option].Sub(subPower)
					results[v1.OptionAbstain] = results[v1.OptionAbstain].Add(subPower)
				}
			}

			return false
		})
		k.deleteVote(ctx, vote.ProposalId, voter)
		return false
	})

	// iterate over the validators again to tally their voting power
	for _, val := range currValidators {
		if len(val.Vote) == 0 {
			continue
		}

		sharesAfterDeductions := val.DelegatorShares.Sub(val.DelegatorDeductions)
		votingPower := sharesAfterDeductions.MulInt(val.BondedTokens).Quo(val.DelegatorShares)

		for _, option := range val.Vote {
			weight, _ := sdk.NewDecFromStr(option.Weight)
			subPower := votingPower.Mul(weight)
			results[option.Option] = results[option.Option].Add(subPower)
		}
		totalVotingPower = totalVotingPower.Add(votingPower)
	}

	tallyParams := k.GetTallyParams(ctx)
	tallyResults = v1.NewTallyResultFromMap(results)

	// TODO: Upgrade the spec to cover all of these cases & remove pseudocode.
	// If there is no staked coins, the proposal fails
	if k.sk.TotalBondedTokens(ctx).IsZero() {
		return false, false, tallyResults
	}

	// If there is not enough quorum of votes, the proposal fails
	percentVoting := totalVotingPower.Quo(sdk.NewDecFromInt(k.sk.TotalBondedTokens(ctx)))
	quorum, _ := sdk.NewDecFromStr(tallyParams.Quorum)
	if percentVoting.LT(quorum) {
		return false, false, tallyResults
	}

	// If no one votes (everyone abstains), proposal fails
	if totalVotingPower.Sub(results[v1.OptionAbstain]).Equal(sdk.ZeroDec()) {
		return false, false, tallyResults
	}

	// If more than 1/3 of voters veto, proposal fails
	vetoThreshold, _ := sdk.NewDecFromStr(tallyParams.VetoThreshold)
	if results[v1.OptionNoWithVeto].Quo(totalVotingPower).GT(vetoThreshold) {
		return false, true, tallyResults
	}

	// If more than 1/2 of non-abstaining voters vote Yes, proposal passes
	threshold, _ := sdk.NewDecFromStr(tallyParams.Threshold)
	if results[v1.OptionYes].Quo(totalVotingPower.Sub(results[v1.OptionAbstain])).GT(threshold) {
		return true, false, tallyResults
	}

	// If more than 1/2 of non-abstaining voters vote No, proposal fails
	return false, false, tallyResults
}
