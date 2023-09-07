package keeper

import (
	"time"

	"github.com/terra-money/alliance/x/alliance/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) RewardDelayTime(ctx sdk.Context) (res time.Duration) {
	params := k.GetParams(ctx)
	return params.RewardDelayTime
}

func (k Keeper) RewardClaimInterval(ctx sdk.Context) (res time.Duration) {
	params := k.GetParams(ctx)
	return params.TakeRateClaimInterval
}

func (k Keeper) LastRewardClaimTime(ctx sdk.Context) (res time.Time) {
	params := k.GetParams(ctx)
	return params.LastTakeRateClaimTime
}

func (k Keeper) SetLastRewardClaimTime(ctx sdk.Context, lastTime time.Time) {
	params := k.GetParams(ctx)
	params.LastTakeRateClaimTime = lastTime
	k.SetParams(ctx, params)
}

func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return
	}
	k.cdc.MustUnmarshal(bz, &params)
	return
}

func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&params)
	store.Set(types.ParamsKey, bz)
}
