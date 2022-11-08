package keeper

import (
	"alliance/x/alliance/types"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"math"
)

// UpdateAllianceAsset updates the alliance asset with new params
// Also saves a snapshot whenever rewards weight changes to make sure delegation reward calculation has reference to
// historical reward rates
func (k Keeper) UpdateAllianceAsset(ctx sdk.Context, newAsset types.AllianceAsset) error {
	asset, found := k.GetAssetByDenom(ctx, newAsset.Denom)
	if !found {
		return types.ErrUnknownAsset
	}

	// Only add a snapshot if reward weight changes
	if !newAsset.RewardWeight.Equal(asset.RewardWeight) {
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
	}

	// Make sure only the take rate and reward weight can be updated
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

// RebalanceBondTokenWeights uses asset reward weights to calculate the expected amount of staking token that has to be
// minted / burned to maintain the right ratio
// It iterates all validators and calculates the expected staked amount based on delegations and delegates/undelegates
// the difference.
func (k Keeper) RebalanceBondTokenWeights(ctx sdk.Context) (err error) {
	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	allianceBondAmount := k.getAllianceBondedAmount(ctx, moduleAddr)

	nativeBondAmount := k.stakingKeeper.TotalBondedTokens(ctx).Sub(allianceBondAmount)
	bondDenom := k.stakingKeeper.BondDenom(ctx)

	assets := k.GetAllAssets(ctx)
	unbondedValidatorShares := sdk.NewDecCoins()
	var bondedValidators []types.AllianceValidator

	// Iterate through all alliance validators to remove those that are unbonded.
	// Unbonded validators will be ignored when rebalancing.
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
			// Ignores assets that were recently added to prevent a small set of stakers from owning too much of the
			// voting power
			if ctx.BlockTime().Before(asset.RewardStartTime) {
				continue
			}
			valShares := validator.ValidatorSharesWithDenom(asset.Denom)
			expectedBondAmountForAsset := asset.RewardWeight.MulInt(nativeBondAmount)

			bondedValidatorShares := asset.TotalValidatorShares.Sub(unbondedValidatorShares.AmountOf(asset.Denom))
			if valShares.IsPositive() && bondedValidatorShares.IsPositive() {
				expectedBondAmount = expectedBondAmount.Add(valShares.Mul(expectedBondAmountForAsset).Quo(bondedValidatorShares))
			}
		}
		if expectedBondAmount.GT(currentBondedAmount) {
			// delegate more tokens to increase the weight
			bondAmount := expectedBondAmount.Sub(currentBondedAmount).TruncateInt()
			err = k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(bondDenom, bondAmount)))
			if err != nil {
				return nil
			}
			_, err = k.stakingKeeper.Delegate(ctx, moduleAddr, bondAmount, stakingtypes.Unbonded, *validator.Validator, true)
			if err != nil {
				return err
			}
		} else if expectedBondAmount.LT(currentBondedAmount) {
			// undelegate more tokens to reduce the weight
			unbondAmount := currentBondedAmount.Sub(expectedBondAmount).TruncateInt()
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
