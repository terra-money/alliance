package keeper

import (
	"alliance/x/alliance/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) UpdateAllianceAsset(ctx sdk.Context, newAsset types.AllianceAsset) error {
	prevAsset, found := k.GetAssetByDenom(ctx, newAsset.Denom)
	if !found {
		return types.ErrUnknownAsset
	}
	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)

	if prevAsset.RewardWeight != newAsset.RewardWeight {
		iter := k.GetAllValidators(ctx)
		for ; iter.Valid(); iter.Next() {
			b := iter.Value()
			var aVal types.Validator
			k.cdc.MustUnmarshal(b, &aVal)
			if !aVal.TotalTokensWithAsset(prevAsset).IsPositive() {
				continue
			}
			valAddr, _ := sdk.ValAddressFromBech32(aVal.ValidatorAddress)
			val, _ := k.stakingKeeper.GetValidator(ctx, valAddr)
			delegation, found := k.stakingKeeper.GetDelegation(ctx, moduleAddr, valAddr)
			if !found {
				return types.ErrZeroDelegations
			}
			currentTokens := types.ConvertNewShareToToken(val.Tokens, val.DelegatorShares, delegation.Shares)
			expectedTokens := newAsset.RewardWeight.MulInt(prevAsset.TotalTokens).TruncateInt()
			if currentTokens.GT(expectedTokens) {
				tokensToRemove := currentTokens.Sub(expectedTokens)
				shares, err := k.stakingKeeper.ValidateUnbondAmount(ctx, moduleAddr, valAddr, tokensToRemove)
				if err != nil {
					return err
				}
				_, err = k.stakingKeeper.Unbond(ctx, moduleAddr, valAddr, shares)
				if err != nil {
					return err
				}
			} else {
				tokensToAdd := expectedTokens.Sub(currentTokens)
				err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(k.stakingKeeper.BondDenom(ctx), tokensToAdd)))
				if err != nil {
					return err
				}
				_, err = k.stakingKeeper.Delegate(ctx, moduleAddr, tokensToAdd, stakingtypes.Unbonded, val, true)
			}
		}
	}

	// Add a snapshot to help with rewards calculation

	// Only allow updating of certain values
	prevAsset.TakeRate = newAsset.TakeRate
	prevAsset.RewardWeight = newAsset.RewardWeight

	k.SetAsset(ctx, prevAsset)
	return nil
}

// SetAsset Does not check if the asset already exists and overwrites it
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

func (k Keeper) GetAssetByDenom(ctx sdk.Context, denom string) (asset types.AllianceAsset, found bool) {
	store := ctx.KVStore(k.storeKey)
	assetKey := types.GetAssetKey(denom)
	b := store.Get(assetKey)

	if b == nil {
		return asset, false
	}

	k.cdc.MustUnmarshal(b, &asset)
	return asset, true
}

func (k Keeper) DeleteAsset(ctx sdk.Context, denom string) {
	store := ctx.KVStore(k.storeKey)
	assetKey := types.GetAssetKey(denom)
	store.Delete(assetKey)
}
