package keeper

import (
	"alliance/x/alliance/types"

	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Delegate(ctx sdk.Context, delAddr sdk.AccAddress, validator stakingtypes.Validator, coin sdk.Coin) (*types.Delegation, error) {
	// Check if asset is whitelisted as an alliance asset
	asset, found := k.GetAssetByDenom(ctx, coin.Denom)
	if !found {
		return nil, status.Errorf(codes.NotFound, "Asset with denom: %s does not exist", coin.Denom)
	}

	// Check and send delegated tokens into the alliance module address
	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, delAddr, types.ModuleName, sdk.NewCoins(coin))
	if err != nil {
		return nil, err
	}

	// Convert delegated tokens into staking tokens using the rewards rate
	tokensToMint := sdk.Coin{
		Denom:  k.stakingKeeper.BondDenom(ctx),
		Amount: asset.ConvertToStake(coin.Amount),
	}
	err = k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(tokensToMint))
	if err != nil {
		return nil, err
	}

	// Claim delegation rewards (if already delegated) before delegating more tokens
	// stakingKeeper.Delegate will actually claim as well but claiming here directly gives us
	// the amount of tokens claimed
	_, found = k.GetDelegation(ctx, delAddr, validator, coin.Denom)
	if found {
		_, err = k.ClaimDelegationRewards(ctx, delAddr, validator, coin.Denom)
		if err != nil {
			return nil, err
		}
	}

	// Delegate stake tokens to validators
	// Delegate would automatically claim rewards into the module address
	_, err = k.stakingKeeper.Delegate(ctx, moduleAddr, tokensToMint.Amount, stakingtypes.Unbonded, validator, true)
	if err != nil {
		return nil, err
	}

	delegation, newDelegationShares := k.upsertDelegationWithNewTokens(ctx, delAddr, validator, coin, asset)

	// Update asset info
	newValidatorShares := types.GetValidatorShares(asset, coin.Amount)
	asset.TotalTokens = asset.TotalTokens.Add(coin.Amount)
	asset.TotalValidatorShares = asset.TotalValidatorShares.Add(newValidatorShares)
	k.SetAsset(ctx, asset)

	// Update validator with tokens and shares
	k.updateValidatorShares(ctx, validator.GetOperator(),
		sdk.NewDecCoins(sdk.NewDecCoinFromDec(coin.Denom, newDelegationShares)),
		sdk.NewDecCoins(sdk.NewDecCoinFromDec(coin.Denom, newValidatorShares)),
		true,
	)

	return &delegation, nil
}

