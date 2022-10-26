package types

import (
	cosmosmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewAsset(denom string, rewardWeight sdk.Dec, takeRate sdk.Dec) AllianceAsset {
	return AllianceAsset{
		Denom:                denom,
		RewardWeight:         rewardWeight,
		TakeRate:             takeRate,
		TotalTokens:          sdk.ZeroInt(),
		TotalValidatorShares: sdk.ZeroDec(),
		TotalStakeTokens:     sdk.ZeroInt(),
	}
}

func (asset AllianceAsset) ConvertToStake(amount cosmosmath.Int) (token cosmosmath.Int) {
	token = asset.RewardWeight.MulInt(amount).TruncateInt()
	return
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

func GetDelegationTokens(del Delegation, val AllianceValidator, asset AllianceAsset) sdk.Coin {
	valTokens := val.TotalTokensWithAsset(asset)
	valShares := val.ValidatorSharesWithDenom(asset.Denom)
	delTokens := ConvertNewShareToToken(valTokens, valShares, del.Shares)
	return sdk.NewCoin(asset.Denom, delTokens)
}

func GetDelegationSharesFromTokens(val AllianceValidator, asset AllianceAsset, token cosmosmath.Int) sdk.Dec {
	valTokens := val.TotalTokensWithAsset(asset)
	valShares := val.ValidatorSharesWithDenom(asset.Denom)
	return ConvertNewTokenToShares(valTokens, valShares, token)
}

func GetValidatorShares(asset AllianceAsset, token cosmosmath.Int) sdk.Dec {
	return ConvertNewTokenToShares(asset.TotalTokens, asset.TotalValidatorShares, token)
}
