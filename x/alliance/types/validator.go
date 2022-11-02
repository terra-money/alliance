package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type AllianceValidator struct {
	*stakingtypes.Validator
	*AllianceValidatorInfo
}

func NewAllianceValidatorInfo() AllianceValidatorInfo {
	return AllianceValidatorInfo{
		GlobalRewardHistory:  RewardHistories{},
		TotalDelegatorShares: sdk.NewDecCoins(),
		ValidatorShares:      sdk.NewDecCoins(),
	}
}

func (v *AllianceValidator) AddShares(delegationShares sdk.DecCoins, validatorShares sdk.DecCoins) {
	v.TotalDelegatorShares = delegationShares.Add(v.TotalDelegatorShares...)
	v.ValidatorShares = validatorShares.Add(v.ValidatorShares...)
}

// ReduceShares handles small inaccuracies when subtracting shares due to rounding errors
func (v *AllianceValidator) ReduceShares(delegationShares sdk.DecCoins, validatorShares sdk.DecCoins) {
	diffs := SubtractDecCoinsWithRounding(v.TotalDelegatorShares, delegationShares)
	v.TotalDelegatorShares = diffs
	diffs = SubtractDecCoinsWithRounding(v.ValidatorShares, validatorShares)
	v.ValidatorShares = diffs
}

func SubtractDecCoinsWithRounding(d1s sdk.DecCoins, d2s sdk.DecCoins) (d3s sdk.DecCoins) {
	d3s = sdk.NewDecCoins(d1s...)
	for _, d2 := range d2s {
		a1 := d1s.AmountOf(d2.Denom)
		if d2.Amount.GT(a1) && d2.Amount.Sub(a1).LT(sdk.OneDec()) {
			d3s = d3s.Sub(sdk.NewDecCoins(sdk.NewDecCoinFromDec(d2.Denom, a1)))
		} else {
			d3s = d3s.Sub(sdk.NewDecCoins(d2))
		}
	}
	return d3s
}

func (v AllianceValidator) TotalSharesWithDenom(denom string) sdk.Dec {
	return sdk.NewDecCoins(v.TotalDelegatorShares...).AmountOf(denom)
}

func (v AllianceValidator) ValidatorSharesWithDenom(denom string) sdk.Dec {
	return sdk.NewDecCoins(v.ValidatorShares...).AmountOf(denom)
}

func (v AllianceValidator) TotalDelegationSharesWithDenom(denom string) sdk.Dec {
	return sdk.NewDecCoins(v.TotalDelegatorShares...).AmountOf(denom)
}

func (v AllianceValidator) TotalTokensWithAsset(asset AllianceAsset) sdk.Int {
	shares := v.ValidatorSharesWithDenom(asset.Denom)
	return ConvertNewShareToDecToken(sdk.NewDecFromInt(asset.TotalTokens), asset.TotalValidatorShares, shares).TruncateInt()
}

func (v AllianceValidator) TotalDecTokensWithAsset(asset AllianceAsset) sdk.Dec {
	shares := v.ValidatorSharesWithDenom(asset.Denom)
	return ConvertNewShareToDecToken(sdk.NewDecFromInt(asset.TotalTokens), asset.TotalValidatorShares, shares)
}