// Redelegate from one validator to another
// Method assumes that all tokens are owned by delegator and has delegations staked with srcVal
func (k Keeper) Redelegate(ctx sdk.Context, delAddr sdk.AccAddress, srcVal stakingtypes.Validator, dstVal stakingtypes.Validator, coin sdk.Coin) (*types.MsgRedelegateResponse, error) {
	asset, found := k.GetAssetByDenom(ctx, coin.Denom)

	if !found {
		return nil, status.Errorf(codes.NotFound, "Asset with denom: %s does not exist", coin.Denom)
	}

	srcDelegation, ok := k.GetDelegation(ctx, delAddr, srcVal, coin.Denom)
	if !ok {
		return nil, stakingtypes.ErrNoDelegatorForAddress
	}

	aVal := k.GetOrCreateValidator(ctx, srcVal.GetOperator())
	updatedDelegationShares, err := k.ValidateDelegatedAmount(srcDelegation, coin, aVal, asset)
	if err != nil {
		return nil, err
	}

	stakeTokens := asset.ConvertToStake(coin.Amount)
	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)

	shares, err := k.stakingKeeper.ValidateUnbondAmount(ctx, moduleAddr, srcVal.GetOperator(), stakeTokens)
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

	completionTime, err := k.stakingKeeper.BeginRedelegation(ctx, moduleAddr, srcVal.GetOperator(), dstVal.GetOperator(), shares)
	if err != nil {
		return nil, err
	}

	// Since all delegations are owned by the module account,
	// we remove redelegation from x/staling here and re-record it in x/alliance to allow transitive re-delegation in x/staking
	// The implication of this is that re-delegations will not be slashed if the src validator is slashed
	// TODO: Add a slashing hook to make sure we handle slashing for redelegations
	k.stakingKeeper.RemoveRedelegation(ctx, stakingtypes.Redelegation{
		DelegatorAddress:    moduleAddr.String(),
		ValidatorSrcAddress: srcVal.OperatorAddress,
		ValidatorDstAddress: dstVal.OperatorAddress,
		Entries:             nil,
	})

	changedValidatorShares := types.GetValidatorShares(asset, coin.Amount)

	// Remove tokens and from from src validator
	k.reduceDelegationShares(ctx, delAddr, srcVal, coin, updatedDelegationShares, srcDelegation)
	k.updateValidatorShares(
		ctx,
		srcVal.GetOperator(),
		sdk.NewDecCoins(sdk.NewDecCoinFromDec(coin.Denom, updatedDelegationShares)),
		sdk.NewDecCoins(sdk.NewDecCoinFromDec(coin.Denom, changedValidatorShares)),
		false,
	)

	// Add tokens and shares to dst validator
	_, newDelegationShares := k.upsertDelegationWithNewTokens(ctx, delAddr, dstVal, coin, asset)
	k.updateValidatorShares(
		ctx,
		dstVal.GetOperator(),
		sdk.NewDecCoins(sdk.NewDecCoinFromDec(coin.Denom, newDelegationShares)),
		sdk.NewDecCoins(sdk.NewDecCoinFromDec(coin.Denom, changedValidatorShares)),
		true,
	)

	k.addRedelegation(ctx, delAddr, srcVal.GetOperator(), dstVal.GetOperator(), coin, completionTime)
	k.queueRedelegation(ctx, delAddr, srcVal.GetOperator(), dstVal.GetOperator(), coin, completionTime)
	return &types.MsgRedelegateResponse{}, nil
}

func (k Keeper) Undelegate(ctx sdk.Context, delAddr sdk.AccAddress, validator stakingtypes.Validator, coin sdk.Coin) error {
	// Query for things needed for undelegation
	asset, found := k.GetAssetByDenom(ctx, coin.Denom)

	if !found {
		return status.Errorf(codes.NotFound, "Asset with denom: %s does not exist", coin.Denom)
	}

	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	delegation, ok := k.GetDelegation(ctx, delAddr, validator, coin.Denom)
	if !ok {
		return stakingtypes.ErrNoDelegatorForAddress
	}

	aVal := k.GetOrCreateValidator(ctx, validator.GetOperator())
	// Calculate how much delegation shares to be undelegated
	delegationSharesToUndelegate, err := k.ValidateDelegatedAmount(delegation, coin, aVal, asset)
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

	// Unbond from x/staking module
	stakeTokens := asset.ConvertToStake(coin.Amount)
	stakeShares, err := k.stakingKeeper.ValidateUnbondAmount(ctx, moduleAddr, validator.GetOperator(), stakeTokens)
	if err != nil {
		return err
	}

	_, err = k.stakingKeeper.Unbond(ctx, moduleAddr, validator.GetOperator(), stakeShares)
	if err != nil {
		return err
	}

	// Remove tokens and shares from src validator
	k.updateValidatorShares(
		ctx,
		validator.GetOperator(),
		sdk.NewDecCoins(sdk.NewDecCoinFromDec(coin.Denom, delegationSharesToUndelegate)),
		sdk.NewDecCoins(sdk.NewDecCoinFromDec(coin.Denom, validatorSharesToRemove)),
		false,
	)

	// Queue undelegation messages to distribute tokens after undelegation completes in the future
	k.queueUndelegation(ctx, delAddr, validator.GetOperator(), coin)
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
func (k Keeper) CompleteUndelegations(ctx sdk.Context) int {
	store := ctx.KVStore(k.storeKey)
	iter := store.Iterator(types.UndelegationQueueKey, types.GetUndelegationQueueKey(ctx.BlockTime()))
	processed := 0
	for ; iter.Valid(); iter.Next() {
		var queued types.QueuedUndelegation
		k.cdc.MustUnmarshal(iter.Value(), &queued)
		for _, undel := range queued.Entries {
			delArr, _ := sdk.AccAddressFromBech32(undel.DelegatorAddress)
			k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, delArr, sdk.NewCoins(undel.Balance))
			processed++
		}
		store.Delete(iter.Key())
	}
	// TODO: Burn stake tokens that got unbonded
	return processed
}

