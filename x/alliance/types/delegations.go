package types

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewDelegation(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, denom string, shares math.LegacyDec, rewardHistory []RewardHistory) Delegation {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return Delegation{
		DelegatorAddress:      delAddr.String(),
		ValidatorAddress:      valAddr.String(),
		Denom:                 denom,
		Shares:                shares,
		RewardHistory:         rewardHistory,
		LastRewardClaimHeight: uint64(sdkCtx.BlockHeight()),
	}
}
