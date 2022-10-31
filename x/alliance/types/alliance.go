package types

func NewRewardRateChangeSnapshot(asset AllianceAsset, val AllianceValidator) RewardRateChangeSnapshot {
	return RewardRateChangeSnapshot{
		PrevRewardWeight: asset.RewardWeight,
		RewardHistories:  val.GlobalRewardHistory,
	}
}
