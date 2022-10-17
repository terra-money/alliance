package keeper

import (
	"alliance/x/alliance/types"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"time"
)

func (k Keeper) Delegate(ctx sdk.Context, delAddr sdk.AccAddress, validator stakingtypes.Validator, coin sdk.Coin) (*types.Delegation, error) {
	asset := k.GetAssetByDenom(ctx, coin.Denom)
	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, delAddr, types.ModuleName, sdk.NewCoins(coin))
	if err != nil {
		return nil, err
	}
	tokensToMint := asset.ConvertToStake(coin.Amount)
	err = k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.Coin{
		Denom:  k.stakingKeeper.BondDenom(ctx),
		Amount: tokensToMint,
	}))
	if err != nil {
		return nil, err
	}
	_, err = k.stakingKeeper.Delegate(ctx, moduleAddr, tokensToMint, stakingtypes.Unbonded, validator, true)
	if err != nil {
		return nil, err
	}
	asset.TotalTokens = asset.TotalTokens.Add(coin.Amount)
	k.SetAsset(ctx, asset)
	delegation := k.upsertDelegationWithNewTokens(ctx, delAddr, validator, coin, asset)
	return &delegation, nil
}

func (k Keeper) upsertDelegationWithNewTokens(ctx sdk.Context, delAddr sdk.AccAddress, validator stakingtypes.Validator, coin sdk.Coin, asset types.AllianceAsset) types.Delegation {
	newShares := convertNewTokenToShares(asset.TotalTokens, asset.TotalShares, coin.Amount)
	return k.upsertDelegationWithNewShares(ctx, delAddr, validator, coin, newShares)
}

func (k Keeper) upsertDelegationWithNewShares(ctx sdk.Context, delAddr sdk.AccAddress, validator stakingtypes.Validator, coin sdk.Coin, shares sdk.Dec) types.Delegation {
	delegation, ok := k.GetDelegation(ctx, delAddr, validator, coin.Denom)
	if !ok {
		delegation = types.Delegation{
			DelegatorAddress: delAddr.String(),
			ValidatorAddress: validator.GetOperator().String(),
			Denom:            coin.Denom,
			Shares:           shares,
		}
	} else {
		delegation.AddShares(shares)
	}
	k.SetDelegation(ctx, delAddr, validator, coin.Denom, delegation)
	return delegation
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

// Redelegate from one validator to another
// Method assumes that all tokens are owned by delegator and has delegations staked with srcVal
func (k Keeper) Redelegate(ctx sdk.Context, delAddr sdk.AccAddress, srcVal stakingtypes.Validator, dstVal stakingtypes.Validator, coin sdk.Coin) (*types.MsgRedelegateResponse, error) {
	asset := k.GetAssetByDenom(ctx, coin.Denom)

	srcDelegation, ok := k.GetDelegation(ctx, delAddr, srcVal, coin.Denom)
	if !ok {
		return nil, stakingtypes.ErrNoDelegatorForAddress
	}
	updatedShares, err := k.ValidateDelegatedAmount(srcDelegation, coin, asset)
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

	completionTime, err := k.stakingKeeper.BeginRedelegation(ctx, moduleAddr, srcVal.GetOperator(), dstVal.GetOperator(), shares)
	if err != nil {
		return nil, err
	}

	// Since all delegations are owned by the module account,
	// we remove redelegation from x/staling here and re-record it in x/alliance to allow transitive re-delegation in x/staking
	// The implication of this is that re-delegations will not be slashed if the src validator is slashed
	// TODO: Update slashing module to make sure we handle slashing here
	k.stakingKeeper.RemoveRedelegation(ctx, stakingtypes.Redelegation{
		DelegatorAddress:    moduleAddr.String(),
		ValidatorSrcAddress: srcVal.OperatorAddress,
		ValidatorDstAddress: dstVal.OperatorAddress,
		Entries:             nil,
	})

	// Reduce shares from src validator
	k.reduceDelegationShares(ctx, delAddr, srcVal, coin, updatedShares, srcDelegation)
	// Add shares to destination validator
	k.upsertDelegationWithNewShares(ctx, delAddr, dstVal, coin, updatedShares)
	k.addRedelegation(ctx, delAddr, srcVal.GetOperator(), dstVal.GetOperator(), coin, completionTime)
	k.queueRedelegation(ctx, delAddr, srcVal.GetOperator(), dstVal.GetOperator(), coin, completionTime)
	return &types.MsgRedelegateResponse{}, nil
}

func (k Keeper) ValidateDelegatedAmount(delegation types.Delegation, coin sdk.Coin, asset types.AllianceAsset) (shares sdk.Dec, err error) {
	shares = convertNewTokenToShares(asset.TotalTokens, asset.TotalShares, coin.Amount)
	if delegation.Shares.LT(shares.TruncateDec()) {
		return sdk.Dec{}, stakingtypes.ErrInsufficientShares
	}
	return shares, nil
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

func (k Keeper) Undelegate() {
	panic("implement me")
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

func convertNewTokenToShares(totalTokens math.Int, totalShares sdk.Dec, newTokens math.Int) (shares sdk.Dec) {
	if totalShares.IsZero() {
		return sdk.NewDecFromInt(newTokens)
	}
	return totalShares.MulInt(newTokens).QuoInt(totalTokens)
}

func convertNewShareToToken(totalTokens math.Int, totalShares sdk.Dec, shares sdk.Dec) (token math.Int) {
	return shares.MulInt(totalTokens).Quo(totalShares).TruncateInt()
}

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

// queueRedelegation Adds a redelegation to a queue to be processed at a later timestamp
// TODO: Handle a max number of entries per timestamp
// TODO: Logic in end block to dequeue and remove redelegations
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
