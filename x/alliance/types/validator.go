package types

import (
	"fmt"
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

// ReduceShares handles small inaccuracies when subtraction shares due to edge case rounding errors
func (v *AllianceValidator) ReduceShares(delegationShares sdk.DecCoins, validatorShares sdk.DecCoins) {
	epsilon := sdk.MustNewDecFromStr("-0.01")
	diffs, hasNeg := sdk.NewDecCoins(v.TotalDelegatorShares...).SafeSub(delegationShares)
	if hasNeg {
		for i, diff := range diffs {
			if diff.IsNegative() {
				if diff.Amount.GTE(epsilon) {
					diffs[i].Amount = sdk.ZeroDec()
				} else {
					panic(fmt.Sprintf("negative shares %s", diff.String()))
				}
			}
		}

	}
	v.TotalDelegatorShares = diffs
	diffs, hasNeg = sdk.NewDecCoins(v.ValidatorShares...).SafeSub(validatorShares)
	if hasNeg {
		for i, diff := range diffs {
			if diff.IsNegative() {
				if diff.Amount.GTE(epsilon) {
					diffs[i].Amount = sdk.ZeroDec()
				} else {
					panic(fmt.Sprintf("negative shares %s", diff.String()))
				}
			}
		}

	}
	v.ValidatorShares = diffs
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
	return ConvertNewShareToToken(asset.TotalTokens, asset.TotalValidatorShares, shares)
}

func (v AllianceValidator) TotalDecTokensWithAsset(asset AllianceAsset) sdk.Dec {
	shares := v.ValidatorSharesWithDenom(asset.Denom)
	return ConvertNewShareToDecToken(asset.TotalTokens, asset.TotalValidatorShares, shares)
}
