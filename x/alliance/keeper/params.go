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

func (k Keeper) GlobalRewardIndices(ctx sdk.Context) (res types.RewardIndices) {
	k.paramstore.Get(ctx, types.GlobalRewardIndices, &res)
	return
}

func (k Keeper) SetGlobalRewardIndex(ctx sdk.Context, index types.RewardIndices) {
	k.paramstore.Set(ctx, types.GlobalRewardIndices, &index)
	return
}
