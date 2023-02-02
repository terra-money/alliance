package types

const (
	EventTypeDelegate               = "alliance_delegate"
	EventTypeUndelegate             = "alliance_undelegate"
	EventTypeRedelegate             = "alliance_redelegate"
	EventTypeClaimDelegationRewards = "alliance_claim_delegation_rewards"

	AttributeKeyValidator      = "validator"
	AttributeKeySrcValidator   = "source_validator"
	AttributeKeyDstValidator   = "destination_validator"
	AttributeKeyCompletionTime = "completion_time"
	AttributeKeyNewShares      = "new_shares"
)
