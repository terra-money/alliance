package keeper

import (
	"alliance/x/alliance/types"
	"fmt"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetAllianceValidator(ctx sdk.Context, valAddr sdk.ValAddress) (types.AllianceValidator, error) {
	val, found := k.stakingKeeper.GetValidator(ctx, valAddr)
	if !found {
		return types.AllianceValidator{}, fmt.Errorf("validator with address %s does not exist", valAddr.String())
	}
	valInfo, found := k.getAllianceValidatorInfo(ctx, valAddr)
	if !found {
		valInfo = k.createAllianceValidatorInfo(ctx, valAddr)
	}
	return types.AllianceValidator{
		&val,
		&valInfo,
	}, nil
}

func (k Keeper) getAllianceValidatorInfo(ctx sdk.Context, valAddr sdk.ValAddress) (types.AllianceValidatorInfo, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetAllianceValidatorInfoKey(valAddr)
	vb := store.Get(key)
	var info types.AllianceValidatorInfo
	if vb == nil {
		return info, false
	} else {
		k.cdc.MustUnmarshal(vb, &info)
		return info, true
	}
}

func (k Keeper) createAllianceValidatorInfo(ctx sdk.Context, valAddr sdk.ValAddress) (val types.AllianceValidatorInfo) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetAllianceValidatorInfoKey(valAddr)
	val = types.NewAllianceValidatorInfo()
	vb := k.cdc.MustMarshal(&val)
	store.Set(key, vb)
	return val
}

func (k Keeper) IterateAllianceValidatorInfo(ctx sdk.Context) storetypes.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.ValidatorInfoKey)
}
