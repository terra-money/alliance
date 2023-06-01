package keeper

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/terra-money/alliance/x/alliance/types"
)

func (k Keeper) GetAllUnbondings(ctx sdk.Context, denom string, delegator sdk.AccAddress) (items []types.UnbondingDelegation, err error) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.UndelegationByValidatorIndexKey) // Start iteration from the beginning
	defer iterator.Close()
	suffix := types.GetUnbondingKeySuffix(denom, delegator)

	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		if len(key) < len(suffix) {
			continue // Skip keys that are shorter than the suffix
		}

		if !bytes.HasSuffix(key, suffix) {
			continue // Skip keys that don't have the desired suffix
		}

		_, time, _, err := types.ParseUndelegationKey(key)
		if err != nil {
			return nil, err
		}

		// Process and append item to the results
		item := types.UnbondingDelegation{
			CompletionTime: time,
		}
		items = append(items, item)
	}

	return items, nil
}
