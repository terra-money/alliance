package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-money/alliance/x/alliance/types"
)

// SlashValidator works by reducing the amount of validator shares for all alliance assets by a `fraction`
// This effectively reallocates tokens from slashed validators to good validators
// On top of slashing currently bonded delegations, we also slash re-delegations and un-delegations
// that are still in the progress of unbonding
func (k Keeper) SlashValidator(ctx context.Context, valAddr sdk.ValAddress, fraction math.LegacyDec) error {
	// Slashing must be checked otherwise we can end up slashing incorrect amounts
	if fraction.LTE(math.LegacyZeroDec()) || fraction.GT(math.LegacyOneDec()) {
		return fmt.Errorf("slashed fraction must be greater than 0 and less than or equal to 1: %d", fraction)
	}

	val, err := k.GetAllianceValidator(ctx, valAddr)
	if err != nil {
		return err
	}
	// slashedValidatorShares accumulates the final validator shares after slashing
	slashedValidatorShares := sdk.NewDecCoins()
	for _, share := range val.ValidatorShares {
		sharesToSlash := share.Amount.Mul(fraction)
		sharesAfterSlashing := sdk.NewDecCoinFromDec(share.Denom, share.Amount.Sub(sharesToSlash))
		slashedValidatorShares = slashedValidatorShares.Add(sharesAfterSlashing)
		asset, found := k.GetAssetByDenom(ctx, share.Denom)
		if !found {
			return types.ErrUnknownAsset
		}
		asset.TotalValidatorShares = asset.TotalValidatorShares.Sub(sharesToSlash)
		if err = k.SetAsset(ctx, asset); err != nil {
			return err
		}
	}
	val.ValidatorShares = slashedValidatorShares
	if err = k.SetValidator(ctx, val); err != nil {
		return err
	}

	err = k.slashRedelegations(ctx, valAddr, fraction)
	if err != nil {
		return err
	}

	err = k.slashUndelegations(ctx, valAddr, fraction)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) slashRedelegations(ctx context.Context, valAddr sdk.ValAddress, fraction math.LegacyDec) error {
	store := k.storeService.OpenKVStore(ctx)
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	// Slash all immature re-delegations
	redelegationIterator := k.IterateRedelegationsBySrcValidator(ctx, valAddr)
	for ; redelegationIterator.Valid(); redelegationIterator.Next() {
		redelegationKey, completion, err := types.ParseRedelegationIndexForRedelegationKey(redelegationIterator.Key())
		if err != nil {
			return err
		}
		// Skip if redelegation is already mature
		if completion.Before(sdkCtx.BlockTime()) {
			continue
		}
		b, err := store.Get(redelegationKey)
		if err != nil {
			return err
		}
		var redelegation types.Redelegation
		k.cdc.MustUnmarshal(b, &redelegation)

		delAddr, err := sdk.AccAddressFromBech32(redelegation.DelegatorAddress)
		if err != nil {
			return err
		}
		dstValAddr, err := sdk.ValAddressFromBech32(redelegation.DstValidatorAddress)
		if err != nil {
			return err
		}
		dstVal, err := k.GetAllianceValidator(ctx, dstValAddr)
		if err != nil {
			return err
		}

		_, err = k.ClaimDelegationRewards(ctx, delAddr, dstVal, redelegation.Balance.Denom)
		if err != nil {
			return err
		}

		delegation, found := k.GetDelegation(ctx, delAddr, dstValAddr, redelegation.Balance.Denom)
		if !found {
			continue
		}

		asset, found := k.GetAssetByDenom(ctx, redelegation.Balance.Denom)
		if !found {
			continue
		}

		// Slash delegation shares
		tokensToSlash := fraction.MulInt(redelegation.Balance.Amount).TruncateInt()
		sharesToSlash, err := k.ValidateDelegatedAmount(delegation, sdk.NewCoin(redelegation.Balance.Denom, tokensToSlash), dstVal, asset)
		if err != nil {
			return err
		}
		dstVal.TotalDelegatorShares = sdk.DecCoins(dstVal.TotalDelegatorShares).Sub(sdk.NewDecCoins(sdk.NewDecCoinFromDec(asset.Denom, sharesToSlash)))
		if err = k.SetValidator(ctx, dstVal); err != nil {
			return err
		}

		delegation.Shares = delegation.Shares.Sub(sharesToSlash)
		if err = k.SetDelegation(ctx, delAddr, dstValAddr, asset.Denom, delegation); err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) slashUndelegations(ctx context.Context, valAddr sdk.ValAddress, fraction math.LegacyDec) error {
	store := k.storeService.OpenKVStore(ctx)
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	// Slash all immature re-delegations
	undelegationIterator := k.IterateUndelegationsBySrcValidator(ctx, valAddr)
	for ; undelegationIterator.Valid(); undelegationIterator.Next() {
		undelegationKey, completion, err := types.ParseUnbondingIndexKeyToUndelegationKey(undelegationIterator.Key())
		if err != nil {
			return err
		}
		// Skip if undelegation is already mature
		if completion.Before(sdkCtx.BlockTime()) {
			continue
		}
		b, err := store.Get(undelegationKey)
		if err != nil {
			return err
		}
		var undelegations types.QueuedUndelegation
		k.cdc.MustUnmarshal(b, &undelegations)

		// Slash undelegations by sending slashed tokens to fee pool
		for _, entry := range undelegations.Entries {
			tokensToSlash := fraction.MulInt(entry.Balance.Amount).TruncateInt()
			entry.Balance = sdk.NewCoin(entry.Balance.Denom, entry.Balance.Amount.Sub(tokensToSlash))
			coinToSlash := sdk.NewCoin(entry.Balance.Denom, tokensToSlash)
			err = k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.feeCollectorName, sdk.NewCoins(coinToSlash))
			if err != nil {
				return err
			}
		}
		b = k.cdc.MustMarshal(&undelegations)
		err = store.Set(undelegationKey, b)
		if err != nil {
			return err
		}
	}
	return nil
}
