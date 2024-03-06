package keeper

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-money/alliance/x/alliance/types"
)

// This method retun all unbonding delegations for a given denom, validator address and delegator address.
// It is the most optimal way to query that data because it uses the indexes that are already in place
// for the unbonding queue and ommits unnecessary checks or data parsings.
func (k Keeper) GetUnbondings(
	ctx sdk.Context,
	denom string,
	delAddr sdk.AccAddress,
	valAddr sdk.ValAddress,
) (unbondingDelegations []types.UnbondingDelegation, err error) {
	// Get the store
	store := ctx.KVStore(k.storeKey)
	// create the iterator with the correct prefix
	prefix := types.GetUndelegationsIndexOrderedByValidatorKey(valAddr)
	// Get the iterator
	iter := sdk.KVStorePrefixIterator(store, prefix)
	defer iter.Close()
	suffix := types.GetPartialUnbondingKeySuffix(denom, delAddr)

	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		// Skip keys that don't have the desired suffix
		if !bytes.HasSuffix(key, suffix) {
			continue
		}

		prefix := types.GetUndelegationsIndexOrderedByValidatorKey(valAddr)
		// Skip keys that don't have the desired suffix
		if !bytes.HasPrefix(key, prefix) {
			continue
		}

		completionTime, err := types.GetTimeFromUndelegationKey(key)
		if err != nil {
			return nil, err
		}
		// Recover the queued undelegation from the store
		b := store.Get(types.GetUndelegationQueueKey(completionTime, delAddr))

		// Parse the model from the bytes
		var queue types.QueuedUndelegation
		err = k.cdc.Unmarshal(b, &queue)
		if err != nil {
			return nil, err
		}

		// Iterate over the entries and append them to the result
		for _, entry := range queue.Entries {
			unbondDelegation := types.UnbondingDelegation{
				ValidatorAddress: entry.ValidatorAddress,
				CompletionTime:   completionTime,
				Amount:           entry.Balance.Amount,
				Denom:            entry.Balance.Denom,
			}
			unbondingDelegations = append(unbondingDelegations, unbondDelegation)
		}
	}

	return unbondingDelegations, err
}

// CompleteUnbondings Go through all queued undelegations and send the tokens to the delAddrs
func (k Keeper) CompleteUnbondings(ctx sdk.Context) error {
	store := ctx.KVStore(k.storeKey)
	iter := k.IterateUndelegationsByCompletionTime(ctx, ctx.BlockTime())
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
			err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, delAddr, sdk.NewCoins(undel.Balance))
			if err != nil {
				return err
			}
			valAddr, err := sdk.ValAddressFromBech32(undel.ValidatorAddress)
			if err != nil {
				return err
			}
			indexKey := types.GetUnbondingIndexKey(valAddr, completionTime, undel.Balance.Denom, delAddr)
			store.Delete(indexKey)
		}
		store.Delete(iter.Key())
	}

	// Burn all "virtual" staking tokens in the module account
	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	coin := k.bankKeeper.GetBalance(ctx, moduleAddr, k.stakingKeeper.BondDenom(ctx))
	if !coin.IsZero() {
		err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(coin))
		if err != nil {
			return err
		}
	}
	return nil
}

// This method retun all in-progress unbondings for a given delegator address
func (k Keeper) GetUnbondingsByDelegator(
	ctx sdk.Context,
	delAddr sdk.AccAddress,
) (unbondingDelegations []types.UnbondingDelegation, err error) {
	// Get and iterate over all alliances
	for _, alliance := range k.GetAllAssets(ctx) {
		// Get the unbonding delegations for the current alliance
		unbondings, err := k.GetUnbondingsByDenomAndDelegator(ctx, alliance.Denom, delAddr)
		if err != nil {
			return nil, err
		}
		unbondingDelegations = append(unbondingDelegations, unbondings...)
	}
	return unbondingDelegations, err
}

// This method retun all unbonding delegations for a given denom and delegator address,
// it is less optimal than GetUnbondings because it has do some data parsing and additional
// checks, plus it returns a larger data set.
func (k Keeper) GetUnbondingsByDenomAndDelegator(
	ctx sdk.Context,
	denom string,
	delAddr sdk.AccAddress,
) (unbondingDelegations []types.UnbondingDelegation, err error) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.UndelegationByValidatorIndexKey)
	defer iter.Close()
	suffix := types.GetPartialUnbondingKeySuffix(denom, delAddr)

	// Iterate over the keys
	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		// Continue to the next iteration if the key is shorter than the suffix
		if len(key) < len(suffix) {
			continue
		}
		// continue to the next iteration if the key doesn't have the desired suffix
		if !bytes.HasSuffix(key, suffix) {
			continue
		}
		// parse the key and get the unbonding completion time
		completionTime, err := types.GetTimeFromUndelegationKey(key)
		if err != nil {
			return nil, err
		}
		// Recover the queued undelegation from the store
		b := store.Get(types.GetUndelegationQueueKey(completionTime, delAddr))

		// Parse the model from the bytes
		var queue types.QueuedUndelegation
		err = k.cdc.Unmarshal(b, &queue)
		if err != nil {
			return nil, err
		}

		// Iterate over the entries and append them to the result
		for _, entry := range queue.Entries {
			unbondDelegation := types.UnbondingDelegation{
				ValidatorAddress: entry.ValidatorAddress,
				CompletionTime:   completionTime,
				Amount:           entry.Balance.Amount,
				Denom:            entry.Balance.Denom,
			}
			unbondingDelegations = append(unbondingDelegations, unbondDelegation)
		}
	}
	return unbondingDelegations, err
}
