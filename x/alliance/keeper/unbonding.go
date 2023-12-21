package keeper

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-money/alliance/x/alliance/types"
)

// BeginUnbondingsForDissolvingAlliances got thought all alliances, check if there is any alliance winding down
// and if so, start the unbonding process for all delegations
func (k Keeper) BeginUnbondingsForDissolvingAlliances(ctx sdk.Context) (err error) {
	assets := k.GetAllAssets(ctx)

	amountOfDelegationsExecuted := 0
	for _, asset := range assets {
		if !asset.IsDissolving {
			continue
		}

		k.IterateDelegations(ctx, func(delegation types.Delegation) (stop bool) {
			// We should begin unbonding in batches of 50 at the time
			// otherwise it can be too expensive to process
			if amountOfDelegationsExecuted == 50 {
				return true
			}

			if delegation.Denom == asset.Denom {
				validator, fail := k.GetAllianceValidator(ctx, sdk.ValAddress(delegation.DelegatorAddress))
				if err != nil {
					err = fail
					return true
				}

				coinsToUndelegate := types.GetDelegationTokensWithShares(delegation.Shares, validator, *asset)
				_, fail = k.Undelegate(ctx, sdk.AccAddress(delegation.DelegatorAddress), validator, coinsToUndelegate)
				if err != nil {
					err = fail
					return true
				}
			}

			return false
		})

		if err != nil {
			return err
		}
	}

	return err
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

// This method retun all unbonding delegations for a given denom, validator address and delegator address.
// It is the most optimal way to query that data because it uses the indexes that are already in place
// for the unbonding queue and ommits unnecessary checks or data parsings.
func (k Keeper) GetUnbondings(
	ctx sdk.Context,
	denom string,
	delAddr sdk.AccAddress,
	valAddr sdk.ValAddress,
) (unbondingDelegations []types.UnbondingDelegation, err error) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.UndelegationByValidatorIndexKey)
	defer iter.Close()
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

func (k Keeper) addUnbondingAmounts(ctx sdk.Context, unbondingDelegations []types.UnbondingDelegation, delAddr sdk.AccAddress) (unbonding []types.UnbondingDelegation) {
	for i := 0; i < len(unbondingDelegations); i++ {
		iter := k.IterateUndelegationsByCompletionTime(ctx, unbondingDelegations[i].CompletionTime)
		defer iter.Close()
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
