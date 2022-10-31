package keeper

import (
	"alliance/x/alliance/types"
	"cosmossdk.io/math"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Delegate(ctx sdk.Context, delAddr sdk.AccAddress, validator types.AllianceValidator, coin sdk.Coin) (*types.Delegation, error) {
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
	delegation, newDelegationShares := k.upsertDelegationWithNewTokens(ctx, delAddr, validator, coin, asset)

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
	return &delegation, nil
}

// Redelegate from one validator to another
// Method assumes that all tokens are owned by delegator and has delegations staked with srcVal
func (k Keeper) Redelegate(ctx sdk.Context, delAddr sdk.AccAddress, srcVal types.AllianceValidator, dstVal types.AllianceValidator, coin sdk.Coin) (*types.MsgRedelegateResponse, error) {
	asset, found := k.GetAssetByDenom(ctx, coin.Denom)

	if !found {
		return nil, status.Errorf(codes.NotFound, "Asset with denom: %s does not exist", coin.Denom)
	}

	srcDelegation, ok := k.GetDelegation(ctx, delAddr, srcVal, coin.Denom)
	if !ok {
		return nil, stakingtypes.ErrNoDelegatorForAddress
	}

	updatedDelegationShares, err := k.ValidateDelegatedAmount(srcDelegation, coin, srcVal, asset)
	if err != nil {
		return nil, err
	}

	// Prevents transitive re-delegations
	// e.g. if A -> B then B -> C is blocked until re-delegation is completed
	iter := k.IterateRedelegations(ctx, delAddr, srcVal.GetOperator(), coin.Denom)
	defer iter.Close()
	if iter.Valid() {
		return nil, stakingtypes.ErrTransitiveRedelegation
	}

	_, found = k.GetDelegation(ctx, delAddr, srcVal, coin.Denom)
	if found {
		_, err = k.ClaimDelegationRewards(ctx, delAddr, srcVal, coin.Denom)
		if err != nil {
			return nil, err
		}
	}
	_, found = k.GetDelegation(ctx, delAddr, dstVal, coin.Denom)
	if found {
		_, err = k.ClaimDelegationRewards(ctx, delAddr, dstVal, coin.Denom)
		if err != nil {
			return nil, err
		}
	}

	completionTime := ctx.BlockHeader().Time.Add(k.stakingKeeper.UnbondingTime(ctx))
	changedValidatorShares := types.GetValidatorShares(asset, coin.Amount)

	// Remove tokens and from from src validator
	k.reduceDelegationShares(ctx, delAddr, srcVal, coin, updatedDelegationShares, srcDelegation)
	k.updateValidatorShares(
		ctx,
		srcVal,
		sdk.NewDecCoins(sdk.NewDecCoinFromDec(coin.Denom, updatedDelegationShares)),
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
	k.queueRedelegation(ctx, delAddr, srcVal.GetOperator(), dstVal.GetOperator(), coin, completionTime)

	k.QueueAssetRebalanceEvent(ctx)
	return &types.MsgRedelegateResponse{}, nil
}

func (k Keeper) Undelegate(ctx sdk.Context, delAddr sdk.AccAddress, validator types.AllianceValidator, coin sdk.Coin) error {
	// Query for things needed for undelegation
	asset, found := k.GetAssetByDenom(ctx, coin.Denom)

	if !found {
		return status.Errorf(codes.NotFound, "Asset with denom: %s does not exist", coin.Denom)
	}

	delegation, ok := k.GetDelegation(ctx, delAddr, validator, coin.Denom)
	if !ok {
		return stakingtypes.ErrNoDelegatorForAddress
	}

	// Calculate how much delegation shares to be undelegated
	delegationSharesToUndelegate, err := k.ValidateDelegatedAmount(delegation, coin, validator, asset)
	if err != nil {
		return err
	}
	validatorSharesToRemove := types.GetValidatorShares(asset, coin.Amount)

	// Claim delegation rewards first
	_, err = k.ClaimDelegationRewards(ctx, delAddr, validator, coin.Denom)
	if err != nil {
		return err
	}

	asset.TotalTokens = asset.TotalTokens.Sub(coin.Amount)
	asset.TotalValidatorShares = asset.TotalValidatorShares.Sub(validatorSharesToRemove)
	k.SetAsset(ctx, asset)
	k.reduceDelegationShares(ctx, delAddr, validator, coin, delegationSharesToUndelegate, delegation)

	// Remove tokens and shares from src validator
	k.updateValidatorShares(
		ctx,
		validator,
		sdk.NewDecCoins(sdk.NewDecCoinFromDec(coin.Denom, delegationSharesToUndelegate)),
		sdk.NewDecCoins(sdk.NewDecCoinFromDec(coin.Denom, validatorSharesToRemove)),
		false,
	)

	// Queue undelegation messages to distribute tokens after undelegation completes in the future
	k.queueUndelegation(ctx, delAddr, validator.GetOperator(), coin)
	k.QueueAssetRebalanceEvent(ctx)
	return nil
}

// CompleteRedelegations Go through the re-delegations queue and remove all that have passed the completion time
func (k Keeper) CompleteRedelegations(ctx sdk.Context) int {
	store := ctx.KVStore(k.storeKey)
	iter := store.Iterator(types.RedelegationQueueKey, types.GetRedelegationQueueKey(ctx.BlockTime()))
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
func (k Keeper) CompleteUndelegations(ctx sdk.Context) (int, error) {
	store := ctx.KVStore(k.storeKey)
	iter := store.Iterator(types.UndelegationQueueKey, types.GetUndelegationQueueKey(ctx.BlockTime()))
	processed := 0
	for ; iter.Valid(); iter.Next() {
		var queued types.QueuedUndelegation
		k.cdc.MustUnmarshal(iter.Value(), &queued)
		for _, undel := range queued.Entries {
			delArr, _ := sdk.AccAddressFromBech32(undel.DelegatorAddress)
			err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, delArr, sdk.NewCoins(undel.Balance))
			if err != nil {
				return 0, err
			}
			processed++
		}
		store.Delete(iter.Key())
	}

	// Burn all "virtual" staking tokens in the module account
	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	coin := k.bankKeeper.GetBalance(ctx, moduleAddr, k.stakingKeeper.BondDenom(ctx))
	if !coin.IsZero() {
		err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(coin))
		if err != nil {
			return processed, err
		}
	}
	return processed, nil
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

func (k Keeper) SetDelegation(ctx sdk.Context, delAddr sdk.AccAddress, validator types.AllianceValidator, denom string, del types.Delegation) {
	key := types.GetDelegationKey(delAddr, validator.GetOperator(), denom)
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
	store := ctx.KVStore(k.storeKey)
	key := types.GetRedelegationKey(delAddr, redel.Balance.Denom, dstValAddr, completion)
	store.Delete(key)
}

func (k Keeper) IterateRedelegations(ctx sdk.Context, delAddr sdk.AccAddress, dstVal sdk.ValAddress, denom string) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	key := types.GetRedelegationsKey(delAddr, denom, dstVal)
	return sdk.KVStorePrefixIterator(store, key)
}

func (k Keeper) IterateRedelegationsByDelegator(ctx sdk.Context, delAddr sdk.AccAddress) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	key := types.GetRedelegationsKeyByDelegator(delAddr)
	return sdk.KVStorePrefixIterator(store, key)
}

