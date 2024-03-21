package keeper

import (
	"context"
	"time"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"

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
func (k Keeper) Delegate(ctx context.Context, delAddr sdk.AccAddress, validator types.AllianceValidator, coin sdk.Coin) (*math.LegacyDec, error) {
	// Check if asset is whitelisted as an alliance asset
	asset, found := k.GetAssetByDenom(ctx, coin.Denom)
	if !found {
		return nil, status.Errorf(codes.NotFound, "asset with denom: %s does not exist in alliance whitelist", coin.Denom)
	}

	// for the AllianceDenomTwo.
	// Check and send delegated tokens into the alliance module address
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, delAddr, types.ModuleName, sdk.NewCoins(coin))
	if err != nil {
		return nil, err
	}

	valAddr, err := validator.GetValAddress()
	if err != nil {
		return nil, err
	}
	// Claim rewards before adding more to a previous delegation
	_, found = k.GetDelegation(ctx, delAddr, valAddr, coin.Denom)
	if found {
		_, err = k.ClaimDelegationRewards(ctx, delAddr, validator, coin.Denom)
		if err != nil {
			return nil, err
		}
	}

	// Create or update a delegation
	_, newDelegationShares, err := k.upsertDelegationWithNewTokens(ctx, delAddr, validator, coin, asset)
	if err != nil {
		return nil, err
	}

	// Update ownership shares in validator and asset
	newValidatorShares := types.GetValidatorShares(asset, coin.Amount)
	asset.TotalTokens = asset.TotalTokens.Add(coin.Amount)
	asset.TotalValidatorShares = asset.TotalValidatorShares.Add(newValidatorShares)
	err = k.SetAsset(ctx, asset)
	if err != nil {
		return nil, err
	}

	err = k.updateValidatorShares(ctx, validator,
		sdk.NewDecCoins(sdk.NewDecCoinFromDec(coin.Denom, newDelegationShares)),
		sdk.NewDecCoins(sdk.NewDecCoinFromDec(coin.Denom, newValidatorShares)),
		true,
	)
	if err != nil {
		return nil, err
	}
	err = k.QueueAssetRebalanceEvent(ctx)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	_ = sdkCtx.EventManager().EmitTypedEvent(
		&types.DelegateAllianceEvent{
			AllianceSender: delAddr.String(),
			Validator:      validator.OperatorAddress,
			Coin:           coin,
			NewShares:      newValidatorShares,
		},
	)

	return &newValidatorShares, nil
}

