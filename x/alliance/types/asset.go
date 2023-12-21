package types

import (
	"time"

	cosmosmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewAllianceAsset(denom string, rewardWeight sdk.Dec, minRewardWeight sdk.Dec, maxRewardWeight sdk.Dec, takeRate sdk.Dec, rewardStartTime time.Time) AllianceAsset {
	return AllianceAsset{
		Denom:        denom,
		RewardWeight: rewardWeight,
		RewardWeightRange: RewardWeightRange{
			Min: minRewardWeight,
			Max: maxRewardWeight,
		},
		TakeRate:             takeRate,
		TotalTokens:          sdk.ZeroInt(),
		TotalValidatorShares: sdk.ZeroDec(),
		RewardStartTime:      rewardStartTime,
		RewardChangeRate:     sdk.OneDec(),
		RewardChangeInterval: time.Duration(0),
		LastRewardChangeTime: rewardStartTime,
		IsInitialized:        false,
		IsDissolving:         false,
	}
}

func ConvertNewTokenToShares(totalTokens sdk.Dec, totalShares sdk.Dec, newTokens cosmosmath.Int) (shares sdk.Dec) {
	if totalShares.IsZero() {
		return sdk.NewDecFromInt(newTokens)
	}
	return totalShares.Quo(totalTokens).MulInt(newTokens)
}

func ConvertNewShareToDecToken(totalTokens sdk.Dec, totalShares sdk.Dec, shares sdk.Dec) (token sdk.Dec) {
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
	delTokens = delTokens.Add(sdk.NewDecWithPrec(1, 6))
	return sdk.NewCoin(asset.Denom, delTokens.TruncateInt())
}

func GetDelegationTokensWithShares(delegatorShares sdk.Dec, val AllianceValidator, asset AllianceAsset) sdk.Coin {
	valTokens := val.TotalTokensWithAsset(asset)
	totalDelegationShares := val.TotalDelegationSharesWithDenom(asset.Denom)
	delTokens := ConvertNewShareToDecToken(valTokens, totalDelegationShares, delegatorShares)

	// We add a small epsilon before rounding down to make sure cases like
	// 9.999999 get round to 10
	delTokens = delTokens.Add(sdk.NewDecWithPrec(1, 6))
	return sdk.NewCoin(asset.Denom, delTokens.TruncateInt())
}

func GetDelegationSharesFromTokens(val AllianceValidator, asset AllianceAsset, token cosmosmath.Int) sdk.Dec {
	valTokens := val.TotalTokensWithAsset(asset)
	totalDelegationShares := val.TotalDelegationSharesWithDenom(asset.Denom)
	if totalDelegationShares.TruncateInt().Equal(sdk.ZeroInt()) {
		return sdk.NewDecFromInt(token)
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
