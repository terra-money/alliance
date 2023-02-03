package keeper

import (
	"time"

	"github.com/terra-money/alliance/x/alliance/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) RewardDelayTime(ctx sdk.Context) (res time.Duration) {
	k.paramstore.Get(ctx, types.RewardDelayTime, &res)
	return
}

func (k Keeper) RewardClaimInterval(ctx sdk.Context) (res time.Duration) {
	k.paramstore.Get(ctx, types.TakeRateClaimInterval, &res)
	return
}

func (k Keeper) LastRewardClaimTime(ctx sdk.Context) (res time.Time) {
	k.paramstore.Get(ctx, types.LastTakeRateClaimTime, &res)
	return
}

func (k Keeper) SetLastRewardClaimTime(ctx sdk.Context, lastTime time.Time) {
	k.paramstore.Set(ctx, types.LastTakeRateClaimTime, &lastTime)
}