// Redelegate from one validator to another
func (k Keeper) Redelegate(ctx context.Context, delAddr sdk.AccAddress, srcVal types.AllianceValidator, dstVal types.AllianceValidator, coin sdk.Coin) (*time.Time, error) {
	if srcVal.Validator.Equal(dstVal.Validator) {
		return nil, status.Errorf(codes.InvalidArgument, "Cannot redelegate to the same validator")
	}

	asset, found := k.GetAssetByDenom(ctx, coin.Denom)
	if !found {
		return nil, status.Errorf(codes.NotFound, "Asset with denom: %s does not exist", coin.Denom)
	}

	srcValAddr, err := srcVal.GetValAddress()
	if err != nil {
		return nil, err
	}

	dstValAddr, err := dstVal.GetValAddress()
	if err != nil {
		return nil, err
	}

	_, found = k.GetDelegation(ctx, delAddr, srcValAddr, coin.Denom)
	if !found {
		return nil, stakingtypes.ErrNoDelegatorForAddress
	}
	_, err = k.ClaimDelegationRewards(ctx, delAddr, srcVal, coin.Denom)
	if err != nil {
		return nil, err
	}
	// re-query delegation since it was updated in `ClaimDelegationRewards`
	srcDelegation, _ := k.GetDelegation(ctx, delAddr, srcValAddr, coin.Denom)

	_, found = k.GetDelegation(ctx, delAddr, dstValAddr, coin.Denom)
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
	if k.HasRedelegation(ctx, delAddr, srcValAddr, coin.Denom) {
		return nil, stakingtypes.ErrTransitiveRedelegation
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	unbondingTime, err := k.stakingKeeper.UnbondingTime(ctx)
	if err != nil {
		return nil, err
	}
	completionTime := sdkCtx.BlockHeader().Time.Add(unbondingTime)
	changedValidatorShares := types.GetValidatorShares(asset, coin.Amount)

	// Remove tokens and shares from src validator
	err = k.reduceDelegationShares(ctx, delAddr, srcVal, coin, delegationSharesToRemove, srcDelegation)
	if err != nil {
		return nil, err
	}
	err = k.updateValidatorShares(
		ctx,
		srcVal,
		sdk.NewDecCoins(sdk.NewDecCoinFromDec(coin.Denom, delegationSharesToRemove)),
		sdk.NewDecCoins(sdk.NewDecCoinFromDec(coin.Denom, changedValidatorShares)),
		false,
	)
	if err != nil {
		return nil, err
	}

	err = k.ClearDustDelegation(ctx, delAddr, srcVal, asset)
	if err != nil {
		return nil, err
	}

	// Add tokens and shares to dst validator
	_, newDelegationShares, err := k.upsertDelegationWithNewTokens(ctx, delAddr, dstVal, coin, asset)
	if err != nil {
		return nil, err
	}
	if err = k.updateValidatorShares(
		ctx,
		dstVal,
		sdk.NewDecCoins(sdk.NewDecCoinFromDec(coin.Denom, newDelegationShares)),
		sdk.NewDecCoins(sdk.NewDecCoinFromDec(coin.Denom, changedValidatorShares)),
		true,
	); err != nil {
		return nil, err
	}

	err = k.addRedelegation(ctx, delAddr, srcValAddr, dstValAddr, coin, completionTime)
	if err != nil {
		return nil, err
	}

	err = k.QueueAssetRebalanceEvent(ctx)
	if err != nil {
		return nil, err
	}

	_ = sdkCtx.EventManager().EmitTypedEvent(
		&types.RedelegateAllianceEvent{
			AllianceSender:       delAddr.String(),
			SourceValidator:      srcVal.OperatorAddress,
			DestinationValidator: dstVal.OperatorAddress,
			Coin:                 coin,
			CompletionTime:       completionTime,
		},
	)

	return &completionTime, nil
}

// Undelegate from a validator
// Staked tokens are only distributed to the delegator after the unbonding period
func (k Keeper) Undelegate(ctx context.Context, delAddr sdk.AccAddress, validator types.AllianceValidator, coin sdk.Coin) (*time.Time, error) {
	asset, found := k.GetAssetByDenom(ctx, coin.Denom)
	if !found {
		return nil, status.Errorf(codes.NotFound, "Asset with denom: %s does not exist", coin.Denom)
	}

	valAddr, err := validator.GetValAddress()
	if err != nil {
		return nil, err
	}

	_, ok := k.GetDelegation(ctx, delAddr, valAddr, coin.Denom)
	if !ok {
		return nil, stakingtypes.ErrNoDelegatorForAddress
	}
	// Claim delegation rewards first
	_, err = k.ClaimDelegationRewards(ctx, delAddr, validator, coin.Denom)
	if err != nil {
		return nil, err
	}

	// Delegation is queried again since it might have been modified when claiming delegation rewards
	delegation, _ := k.GetDelegation(ctx, delAddr, valAddr, coin.Denom)

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
	err = k.SetAsset(ctx, asset)
	if err != nil {
		return nil, err
	}

	// Remove shares from the delegation
	err = k.reduceDelegationShares(ctx, delAddr, validator, coin, delegationSharesToUndelegate, delegation)
	if err != nil {
		return nil, err
	}

	// Remove tokens and shares from src validator
	err = k.updateValidatorShares(
		ctx,
		validator,
		sdk.NewDecCoins(sdk.NewDecCoinFromDec(coin.Denom, delegationSharesToUndelegate)),
		sdk.NewDecCoins(sdk.NewDecCoinFromDec(coin.Denom, validatorSharesToRemove)),
		false,
	)
	if err != nil {
		return nil, err
	}

	err = k.ClearDustDelegation(ctx, delAddr, validator, asset)
	if err != nil {
		return nil, err
	}

	// Queue undelegation messages to distribute tokens after undelegation completes in the future
	completionTime, err := k.queueUndelegation(ctx, delAddr, valAddr, coin)
	if err != nil {
		return nil, err
	}
	err = k.QueueAssetRebalanceEvent(ctx)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	_ = sdkCtx.EventManager().EmitTypedEvent(
		&types.UndelegateAllianceEvent{
			AllianceSender: delAddr.String(),
			Validator:      validator.OperatorAddress,
			Coin:           coin,
			CompletionTime: completionTime,
		},
	)

	return &completionTime, nil
}

// CompleteRedelegations Go through the re-delegations queue and remove all that have passed the completion time
func (k Keeper) CompleteRedelegations(ctx sdk.Context) (deleted int) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iter := store.Iterator(types.RedelegationQueueKey, types.GetRedelegationQueueKey(ctx.BlockTime()))
	defer iter.Close() //nolint:errcheck,nolintlint
	for ; iter.Valid(); iter.Next() {
		completion := types.ParseRedelegationQueueKey(iter.Key())
		var queued types.QueuedRedelegation
		k.cdc.MustUnmarshal(iter.Value(), &queued)
		for _, redel := range queued.Entries {
			err := k.DeleteRedelegation(ctx, *redel, completion)
			if err != nil {
				return deleted
			}
			deleted++
		}
		store.Delete(iter.Key())
	}
	return deleted
}

