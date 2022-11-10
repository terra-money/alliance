package types

func NewRewardWeightChangeSnapshot(asset AllianceAsset, val AllianceValidator) RewardWeightChangeSnapshot {
	return RewardWeightChangeSnapshot{
		PrevRewardWeight: asset.RewardWeight,
		RewardHistories:  val.GlobalRewardHistory,
	}
}
