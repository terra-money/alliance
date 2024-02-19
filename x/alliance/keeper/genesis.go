package keeper

import (
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-money/alliance/x/alliance/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, g *types.GenesisState) []abci.ValidatorUpdate {
	if err := k.SetParams(ctx, g.Params); err != nil {
		panic(err)
	}
	for _, asset := range g.Assets {
		if err := sdk.ValidateDenom(asset.Denom); err != nil {
			panic(err)
		}
		if err := k.SetAsset(ctx, asset); err != nil {
			panic(err)
		}

	}

	for _, val := range g.ValidatorInfos {
		valAddr, _ := sdk.ValAddressFromBech32(val.ValidatorAddress)
		if err := k.SetValidatorInfo(ctx, valAddr, val.Validator); err != nil {
			panic(err)
		}
	}

	for _, delegation := range g.Delegations {
		delAddr, _ := sdk.AccAddressFromBech32(delegation.DelegatorAddress)
		valAddr, _ := sdk.ValAddressFromBech32(delegation.ValidatorAddress)
		if err := k.SetDelegation(ctx, delAddr, valAddr, delegation.Denom, delegation); err != nil {
			panic(err)
		}
	}

	for _, redelegationState := range g.Redelegations {
		delAddr, _ := sdk.AccAddressFromBech32(redelegationState.Redelegation.DelegatorAddress)
		srcValAddr, _ := sdk.ValAddressFromBech32(redelegationState.Redelegation.SrcValidatorAddress)
		dstValAddr, _ := sdk.ValAddressFromBech32(redelegationState.Redelegation.DstValidatorAddress)

		if err := k.addRedelegation(ctx, delAddr, srcValAddr, dstValAddr, redelegationState.Redelegation.Balance, redelegationState.CompletionTime); err != nil {
			panic(err)
		}
		if err := k.queueRedelegation(ctx, delAddr, srcValAddr, dstValAddr, redelegationState.Redelegation.Balance, redelegationState.CompletionTime); err != nil {
			panic(err)
		}
	}

	for _, undelegationState := range g.Undelegations {
		if len(undelegationState.Undelegation.Entries) == 0 {
			continue
		}
		delAddr, _ := sdk.AccAddressFromBech32(undelegationState.Undelegation.Entries[0].DelegatorAddress)
		if err := k.setQueuedUndelegations(ctx, undelegationState.CompletionTime, delAddr, undelegationState.Undelegation); err != nil {
			panic(err)
		}
		for _, undelegation := range undelegationState.Undelegation.Entries {
			valAddr, _ := sdk.ValAddressFromBech32(undelegation.ValidatorAddress)
			if err := k.setUnbondingIndexByVal(ctx, valAddr, undelegationState.CompletionTime, delAddr, undelegation.Balance.Denom); err != nil {
				panic(err)
			}
		}
	}

	for _, rewardWeightSnapshot := range g.RewardWeightChangeSnaphots {
		valAddr, _ := sdk.ValAddressFromBech32(rewardWeightSnapshot.Validator)
		if err := k.setRewardWeightChangeSnapshot(ctx, rewardWeightSnapshot.Denom, valAddr, rewardWeightSnapshot.Height, rewardWeightSnapshot.Snapshot); err != nil {
			panic(err)
		}
	}

	return []abci.ValidatorUpdate{}
}

func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	state := types.GenesisState{}
	assets := k.GetAllAssets(ctx)
	for _, asset := range assets {
		state.Assets = append(state.Assets, *asset)
	}

	if err := k.IterateAllianceValidatorInfo(ctx, func(valAddr sdk.ValAddress, info types.AllianceValidatorInfo) (stop bool) {
		state.ValidatorInfos = append(state.ValidatorInfos, types.ValidatorInfoState{
			ValidatorAddress: valAddr.String(),
			Validator:        info,
		})
		return false
	}); err != nil {
		panic(err)
	}

	if err := k.IterateDelegations(ctx, func(d types.Delegation) (stop bool) {
		state.Delegations = append(state.Delegations, d)
		return false
	}); err != nil {
		panic(err)
	}

	k.IterateRedelegations(ctx, func(r types.Redelegation, completionTime time.Time) (stop bool) {
		state.Redelegations = append(state.Redelegations, types.RedelegationState{
			CompletionTime: completionTime,
			Redelegation:   r,
		})
		return false
	})

	k.IterateUndelegations(ctx, func(u types.QueuedUndelegation, completionTime time.Time) (stop bool) {
		state.Undelegations = append(state.Undelegations, types.UndelegationState{
			CompletionTime: completionTime,
			Undelegation:   u,
		})
		return false
	})

	k.IterateAllWeightChangeSnapshot(ctx, func(denom string, valAddr sdk.ValAddress, height uint64, snapshot types.RewardWeightChangeSnapshot) bool {
		state.RewardWeightChangeSnaphots = append(state.RewardWeightChangeSnaphots, types.RewardWeightChangeSnapshotState{
			Height:    height,
			Validator: valAddr.String(),
			Denom:     denom,
			Snapshot:  snapshot,
		})
		return false
	})

	state.Params = types.Params{
		RewardDelayTime:       k.RewardDelayTime(ctx),
		TakeRateClaimInterval: k.RewardClaimInterval(ctx),
		LastTakeRateClaimTime: k.LastRewardClaimTime(ctx),
	}

	return &state
}