func (k Keeper) GetDelegation(
	ctx context.Context,
	delAddr sdk.AccAddress,
	valAddr sdk.ValAddress,
	denom string,
) (d types.Delegation, found bool) {
	key := types.GetDelegationKey(delAddr, valAddr, denom)
	b, err := k.storeService.OpenKVStore(ctx).Get(key)
	if b == nil || err != nil {
		return d, false
	}
	k.cdc.MustUnmarshal(b, &d)
	return d, true
}

func (k Keeper) SetDelegation(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, denom string, del types.Delegation) error {
	key := types.GetDelegationKey(delAddr, valAddr, denom)
	b := k.cdc.MustMarshal(&del)
	return k.storeService.OpenKVStore(ctx).Set(key, b)
}

func (k Keeper) DeleteRedelegation(ctx context.Context, redel types.Redelegation, completion time.Time) error {
	delAddr, err := sdk.AccAddressFromBech32(redel.DelegatorAddress)
	if err != nil {
		return err
	}
	dstValAddr, err := sdk.ValAddressFromBech32(redel.DstValidatorAddress)
	if err != nil {
		return err
	}
	srcValAddr, err := sdk.ValAddressFromBech32(redel.SrcValidatorAddress)
	if err != nil {
		return err
	}
	store := k.storeService.OpenKVStore(ctx)
	key := types.GetRedelegationKey(delAddr, redel.Balance.Denom, dstValAddr, completion)
	err = store.Delete(key)
	if err != nil {
		return err
	}
	indexKey := types.GetRedelegationIndexKey(srcValAddr, completion, redel.Balance.Denom, dstValAddr, delAddr)
	return store.Delete(indexKey)
}

func (k Keeper) IterateDelegations(ctx context.Context, cb func(d types.Delegation) (stop bool)) (err error) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iter := storetypes.KVStorePrefixIterator(store, types.DelegationKey)
	defer iter.Close() //nolint:errcheck,nolintlint
	for ; iter.Valid(); iter.Next() {
		var delegation types.Delegation
		k.cdc.MustUnmarshal(iter.Value(), &delegation)
		if cb(delegation) {
			break
		}
	}
	return nil
}

func (k Keeper) HasRedelegation(ctx context.Context, delAddr sdk.AccAddress, dstVal sdk.ValAddress, denom string) bool {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetRedelegationsKey(delAddr, denom, dstVal)
	iter := storetypes.KVStorePrefixIterator(store, key)
	defer iter.Close() //nolint:errcheck,nolintlint
	return iter.Valid()
}

func (k Keeper) IterateRedelegations(ctx context.Context, cb func(redelegation types.Redelegation, completionTime time.Time) (stop bool)) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iter := storetypes.KVStorePrefixIterator(store, types.RedelegationKey)
	defer iter.Close() //nolint:errcheck,nolintlint
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

func (k Keeper) IterateRedelegationsByDelegator(ctx context.Context, delAddr sdk.AccAddress) storetypes.Iterator {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetRedelegationsKeyByDelegator(delAddr)
	return storetypes.KVStorePrefixIterator(store, key)
}

func (k Keeper) IterateRedelegationsBySrcValidator(ctx context.Context, srcValAddr sdk.ValAddress) storetypes.Iterator {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	prefix := types.GetRedelegationsIndexOrderedByValidatorKey(srcValAddr)
	return storetypes.KVStorePrefixIterator(store, prefix)
}

