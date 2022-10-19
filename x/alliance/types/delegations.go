package types

import sdk "github.com/cosmos/cosmos-sdk/types"

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

func NewValidator(valAddr sdk.ValAddress) Validator {
	return Validator{
		ValidatorAddress: valAddr.String(),
		RewardIndices:    RewardIndices{},
		TotalTokens:      sdk.NewCoins(),
		TotalShares:      sdk.NewDecCoins(),
	}
}

func (v *Validator) AddTokens(coins sdk.Coins) {
	v.TotalTokens = coins.Add(v.TotalTokens...)
}

func (v *Validator) ReduceTokens(coins sdk.Coins) {
	v.TotalTokens = sdk.NewCoins(v.TotalTokens...).Sub(coins...)
}

func (v *Validator) AddShares(shares sdk.DecCoins) {
	v.TotalShares = shares.Add(v.TotalShares...)
}

func (v *Validator) ReduceShares(shares sdk.DecCoins) {
	v.TotalShares = sdk.NewDecCoins(v.TotalShares...).Sub(shares)
}

func (v Validator) TotalSharesWithDenom(denom string) sdk.Dec {
	return sdk.NewDecCoins(v.TotalShares...).AmountOf(denom)
}

func (v Validator) TotalTokensWithDenom(denom string) sdk.Int {
	return sdk.NewCoins(v.TotalTokens...).AmountOf(denom)
}
