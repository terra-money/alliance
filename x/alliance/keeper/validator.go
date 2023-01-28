package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/terra-money/alliance/x/alliance/types"
)

func (k Keeper) GetAllianceValidator(ctx sdk.Context, valAddr sdk.ValAddress) (types.AllianceValidator, error) {
	val, found := k.stakingKeeper.GetValidator(ctx, valAddr)
	if !found {
		return types.AllianceValidator{}, fmt.Errorf("validator with address %s does not exist", valAddr.String())
	}
	valInfo, found := k.GetAllianceValidatorInfo(ctx, valAddr)
	if !found {
		valInfo = k.createAllianceValidatorInfo(ctx, valAddr)
	}
	return types.AllianceValidator{
		Validator:             &val,
		AllianceValidatorInfo: &valInfo,
	}, nil
}

func (k Keeper) GetAllianceValidatorInfo(ctx sdk.Context, valAddr sdk.ValAddress) (types.AllianceValidatorInfo, bool) {
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

func (k Keeper) IterateAllianceValidatorInfo(ctx sdk.Context, cb func(valAddr sdk.ValAddress, info types.AllianceValidatorInfo) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.ValidatorInfoKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var info types.AllianceValidatorInfo
		b := iter.Value()
		k.cdc.MustUnmarshal(b, &info)
		valAddr := types.ParseAllianceValidatorKey(iter.Key())
		if cb(valAddr, info) {
			return
		}
	}
}

func (k Keeper) GetAllAllianceValidatorInfo(ctx sdk.Context) []types.AllianceValidatorInfo {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.ValidatorInfoKey)
	defer iter.Close()
	var infos []types.AllianceValidatorInfo
	for ; iter.Valid(); iter.Next() {
		b := iter.Value()
		var info types.AllianceValidatorInfo
		k.cdc.UnmarshalInterface(b, &info)
		infos = append(infos, info)
	}
	return infos
}
