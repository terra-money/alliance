package types

import (
	"cosmossdk.io/math"
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
func (v *AllianceValidator) ReduceShares(delegationShares sdk.DecCoins, validatorShares sdk.DecCoins) {
	newDelegatorShares := SubtractDecCoinsWithRounding(v.TotalDelegatorShares, delegationShares)
	v.TotalDelegatorShares = newDelegatorShares
	newValidatorShares := SubtractDecCoinsWithRounding(v.ValidatorShares, validatorShares)
	v.ValidatorShares = newValidatorShares
}

func SubtractDecCoinsWithRounding(d1s sdk.DecCoins, d2s sdk.DecCoins) sdk.DecCoins {
	d1Copy := sdk.NewDecCoins(d1s...)
	for _, d2 := range d2s {
		a1 := d1s.AmountOf(d2.Denom)
		a2 := d2.Amount
		if a2.GT(a1) && a2.Sub(a1).LT(math.LegacyOneDec()) {
			d1Copy = d1Copy.Sub(sdk.NewDecCoins(sdk.NewDecCoinFromDec(d2.Denom, a1)))
		} else {
			d1Copy = d1Copy.Sub(sdk.NewDecCoins(d2))
		}
	}
	return d1Copy
}

func (v AllianceValidator) ValidatorSharesWithDenom(denom string) math.LegacyDec {
	// This is used instead of coins.AmountOf to reduce the need for regex matching to speed up the query
	for _, c := range v.ValidatorShares {
		if c.Denom == denom {
			return c.Amount
		}
	}
	return math.LegacyZeroDec()
}

func (v AllianceValidator) TotalDelegationSharesWithDenom(denom string) math.LegacyDec {
	return sdk.DecCoins(v.TotalDelegatorShares).AmountOf(denom)
}

func (v AllianceValidator) TotalTokensWithAsset(asset AllianceAsset) math.LegacyDec {
	shares := v.ValidatorSharesWithDenom(asset.Denom)
	dec := ConvertNewShareToDecToken(math.LegacyNewDecFromInt(asset.TotalTokens), asset.TotalValidatorShares, shares)
	return dec
}

func GetValidatorShares(asset AllianceAsset, token math.Int) math.LegacyDec {
	return ConvertNewTokenToShares(math.LegacyNewDecFromInt(asset.TotalTokens), asset.TotalValidatorShares, token)
}

func (v AllianceValidator) GetValAddress() (sdk.ValAddress, error) {
	return sdk.ValAddressFromBech32(v.GetOperator())
}
