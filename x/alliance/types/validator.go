package types

import (
	cosmosmath "cosmossdk.io/math"
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
	v.TotalDelegatorShares = sdk.DecCoins(v.TotalDelegatorShares).Add(delegationShares...)
	v.ValidatorShares = sdk.DecCoins(v.ValidatorShares).Add(validatorShares...)
}

// ReduceShares handles small inaccuracies (~ < 1) when subtracting shares due to rounding errors
func (v *AllianceValidator) ReduceShares(delegationShares sdk.DecCoins, validatorShares sdk.DecCoins) error {
	newDelegatorShares, err := SubtractDecCoinsWithRounding(v.TotalDelegatorShares, delegationShares)
	if err != nil {
		return err
	}
	v.TotalDelegatorShares = newDelegatorShares

	newValidatorShares, err := SubtractDecCoinsWithRounding(v.ValidatorShares, validatorShares)
	if err != nil {
		return err
	}
	v.ValidatorShares = newValidatorShares

	return nil
}

func SubtractDecCoinsWithRounding(d1s sdk.DecCoins, d2s sdk.DecCoins) (sdk.DecCoins, error) {
	d1Copy := sdk.NewDecCoins(d1s...)
	for _, d2 := range d2s {
		a1 := d1s.AmountOf(d2.Denom)
		a2 := d2.Amount
		// check if the result of the SafeSub is negative ...
		isNegativeResult := false

		if a2.GT(a1) && a2.Sub(a1).LT(sdk.OneDec()) {
			d1Copy, isNegativeResult = d1Copy.SafeSub(sdk.NewDecCoins(sdk.NewDecCoinFromDec(d2.Denom, a1)))
		} else {
			d1Copy, isNegativeResult = d1Copy.SafeSub(sdk.NewDecCoins(d2))
		}

		// ... if the SafeSub returns negative an error should be returned
		if isNegativeResult {
			return nil, ErrInsufficientShares
		}
	}
	return d1Copy, nil
}

func (v AllianceValidator) TotalSharesWithDenom(denom string) sdk.Dec {
	return sdk.DecCoins(v.TotalDelegatorShares).AmountOf(denom)
}

func (v AllianceValidator) ValidatorSharesWithDenom(denom string) sdk.Dec {
	// This is used instead of coins.AmountOf to reduce the need for regex matching to speed up the query
	for _, c := range v.ValidatorShares {
		if c.Denom == denom {
			return c.Amount
		}
	}
	return sdk.ZeroDec()
}

func (v AllianceValidator) TotalDelegationSharesWithDenom(denom string) sdk.Dec {
	return sdk.DecCoins(v.TotalDelegatorShares).AmountOf(denom)
}

func (v AllianceValidator) TotalTokensWithAsset(asset AllianceAsset) sdk.Dec {
	shares := v.ValidatorSharesWithDenom(asset.Denom)
	dec := ConvertNewShareToDecToken(sdk.NewDecFromInt(asset.TotalTokens), asset.TotalValidatorShares, shares)
	return dec
}

func (v AllianceValidator) TotalDecTokensWithAsset(asset AllianceAsset) sdk.Dec {
	shares := v.ValidatorSharesWithDenom(asset.Denom)
	return ConvertNewShareToDecToken(sdk.NewDecFromInt(asset.TotalTokens), asset.TotalValidatorShares, shares)
}

func GetValidatorShares(asset AllianceAsset, token cosmosmath.Int) sdk.Dec {
	return ConvertNewTokenToShares(sdk.NewDecFromInt(asset.TotalTokens), asset.TotalValidatorShares, token)
}
