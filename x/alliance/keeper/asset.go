package keeper

import (
	"alliance/x/alliance/types"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"math"
)

func (k Keeper) UpdateAllianceAsset(ctx sdk.Context, newAsset types.AllianceAsset) error {
	asset, found := k.GetAssetByDenom(ctx, newAsset.Denom)
	if !found {
		return types.ErrUnknownAsset
	}

	valIter := k.IterateAllianceValidatorInfo(ctx)
	defer valIter.Close()
	for ; valIter.Valid(); valIter.Next() {
		valAddr := types.ParseAllianceValidatorKey(valIter.Key())
		validator, err := k.GetAllianceValidator(ctx, valAddr)
		if err != nil {
			return err
		}
		_, err = k.ClaimValidatorRewards(ctx, validator)
		if err != nil {
			return err
		}
		k.SetRewardRatesChangeSnapshot(ctx, asset, validator)
	}

	asset.TakeRate = newAsset.TakeRate
	asset.RewardWeight = newAsset.RewardWeight
	k.SetAsset(ctx, asset)
	return nil
}

func (k Keeper) RebalanceHook(ctx sdk.Context) error {
	if k.ConsumeAssetRebalanceEvent(ctx) {
		return k.RebalanceBondTokenWeights(ctx)
	}
	return nil
}

func (k Keeper) RebalanceBondTokenWeights(ctx sdk.Context) (err error) {
	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	allianceBondAmount := k.getAllianceBondedAmount(ctx, moduleAddr)

	nativeBondAmount := k.stakingKeeper.TotalBondedTokens(ctx).Sub(allianceBondAmount)
	bondDenom := k.stakingKeeper.BondDenom(ctx)

	assets := k.GetAllAssets(ctx)
	unbondedValidatorShares := sdk.NewDecCoins()
	var bondedValidators []types.AllianceValidator

	iter := k.IterateAllianceValidatorInfo(ctx)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		valAddr := types.ParseAllianceValidatorKey(iter.Key())
		validator, err := k.GetAllianceValidator(ctx, valAddr)
		if err != nil {
			return err
		}
		if validator.IsBonded() {
			bondedValidators = append(bondedValidators, validator)
		} else {
			unbondedValidatorShares = unbondedValidatorShares.Add(validator.ValidatorShares...)
		}
	}

	for _, validator := range bondedValidators {
		currentBondedAmount := sdk.NewDec(0)
		delegation, found := k.stakingKeeper.GetDelegation(ctx, moduleAddr, validator.GetOperator())
		if found {
			currentBondedAmount = validator.TokensFromShares(delegation.GetShares())
		}

		expectedBondAmount := sdk.ZeroDec()
		for _, asset := range assets {
			if ctx.BlockTime().Before(asset.RewardStartTime) {
				continue
			}
			valShares := validator.ValidatorSharesWithDenom(asset.Denom)
			expectedBondAmountForAsset := asset.RewardWeight.MulInt(nativeBondAmount).TruncateInt()

			// Accumulate expected tokens staked by adding up all expected tokens from each alliance asset
			if valShares.IsPositive() {
				expectedBondAmount = expectedBondAmount.Add(valShares.MulInt(expectedBondAmountForAsset).Quo(asset.TotalValidatorShares.Sub(unbondedValidatorShares.AmountOf(asset.Denom))))
			}
		}
		if expectedBondAmount.GT(currentBondedAmount) {
			// add
			bondAmount := expectedBondAmount.Sub(currentBondedAmount).RoundInt()
			err = k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(bondDenom, bondAmount)))
			if err != nil {
				return nil
			}
			_, err = k.stakingKeeper.Delegate(ctx, moduleAddr, bondAmount, stakingtypes.Unbonded, *validator.Validator, true)
			if err != nil {
				return err
			}
		} else if expectedBondAmount.LT(currentBondedAmount) {
			// sub
			unbondAmount := currentBondedAmount.Sub(expectedBondAmount).RoundInt()
			sharesToUnbond, err := k.stakingKeeper.ValidateUnbondAmount(ctx, moduleAddr, validator.GetOperator(), unbondAmount)
			if err != nil {
				return err
			}
			tokensToBurn, err := k.stakingKeeper.Unbond(ctx, moduleAddr, validator.GetOperator(), sharesToUnbond)
			if err != nil {
				return err
			}
			err = k.bankKeeper.BurnCoins(ctx, stakingtypes.BondedPoolName, sdk.NewCoins(sdk.NewCoin(bondDenom, tokensToBurn)))
			if err != nil {
				return err
			}

		}
	}
	return nil
}

