package types

func NewRewardWeightChangeSnapshot(asset AllianceAsset, val AllianceValidator) RewardWeightChangeSnapshot {
	assetRewardHistory := RewardHistories(val.GlobalRewardHistory).GetIndexByAlliance(asset.Denom)
	return RewardWeightChangeSnapshot{
		PrevRewardWeight: asset.RewardWeight,
		RewardHistories:  assetRewardHistory,
	}
}
