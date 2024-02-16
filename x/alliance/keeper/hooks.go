package keeper

import (
	"context"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type Hooks struct {
	k Keeper
}

var _ stakingtypes.StakingHooks = Hooks{}

func (h Hooks) AfterValidatorCreated(_ context.Context, _ sdk.ValAddress) error {
	return nil
}

func (h Hooks) BeforeValidatorModified(_ context.Context, _ sdk.ValAddress) error {
	return nil
}

func (h Hooks) AfterValidatorRemoved(ctx context.Context, _ sdk.ConsAddress, valAddr sdk.ValAddress) error {
	h.k.DeleteValidatorInfo(ctx, valAddr)
	h.k.QueueAssetRebalanceEvent(ctx)
	return nil
}

func (h Hooks) AfterValidatorBonded(ctx context.Context, _ sdk.ConsAddress, _ sdk.ValAddress) error {
	h.k.QueueAssetRebalanceEvent(ctx)
	return nil
}

func (h Hooks) AfterValidatorBeginUnbonding(ctx context.Context, _ sdk.ConsAddress, _ sdk.ValAddress) error {
	h.k.QueueAssetRebalanceEvent(ctx)
	return nil
}

func (h Hooks) BeforeDelegationCreated(_ context.Context, _ sdk.AccAddress, _ sdk.ValAddress) error {
	return nil
}

func (h Hooks) BeforeDelegationSharesModified(_ context.Context, _ sdk.AccAddress, _ sdk.ValAddress) error {
	return nil
}

func (h Hooks) BeforeDelegationRemoved(_ context.Context, _ sdk.AccAddress, _ sdk.ValAddress) error {
	return nil
}

func (h Hooks) AfterDelegationModified(ctx context.Context, _ sdk.AccAddress, _ sdk.ValAddress) error {
	h.k.QueueAssetRebalanceEvent(ctx)
	return nil
}

func (h Hooks) BeforeValidatorSlashed(ctx context.Context, valAddr sdk.ValAddress, fraction math.LegacyDec) error {
	err := h.k.SlashValidator(ctx, valAddr, fraction)
	if err != nil {
		return err
	}
	h.k.QueueAssetRebalanceEvent(ctx)
	return nil
}

func (h Hooks) AfterValidatorSlashed(_ context.Context, _ sdk.ConsAddress, _ sdk.ValAddress, _ math.LegacyDec) {
}

func (h Hooks) AfterUnbondingInitiated(_ context.Context, _ uint64) error {
	return nil
}
