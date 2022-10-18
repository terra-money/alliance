package keeper

import (
	"alliance/x/alliance/types"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) RewardDelayTime(ctx sdk.Context) (res time.Duration) {
	k.paramstore.Get(ctx, types.RewardDelayTime, &res)
	return
}

func (k Keeper) GlobalRewardIndex(ctx sdk.Context) (res sdk.Dec) {
	k.paramstore.Get(ctx, types.GlobalRewardIndex, &res)
	return
}

func (k Keeper) SetGlobalRewardIndex(ctx sdk.Context, index sdk.Dec) {
	k.paramstore.Set(ctx, types.GlobalRewardIndex, &index)
	return
}
