package keeper

import (
	"time"

	"github.com/terra-money/alliance/x/alliance/types"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Delegate is the entry point for delegators to delegate alliance assets to validators
// Voting power is not immediately accured to the delegators in this method and a flag is set to rebalance voting power
// at the end of the block. This improves performance since rebalancing only needs to happen once regardless of how many
// delegations are made in a single block
func (k Keeper) Delegate(ctx sdk.Context, delAddr sdk.AccAddress, validator types.AllianceValidator, coin sdk.Coin) (*sdk.Dec, error) {
	// Check if asset is whitelisted as an alliance asset
	asset, found := k.GetAssetByDenom(ctx, coin.Denom)
	if !found {
		return nil, status.Errorf(codes.NotFound, "asset with denom: %s does not exist in alliance whitelist", coin.Denom)
	}

	// Check and send delegated tokens into the alliance module address
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, delAddr, types.ModuleName, sdk.NewCoins(coin))
	if err != nil {
		return nil, err
	}

	// Claim rewards before adding more to a previous delegation
	_, found = k.GetDelegation(ctx, delAddr, validator, coin.Denom)
	if found {
		_, err = k.ClaimDelegationRewards(ctx, delAddr, validator, coin.Denom)
		if err != nil {
			return nil, err
		}
	}

	// Create or update a delegation
	_, newDelegationShares := k.upsertDelegationWithNewTokens(ctx, delAddr, validator, coin, asset)

	// Update ownership shares in validator and asset
	newValidatorShares := types.GetValidatorShares(asset, coin.Amount)
	asset.TotalTokens = asset.TotalTokens.Add(coin.Amount)
	asset.TotalValidatorShares = asset.TotalValidatorShares.Add(newValidatorShares)
	k.SetAsset(ctx, asset)
	k.updateValidatorShares(ctx, validator,
		sdk.NewDecCoins(sdk.NewDecCoinFromDec(coin.Denom, newDelegationShares)),
		sdk.NewDecCoins(sdk.NewDecCoinFromDec(coin.Denom, newValidatorShares)),
		true,
	)
	k.QueueAssetRebalanceEvent(ctx)
	return &newValidatorShares, nil
}

// Redelegate from one validator to another
func (k Keeper) Redelegate(ctx sdk.Context, delAddr sdk.AccAddress, srcVal types.AllianceValidator, dstVal types.AllianceValidator, coin sdk.Coin) (*time.Time, error) {
	if srcVal.Validator.Equal(dstVal.Validator) {
		return nil, status.Errorf(codes.InvalidArgument, "Cannot redelegate to the same validator")
	}

	asset, found := k.GetAssetByDenom(ctx, coin.Denom)
	if !found {
		return nil, status.Errorf(codes.NotFound, "Asset with denom: %s does not exist", coin.Denom)
	}

	_, found = k.GetDelegation(ctx, delAddr, srcVal, coin.Denom)
	if !found {
		return nil, stakingtypes.ErrNoDelegatorForAddress
	}
	_, err := k.ClaimDelegationRewards(ctx, delAddr, srcVal, coin.Denom)
	if err != nil {
		return nil, err
	}
	// re-query delegation since it was updated in `ClaimDelegationRewards`
	srcDelegation, _ := k.GetDelegation(ctx, delAddr, srcVal, coin.Denom)

	_, found = k.GetDelegation(ctx, delAddr, dstVal, coin.Denom)
	if found {
		_, err = k.ClaimDelegationRewards(ctx, delAddr, dstVal, coin.Denom)
		if err != nil {
			return nil, err
		}
	}

	delegationSharesToRemove, err := k.ValidateDelegatedAmount(srcDelegation, coin, srcVal, asset)
	if err != nil {
		return nil, err
	}

	// Now we have shares we want to re-delegate, we re-calculate how many tokens to actually re-delegate
	// Directly using the input amount can result in un-delegating more tokens than expected due to rounding issues
	coinsToRedelegate := types.GetDelegationTokensWithShares(delegationSharesToRemove, srcVal, asset)
	if coin.Amount.GT(coinsToRedelegate.Amount) {
		return nil, types.ErrInsufficientTokens.Wrapf("wanted %s but have %s", coin.Amount, coinsToRedelegate.Amount)
	}

	// Prevents transitive re-delegations
	// e.g. if a redelegation from A -> B is made before another request from B -> C
	// the latter is blocked until the first redelegation is mature (time > unbonding time)
	if k.HasRedelegation(ctx, delAddr, srcVal.GetOperator(), coin.Denom) {
		return nil, stakingtypes.ErrTransitiveRedelegation
	}

	completionTime := ctx.BlockHeader().Time.Add(k.stakingKeeper.UnbondingTime(ctx))
	changedValidatorShares := types.GetValidatorShares(asset, coin.Amount)

	// Remove tokens and shares from src validator
	k.reduceDelegationShares(ctx, delAddr, srcVal, coin, delegationSharesToRemove, srcDelegation)
	k.updateValidatorShares(
		ctx,
		srcVal,
		sdk.NewDecCoins(sdk.NewDecCoinFromDec(coin.Denom, delegationSharesToRemove)),
		sdk.NewDecCoins(sdk.NewDecCoinFromDec(coin.Denom, changedValidatorShares)),
		false,
	)

	// Add tokens and shares to dst validator
	_, newDelegationShares := k.upsertDelegationWithNewTokens(ctx, delAddr, dstVal, coin, asset)
	k.updateValidatorShares(
		ctx,
		dstVal,
		sdk.NewDecCoins(sdk.NewDecCoinFromDec(coin.Denom, newDelegationShares)),
		sdk.NewDecCoins(sdk.NewDecCoinFromDec(coin.Denom, changedValidatorShares)),
		true,
	)

	k.addRedelegation(ctx, delAddr, srcVal.GetOperator(), dstVal.GetOperator(), coin, completionTime)

	k.QueueAssetRebalanceEvent(ctx)

	return &completionTime, nil
}

