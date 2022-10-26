package keeper

import (
	"alliance/x/alliance/types"
	cosmosmath "cosmossdk.io/math"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (k Keeper) UpdateAllianceAsset(ctx sdk.Context, newAsset types.AllianceAsset) error {
	//prevAsset, found := k.GetAssetByDenom(ctx, newAsset.Denom)
	//if !found {
	//	return types.ErrUnknownAsset
	//}
	//moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	//
	//if prevAsset.RewardWeight != newAsset.RewardWeight {
	//	iter := k.IterateAllianceValidatorInfo(ctx)
	//	for ; iter.Valid(); iter.Next() {
	//		b := iter.Value()
	//		var aVal types.Validator
	//		k.cdc.MustUnmarshal(b, &aVal)
	//		if !aVal.TotalTokensWithAsset(prevAsset).IsPositive() {
	//			continue
	//		}
	//		valAddr, _ := sdk.ValAddressFromBech32(aVal.ValidatorAddress)
	//		val, _ := k.stakingKeeper.GetValidator(ctx, valAddr)
	//		delegation, found := k.stakingKeeper.GetDelegation(ctx, moduleAddr, valAddr)
	//		if !found {
	//			return types.ErrZeroDelegations
	//		}
	//		currentTokens := types.ConvertNewShareToToken(val.Tokens, val.DelegatorShares, delegation.Shares)
	//		expectedTokens := newAsset.RewardWeight.MulInt(prevAsset.TotalTokens).TruncateInt()
	//		if currentTokens.GT(expectedTokens) {
	//			tokensToRemove := currentTokens.Sub(expectedTokens)
	//			shares, err := k.stakingKeeper.ValidateUnbondAmount(ctx, moduleAddr, valAddr, tokensToRemove)
	//			if err != nil {
	//				return err
	//			}
	//			_, err = k.stakingKeeper.Unbond(ctx, moduleAddr, valAddr, shares)
	//			if err != nil {
	//				return err
	//			}
	//		} else {
	//			tokensToAdd := expectedTokens.Sub(currentTokens)
	//			err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(k.stakingKeeper.BondDenom(ctx), tokensToAdd)))
	//			if err != nil {
	//				return err
	//			}
	//			_, err = k.stakingKeeper.Delegate(ctx, moduleAddr, tokensToAdd, stakingtypes.Unbonded, val, true)
	//		}
	//	}
	//}
	//
	//// Add a snapshot to help with rewards calculation
	//
	//// Only allow updating of certain values
	//prevAsset.TakeRate = newAsset.TakeRate
	//prevAsset.RewardWeight = newAsset.RewardWeight
	//
	//k.SetAsset(ctx, prevAsset)
	return nil
}

func (k Keeper) RebalanceInternalStakeWeights(ctx sdk.Context) error {
	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	allianceStakedAmount := k.stakingKeeper.GetDelegatorBonded(ctx, moduleAddr)
	nativeStakedAmount := k.stakingKeeper.TotalBondedTokens(ctx).Sub(allianceStakedAmount)
	bondDenom := k.stakingKeeper.BondDenom(ctx)

	assets := k.GetAllAssets(ctx)

	for iter := k.IterateAllianceValidatorInfo(ctx); iter.Valid(); iter.Next() {
		valAddr := types.ParseAllianceValidatorKey(iter.Key())
		validator, err := k.GetAllianceValidator(ctx, valAddr)
		if err != nil {
			return err
		}
		actualStakeAmount := sdk.NewDec(0)
		delegation, found := k.stakingKeeper.GetDelegation(ctx, moduleAddr, valAddr)
		if found {
			actualStakeAmount = validator.TokensFromShares(delegation.GetShares())
		}

		expectedStakeAmount := sdk.ZeroDec()
		for _, asset := range assets {
			valShares := validator.ValidatorSharesWithDenom(asset.Denom)
			totalStakeTokens := asset.RewardWeight.MulInt(nativeStakedAmount).TruncateInt()

			// Accumulate expected tokens staked by adding up all expected tokens from each alliance asset
			if valShares.IsPositive() {
				expectedStakeAmount = expectedStakeAmount.Add(valShares.MulInt(totalStakeTokens).Quo(asset.TotalValidatorShares))
			}

			// Update total staked tokens if we are handling this alliance token for the first time
			if asset.TotalStakeTokens.IsZero() && asset.TotalValidatorShares.IsPositive() {
				asset.TotalStakeTokens = totalStakeTokens
				k.SetAsset(ctx, *asset)
			}
		}
		if expectedStakeAmount.GT(actualStakeAmount) {
			// add
			bondAmount := expectedStakeAmount.Sub(actualStakeAmount).TruncateInt()
			err = k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(bondDenom, bondAmount)))
			if err != nil {
				return nil
			}
			_, err = k.stakingKeeper.Delegate(ctx, moduleAddr, bondAmount, stakingtypes.Unbonded, *validator.Validator, true)
			if err != nil {
				return err
			}
		} else if expectedStakeAmount.LT(actualStakeAmount) {
			// sub
			unbondAmount := actualStakeAmount.Sub(expectedStakeAmount)
			sharesToUnbond, err := k.stakingKeeper.ValidateUnbondAmount(ctx, moduleAddr, validator.GetOperator(), unbondAmount.TruncateInt())
			if err != nil {
				return err
			}
			tokensToBurn, err := k.stakingKeeper.Unbond(ctx, moduleAddr, validator.GetOperator(), sharesToUnbond)
			if err != nil {
				return err
			}
			k.bankKeeper.BurnCoins(ctx, stakingtypes.BondedPoolName, sdk.NewCoins(sdk.NewCoin(bondDenom, tokensToBurn)))
		}
	}
	return nil
}

// SetAsset Does not check if the asset already exists and overwrites it
func (k Keeper) SetAsset(ctx sdk.Context, asset types.AllianceAsset) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&asset)
	store.Set(types.GetAssetKey(asset.Denom), b)
}

func (k Keeper) GetAllAssets(ctx sdk.Context) (assets []*types.AllianceAsset) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.AssetKey)
	defer iter.Close()

	for iter.Valid() {
		var asset types.AllianceAsset
		b := iter.Value()
		k.cdc.MustUnmarshal(b, &asset)
		assets = append(assets, &asset)
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

func (k Keeper) QueueAssetRebalanceEvent(ctx sdk.Context, denom string) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetAssetRebalanceQueueKeyByDenom(denom)
	store.Set(key, []byte{0x00})
}

func (k Keeper) IterateAssetRebalanceQueue(ctx sdk.Context) storetypes.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.AssetRebalanceQueueKey)
}

func (k Keeper) DeleteAssetRebalanceEvent(ctx sdk.Context, denom string) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetAssetRebalanceQueueKeyByDenom(denom)
	store.Delete(key)
}

func (k Keeper) GetTotalStakedAmount(ctx sdk.Context) cosmosmath.Int {
	assets := k.GetAllAssets(ctx)
	amount := sdk.ZeroInt()
	for _, asset := range assets {
		amount = amount.Add(asset.TotalTokens)
	}
	return amount
}
