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

func (v *AllianceValidator) ReduceShares(delegationShares sdk.DecCoins, validatorShares sdk.DecCoins) {
	v.TotalDelegatorShares = sdk.NewDecCoins(v.TotalDelegatorShares...).Sub(delegationShares)
	v.ValidatorShares = sdk.NewDecCoins(v.ValidatorShares...).Sub(validatorShares)
}

func (v AllianceValidator) TotalSharesWithDenom(denom string) sdk.Dec {
	return sdk.NewDecCoins(v.TotalDelegatorShares...).AmountOf(denom)
}

func (v AllianceValidator) ValidatorSharesWithDenom(denom string) sdk.Dec {
	return sdk.NewDecCoins(v.ValidatorShares...).AmountOf(denom)
}

func (v AllianceValidator) TotalTokensWithAsset(asset AllianceAsset) sdk.Int {
	shares := v.ValidatorSharesWithDenom(asset.Denom)
	return ConvertNewShareToToken(asset.TotalTokens, asset.TotalValidatorShares, shares)
}
