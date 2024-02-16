package types

import (
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewAllianceAsset(denom string, rewardWeight math.LegacyDec, minRewardWeight math.LegacyDec, maxRewardWeight math.LegacyDec, takeRate math.LegacyDec, rewardStartTime time.Time) AllianceAsset {
	return AllianceAsset{
		Denom:        denom,
		RewardWeight: rewardWeight,
		RewardWeightRange: RewardWeightRange{
			Min: minRewardWeight,
			Max: maxRewardWeight,
		},
		TakeRate:             takeRate,
		TotalTokens:          math.ZeroInt(),
		TotalValidatorShares: math.LegacyZeroDec(),
		RewardStartTime:      rewardStartTime,
		RewardChangeRate:     math.LegacyOneDec(),
		RewardChangeInterval: time.Duration(0),
		LastRewardChangeTime: rewardStartTime,
		IsInitialized:        false,
	}
}

func ConvertNewTokenToShares(totalTokens math.LegacyDec, totalShares math.LegacyDec, newTokens math.Int) (shares math.LegacyDec) {
	if totalShares.IsZero() {
		return math.LegacyNewDecFromInt(newTokens)
	}
	return totalShares.Quo(totalTokens).MulInt(newTokens)
}

func ConvertNewShareToDecToken(totalTokens math.LegacyDec, totalShares math.LegacyDec, shares math.LegacyDec) (token math.LegacyDec) {
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
	delTokens = delTokens.Add(math.LegacyNewDecWithPrec(1, 6))
	return sdk.NewCoin(asset.Denom, delTokens.TruncateInt())
}

func GetDelegationTokensWithShares(delegatorShares math.LegacyDec, val AllianceValidator, asset AllianceAsset) sdk.Coin {
	valTokens := val.TotalTokensWithAsset(asset)
	totalDelegationShares := val.TotalDelegationSharesWithDenom(asset.Denom)
	delTokens := ConvertNewShareToDecToken(valTokens, totalDelegationShares, delegatorShares)

	// We add a small epsilon before rounding down to make sure cases like
	// 9.999999 get round to 10
	delTokens = delTokens.Add(math.LegacyNewDecWithPrec(1, 6))
	return sdk.NewCoin(asset.Denom, delTokens.TruncateInt())
}

func GetDelegationSharesFromTokens(val AllianceValidator, asset AllianceAsset, token math.Int) math.LegacyDec {
	valTokens := val.TotalTokensWithAsset(asset)
	totalDelegationShares := val.TotalDelegationSharesWithDenom(asset.Denom)
	if totalDelegationShares.TruncateInt().Equal(math.ZeroInt()) {
		return math.LegacyNewDecFromInt(token)
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