// Undelegate from a validator
// Staked tokens are only distributed to the delegator after the unbonding period
func (k Keeper) Undelegate(ctx sdk.Context, delAddr sdk.AccAddress, validator types.AllianceValidator, coin sdk.Coin) (*time.Time, error) {
	asset, found := k.GetAssetByDenom(ctx, coin.Denom)
	if !found {
		return nil, status.Errorf(codes.NotFound, "Asset with denom: %s does not exist", coin.Denom)
	}

	_, ok := k.GetDelegation(ctx, delAddr, validator, coin.Denom)
	if !ok {
		return nil, stakingtypes.ErrNoDelegatorForAddress
	}
	// Claim delegation rewards first
	_, err := k.ClaimDelegationRewards(ctx, delAddr, validator, coin.Denom)
	if err != nil {
		return nil, err
	}

	// Delegation is queried again since it might have been modified when claiming delegation rewards
	delegation, _ := k.GetDelegation(ctx, delAddr, validator, coin.Denom)

	// Calculate how much delegation shares to be undelegated taking into account rounding issues
	delegationSharesToUndelegate, err := k.ValidateDelegatedAmount(delegation, coin, validator, asset)
	if err != nil {
		return nil, err
	}

	// Now we have shares we want to un-delegate, we re-calculate how many tokens to actually un-delegate
	// Directly using the input amount can result in un-delegating more tokens than expected due to rounding issues
	coinsToUndelegate := types.GetDelegationTokensWithShares(delegationSharesToUndelegate, validator, asset)
	if coin.Amount.GT(coinsToUndelegate.Amount) {
		return nil, types.ErrInsufficientTokens.Wrapf("wanted %s but have %s", coin.Amount, coinsToUndelegate.Amount)
	}
	validatorSharesToRemove := types.GetValidatorShares(asset, coin.Amount)

	// Remove tokens and shares from the alliance asset
	asset.TotalTokens = asset.TotalTokens.Sub(coin.Amount)
	asset.TotalValidatorShares = asset.TotalValidatorShares.Sub(validatorSharesToRemove)
	k.SetAsset(ctx, asset)

	// Remove shares from the delegation
	k.reduceDelegationShares(ctx, delAddr, validator, coin, delegationSharesToUndelegate, delegation)

	// Remove tokens and shares from src validator
	k.updateValidatorShares(
		ctx,
		validator,
		sdk.NewDecCoins(sdk.NewDecCoinFromDec(coin.Denom, delegationSharesToUndelegate)),
		sdk.NewDecCoins(sdk.NewDecCoinFromDec(coin.Denom, validatorSharesToRemove)),
		false,
	)

	// When there are no more tokens recorded in the asset, clear all share records that might remain
	// from rounding errors to prevent dust amounts from staying in the stores
	if asset.TotalTokens.IsZero() {
		k.ResetAssetAndValidators(ctx, asset)
	}

	// Queue undelegation messages to distribute tokens after undelegation completes in the future
	completionTime := k.queueUndelegation(ctx, delAddr, validator.GetOperator(), coin)
	k.QueueAssetRebalanceEvent(ctx)
	return &completionTime, nil
}

