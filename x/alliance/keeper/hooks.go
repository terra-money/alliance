package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type Hooks struct {
	k Keeper
}

var _ stakingtypes.StakingHooks = Hooks{}

func (h Hooks) AfterValidatorCreated(_ sdk.Context, _ sdk.ValAddress) error {
	return nil
}

func (h Hooks) BeforeValidatorModified(_ sdk.Context, _ sdk.ValAddress) error {
	return nil
}

func (h Hooks) AfterValidatorRemoved(ctx sdk.Context, _ sdk.ConsAddress, valAddr sdk.ValAddress) error {
	h.k.DeleteValidatorInfo(ctx, valAddr)
	h.k.QueueAssetRebalanceEvent(ctx)
	return nil
}

func (h Hooks) AfterValidatorBonded(ctx sdk.Context, _ sdk.ConsAddress, _ sdk.ValAddress) error {
	h.k.QueueAssetRebalanceEvent(ctx)
	return nil
}

func (h Hooks) AfterValidatorBeginUnbonding(ctx sdk.Context, _ sdk.ConsAddress, _ sdk.ValAddress) error {
	h.k.QueueAssetRebalanceEvent(ctx)
	return nil
}

func (h Hooks) BeforeDelegationCreated(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress) error {
	return nil
}

func (h Hooks) BeforeDelegationSharesModified(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress) error {
	return nil
}

func (h Hooks) BeforeDelegationRemoved(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress) error {
	return nil
}

func (h Hooks) AfterDelegationModified(ctx sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress) error {
	h.k.QueueAssetRebalanceEvent(ctx)
	return nil
}

func (h Hooks) BeforeValidatorSlashed(ctx sdk.Context, valAddr sdk.ValAddress, fraction sdk.Dec) error {
	err := h.k.SlashValidator(ctx, valAddr, fraction)
	if err != nil {
		return err
	}
	h.k.QueueAssetRebalanceEvent(ctx)
	return nil
}

func (h Hooks) AfterValidatorSlashed(_ sdk.Context, _ sdk.ConsAddress, _ sdk.ValAddress, _ sdk.Dec) {
}

func (h Hooks) AfterUnbondingInitiated(_ sdk.Context, _ uint64) error {
	return nil
}
