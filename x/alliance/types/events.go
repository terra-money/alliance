package types

const (
	EventTypeDelegate               = "delegate"
	EventTypeUndelegate             = "undelegate"
	EventTypeRedelegate             = "redelegate"
	EventTypeClaimDelegationRewards = "claim_delegation_rewards"

	AttributeKeyValidator         = "validator"
	AttributeKeyCommissionRate    = "commission_rate"
	AttributeKeyMinSelfDelegation = "min_self_delegation"
	AttributeKeySrcValidator      = "source_validator"
	AttributeKeyDstValidator      = "destination_validator"
	AttributeKeyDelegator         = "delegator"
	AttributeKeyCreationHeight    = "creation_height"
	AttributeKeyCompletionTime    = "completion_time"
	AttributeKeyNewShares         = "new_shares"
)