// CompleteRedelegations Go through the re-delegations queue and remove all that have passed the completion time
func (k Keeper) CompleteRedelegations(ctx sdk.Context) int {
	store := ctx.KVStore(k.storeKey)
	iter := store.Iterator(types.RedelegationQueueKey, types.GetRedelegationQueueKey(ctx.BlockTime()))
	defer iter.Close()
	deleted := 0
	for ; iter.Valid(); iter.Next() {
		completion := types.ParseRedelegationQueueKey(iter.Key())
		var queued types.QueuedRedelegation
		k.cdc.MustUnmarshal(iter.Value(), &queued)
		for _, redel := range queued.Entries {
			k.DeleteRedelegation(ctx, *redel, completion)
			deleted++
		}
		store.Delete(iter.Key())
	}
	return deleted
}

// CompleteUndelegations Go through all queued undelegations and send the tokens to the delegators
func (k Keeper) CompleteUndelegations(ctx sdk.Context) error {
	store := ctx.KVStore(k.storeKey)
	iter := k.IterateUndelegationsByCompletionTime(ctx, ctx.BlockTime())
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var queued types.QueuedUndelegation
		completionTime, err := types.ParseUndelegationQueueKeyForCompletionTime(iter.Key())
		if err != nil {
			return err
		}
		k.cdc.MustUnmarshal(iter.Value(), &queued)
		for _, undel := range queued.Entries {
			delAddr, err := sdk.AccAddressFromBech32(undel.DelegatorAddress)
			if err != nil {
				return err
			}
			err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, delAddr, sdk.NewCoins(undel.Balance))
			if err != nil {
				return err
			}
			valAddr, err := sdk.ValAddressFromBech32(undel.ValidatorAddress)
			if err != nil {
				return err
			}
			indexKey := types.GetUnbondingIndexKey(valAddr, completionTime, undel.Balance.Denom, delAddr)
			store.Delete(indexKey)
		}
		store.Delete(iter.Key())
	}

	// Burn all "virtual" staking tokens in the module account
	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	coin := k.bankKeeper.GetBalance(ctx, moduleAddr, k.stakingKeeper.BondDenom(ctx))
	if !coin.IsZero() {
		err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(coin))
		if err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) GetDelegation(ctx sdk.Context, delAddr sdk.AccAddress, validator types.AllianceValidator, denom string) (d types.Delegation, found bool) {
	key := types.GetDelegationKey(delAddr, validator.GetOperator(), denom)
	b := ctx.KVStore(k.storeKey).Get(key)
	if b == nil {
		return d, false
	}
	k.cdc.MustUnmarshal(b, &d)
	return d, true
}

func (k Keeper) SetDelegation(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, denom string, del types.Delegation) {
	key := types.GetDelegationKey(delAddr, valAddr, denom)
	b := k.cdc.MustMarshal(&del)
	ctx.KVStore(k.storeKey).Set(key, b)
}

func (k Keeper) DeleteRedelegation(ctx sdk.Context, redel types.Redelegation, completion time.Time) {
	delAddr, err := sdk.AccAddressFromBech32(redel.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	dstValAddr, err := sdk.ValAddressFromBech32(redel.DstValidatorAddress)
	if err != nil {
		panic(err)
	}
	srcValAddr, err := sdk.ValAddressFromBech32(redel.SrcValidatorAddress)
	if err != nil {
		panic(err)
	}
	store := ctx.KVStore(k.storeKey)
	key := types.GetRedelegationKey(delAddr, redel.Balance.Denom, dstValAddr, completion)
	store.Delete(key)
	indexKey := types.GetRedelegationIndexKey(srcValAddr, completion, redel.Balance.Denom, dstValAddr, delAddr)
	store.Delete(indexKey)
}

func (k Keeper) IterateDelegations(ctx sdk.Context, cb func(d types.Delegation) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.DelegationKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var delegation types.Delegation
		k.cdc.MustUnmarshal(iter.Value(), &delegation)
		if cb(delegation) {
			break
		}
	}
}

