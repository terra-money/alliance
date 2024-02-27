package keeper

import (
	"bytes"
	"context"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-money/alliance/x/alliance/types"
)

// BeginUnbondingsForDissolvingAlliances iterates over the alliances,
// if there is any alliance being dissolved:
//   - check if the AllianceDissolutionTime has been set and if it is in the past compared to
//     the current block time the alliance will be deleted.
//   - if the alliance is not initialized and is being dissolved it means that
//     all alliance unbondings have been executed and no further comput needs to be done.
//   - else iterate the delegations and begin the unbonding for all assets,
//     setting always the last unbonding period for the latest unbonding executed
//     as the AllianceDissolutionTime, so that we can keep track of when to delete the alliance.
func (k Keeper) BeginUnbondingsForDissolvingAlliances(ctx sdk.Context) (err error) {
	// Iterate over all alliances...
	for _, asset := range k.GetAllAssets(ctx) {
		// Check if any alliance is dissolving.
		// ¡NOTE!: It's important to continue the loop if the alliance is not dissolving.
		// because there is logic down the linde that depends on dissolving alliances.
		if !asset.IsDissolving {
			continue
		}
		assert := *asset
		// Variable that keeps track of how many unbondings have been executed
		amountOfUndelegationsExecuted := 0

		// In theory if an alliance is dissolving it should have a dissolution time set
		// and if the dissolution time is in the past we should delete the alliance
		// because it means that the unbonding period for all the delegations is completed
		// and funds have been sent back to the delegators.
		if assert.AllianceDissolutionTime != nil {
			if ctx.BlockTime().After(*assert.AllianceDissolutionTime) {
				err := k.DeleteAsset(ctx, assert)
				if err != nil {
					return err
				}
				continue
			}
		}
		// If the alliance is not initialized and is being dissolved it means that
		// all alliance unbondings have been executed and no further comput needs to be done.
		if !assert.IsInitialized && assert.IsDissolving {
			continue
		}

		// Iterate over all the delegations
		err := k.IterateDelegations(ctx, func(delegation types.Delegation) (stop bool) {
			// if the delegation being checked is not for the current alliance we should
			// continue to the next delegation without spending more comput time.
			if delegation.Denom != assert.Denom {
				return false
			}

			// We should begin unbonding in batches of 50 at the time
			// otherwise it can be too expensive to process and blocks
			// can take too long to be process.
			// TODO: make of this a module parameter at some point.
			if amountOfUndelegationsExecuted == 50 {
				return true
			}

			// Parse delegator address
			delAddr, fail := sdk.AccAddressFromBech32(delegation.DelegatorAddress)
			if fail != nil {
				err = fail
				return true
			}

			// Parse validator address
			vaAddr, fail := sdk.ValAddressFromBech32(delegation.ValidatorAddress)
			if fail != nil {
				err = fail
				return true
			}

			// Get the alliances alliance validator info
			validator, fail := k.GetAllianceValidator(ctx, vaAddr)
			if fail != nil {
				err = fail
				return true
			}

			// Calculate the amount of coins to unbond
			coinsToUnbond := types.GetDelegationTokens(delegation, validator, assert)

			// Execute the unbonding
			time, fail := k.Undelegate(ctx, delAddr, validator, coinsToUnbond)
			if fail != nil {
				err = fail
				return true
			}

			// Set the last unbonding time as the alliance dissolution time
			asset.AllianceDissolutionTime = time
			fail = k.SetAsset(ctx, assert)
			if fail != nil {
				err = fail
				return true
			}

			// Increment the amount of unbondings executed so at the begining of the
			// loop we can check if we have executed 50 unbondings and stop the loop.
			amountOfUndelegationsExecuted++

			return false
		})
		if err != nil {
			return err
		}

		// ¡NOTE! If no unbondings have been executed but the alliance is in dissolving state,
		// we should set the alliance as not initialized because it could be that the alliance
		// the alliance has never had delegations or that all the delegations have been unbonded already,
		// so going back to the line 52 of this file we avoid doing unnecessary comput.
		if amountOfUndelegationsExecuted == 0 {
			assert.IsInitialized = false
			err = k.SetAsset(ctx, assert)
			if err != nil {
				return err
			}
		}
	}

	return err
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
	// Get the store
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	// create the iterator with the correct prefix
	prefix := types.GetUndelegationsIndexOrderedByValidatorKey(valAddr)
	// Get the iterator
	iter := storetypes.KVStorePrefixIterator(store, prefix)
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

// This method retun all in-progress unbondings for a given delegator address
func (k Keeper) GetUnbondingsByDelegator(
	ctx context.Context,
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
	ctx context.Context,
	denom string,
	delAddr sdk.AccAddress,
) (unbondingDelegations []types.UnbondingDelegation, err error) {
	// Get the store
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	// create the iterator with the correct prefix
	iter := storetypes.KVStorePrefixIterator(store, types.UndelegationByValidatorIndexKey)
	defer iter.Close() //nolint:errcheck,nolintlint
	// Get the suffix key
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
