package keeper

import (
	"context"
	storetypes "cosmossdk.io/store/types"
	"fmt"
	"github.com/cosmos/cosmos-sdk/runtime"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-money/alliance/x/alliance/types"
)

func (k Keeper) GetAllianceValidator(ctx context.Context, valAddr sdk.ValAddress) (types.AllianceValidator, error) {
	val, err := k.stakingKeeper.GetValidator(ctx, valAddr)
	if err != nil {
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

func (k Keeper) GetAllianceValidatorInfo(ctx context.Context, valAddr sdk.ValAddress) (types.AllianceValidatorInfo, bool) {
	store := k.storeService.OpenKVStore(ctx)
	key := types.GetAllianceValidatorInfoKey(valAddr)
	vb, err := store.Get(key)
	var info types.AllianceValidatorInfo
	if vb == nil || err != nil {
		return info, false
	}
	k.cdc.MustUnmarshal(vb, &info)
	return info, true
}

func (k Keeper) createAllianceValidatorInfo(ctx context.Context, valAddr sdk.ValAddress) (val types.AllianceValidatorInfo) {
	store := k.storeService.OpenKVStore(ctx)
	key := types.GetAllianceValidatorInfoKey(valAddr)
	val = types.NewAllianceValidatorInfo()
	vb := k.cdc.MustMarshal(&val)
	store.Set(key, vb)
	return val
}

func (k Keeper) IterateAllianceValidatorInfo(ctx context.Context, cb func(valAddr sdk.ValAddress, info types.AllianceValidatorInfo) (stop bool)) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iter := storetypes.KVStorePrefixIterator(store, types.ValidatorInfoKey)
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

func (k Keeper) GetAllAllianceValidatorInfo(ctx context.Context) []types.AllianceValidatorInfo {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iter := storetypes.KVStorePrefixIterator(store, types.ValidatorInfoKey)
	defer iter.Close()
	var infos []types.AllianceValidatorInfo
	for ; iter.Valid(); iter.Next() {
		b := iter.Value()
		var info types.AllianceValidatorInfo
		k.cdc.UnmarshalInterface(b, &info) //nolint:errcheck
		infos = append(infos, info)
	}
	return infos
}

func (k Keeper) DeleteValidatorInfo(ctx context.Context, valAddr sdk.ValAddress) {
	store := k.storeService.OpenKVStore(ctx)
	key := types.GetAllianceValidatorInfoKey(valAddr)
	store.Delete(key)
}