func (k Keeper) IterateUndelegationsBySrcValidator(ctx context.Context, valAddr sdk.ValAddress) storetypes.Iterator {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	prefix := types.GetUndelegationsIndexOrderedByValidatorKey(valAddr)
	return storetypes.KVStorePrefixIterator(store, prefix)
}

func (k Keeper) IterateUndelegationsByCompletionTime(ctx context.Context, completionTime time.Time) storetypes.Iterator {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	return store.Iterator(types.UndelegationQueueKey, types.GetUndelegationQueueKeyByTime(completionTime))
}

func (k Keeper) IterateUndelegations(ctx context.Context, cb func(undelegation types.QueuedUndelegation, completionTime time.Time) (stop bool)) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iter := storetypes.KVStorePrefixIterator(store, types.UndelegationQueueKey)
	defer iter.Close() //nolint:errcheck,nolintlint
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

func (k Keeper) SetValidator(ctx context.Context, val types.AllianceValidator) error {
	store := k.storeService.OpenKVStore(ctx)
	valAddr, err := val.GetValAddress()
	if err != nil {
		return err
	}
	key := types.GetAllianceValidatorInfoKey(valAddr)
	vb := k.cdc.MustMarshal(val.AllianceValidatorInfo)
	return store.Set(key, vb)
}

func (k Keeper) SetValidatorInfo(ctx context.Context, valAddr sdk.ValAddress, val types.AllianceValidatorInfo) error {
	store := k.storeService.OpenKVStore(ctx)
	key := types.GetAllianceValidatorInfoKey(valAddr)
	vb := k.cdc.MustMarshal(&val)
	return store.Set(key, vb)
}

// ValidateDelegatedAmount returns the amount of shares for a given coin that is staked
// Returns the number of shares that represents the amount of staked tokens that was requested
func (k Keeper) ValidateDelegatedAmount(delegation types.Delegation, coin sdk.Coin, val types.AllianceValidator, asset types.AllianceAsset) (shares math.LegacyDec, err error) {
	delegationSharesToUpdate := types.GetDelegationSharesFromTokens(val, asset, coin.Amount)
	// Account for rounding in which shares for a full withdraw is slightly more or less than the number of shares recorded
	// Withdraw all in that case
	// 1e6 of margin should be enough to handle realistic rounding issues caused by using the fix-point math.
	if delegation.Shares.Sub(delegationSharesToUpdate).Abs().LT(types.Rounder) {
		return delegation.Shares, nil
	}

	if delegation.Shares.LT(delegationSharesToUpdate.TruncateDec()) {
		return math.LegacyDec{}, stakingtypes.ErrInsufficientShares
	}
	// Cap the shares at the delegation's shares. Shares being greater could occur
	// due to rounding, however we don't want to truncate the shares or take the
	// minimum because we want to allow for the full withdraw of shares from a
	// delegation.
	// This logic is similar to that found in x/staking
	if delegationSharesToUpdate.GT(delegation.Shares) {
		delegationSharesToUpdate = delegation.Shares
	}

	return delegationSharesToUpdate, nil
}

