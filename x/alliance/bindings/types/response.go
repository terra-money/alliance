package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type AllianceResponse struct {
	Denom                string            `json:"denom"`
	RewardWeight         string            `json:"reward_weight"`
	TakeRate             string            `json:"take_rate"`
	TotalTokens          string            `json:"total_tokens"`
	TotalValidatorShares string            `json:"total_validator_shares"`
	RewardStartTime      uint64            `json:"reward_start_time"`
	RewardChangeRate     string            `json:"reward_change_rate"`
	LastRewardChangeTime uint64            `json:"last_reward_change_time"`
	RewardWeightRange    RewardWeightRange `json:"reward_weight_range"`
	IsInitialized        bool              `json:"is_initialized"`
}

type RewardWeightRange struct {
	Min string `json:"min"`
	Max string `json:"max"`
}

type DelegationResponse struct {
	Delegator string `json:"delegator"`
	Validator string `json:"validator"`
	Denom     string `json:"denom"`
	Amount    string `json:"amount"`
}

type DelegationRewardsResponse struct {
	Rewards sdk.Coins `json:"rewards"`
}
