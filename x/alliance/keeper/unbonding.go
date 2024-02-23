package keeper

import (
	"bytes"
	"context"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-money/alliance/x/alliance/types"
)

// CompleteUnbondings Go through all queued undelegations and send the tokens to the delAddrs
func (k Keeper) CompleteUnbondings(ctx context.Context) error {
	store := k.storeService.OpenKVStore(ctx)
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	iter := k.IterateUndelegationsByCompletionTime(ctx, sdkCtx.BlockTime())
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var queued types.QueuedUndelegation
		completionTime, err := types.ParseUndelegationQueueKeyForCompletionTime(iter.Key())
		if err != nil {
			return err
		}
		k.cdc.MustUnmarshal(iter.Value(), &queued)
		for _, undel := range queued.Entries {
			delAddr, err := sdk.AccAddressFromBech32(undel.DelegatorAddress)
			if err != nil {
				return err
			}
			if err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, delAddr, sdk.NewCoins(undel.Balance)); err != nil {
				return err
			}
			valAddr, err := sdk.ValAddressFromBech32(undel.ValidatorAddress)
			if err != nil {
				return err
			}
			indexKey := types.GetUnbondingIndexKey(valAddr, completionTime, undel.Balance.Denom, delAddr)
			if err = store.Delete(indexKey); err != nil {
				return err
			}
		}
		if err := store.Delete(iter.Key()); err != nil {
			return err
		}
	}

	// Burn all "virtual" staking tokens in the module account
	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	bondDenom, err := k.stakingKeeper.BondDenom(ctx)
	if err != nil {
		return err
	}
	coin := k.bankKeeper.GetBalance(ctx, moduleAddr, bondDenom)
	if !coin.IsZero() {
		err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(coin))
		if err != nil {
			return err
		}
	}
	return nil
}

// This method retun all unbonding delegations for a given denom, validator address and delegator address.
// It is the most optimal way to query that data because it uses the indexes that are already in place
// for the unbonding queue and ommits unnecessary checks or data parsings.
func (k Keeper) GetUnbondings(
	ctx context.Context,
	denom string,
	delAddr sdk.AccAddress,
	valAddr sdk.ValAddress,
) (unbondingDelegations []types.UnbondingDelegation, err error) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iter := storetypes.KVStorePrefixIterator(store, types.UndelegationByValidatorIndexKey)
	defer iter.Close() //nolint:errcheck,nolintlint
	suffix := types.GetPartialUnbondingKeySuffix(denom, delAddr)
	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		if len(key) < len(suffix) {
			continue // Skip keys that are shorter than the suffix
		}

		prefix := types.GetUndelegationsIndexOrderedByValidatorKey(valAddr)
		if !bytes.HasPrefix(key, prefix) {
			continue // Skip keys that don't have the desired suffix
		}

		if !bytes.HasSuffix(key, suffix) {
			continue // Skip keys that don't have the desired suffix
		}

		_, unbondingCompletionTime, err := types.PartiallyParseUndelegationKeyBytes(key)
		if err != nil {
			return nil, err
		}
		// Process and append item to the results
		unbondDelegation := types.UnbondingDelegation{
			ValidatorAddress: valAddr.String(),
			CompletionTime:   unbondingCompletionTime,
		}
		unbondingDelegations = append(unbondingDelegations, unbondDelegation)
	}

	unbondingDelegations = k.addUnbondingAmounts(ctx, unbondingDelegations, delAddr)

	return unbondingDelegations, err
}

// This method retun all in-progress unbondings for a given delegator address
// it is less optimal than GetUnbondingsByDenomAndDelegator because it
// has to iterate over all alliances to get the list of all assets
func (k Keeper) GetUnbondingsByDelegator(
	ctx context.Context,
	delAddr sdk.AccAddress,
) (unbondingDelegations []types.UnbondingDelegation, err error) {
	// Retrieve all Aliances to get the list of all assets
	alliances := k.GetAllAssets(ctx)

	for _, alliance := range alliances {
		// Get the unbonding delegations for the current alliance
		unbondingDelegations, err = k.GetUnbondingsByDenomAndDelegator(ctx, alliance.Denom, delAddr)
		if err != nil {
			return nil, err
		}
	}
	return unbondingDelegations, err
}

// This method retun all unbonding delegations for a given denom and delegator address,
// it is less optimal than GetUnbondings because it has do some data parsing and additional
// checks, plus it returns a larger data set.
func (k Keeper) GetUnbondingsByDenomAndDelegator(
	ctx context.Context,
	denom string,
	delAddr sdk.AccAddress,
) (unbondingDelegations []types.UnbondingDelegation, err error) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iter := storetypes.KVStorePrefixIterator(store, types.UndelegationByValidatorIndexKey)
	defer iter.Close() //nolint:errcheck,nolintlint
	suffix := types.GetPartialUnbondingKeySuffix(denom, delAddr)

	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		if len(key) < len(suffix) {
			continue // Skip keys that are shorter than the suffix
		}

		if !bytes.HasSuffix(key, suffix) {
			continue // Skip keys that don't have the desired suffix
		}

		valAddr, unbondingCompletionTime, err := types.PartiallyParseUndelegationKeyBytes(key)
		if err != nil {
			return nil, err
		}
		// Process and append item to the results
		unbondDelegation := types.UnbondingDelegation{
			ValidatorAddress: valAddr.String(),
			CompletionTime:   unbondingCompletionTime,
		}
		unbondingDelegations = append(unbondingDelegations, unbondDelegation)
	}

	unbondingDelegations = k.addUnbondingAmounts(ctx, unbondingDelegations, delAddr)

	return unbondingDelegations, err
}

func (k Keeper) addUnbondingAmounts(ctx context.Context, unbondingDelegations []types.UnbondingDelegation, delAddr sdk.AccAddress) (unbonding []types.UnbondingDelegation) {
	for i := 0; i < len(unbondingDelegations); i++ {
		iter := k.IterateUndelegationsByCompletionTime(ctx, unbondingDelegations[i].CompletionTime)
		defer iter.Close() //nolint:errcheck,nolintlint
		for ; iter.Valid(); iter.Next() {
			var queued types.QueuedUndelegation
			k.cdc.MustUnmarshal(iter.Value(), &queued)

			for _, undel := range queued.Entries {
				if undel.DelegatorAddress != delAddr.String() {
					continue
				}

				if undel.ValidatorAddress != unbondingDelegations[i].ValidatorAddress {
					continue
				}

				unbondingDelegations[i].Amount = undel.Balance.Amount
			}
		}
	}

	for _, unbondingDelegation := range unbondingDelegations {
		if unbondingDelegation.Amount.IsNil() {
			continue
		}

		unbonding = append(unbonding, unbondingDelegation)
	}

	return unbonding
}
