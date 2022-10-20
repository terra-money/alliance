package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewDelegation(delAddr sdk.AccAddress, valAddr sdk.ValAddress, denom string, shares sdk.Dec, rewardIndices []RewardIndex) Delegation {
	return Delegation{
		DelegatorAddress: delAddr.String(),
		ValidatorAddress: valAddr.String(),
		Denom:            denom,
		Shares:           shares,
		RewardIndices:    rewardIndices,
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

func NewValidator(valAddr sdk.ValAddress) Validator {
	return Validator{
		ValidatorAddress: valAddr.String(),
		RewardIndices:    RewardIndices{},
		TotalShares:      sdk.NewDecCoins(),
		ValidatorShares:  sdk.NewDecCoins(),
	}
}

func (v *Validator) AddShares(delegationShares sdk.DecCoins, validatorShares sdk.DecCoins) {
	v.TotalShares = delegationShares.Add(v.TotalShares...)
	v.ValidatorShares = validatorShares.Add(v.ValidatorShares...)
}

func (v *Validator) ReduceShares(delegationShares sdk.DecCoins, validatorShares sdk.DecCoins) {
	v.TotalShares = sdk.NewDecCoins(v.TotalShares...).Sub(delegationShares)
	v.ValidatorShares = sdk.NewDecCoins(v.ValidatorShares...).Sub(validatorShares)
}

func (v Validator) TotalSharesWithDenom(denom string) sdk.Dec {
	return sdk.NewDecCoins(v.TotalShares...).AmountOf(denom)
}

func (v Validator) ValidatorSharesWithDenom(denom string) sdk.Dec {
	return sdk.NewDecCoins(v.ValidatorShares...).AmountOf(denom)
}

func (v Validator) TotalTokensWithAsset(asset AllianceAsset) sdk.Int {
	shares := v.ValidatorSharesWithDenom(asset.Denom)
	return ConvertNewShareToToken(asset.TotalTokens, asset.TotalValidatorShares, shares)
}