func (k Keeper) SlashValidator(ctx sdk.Context, valAddr sdk.ValAddress, fraction sdk.Dec) error {
	val, err := k.GetAllianceValidator(ctx, valAddr)
	if err != nil {
		return err
	}
	slashedValidatorShares := sdk.NewDecCoins()
	for _, share := range val.ValidatorShares {
		sharesToSlash := share.Amount.Mul(fraction)
		slashedValidatorShares = append(slashedValidatorShares, sdk.NewDecCoinFromDec(share.Denom, share.Amount.Sub(sharesToSlash)))
		asset, found := k.GetAssetByDenom(ctx, share.Denom)
		if !found {
			return types.ErrUnknownAsset
		}
		asset.TotalValidatorShares = asset.TotalValidatorShares.Sub(sharesToSlash)
		k.SetAsset(ctx, asset)
	}
	val.ValidatorShares = slashedValidatorShares
	k.SetValidator(ctx, val)

	err = k.SlashRedelegations(ctx, valAddr, fraction)
	if err != nil {
		return err
	}

	err = k.SlashUndelegations(ctx, valAddr, fraction)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) SlashRedelegations(ctx sdk.Context, valAddr sdk.ValAddress, fraction sdk.Dec) error {
	store := ctx.KVStore(k.storeKey)
	// Slash all immature re-delegations
	redelegationIterator := k.IterateRedelegationsBySrcValidator(ctx, valAddr)
	for ; redelegationIterator.Valid(); redelegationIterator.Next() {
		redelegationKey, completion, err := types.ParseRedelegationIndexForRedelegationKey(redelegationIterator.Key())
		if err != nil {
			return err
		}
		// Skip if redelegation is already mature
		if completion.Before(ctx.BlockTime()) {
			continue
		}
		b := store.Get(redelegationKey)
		var redelegation types.Redelegation
		k.cdc.MustUnmarshal(b, &redelegation)

		delAddr, err := sdk.AccAddressFromBech32(redelegation.DelegatorAddress)
		if err != nil {
			return err
		}
		dstValAddr, err := sdk.ValAddressFromBech32(redelegation.DstValidatorAddress)
		if err != nil {
			return err
		}
		dstVal, err := k.GetAllianceValidator(ctx, dstValAddr)
		if err != nil {
			return err
		}

		delegation, found := k.GetDelegation(ctx, delAddr, dstVal, redelegation.Balance.Denom)
		if !found {
			continue
		}

		asset, found := k.GetAssetByDenom(ctx, redelegation.Balance.Denom)
		if !found {
			continue
		}

		// Slash delegation shares
		tokensToSlash := fraction.MulInt(redelegation.Balance.Amount).TruncateInt()
		sharesToSlash, err := k.ValidateDelegatedAmount(delegation, sdk.NewCoin(redelegation.Balance.Denom, tokensToSlash), dstVal, asset)
		if err != nil {
			return err
		}
		dstVal.TotalDelegatorShares = sdk.NewDecCoins(dstVal.TotalDelegatorShares...).Sub(sdk.NewDecCoins(sdk.NewDecCoinFromDec(asset.Denom, sharesToSlash)))
		k.SetValidator(ctx, dstVal)

		delegation.Shares = delegation.Shares.Sub(sharesToSlash)
		k.SetDelegation(ctx, delAddr, dstVal, asset.Denom, delegation)
	}
	return nil
}

func (k Keeper) SlashUndelegations(ctx sdk.Context, valAddr sdk.ValAddress, fraction sdk.Dec) error {
	store := ctx.KVStore(k.storeKey)
	// Slash all immature re-delegations
	undelegationIterator := k.IterateUndelegationsBySrcValidator(ctx, valAddr)
	for ; undelegationIterator.Valid(); undelegationIterator.Next() {
		undelegationKey, completion, err := types.ParseUnbondingIndexKeyToUndelegationKey(undelegationIterator.Key())
		if err != nil {
			return err
		}
		// Skip if undelegation is already mature
		if completion.Before(ctx.BlockTime()) {
			continue
		}
		b := store.Get(undelegationKey)
		var undelegations types.QueuedUndelegation
		k.cdc.MustUnmarshal(b, &undelegations)

		// Slash undelegations by sending slashed tokens to fee pool
		for _, entry := range undelegations.Entries {
			tokensToSlash := fraction.MulInt(entry.Balance.Amount).TruncateInt()
			entry.Balance = sdk.NewCoin(entry.Balance.Denom, entry.Balance.Amount.Sub(tokensToSlash))
			coinToSlash := sdk.NewCoin(entry.Balance.Denom, tokensToSlash)
			err = k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, authtypes.FeeCollectorName, sdk.NewCoins(coinToSlash))
			if err != nil {
				return err
			}
		}
		b = k.cdc.MustMarshal(&undelegations)
		store.Set(undelegationKey, b)
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

func (k Keeper) QueueAssetRebalanceEvent(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	key := types.AssetRebalanceQueueKey
	store.Set(key, []byte{0x00})
}

func (k Keeper) ConsumeAssetRebalanceEvent(ctx sdk.Context) bool {
	store := ctx.KVStore(k.storeKey)
	key := types.AssetRebalanceQueueKey
	b := store.Get(key)
	if b == nil {
		return false
	}
	store.Delete(key)
	return true
}

func (k Keeper) SetRewardRatesChangeSnapshot(ctx sdk.Context, asset types.AllianceAsset, val types.AllianceValidator) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetRewardRateChangeSnapshotKey(asset.Denom, val.GetOperator(), uint64(ctx.BlockHeight()))
	snapshot := types.NewRewardRateChangeSnapshot(asset, val)
	b := k.cdc.MustMarshal(&snapshot)
	store.Set(key, b)
}

func (k Keeper) IterateRewardRatesChangeSnapshot(ctx sdk.Context, denom string, valAddr sdk.ValAddress, lastClaimHeight uint64) store.Iterator {
	store := ctx.KVStore(k.storeKey)
	key := types.GetRewardRateChangeSnapshotKey(denom, valAddr, lastClaimHeight)
	end := types.GetRewardRateChangeSnapshotKey(denom, valAddr, math.MaxUint64)
	return store.Iterator(key, end)
}
