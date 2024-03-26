package types

import (
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewAllianceAsset(denom string, rewardWeight sdkmath.LegacyDec, minRewardWeight sdkmath.LegacyDec, maxRewardWeight sdkmath.LegacyDec, takeRate sdkmath.LegacyDec, rewardStartTime time.Time) AllianceAsset {
	return AllianceAsset{
		Denom:        denom,
		RewardWeight: rewardWeight,
		RewardWeightRange: RewardWeightRange{
			Min: minRewardWeight,
			Max: maxRewardWeight,
		},
		TakeRate:             takeRate,
		TotalTokens:          sdkmath.ZeroInt(),
		TotalValidatorShares: sdkmath.LegacyZeroDec(),
		RewardStartTime:      rewardStartTime,
		RewardChangeRate:     sdkmath.LegacyOneDec(),
		RewardChangeInterval: time.Duration(0),
		LastRewardChangeTime: rewardStartTime,
		IsInitialized:        false,
	}
}

func ConvertNewTokenToShares(totalTokens sdkmath.LegacyDec, totalShares sdkmath.LegacyDec, newTokens sdkmath.Int) (shares sdkmath.LegacyDec) {
	if totalShares.IsZero() {
		return sdkmath.LegacyNewDecFromInt(newTokens)
	}
	return totalShares.Quo(totalTokens).MulInt(newTokens)
}

func ConvertNewShareToDecToken(totalTokens sdkmath.LegacyDec, totalShares sdkmath.LegacyDec, shares sdkmath.LegacyDec) (token sdkmath.LegacyDec) {
	if totalShares.IsZero() {
		return totalTokens
	}
	return shares.Quo(totalShares).Mul(totalTokens)
}

func GetDelegationTokens(del Delegation, val AllianceValidator, asset AllianceAsset) sdk.Coin {
	valTokens := val.TotalTokensWithAsset(asset)
	totalDelegationShares := val.TotalDelegationSharesWithDenom(asset.Denom)
	delTokens := ConvertNewShareToDecToken(valTokens, totalDelegationShares, del.Shares)

	// We add a small epsilon before rounding down to make sure cases like
	// 9.999999 get round to 10
	delTokens = delTokens.Add(sdkmath.LegacyNewDecWithPrec(1, 6))
	return sdk.NewCoin(asset.Denom, delTokens.TruncateInt())
}

func GetDelegationTokensWithShares(delegatorShares sdkmath.LegacyDec, val AllianceValidator, asset AllianceAsset) sdk.Coin {
	valTokens := val.TotalTokensWithAsset(asset)
	totalDelegationShares := val.TotalDelegationSharesWithDenom(asset.Denom)
	delTokens := ConvertNewShareToDecToken(valTokens, totalDelegationShares, delegatorShares)

	// We add a small epsilon before rounding down to make sure cases like
	// 9.999999 get round to 10
	delTokens = delTokens.Add(sdkmath.LegacyNewDecWithPrec(1, 6))
	return sdk.NewCoin(asset.Denom, delTokens.TruncateInt())
}

func GetDelegationSharesFromTokens(val AllianceValidator, asset AllianceAsset, token sdkmath.Int) sdkmath.LegacyDec {
	valTokens := val.TotalTokensWithAsset(asset)
	totalDelegationShares := val.TotalDelegationSharesWithDenom(asset.Denom)
	if totalDelegationShares.TruncateInt().Equal(sdkmath.ZeroInt()) {
		return sdkmath.LegacyNewDecFromInt(token)
	}
	return ConvertNewTokenToShares(valTokens, totalDelegationShares, token)
}

func (a AllianceAsset) HasPositiveDecay() bool {
	return a.RewardChangeInterval > 0 && a.RewardChangeRate.IsPositive()
}

// RewardsStarted helper function to check if rewards for the alliance has started
func (a AllianceAsset) RewardsStarted(blockTime time.Time) bool {
	return blockTime.After(a.RewardStartTime) || blockTime.Equal(a.RewardStartTime)
}
