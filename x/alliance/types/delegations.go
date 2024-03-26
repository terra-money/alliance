package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewDelegation(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, denom string, shares sdkmath.LegacyDec, rewardHistory []RewardHistory) Delegation {
	return Delegation{
		DelegatorAddress:      delAddr.String(),
		ValidatorAddress:      valAddr.String(),
		Denom:                 denom,
		Shares:                shares,
		RewardHistory:         rewardHistory,
		LastRewardClaimHeight: uint64(ctx.BlockHeight()),
	}
}
