package keeper

import (
	"alliance/x/alliance/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetAllValidators(ctx sdk.Context) storetypes.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.ValidatorKey)
}