func (k Keeper) HasRedelegation(ctx sdk.Context, delAddr sdk.AccAddress, dstVal sdk.ValAddress, denom string) bool {
	store := ctx.KVStore(k.storeKey)
	key := types.GetRedelegationsKey(delAddr, denom, dstVal)
	iter := sdk.KVStorePrefixIterator(store, key)
	defer iter.Close()
	return iter.Valid()
}

func (k Keeper) IterateRedelegations(ctx sdk.Context, cb func(redelegation types.Redelegation, completionTime time.Time) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.RedelegationKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var redelegation types.Redelegation
		b := iter.Value()
		k.cdc.MustUnmarshal(b, &redelegation)
		completionTime := types.ParseRedelegationKeyForCompletionTime(iter.Key())
		if cb(redelegation, completionTime) {
			return
		}
	}
}

func (k Keeper) IterateRedelegationsByDelegator(ctx sdk.Context, delAddr sdk.AccAddress) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	key := types.GetRedelegationsKeyByDelegator(delAddr)
	return sdk.KVStorePrefixIterator(store, key)
}

func (k Keeper) IterateRedelegationsBySrcValidator(ctx sdk.Context, srcValAddr sdk.ValAddress) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	prefix := types.GetRedelegationsIndexOrderedByValidatorKey(srcValAddr)
	return sdk.KVStorePrefixIterator(store, prefix)
}

func (k Keeper) IterateUndelegationsBySrcValidator(ctx sdk.Context, valAddr sdk.ValAddress) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	prefix := types.GetUndelegationsIndexOrderedByValidatorKey(valAddr)
	return sdk.KVStorePrefixIterator(store, prefix)
}

func (k Keeper) IterateUndelegationsByCompletionTime(ctx sdk.Context, completionTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator(types.UndelegationQueueKey, types.GetUndelegationQueueKeyByTime(completionTime))
}

func (k Keeper) IterateUndelegations(ctx sdk.Context, cb func(undelegation types.QueuedUndelegation, completionTime time.Time) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.UndelegationQueueKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var undelegation types.QueuedUndelegation
		b := iter.Value()
		k.cdc.MustUnmarshal(b, &undelegation)
		completionTime, _ := types.ParseUndelegationQueueKeyForCompletionTime(iter.Key())
		if cb(undelegation, completionTime) {
			return
		}
	}
}

func (k Keeper) SetValidator(ctx sdk.Context, val types.AllianceValidator) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetAllianceValidatorInfoKey(val.GetOperator())
	vb := k.cdc.MustMarshal(val.AllianceValidatorInfo)
	store.Set(key, vb)
}

func (k Keeper) SetValidatorInfo(ctx sdk.Context, valAddr sdk.ValAddress, val types.AllianceValidatorInfo) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetAllianceValidatorInfoKey(valAddr)
	vb := k.cdc.MustMarshal(&val)
	store.Set(key, vb)
}

// ValidateDelegatedAmount returns the amount of shares for a given coin that is staked
// Returns the number of shares that represents the amount of staked tokens that was requested
func (k Keeper) ValidateDelegatedAmount(delegation types.Delegation, coin sdk.Coin, val types.AllianceValidator, asset types.AllianceAsset) (shares sdk.Dec, err error) {
	delegationSharesToUpdate := types.GetDelegationSharesFromTokens(val, asset, coin.Amount)
	if delegation.Shares.LT(delegationSharesToUpdate.TruncateDec()) {
		return sdk.Dec{}, stakingtypes.ErrInsufficientShares
	}

	// Account for rounding in which shares for a full withdraw is slightly more or less than the number of shares recorded
	// Withdraw all in that case
	if delegation.Shares.Sub(delegationSharesToUpdate).Abs().LT(sdk.OneDec()) {
		delegationSharesToUpdate = delegation.Shares
	}

	return delegationSharesToUpdate, nil
}

