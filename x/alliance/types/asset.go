package types

import (
	cosmosmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

func NewAllianceAsset(denom string, rewardWeight sdk.Dec, takeRate sdk.Dec, rewardStartTime time.Time) AllianceAsset {
	return AllianceAsset{
		Denom:                denom,
		RewardWeight:         rewardWeight,
		TakeRate:             takeRate,
		TotalTokens:          sdk.ZeroInt(),
		TotalValidatorShares: sdk.ZeroDec(),
		RewardStartTime:      rewardStartTime,
	}
}

func ConvertNewTokenToShares(totalTokens cosmosmath.Int, totalShares sdk.Dec, newTokens cosmosmath.Int) (shares sdk.Dec) {
	// TODO: Verify this logic when totalShares != 0 or totalTotals != 0
	if totalShares.IsZero() || totalTokens.IsZero() {
		return sdk.NewDecFromInt(newTokens)
	}
	return totalShares.MulInt(newTokens).QuoInt(totalTokens)
}

func ConvertNewShareToToken(totalTokens cosmosmath.Int, totalShares sdk.Dec, shares sdk.Dec) (token cosmosmath.Int) {
	if totalShares.IsZero() {
		return totalTokens
	}
	return shares.MulInt(totalTokens).Quo(totalShares).TruncateInt()
}

func ConvertNewShareToDecToken(totalTokens cosmosmath.Int, totalShares sdk.Dec, shares sdk.Dec) (token sdk.Dec) {
	if totalShares.IsZero() {
		return sdk.NewDecFromInt(totalTokens)
	}
	return shares.MulInt(totalTokens).Quo(totalShares)
}

func GetDelegationTokens(del Delegation, val AllianceValidator, asset AllianceAsset) sdk.Coin {
	valTokens := val.TotalTokensWithAsset(asset)
	totalDelegationShares := val.TotalDelegationSharesWithDenom(asset.Denom)
	delTokens := ConvertNewShareToToken(valTokens, totalDelegationShares, del.Shares)
	return sdk.NewCoin(asset.Denom, delTokens)
}

func GetDelegationSharesFromTokens(val AllianceValidator, asset AllianceAsset, token cosmosmath.Int) sdk.Dec {
	valTokens := val.TotalTokensWithAsset(asset)
	totalDelegationShares := val.TotalDelegationSharesWithDenom(asset.Denom)
	return ConvertNewTokenToShares(valTokens, totalDelegationShares, token)
}

func GetValidatorShares(asset AllianceAsset, token cosmosmath.Int) sdk.Dec {
	return ConvertNewTokenToShares(asset.TotalTokens, asset.TotalValidatorShares, token)
}
