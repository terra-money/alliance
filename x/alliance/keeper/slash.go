package keeper

import (
	"alliance/x/alliance/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func (k Keeper) SlashValidator(ctx sdk.Context, valAddr sdk.ValAddress, fraction sdk.Dec) error {
	val, err := k.GetAllianceValidator(ctx, valAddr)
	if err != nil {
		return err
	}
	slashedValidatorShares := sdk.NewDecCoins()
	for _, share := range val.ValidatorShares {
		sharesToSlash := share.Amount.Mul(fraction)
		slashedValidatorShares = slashedValidatorShares.Add(sdk.NewDecCoinFromDec(share.Denom, share.Amount.Sub(sharesToSlash)))
		asset, found := k.GetAssetByDenom(ctx, share.Denom)
		if !found {
			return types.ErrUnknownAsset
		}
		asset.TotalValidatorShares = asset.TotalValidatorShares.Sub(sharesToSlash)
		k.SetAsset(ctx, asset)
	}
	val.ValidatorShares = slashedValidatorShares
	k.SetValidator(ctx, val)

	err = k.SlashRedelegations(ctx, valAddr, fraction)
	if err != nil {
		return err
	}

	err = k.SlashUndelegations(ctx, valAddr, fraction)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) SlashRedelegations(ctx sdk.Context, valAddr sdk.ValAddress, fraction sdk.Dec) error {
	store := ctx.KVStore(k.storeKey)
	// Slash all immature re-delegations
	redelegationIterator := k.IterateRedelegationsBySrcValidator(ctx, valAddr)
	for ; redelegationIterator.Valid(); redelegationIterator.Next() {
		redelegationKey, completion, err := types.ParseRedelegationIndexForRedelegationKey(redelegationIterator.Key())
		if err != nil {
			return err
		}
		// Skip if redelegation is already mature
		if completion.Before(ctx.BlockTime()) {
			continue
		}
		b := store.Get(redelegationKey)
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

		delegation, found := k.GetDelegation(ctx, delAddr, dstVal, redelegation.Balance.Denom)
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
		dstVal.TotalDelegatorShares = sdk.NewDecCoins(dstVal.TotalDelegatorShares...).Sub(sdk.NewDecCoins(sdk.NewDecCoinFromDec(asset.Denom, sharesToSlash)))
		k.SetValidator(ctx, dstVal)

		delegation.Shares = delegation.Shares.Sub(sharesToSlash)
		k.SetDelegation(ctx, delAddr, dstVal, asset.Denom, delegation)
	}
	return nil
}

func (k Keeper) SlashUndelegations(ctx sdk.Context, valAddr sdk.ValAddress, fraction sdk.Dec) error {
	store := ctx.KVStore(k.storeKey)
	// Slash all immature re-delegations
	undelegationIterator := k.IterateUndelegationsBySrcValidator(ctx, valAddr)
	for ; undelegationIterator.Valid(); undelegationIterator.Next() {
		undelegationKey, completion, err := types.ParseUnbondingIndexKeyToUndelegationKey(undelegationIterator.Key())
		if err != nil {
			return err
		}
		// Skip if undelegation is already mature
		if completion.Before(ctx.BlockTime()) {
			continue
		}
		b := store.Get(undelegationKey)
		var undelegations types.QueuedUndelegation
		k.cdc.MustUnmarshal(b, &undelegations)

		// Slash undelegations by sending slashed tokens to fee pool
		for _, entry := range undelegations.Entries {
			tokensToSlash := fraction.MulInt(entry.Balance.Amount).TruncateInt()
			entry.Balance = sdk.NewCoin(entry.Balance.Denom, entry.Balance.Amount.Sub(tokensToSlash))
			coinToSlash := sdk.NewCoin(entry.Balance.Denom, tokensToSlash)
			err = k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, authtypes.FeeCollectorName, sdk.NewCoins(coinToSlash))
			if err != nil {
				return err
			}
		}
		b = k.cdc.MustMarshal(&undelegations)
		store.Set(undelegationKey, b)
	}
	return nil
}