// addRedelegation adds a redelegation entry to be used to prevent premature re-delegations
func (k Keeper) addRedelegation(ctx sdk.Context, delAddr sdk.AccAddress, srcVal sdk.ValAddress, dstVal sdk.ValAddress, coin sdk.Coin, completionTime time.Time) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetRedelegationKey(delAddr, coin.Denom, dstVal, completionTime)
	b := store.Get(key)
	var redelegation types.Redelegation
	if b == nil {
		redelegation = types.Redelegation{
			DelegatorAddress:    delAddr.String(),
			SrcValidatorAddress: srcVal.String(),
			DstValidatorAddress: dstVal.String(),
			Balance:             coin,
		}
	} else {
		k.cdc.MustUnmarshal(b, &redelegation)
		redelegation.Balance = redelegation.Balance.Add(coin)
	}
	b = k.cdc.MustMarshal(&redelegation)
	store.Set(key, b)

	// Add another entry as an index to retrieve redelegations by validator
	indexKey := types.GetRedelegationIndexKey(srcVal, completionTime, coin.Denom, dstVal, delAddr)
	store.Set(indexKey, []byte{})
	k.queueRedelegation(ctx, delAddr, srcVal, dstVal, coin, completionTime)
}

// queueRedelegation adds a redelegation queue object that is indexed by completion time
// This is used to track and clean-up redelegation events once they are mature
// Use addRedelegations instead
func (k Keeper) queueRedelegation(ctx sdk.Context, delAddr sdk.AccAddress, srcVal sdk.ValAddress, dstVal sdk.ValAddress, coin sdk.Coin, completionTime time.Time) {
	store := ctx.KVStore(k.storeKey)
	queueKey := types.GetRedelegationQueueKey(completionTime)
	b := store.Get(queueKey)
	var queuedDelegations types.QueuedRedelegation
	if b == nil {
		queuedDelegations = types.QueuedRedelegation{
			Entries: []*types.Redelegation{
				{
					DelegatorAddress:    delAddr.String(),
					SrcValidatorAddress: srcVal.String(),
					DstValidatorAddress: dstVal.String(),
					Balance:             coin,
				},
			},
		}
	} else {
		k.cdc.MustUnmarshal(b, &queuedDelegations)
		queuedDelegations.Entries = append(queuedDelegations.Entries, &types.Redelegation{
			DelegatorAddress:    delAddr.String(),
			SrcValidatorAddress: srcVal.String(),
			DstValidatorAddress: dstVal.String(),
			Balance:             coin,
		})
	}
	b = k.cdc.MustMarshal(&queuedDelegations)
	store.Set(queueKey, b)
}

// queueUndelegation adds a undelegation queue object that is indexed by completion time
// This is used to track and clean-up undelegation events once they are mature
func (k Keeper) queueUndelegation(ctx sdk.Context, delAddr sdk.AccAddress, val sdk.ValAddress, coin sdk.Coin) time.Time {
	store := ctx.KVStore(k.storeKey)
	completionTime := ctx.BlockTime().Add(k.stakingKeeper.UnbondingTime(ctx))
	queueKey := types.GetUndelegationQueueKey(completionTime, delAddr)
	b := store.Get(queueKey)
	var queue types.QueuedUndelegation
	if b == nil {
		queue = types.QueuedUndelegation{
			Entries: []*types.Undelegation{
				{
					DelegatorAddress: delAddr.String(),
					ValidatorAddress: val.String(),
					Balance:          coin,
				},
			},
		}
	} else {
		k.cdc.MustUnmarshal(b, &queue)
		queue.Entries = append(queue.Entries, &types.Undelegation{
			DelegatorAddress: delAddr.String(),
			ValidatorAddress: val.String(),
			Balance:          coin,
		})
	}
	k.setQueuedUndelegations(ctx, completionTime, delAddr, queue)
	k.setUnbondingIndexByVal(ctx, val, completionTime, delAddr, coin.Denom)
	return completionTime
}

func (k Keeper) setQueuedUndelegations(ctx sdk.Context, completionTime time.Time, delAddr sdk.AccAddress, queuedDelegations types.QueuedUndelegation) {
	store := ctx.KVStore(k.storeKey)
	queueKey := types.GetUndelegationQueueKey(completionTime, delAddr)
	b := k.cdc.MustMarshal(&queuedDelegations)
	store.Set(queueKey, b)
}

