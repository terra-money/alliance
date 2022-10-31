package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewDelegation(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, denom string, shares sdk.Dec, rewardHistory []RewardHistory) Delegation {
	return Delegation{
		DelegatorAddress:      delAddr.String(),
		ValidatorAddress:      valAddr.String(),
		Denom:                 denom,
		Shares:                shares,
		RewardHistory:         rewardHistory,
		LastRewardClaimHeight: uint64(ctx.BlockHeight()),
	}
}

// ReduceShares
func (d *Delegation) ReduceShares(shares sdk.Dec) {
	if d.Shares.LTE(shares) {
		d.Shares = sdk.ZeroDec()
	} else {
		d.Shares = d.Shares.Sub(shares)
	}
}

// AddShares
func (d *Delegation) AddShares(shares sdk.Dec) {
	d.Shares = d.Shares.Add(shares)
}