// addRedelegation adds a redelegation entry to be used to prevent premature re-delegations
func (k Keeper) addRedelegation(ctx context.Context, delAddr sdk.AccAddress, srcVal sdk.ValAddress, dstVal sdk.ValAddress, coin sdk.Coin, completionTime time.Time) (err error) {
	store := k.storeService.OpenKVStore(ctx)
	key := types.GetRedelegationKey(delAddr, coin.Denom, dstVal, completionTime)
	b, err := store.Get(key)
	var redelegation types.Redelegation
	if b == nil || err != nil {
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
	if err = store.Set(key, b); err != nil {
		return err
	}

	// Add another entry as an index to retrieve redelegations by validator
	indexKey := types.GetRedelegationIndexKey(srcVal, completionTime, coin.Denom, dstVal, delAddr)
	if err = store.Set(indexKey, []byte{}); err != nil {
		return err
	}
	return k.queueRedelegation(ctx, delAddr, srcVal, dstVal, coin, completionTime)
}

// queueRedelegation adds a redelegation queue object that is indexed by completion time
// This is used to track and clean-up redelegation events once they are mature
// Use addRedelegations instead
func (k Keeper) queueRedelegation(ctx context.Context, delAddr sdk.AccAddress, srcVal sdk.ValAddress, dstVal sdk.ValAddress, coin sdk.Coin, completionTime time.Time) error {
	store := k.storeService.OpenKVStore(ctx)
	queueKey := types.GetRedelegationQueueKey(completionTime)
	b, err := store.Get(queueKey)
	var queuedDelegations types.QueuedRedelegation
	if b == nil || err != nil {
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
	return store.Set(queueKey, b)
}

// queueUndelegation adds an undelegation queue object that is indexed by completion time
// This is used to track and clean-up undelegation events once they are mature
func (k Keeper) queueUndelegation(ctx context.Context, delAddr sdk.AccAddress, val sdk.ValAddress, coin sdk.Coin) (time.Time, error) {
	store := k.storeService.OpenKVStore(ctx)
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	unbondingTime, err := k.stakingKeeper.UnbondingTime(ctx)
	if err != nil {
		return time.Time{}, err
	}
	completionTime := sdkCtx.BlockTime().Add(unbondingTime)
	queueKey := types.GetUndelegationQueueKey(completionTime, delAddr)
	b, err := store.Get(queueKey)
	var queue types.QueuedUndelegation
	if b == nil || err != nil {
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
	err = k.setQueuedUndelegations(ctx, completionTime, delAddr, queue)
	if err != nil {
		return time.Time{}, err
	}
	err = k.setUnbondingIndexByVal(ctx, val, completionTime, delAddr, coin.Denom)
	if err != nil {
		return time.Time{}, err
	}
	return completionTime, nil
}

func (k Keeper) setQueuedUndelegations(ctx context.Context, completionTime time.Time, delAddr sdk.AccAddress, queuedDelegations types.QueuedUndelegation) error {
	store := k.storeService.OpenKVStore(ctx)
	queueKey := types.GetUndelegationQueueKey(completionTime, delAddr)
	b := k.cdc.MustMarshal(&queuedDelegations)
	return store.Set(queueKey, b)
}

func (k Keeper) setUnbondingIndexByVal(ctx context.Context, valAddr sdk.ValAddress, completionTime time.Time, delAddr sdk.AccAddress, denom string) error {
	store := k.storeService.OpenKVStore(ctx)
	indexKey := types.GetUnbondingIndexKey(valAddr, completionTime, denom, delAddr)
	return store.Set(indexKey, []byte{})
}

func (k Keeper) upsertDelegationWithNewTokens(ctx context.Context, delAddr sdk.AccAddress, validator types.AllianceValidator, coin sdk.Coin, asset types.AllianceAsset) (types.Delegation, math.LegacyDec, error) { //nolint:unparam // may wish to investigate
	newShares := types.GetDelegationSharesFromTokens(validator, asset, coin.Amount)
	valAddr, err := validator.GetValAddress()
	if err != nil {
		return types.Delegation{}, math.LegacyDec{}, err
	}
	delegation, found := k.GetDelegation(ctx, delAddr, valAddr, coin.Denom)
	if !found {
		delegation = types.NewDelegation(ctx, delAddr, valAddr, coin.Denom, newShares, validator.GlobalRewardHistory)
	} else {
		delegation.Shares = delegation.Shares.Add(newShares)
	}
	err = k.SetDelegation(ctx, delAddr, valAddr, coin.Denom, delegation)
	return delegation, newShares, err
}

// reduceDelegationShares
// If shares after reduction = 0, delegation will be deleted
func (k Keeper) reduceDelegationShares(ctx context.Context, delAddr sdk.AccAddress, validator types.AllianceValidator, coin sdk.Coin, shares math.LegacyDec, delegation types.Delegation) error {
	delegation.Shares = delegation.Shares.Sub(shares)
	valAddr, err := validator.GetValAddress()
	if err != nil {
		return err
	}
	if delegation.Shares.IsZero() {
		store := k.storeService.OpenKVStore(ctx)
		key := types.GetDelegationKey(delAddr, valAddr, coin.Denom)
		return store.Delete(key)
	}
	return k.SetDelegation(ctx, delAddr, valAddr, coin.Denom, delegation)
}

func (k Keeper) updateValidatorShares(ctx context.Context, validator types.AllianceValidator, delegationShares sdk.DecCoins, validatorShares sdk.DecCoins, isAdd bool) error {
	if isAdd {
		validator.AddShares(delegationShares, validatorShares)
	} else {
		validator.ReduceShares(delegationShares, validatorShares)
	}
	return k.SetValidator(ctx, validator)
}

// GetAllianceBondedAmount returns the total amount of bonded native tokens that are not in the
// unbonding pool
func (k Keeper) GetAllianceBondedAmount(ctx context.Context, delegator sdk.AccAddress) (math.Int, error) {
	bonded := math.LegacyZeroDec()
	err := k.stakingKeeper.IterateDelegatorDelegations(ctx, delegator, func(delegation stakingtypes.Delegation) bool {
		validatorAddr, err := sdk.ValAddressFromBech32(delegation.ValidatorAddress)
		if err != nil {
			panic(err) // shouldn't happen
		}
		validator, err := k.stakingKeeper.GetValidator(ctx, validatorAddr)
		if err == nil && validator.IsBonded() {
			shares := delegation.Shares
			tokens := validator.TokensFromSharesTruncated(shares)
			bonded = bonded.Add(tokens)
		}
		return false
	})
	if err != nil {
		return math.Int{}, err
	}
	return bonded.TruncateInt(), nil
}

// ResetAssetAndValidators
// When an asset has no more tokens being delegated, go through all validators and set
// validator shares to zero
func (k Keeper) ResetAssetAndValidators(ctx context.Context, asset types.AllianceAsset) error {
	// When there are no more tokens recorded in the asset, clear all share records that might remain
	// from rounding errors to prevent dust amounts from staying in the stores
	if !asset.TotalTokens.IsZero() {
		return nil
	}
	var err error
	err = k.IterateAllianceValidatorInfo(ctx, func(valAddr sdk.ValAddress, info types.AllianceValidatorInfo) (stop bool) {
		updatedShares := sdk.NewDecCoins()
		for _, share := range info.ValidatorShares {
			if share.Denom != asset.Denom {
				updatedShares = append(updatedShares, share)
			}
		}
		info.ValidatorShares = updatedShares
		if err = k.SetValidatorInfo(ctx, valAddr, info); err != nil {
			return true
		}
		return false
	})
	if err != nil {
		return err
	}
	asset.TotalValidatorShares = math.LegacyZeroDec()
	return k.SetAsset(ctx, asset)
}

func (k Keeper) ClearDustDelegation(ctx context.Context, delAddr sdk.AccAddress, validator types.AllianceValidator, asset types.AllianceAsset) error {
	delegatorSharesToRemove := sdk.NewDecCoinFromDec(asset.Denom, math.LegacyZeroDec())
	validatorSharesToRemove := sdk.NewDecCoinFromDec(asset.Denom, math.LegacyZeroDec())
	valAddr, err := validator.GetValAddress()
	if err != nil {
		return err
	}

	delegation, found := k.GetDelegation(ctx, delAddr, valAddr, asset.Denom)
	// If not found then the delegation has already been deleted, do nothing else
	if found {
		tokensLeft := types.GetDelegationTokensWithShares(delegation.Shares, validator, asset)
		// If there are no tokens that can be claimed by the delegation, delete the delegation
		if tokensLeft.IsZero() {
			store := k.storeService.OpenKVStore(ctx)
			delAddr := sdk.MustAccAddressFromBech32(delegation.DelegatorAddress) // acc address should always be valid here
			key := types.GetDelegationKey(delAddr, valAddr, asset.Denom)
			err = store.Delete(key)
			if err != nil {
				return err
			}

			delegatorSharesToRemove = sdk.NewDecCoinFromDec(asset.Denom, delegation.Shares)
		}
	}

	validatorTokensLeft := validator.TotalTokensWithAsset(asset)
	if validatorTokensLeft.IsZero() {
		validatorSharesToRemove = sdk.NewDecCoinFromDec(asset.Denom, validator.ValidatorSharesWithDenom(asset.Denom))
	}

	// Reduce the dust shares from validator to make sure everything adds up
	validator.ReduceShares(sdk.NewDecCoins(delegatorSharesToRemove), sdk.NewDecCoins(validatorSharesToRemove))
	if err = k.SetValidator(ctx, validator); err != nil {
		return err
	}

	err = k.ResetAssetAndValidators(ctx, asset)
	if err != nil {
		return err
	}
	return nil
}
