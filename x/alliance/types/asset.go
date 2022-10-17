package types

import cosmosmath "cosmossdk.io/math"

func (asset AllianceAsset) ConvertToStake(amount cosmosmath.Int) (token cosmosmath.Int) {
	token = asset.RewardWeight.MulInt(amount).TruncateInt()
	return
}
