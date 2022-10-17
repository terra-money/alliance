package keeper

import (
	"alliance/x/alliance/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// setAsset Does not check if the asset already exists and overwrites it
func (k Keeper) SetAsset(ctx sdk.Context, asset types.AllianceAsset) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&asset)
	store.Set(types.GetAssetKey(asset.Denom), b)
}

func (k Keeper) GetAllAssets(ctx sdk.Context) (assets []types.AllianceAsset) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.AssetKey)
	defer iter.Close()

	for iter.Valid() {
		var asset types.AllianceAsset
		b := iter.Value()
		k.cdc.MustUnmarshal(b, &asset)
		assets = append(assets, asset)
		iter.Next()
	}
	return assets
}

func (k Keeper) GetAssetByDenom(ctx sdk.Context, denom string) (asset types.AllianceAsset) {
	store := ctx.KVStore(k.storeKey)
	assetKey := types.GetAssetKey(denom)
	b := store.Get(assetKey)
	k.cdc.MustUnmarshal(b, &asset)
	return
}

func (k Keeper) AddAsset() {
	panic("implement me")
}

func (k Keeper) RemoveAsset() {
	panic("implement me")
}

func (k Keeper) UpdateAsset() {
	panic("implement me")
}