func (k Keeper) SetValidator(ctx sdk.Context, val types.AllianceValidator) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetAllianceValidatorInfoKey(val.GetOperator())
	vb := k.cdc.MustMarshal(val.AllianceValidatorInfo)
	store.Set(key, vb)
}

func (k Keeper) ValidateDelegatedAmount(delegation types.Delegation, coin sdk.Coin, val types.AllianceValidator, asset types.AllianceAsset) (shares sdk.Dec, err error) {
	delegationShares := types.GetDelegationSharesFromTokens(val, asset, coin.Amount)
	if delegation.Shares.LT(delegationShares.TruncateDec()) {
		return sdk.Dec{}, stakingtypes.ErrInsufficientShares
	}
	return delegationShares, nil
}

// queueRedelegation Adds a redelegation to a queue to be processed at a later timestamp
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
}

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

func (k Keeper) queueUndelegation(ctx sdk.Context, delAddr sdk.AccAddress, val sdk.ValAddress, coin sdk.Coin) {
	store := ctx.KVStore(k.storeKey)
	completionTime := ctx.BlockTime().Add(k.stakingKeeper.UnbondingTime(ctx))
	queueKey := types.GetUndelegationQueueKey(completionTime)
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
	b = k.cdc.MustMarshal(&queue)
	store.Set(queueKey, b)
}

func (k Keeper) upsertDelegationWithNewTokens(ctx sdk.Context, delAddr sdk.AccAddress, validator types.AllianceValidator, coin sdk.Coin, asset types.AllianceAsset) (types.Delegation, sdk.Dec) {
	newShares := types.GetDelegationSharesFromTokens(validator, asset, coin.Amount)
	delegation, ok := k.GetDelegation(ctx, delAddr, validator, coin.Denom)
	latestClaimHistory := validator.GlobalRewardHistory
	if !ok {
		delegation = types.NewDelegation(ctx, delAddr, validator.GetOperator(), coin.Denom, newShares, latestClaimHistory)
	} else {
		delegation.AddShares(newShares)
	}
	k.SetDelegation(ctx, delAddr, validator, coin.Denom, delegation)
	return delegation, newShares
}

// reduceDelegationShares
// If shares after reduction = 0, delegation will be deleted
func (k Keeper) reduceDelegationShares(ctx sdk.Context, delAddr sdk.AccAddress, validator types.AllianceValidator, coin sdk.Coin, shares sdk.Dec, delegation types.Delegation) {
	delegation.ReduceShares(shares)
	store := ctx.KVStore(k.storeKey)
	key := types.GetDelegationKey(delAddr, validator.GetOperator(), coin.Denom)
	if delegation.Shares.IsZero() {
		store.Delete(key)
	} else {
		b := k.cdc.MustMarshal(&delegation)
		ctx.KVStore(k.storeKey).Set(key, b)
		store.Set(key, b)
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

// getAllianceBondedAmount returns the total amount of bonded native tokens that are not in the
// unbonding pool
func (k Keeper) getAllianceBondedAmount(ctx sdk.Context, delegator sdk.AccAddress) math.Int {
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
	return bonded.RoundInt()
}
