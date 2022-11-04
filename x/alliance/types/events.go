package types

const (
	EventTypeDelegate               = "delegate"
	EventTypeUndelegate             = "undelegate"
	EventTypeRedelegate             = "redelegate"
	EventTypeClaimDelegationRewards = "claim_delegation_rewards"

	AttributeKeyValidator      = "validator"
	AttributeKeySrcValidator   = "source_validator"
	AttributeKeyDstValidator   = "destination_validator"
	AttributeKeyCompletionTime = "completion_time"
	AttributeKeyNewShares      = "new_shares"
)