func (k Keeper) setUnbondingIndexByVal(ctx sdk.Context, valAddr sdk.ValAddress, completionTime time.Time, delAddr sdk.AccAddress, denom string) {
	store := ctx.KVStore(k.storeKey)
	indexKey := types.GetUnbondingIndexKey(valAddr, completionTime, denom, delAddr)
	store.Set(indexKey, []byte{})
}

func (k Keeper) upsertDelegationWithNewTokens(ctx sdk.Context, delAddr sdk.AccAddress, validator types.AllianceValidator, coin sdk.Coin, asset types.AllianceAsset) (types.Delegation, sdk.Dec) { //nolint:unparam // may wish to investigate
	newShares := types.GetDelegationSharesFromTokens(validator, asset, coin.Amount)
	delegation, found := k.GetDelegation(ctx, delAddr, validator, coin.Denom)
	if !found {
		delegation = types.NewDelegation(ctx, delAddr, validator.GetOperator(), coin.Denom, newShares, validator.GlobalRewardHistory)
	} else {
		delegation.Shares = delegation.Shares.Add(newShares)
	}
	k.SetDelegation(ctx, delAddr, validator.GetOperator(), coin.Denom, delegation)
	return delegation, newShares
}

// reduceDelegationShares
// If shares after reduction = 0, delegation will be deleted
func (k Keeper) reduceDelegationShares(ctx sdk.Context, delAddr sdk.AccAddress, validator types.AllianceValidator, coin sdk.Coin, shares sdk.Dec, delegation types.Delegation) {
	delegation.Shares = delegation.Shares.Sub(shares)
	if delegation.Shares.IsZero() {
		store := ctx.KVStore(k.storeKey)
		key := types.GetDelegationKey(delAddr, validator.GetOperator(), coin.Denom)
		store.Delete(key)
	} else {
		k.SetDelegation(ctx, delAddr, validator.GetOperator(), coin.Denom, delegation)
	}
}

func (k Keeper) updateValidatorShares(ctx sdk.Context, validator types.AllianceValidator, delegationShares sdk.DecCoins, validatorShares sdk.DecCoins, isAdd bool) {
	if isAdd {
		validator.AddShares(delegationShares, validatorShares)
	} else {
		validator.ReduceShares(delegationShares, validatorShares)
	}
	k.SetValidator(ctx, validator)
}

// GetAllianceBondedAmount returns the total amount of bonded native tokens that are not in the
// unbonding pool
func (k Keeper) GetAllianceBondedAmount(ctx sdk.Context, delegator sdk.AccAddress) math.Int {
	bonded := sdk.ZeroDec()
	k.stakingKeeper.IterateDelegatorDelegations(ctx, delegator, func(delegation stakingtypes.Delegation) bool {
		validatorAddr, err := sdk.ValAddressFromBech32(delegation.ValidatorAddress)
		if err != nil {
			panic(err) // shouldn't happen
		}
		validator, found := k.stakingKeeper.GetValidator(ctx, validatorAddr)
		if found && validator.IsBonded() {
			shares := delegation.Shares
			tokens := validator.TokensFromSharesTruncated(shares)
			bonded = bonded.Add(tokens)
		}
		return false
	})
	return bonded.TruncateInt()
}

// ResetAssetAndValidators
// When an asset has no more tokens being delegated, go through all validators and set
// validator shares to zero
func (k Keeper) ResetAssetAndValidators(ctx sdk.Context, asset types.AllianceAsset) {
	k.IterateAllianceValidatorInfo(ctx, func(valAddr sdk.ValAddress, info types.AllianceValidatorInfo) (stop bool) {
		updatedShares := sdk.NewDecCoins()
		for _, share := range info.ValidatorShares {
			if share.Denom != asset.Denom {
				updatedShares = append(updatedShares, share)
			}
		}
		info.ValidatorShares = updatedShares
		k.SetValidatorInfo(ctx, valAddr, info)
		return false
	})
	asset.TotalValidatorShares = sdk.ZeroDec()
	k.SetAsset(ctx, asset)
}
