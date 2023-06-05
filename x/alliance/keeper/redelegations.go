package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/terra-money/alliance/x/alliance/types"
)

// This method retun all redelegations delegations for a given denom and delegator address.
func (k Keeper) GetRedelegations(
	ctx sdk.Context,
	denom string,
	delAddr sdk.AccAddress,
) (redelegationEntries []types.RedelegationEntry, err error) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetRedelegationsKeyByDelegatorAndDenom(delAddr, denom))

	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var redelegation types.Redelegation
		b := iter.Value()
		k.cdc.MustUnmarshal(b, &redelegation)
		completionTime := types.ParseRedelegationKeyForCompletionTime(iter.Key())

		redelegationEntry := types.RedelegationEntry{
			DelegatorAddress:    redelegation.DelegatorAddress,
			SrcValidatorAddress: redelegation.SrcValidatorAddress,
			DstValidatorAddress: redelegation.DstValidatorAddress,
			Balance:             redelegation.Balance,
			CompletionTime:      completionTime,
		}
		redelegationEntries = append(redelegationEntries, redelegationEntry)
	}

	return redelegationEntries, err
}
