package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/terra-money/alliance/x/alliance/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, g *types.GenesisState) []abci.ValidatorUpdate {
	k.SetParams(ctx, g.Params)
	for _, asset := range g.Assets {
		if err := sdk.ValidateDenom(asset.Denom); err != nil {
			panic(err)
		}
		k.SetAsset(ctx, asset)
	}

	for _, val := range g.ValidatorInfos {
		valAddr, _ := sdk.ValAddressFromBech32(val.ValidatorAddress)
		k.SetValidatorInfo(ctx, valAddr, val.Validator)
	}

	for _, delegation := range g.Delegations {
		delAddr, _ := sdk.AccAddressFromBech32(delegation.DelegatorAddress)
		valAddr, _ := sdk.ValAddressFromBech32(delegation.ValidatorAddress)
		k.SetDelegation(ctx, delAddr, valAddr, delegation.Denom, delegation)
	}

	for _, redelegationState := range g.Redelegations {
		delAddr, _ := sdk.AccAddressFromBech32(redelegationState.Redelegation.DelegatorAddress)
		srcValAddr, _ := sdk.ValAddressFromBech32(redelegationState.Redelegation.SrcValidatorAddress)
		dstValAddr, _ := sdk.ValAddressFromBech32(redelegationState.Redelegation.DstValidatorAddress)

		k.addRedelegation(ctx, delAddr, srcValAddr, dstValAddr, redelegationState.Redelegation.Balance, redelegationState.CompletionTime)
		k.queueRedelegation(ctx, delAddr, srcValAddr, dstValAddr, redelegationState.Redelegation.Balance, redelegationState.CompletionTime)
	}

	for _, undelegationState := range g.Undelegations {
		if len(undelegationState.Undelegation.Entries) == 0 {
			continue
		}
		delAddr, _ := sdk.AccAddressFromBech32(undelegationState.Undelegation.Entries[0].DelegatorAddress)
		k.setQueuedUndelegations(ctx, undelegationState.CompletionTime, delAddr, undelegationState.Undelegation)
		for _, undelegation := range undelegationState.Undelegation.Entries {
			valAddr, _ := sdk.ValAddressFromBech32(undelegation.ValidatorAddress)
			k.setUnbondingIndexByVal(ctx, valAddr, undelegationState.CompletionTime, delAddr, undelegation.Balance.Denom)
		}
	}

	for _, rewardWeightSnapshot := range g.RewardWeightChangeSnaphots {
		valAddr, _ := sdk.ValAddressFromBech32(rewardWeightSnapshot.Validator)
		k.setRewardWeightChangeSnapshot(ctx, rewardWeightSnapshot.Denom, valAddr, rewardWeightSnapshot.Height, rewardWeightSnapshot.Snapshot)
	}

	return []abci.ValidatorUpdate{}
}

func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	state := types.GenesisState{}
	assets := k.GetAllAssets(ctx)
	for _, asset := range assets {
		state.Assets = append(state.Assets, *asset)
	}

	k.IterateAllianceValidatorInfo(ctx, func(valAddr sdk.ValAddress, info types.AllianceValidatorInfo) (stop bool) {
		state.ValidatorInfos = append(state.ValidatorInfos, types.ValidatorInfoState{
			ValidatorAddress: valAddr.String(),
			Validator:        info,
		})
		return false
	})

	k.IterateDelegations(ctx, func(d types.Delegation) (stop bool) {
		state.Delegations = append(state.Delegations, d)
		return false
	})

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

func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}