func (k Keeper) GetDelegation(ctx sdk.Context, delAddr sdk.AccAddress, validator stakingtypes.Validator, denom string) (d types.Delegation, found bool) {
	key := types.GetDelegationKey(delAddr, validator.GetOperator(), denom)
	b := ctx.KVStore(k.storeKey).Get(key)
	if b == nil {
		return d, false
	}
	k.cdc.MustUnmarshal(b, &d)
	return d, true
}

func (k Keeper) SetDelegation(ctx sdk.Context, delAddr sdk.AccAddress, validator stakingtypes.Validator, denom string, del types.Delegation) {
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

func (k Keeper) GetOrCreateValidator(ctx sdk.Context, valAddr sdk.ValAddress) (val types.Validator) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetValidatorKey(valAddr)
	vb := store.Get(key)
	if vb == nil {
		val = types.NewValidator(valAddr)
		vb = k.cdc.MustMarshal(&val)
		store.Set(key, vb)
	} else {
		k.cdc.MustUnmarshal(vb, &val)
	}
	return
}

func (k Keeper) SetValidator(ctx sdk.Context, valAddr sdk.ValAddress, val types.Validator) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetValidatorKey(valAddr)
	vb := k.cdc.MustMarshal(&val)
	store.Set(key, vb)
}

func (k Keeper) ValidateDelegatedAmount(delegation types.Delegation, coin sdk.Coin, aVal types.Validator, asset types.AllianceAsset) (shares sdk.Dec, err error) {
	delegationShares := types.GetDelegationSharesFromTokens(aVal, asset, coin.Amount)
	if delegation.Shares.LT(delegationShares.TruncateDec()) {
		return sdk.Dec{}, stakingtypes.ErrInsufficientShares
	}
	return delegationShares, nil
}

// queueRedelegation Adds a redelegation to a queue to be processed at a later timestamp
// TODO: Handle a max number of entries per timestamp
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

func (k Keeper) upsertDelegationWithNewTokens(ctx sdk.Context, delAddr sdk.AccAddress, validator stakingtypes.Validator, coin sdk.Coin, asset types.AllianceAsset) (types.Delegation, sdk.Dec) {
	aVal := k.GetOrCreateValidator(ctx, validator.GetOperator())
	newShares := types.GetDelegationSharesFromTokens(aVal, asset, coin.Amount)

	delegation, ok := k.GetDelegation(ctx, delAddr, validator, coin.Denom)
	globalRewardIndices := aVal.RewardIndices
	if !ok {
		delegation = types.NewDelegation(delAddr, validator.GetOperator(), coin.Denom, newShares, globalRewardIndices)
	} else {
		delegation.AddShares(newShares)
	}
	k.SetDelegation(ctx, delAddr, validator, coin.Denom, delegation)
	return delegation, newShares
}

// reduceDelegationShares
// If shares after reduction = 0, delegation will be deleted
func (k Keeper) reduceDelegationShares(ctx sdk.Context, delAddr sdk.AccAddress, validator stakingtypes.Validator, coin sdk.Coin, shares sdk.Dec, delegation types.Delegation) {
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

func (k Keeper) updateValidatorShares(ctx sdk.Context, valAddr sdk.ValAddress, delegationShares sdk.DecCoins, validatorShares sdk.DecCoins, isAdd bool) {
	aVal := k.GetOrCreateValidator(ctx, valAddr)
	if isAdd {
		aVal.AddShares(delegationShares, validatorShares)
	} else {
		aVal.ReduceShares(delegationShares, validatorShares)
	}
	k.SetValidator(ctx, valAddr, aVal)
}
