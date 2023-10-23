<!-- This file is auto-generated. Please do not modify it yourself. -->
# Protobuf Documentation
<a name="top"></a>

## Table of Contents

- [alliance/alliance/params.proto](#alliance/alliance/params.proto)
    - [Params](#alliance.alliance.Params)
    - [RewardHistory](#alliance.alliance.RewardHistory)
  
- [alliance/alliance/alliance.proto](#alliance/alliance/alliance.proto)
    - [AllianceAsset](#alliance.alliance.AllianceAsset)
    - [RewardWeightChangeSnapshot](#alliance.alliance.RewardWeightChangeSnapshot)
    - [RewardWeightRange](#alliance.alliance.RewardWeightRange)
  
- [alliance/alliance/delegations.proto](#alliance/alliance/delegations.proto)
    - [AllianceValidatorInfo](#alliance.alliance.AllianceValidatorInfo)
    - [Delegation](#alliance.alliance.Delegation)
    - [QueuedUndelegation](#alliance.alliance.QueuedUndelegation)
    - [Undelegation](#alliance.alliance.Undelegation)
  
- [alliance/alliance/events.proto](#alliance/alliance/events.proto)
    - [ClaimAllianceRewardsEvent](#alliance.alliance.ClaimAllianceRewardsEvent)
    - [DeductAllianceAssetsEvent](#alliance.alliance.DeductAllianceAssetsEvent)
    - [DelegateAllianceEvent](#alliance.alliance.DelegateAllianceEvent)
    - [RedelegateAllianceEvent](#alliance.alliance.RedelegateAllianceEvent)
    - [UndelegateAllianceEvent](#alliance.alliance.UndelegateAllianceEvent)
  
- [alliance/alliance/redelegations.proto](#alliance/alliance/redelegations.proto)
    - [QueuedRedelegation](#alliance.alliance.QueuedRedelegation)
    - [Redelegation](#alliance.alliance.Redelegation)
    - [RedelegationEntry](#alliance.alliance.RedelegationEntry)
  
- [alliance/alliance/genesis.proto](#alliance/alliance/genesis.proto)
    - [GenesisState](#alliance.alliance.GenesisState)
    - [RedelegationState](#alliance.alliance.RedelegationState)
    - [RewardWeightChangeSnapshotState](#alliance.alliance.RewardWeightChangeSnapshotState)
    - [UndelegationState](#alliance.alliance.UndelegationState)
    - [ValidatorInfoState](#alliance.alliance.ValidatorInfoState)
  
- [alliance/alliance/gov.proto](#alliance/alliance/gov.proto)
    - [MsgCreateAllianceProposal](#alliance.alliance.MsgCreateAllianceProposal)
    - [MsgDeleteAllianceProposal](#alliance.alliance.MsgDeleteAllianceProposal)
    - [MsgUpdateAllianceProposal](#alliance.alliance.MsgUpdateAllianceProposal)
  
- [alliance/alliance/unbonding.proto](#alliance/alliance/unbonding.proto)
    - [UnbondingDelegation](#alliance.alliance.UnbondingDelegation)
  
- [alliance/alliance/query.proto](#alliance/alliance/query.proto)
    - [DelegationResponse](#alliance.alliance.DelegationResponse)
    - [QueryAllAllianceValidatorsRequest](#alliance.alliance.QueryAllAllianceValidatorsRequest)
    - [QueryAllAlliancesDelegationsRequest](#alliance.alliance.QueryAllAlliancesDelegationsRequest)
    - [QueryAllianceDelegationRequest](#alliance.alliance.QueryAllianceDelegationRequest)
    - [QueryAllianceDelegationResponse](#alliance.alliance.QueryAllianceDelegationResponse)
    - [QueryAllianceDelegationRewardsRequest](#alliance.alliance.QueryAllianceDelegationRewardsRequest)
    - [QueryAllianceDelegationRewardsResponse](#alliance.alliance.QueryAllianceDelegationRewardsResponse)
    - [QueryAllianceRedelegationsRequest](#alliance.alliance.QueryAllianceRedelegationsRequest)
    - [QueryAllianceRedelegationsResponse](#alliance.alliance.QueryAllianceRedelegationsResponse)
    - [QueryAllianceRequest](#alliance.alliance.QueryAllianceRequest)
    - [QueryAllianceResponse](#alliance.alliance.QueryAllianceResponse)
    - [QueryAllianceUnbondingsByDenomAndDelegatorRequest](#alliance.alliance.QueryAllianceUnbondingsByDenomAndDelegatorRequest)
    - [QueryAllianceUnbondingsByDenomAndDelegatorResponse](#alliance.alliance.QueryAllianceUnbondingsByDenomAndDelegatorResponse)
    - [QueryAllianceUnbondingsRequest](#alliance.alliance.QueryAllianceUnbondingsRequest)
    - [QueryAllianceUnbondingsResponse](#alliance.alliance.QueryAllianceUnbondingsResponse)
    - [QueryAllianceValidatorRequest](#alliance.alliance.QueryAllianceValidatorRequest)
    - [QueryAllianceValidatorResponse](#alliance.alliance.QueryAllianceValidatorResponse)
    - [QueryAllianceValidatorsResponse](#alliance.alliance.QueryAllianceValidatorsResponse)
    - [QueryAlliancesDelegationByValidatorRequest](#alliance.alliance.QueryAlliancesDelegationByValidatorRequest)
    - [QueryAlliancesDelegationsRequest](#alliance.alliance.QueryAlliancesDelegationsRequest)
    - [QueryAlliancesDelegationsResponse](#alliance.alliance.QueryAlliancesDelegationsResponse)
    - [QueryAlliancesRequest](#alliance.alliance.QueryAlliancesRequest)
    - [QueryAlliancesResponse](#alliance.alliance.QueryAlliancesResponse)
    - [QueryIBCAllianceDelegationRequest](#alliance.alliance.QueryIBCAllianceDelegationRequest)
    - [QueryIBCAllianceDelegationRewardsRequest](#alliance.alliance.QueryIBCAllianceDelegationRewardsRequest)
    - [QueryIBCAllianceRequest](#alliance.alliance.QueryIBCAllianceRequest)
    - [QueryParamsRequest](#alliance.alliance.QueryParamsRequest)
    - [QueryParamsResponse](#alliance.alliance.QueryParamsResponse)
  
    - [Query](#alliance.alliance.Query)
  
- [alliance/alliance/tx.proto](#alliance/alliance/tx.proto)
    - [MsgClaimDelegationRewards](#alliance.alliance.MsgClaimDelegationRewards)
    - [MsgClaimDelegationRewardsResponse](#alliance.alliance.MsgClaimDelegationRewardsResponse)
    - [MsgCreateAlliance](#alliance.alliance.MsgCreateAlliance)
    - [MsgCreateAllianceResponse](#alliance.alliance.MsgCreateAllianceResponse)
    - [MsgDelegate](#alliance.alliance.MsgDelegate)
    - [MsgDelegateResponse](#alliance.alliance.MsgDelegateResponse)
    - [MsgDeleteAlliance](#alliance.alliance.MsgDeleteAlliance)
    - [MsgDeleteAllianceResponse](#alliance.alliance.MsgDeleteAllianceResponse)
    - [MsgRedelegate](#alliance.alliance.MsgRedelegate)
    - [MsgRedelegateResponse](#alliance.alliance.MsgRedelegateResponse)
    - [MsgUndelegate](#alliance.alliance.MsgUndelegate)
    - [MsgUndelegateResponse](#alliance.alliance.MsgUndelegateResponse)
    - [MsgUpdateAlliance](#alliance.alliance.MsgUpdateAlliance)
    - [MsgUpdateAllianceResponse](#alliance.alliance.MsgUpdateAllianceResponse)
    - [MsgUpdateParams](#alliance.alliance.MsgUpdateParams)
    - [MsgUpdateParamsResponse](#alliance.alliance.MsgUpdateParamsResponse)
  
    - [Msg](#alliance.alliance.Msg)
  
- [Scalar Value Types](#scalar-value-types)



<a name="alliance/alliance/params.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## alliance/alliance/params.proto



<a name="alliance.alliance.Params"></a>

### Params



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `reward_delay_time` | [google.protobuf.Duration](#google.protobuf.Duration) |  |  |
| `take_rate_claim_interval` | [google.protobuf.Duration](#google.protobuf.Duration) |  | Time interval between consecutive applications of `take_rate` |
| `last_take_rate_claim_time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | Last application of `take_rate` on assets |






<a name="alliance.alliance.RewardHistory"></a>

### RewardHistory



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [string](#string) |  |  |
| `index` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="alliance/alliance/alliance.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## alliance/alliance/alliance.proto



<a name="alliance.alliance.AllianceAsset"></a>

### AllianceAsset
key: denom value: AllianceAsset


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [string](#string) |  | Denom of the asset. It could either be a native token or an IBC token |
| `reward_weight` | [string](#string) |  | The reward weight specifies the ratio of rewards that will be given to each alliance asset It does not need to sum to 1. rate = weight / total_weight Native asset is always assumed to have a weight of 1.s |
| `take_rate` | [string](#string) |  | A positive take rate is used for liquid staking derivatives. It defines an rate that is applied per take_rate_interval that will be redirected to the distribution rewards pool |
| `total_tokens` | [string](#string) |  |  |
| `total_validator_shares` | [string](#string) |  |  |
| `reward_start_time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |
| `reward_change_rate` | [string](#string) |  |  |
| `reward_change_interval` | [google.protobuf.Duration](#google.protobuf.Duration) |  |  |
| `last_reward_change_time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |
| `reward_weight_range` | [RewardWeightRange](#alliance.alliance.RewardWeightRange) |  | set a bound of weight range to limit how much reward weights can scale. |
| `is_initialized` | [bool](#bool) |  | flag to check if an asset has completed the initialization process after the reward delay |






<a name="alliance.alliance.RewardWeightChangeSnapshot"></a>

### RewardWeightChangeSnapshot



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `prev_reward_weight` | [string](#string) |  |  |
| `reward_histories` | [RewardHistory](#alliance.alliance.RewardHistory) | repeated |  |






<a name="alliance.alliance.RewardWeightRange"></a>

### RewardWeightRange



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `min` | [string](#string) |  |  |
| `max` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="alliance/alliance/delegations.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## alliance/alliance/delegations.proto



<a name="alliance.alliance.AllianceValidatorInfo"></a>

### AllianceValidatorInfo



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `global_reward_history` | [RewardHistory](#alliance.alliance.RewardHistory) | repeated |  |
| `total_delegator_shares` | [cosmos.base.v1beta1.DecCoin](#cosmos.base.v1beta1.DecCoin) | repeated |  |
| `validator_shares` | [cosmos.base.v1beta1.DecCoin](#cosmos.base.v1beta1.DecCoin) | repeated |  |






<a name="alliance.alliance.Delegation"></a>

### Delegation



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  | delegator_address is the bech32-encoded address of the delegator. |
| `validator_address` | [string](#string) |  | validator_address is the bech32-encoded address of the validator. |
| `denom` | [string](#string) |  | denom of token staked |
| `shares` | [string](#string) |  | shares define the delegation shares received. |
| `reward_history` | [RewardHistory](#alliance.alliance.RewardHistory) | repeated |  |
| `last_reward_claim_height` | [uint64](#uint64) |  |  |






<a name="alliance.alliance.QueuedUndelegation"></a>

### QueuedUndelegation



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `entries` | [Undelegation](#alliance.alliance.Undelegation) | repeated |  |






<a name="alliance.alliance.Undelegation"></a>

### Undelegation



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  |  |
| `validator_address` | [string](#string) |  |  |
| `balance` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="alliance/alliance/events.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## alliance/alliance/events.proto



<a name="alliance.alliance.ClaimAllianceRewardsEvent"></a>

### ClaimAllianceRewardsEvent



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `allianceSender` | [string](#string) |  |  |
| `validator` | [string](#string) |  |  |
| `coins` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |






<a name="alliance.alliance.DeductAllianceAssetsEvent"></a>

### DeductAllianceAssetsEvent



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `coins` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |






<a name="alliance.alliance.DelegateAllianceEvent"></a>

### DelegateAllianceEvent



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `allianceSender` | [string](#string) |  |  |
| `validator` | [string](#string) |  |  |
| `coin` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |
| `newShares` | [string](#string) |  |  |






<a name="alliance.alliance.RedelegateAllianceEvent"></a>

### RedelegateAllianceEvent



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `allianceSender` | [string](#string) |  |  |
| `sourceValidator` | [string](#string) |  |  |
| `destinationValidator` | [string](#string) |  |  |
| `coin` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |
| `completionTime` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |






<a name="alliance.alliance.UndelegateAllianceEvent"></a>

### UndelegateAllianceEvent



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `allianceSender` | [string](#string) |  |  |
| `validator` | [string](#string) |  |  |
| `coin` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |
| `completionTime` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="alliance/alliance/redelegations.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## alliance/alliance/redelegations.proto



<a name="alliance.alliance.QueuedRedelegation"></a>

### QueuedRedelegation
Used internally to keep track of redelegations


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `entries` | [Redelegation](#alliance.alliance.Redelegation) | repeated |  |






<a name="alliance.alliance.Redelegation"></a>

### Redelegation



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  | internal or external user address |
| `src_validator_address` | [string](#string) |  | redelegation source validator |
| `dst_validator_address` | [string](#string) |  | redelegation destination validator |
| `balance` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  | amount to redelegate |






<a name="alliance.alliance.RedelegationEntry"></a>

### RedelegationEntry
Used on QueryServer


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  | internal or external user address |
| `src_validator_address` | [string](#string) |  | redelegation source validator |
| `dst_validator_address` | [string](#string) |  | redelegation destination validator |
| `balance` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  | amount to redelegate |
| `completion_time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | completion_time defines the unix time for redelegation completion. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="alliance/alliance/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## alliance/alliance/genesis.proto



<a name="alliance.alliance.GenesisState"></a>

### GenesisState
GenesisState defines the module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#alliance.alliance.Params) |  |  |
| `assets` | [AllianceAsset](#alliance.alliance.AllianceAsset) | repeated |  |
| `validator_infos` | [ValidatorInfoState](#alliance.alliance.ValidatorInfoState) | repeated |  |
| `reward_weight_change_snaphots` | [RewardWeightChangeSnapshotState](#alliance.alliance.RewardWeightChangeSnapshotState) | repeated |  |
| `delegations` | [Delegation](#alliance.alliance.Delegation) | repeated |  |
| `redelegations` | [RedelegationState](#alliance.alliance.RedelegationState) | repeated |  |
| `undelegations` | [UndelegationState](#alliance.alliance.UndelegationState) | repeated |  |






<a name="alliance.alliance.RedelegationState"></a>

### RedelegationState



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `completion_time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |
| `redelegation` | [Redelegation](#alliance.alliance.Redelegation) |  |  |






<a name="alliance.alliance.RewardWeightChangeSnapshotState"></a>

### RewardWeightChangeSnapshotState



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `height` | [uint64](#uint64) |  |  |
| `validator` | [string](#string) |  |  |
| `denom` | [string](#string) |  |  |
| `snapshot` | [RewardWeightChangeSnapshot](#alliance.alliance.RewardWeightChangeSnapshot) |  |  |






<a name="alliance.alliance.UndelegationState"></a>

### UndelegationState



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `completion_time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |
| `undelegation` | [QueuedUndelegation](#alliance.alliance.QueuedUndelegation) |  |  |






<a name="alliance.alliance.ValidatorInfoState"></a>

### ValidatorInfoState



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_address` | [string](#string) |  |  |
| `validator` | [AllianceValidatorInfo](#alliance.alliance.AllianceValidatorInfo) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="alliance/alliance/gov.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## alliance/alliance/gov.proto



<a name="alliance.alliance.MsgCreateAllianceProposal"></a>

### MsgCreateAllianceProposal



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  | the title of the update proposal |
| `description` | [string](#string) |  | the description of the proposal |
| `denom` | [string](#string) |  | Denom of the asset. It could either be a native token or an IBC token |
| `reward_weight` | [string](#string) |  | The reward weight specifies the ratio of rewards that will be given to each alliance asset It does not need to sum to 1. rate = weight / total_weight Native asset is always assumed to have a weight of 1. |
| `take_rate` | [string](#string) |  | A positive take rate is used for liquid staking derivatives. It defines an annualized reward rate that will be redirected to the distribution rewards pool |
| `reward_change_rate` | [string](#string) |  |  |
| `reward_change_interval` | [google.protobuf.Duration](#google.protobuf.Duration) |  |  |
| `reward_weight_range` | [RewardWeightRange](#alliance.alliance.RewardWeightRange) |  | set a bound of weight range to limit how much reward weights can scale. |






<a name="alliance.alliance.MsgDeleteAllianceProposal"></a>

### MsgDeleteAllianceProposal



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  | the title of the update proposal |
| `description` | [string](#string) |  | the description of the proposal |
| `denom` | [string](#string) |  |  |






<a name="alliance.alliance.MsgUpdateAllianceProposal"></a>

### MsgUpdateAllianceProposal



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  | the title of the update proposal |
| `description` | [string](#string) |  | the description of the proposal |
| `denom` | [string](#string) |  | Denom of the asset. It could either be a native token or an IBC token |
| `reward_weight` | [string](#string) |  | The reward weight specifies the ratio of rewards that will be given to each alliance asset It does not need to sum to 1. rate = weight / total_weight Native asset is always assumed to have a weight of 1. |
| `take_rate` | [string](#string) |  |  |
| `reward_change_rate` | [string](#string) |  |  |
| `reward_change_interval` | [google.protobuf.Duration](#google.protobuf.Duration) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="alliance/alliance/unbonding.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## alliance/alliance/unbonding.proto



<a name="alliance.alliance.UnbondingDelegation"></a>

### UnbondingDelegation
UnbondingDelegation defines an unbonding object with relevant metadata.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `completion_time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | completion_time is the unix time for unbonding completion. |
| `validator_address` | [string](#string) |  | validator_address is the bech32-encoded address of the validator. |
| `amount` | [string](#string) |  | amount defines the tokens to receive at completion. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="alliance/alliance/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## alliance/alliance/query.proto



<a name="alliance.alliance.DelegationResponse"></a>

### DelegationResponse
DelegationResponse is equivalent to Delegation except that it contains a
balance in addition to shares which is more suitable for client responses.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegation` | [Delegation](#alliance.alliance.Delegation) |  |  |
| `balance` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |






<a name="alliance.alliance.QueryAllAllianceValidatorsRequest"></a>

### QueryAllAllianceValidatorsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  |  |






<a name="alliance.alliance.QueryAllAlliancesDelegationsRequest"></a>

### QueryAllAlliancesDelegationsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  |  |






<a name="alliance.alliance.QueryAllianceDelegationRequest"></a>

### QueryAllianceDelegationRequest
AllianceDelegation


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_addr` | [string](#string) |  |  |
| `validator_addr` | [string](#string) |  |  |
| `denom` | [string](#string) |  |  |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  |  |






<a name="alliance.alliance.QueryAllianceDelegationResponse"></a>

### QueryAllianceDelegationResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegation` | [DelegationResponse](#alliance.alliance.DelegationResponse) |  |  |






<a name="alliance.alliance.QueryAllianceDelegationRewardsRequest"></a>

### QueryAllianceDelegationRewardsRequest
AllianceDelegation


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_addr` | [string](#string) |  |  |
| `validator_addr` | [string](#string) |  |  |
| `denom` | [string](#string) |  |  |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  |  |






<a name="alliance.alliance.QueryAllianceDelegationRewardsResponse"></a>

### QueryAllianceDelegationRewardsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `rewards` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |






<a name="alliance.alliance.QueryAllianceRedelegationsRequest"></a>

### QueryAllianceRedelegationsRequest
Redelegations


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [string](#string) |  |  |
| `delegator_addr` | [string](#string) |  |  |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  |  |






<a name="alliance.alliance.QueryAllianceRedelegationsResponse"></a>

### QueryAllianceRedelegationsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `redelegations` | [RedelegationEntry](#alliance.alliance.RedelegationEntry) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  |  |






<a name="alliance.alliance.QueryAllianceRequest"></a>

### QueryAllianceRequest
Alliance


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [string](#string) |  |  |






<a name="alliance.alliance.QueryAllianceResponse"></a>

### QueryAllianceResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `alliance` | [AllianceAsset](#alliance.alliance.AllianceAsset) |  |  |






<a name="alliance.alliance.QueryAllianceUnbondingsByDenomAndDelegatorRequest"></a>

### QueryAllianceUnbondingsByDenomAndDelegatorRequest
AllianceDelegation


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [string](#string) |  |  |
| `delegator_addr` | [string](#string) |  |  |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  |  |






<a name="alliance.alliance.QueryAllianceUnbondingsByDenomAndDelegatorResponse"></a>

### QueryAllianceUnbondingsByDenomAndDelegatorResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `unbondings` | [UnbondingDelegation](#alliance.alliance.UnbondingDelegation) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  |  |






<a name="alliance.alliance.QueryAllianceUnbondingsRequest"></a>

### QueryAllianceUnbondingsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [string](#string) |  |  |
| `delegator_addr` | [string](#string) |  |  |
| `validator_addr` | [string](#string) |  |  |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  |  |






<a name="alliance.alliance.QueryAllianceUnbondingsResponse"></a>

### QueryAllianceUnbondingsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `unbondings` | [UnbondingDelegation](#alliance.alliance.UnbondingDelegation) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  |  |






<a name="alliance.alliance.QueryAllianceValidatorRequest"></a>

### QueryAllianceValidatorRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_addr` | [string](#string) |  |  |






<a name="alliance.alliance.QueryAllianceValidatorResponse"></a>

### QueryAllianceValidatorResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_addr` | [string](#string) |  |  |
| `total_delegation_shares` | [cosmos.base.v1beta1.DecCoin](#cosmos.base.v1beta1.DecCoin) | repeated |  |
| `validator_shares` | [cosmos.base.v1beta1.DecCoin](#cosmos.base.v1beta1.DecCoin) | repeated |  |
| `total_staked` | [cosmos.base.v1beta1.DecCoin](#cosmos.base.v1beta1.DecCoin) | repeated |  |






<a name="alliance.alliance.QueryAllianceValidatorsResponse"></a>

### QueryAllianceValidatorsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validators` | [QueryAllianceValidatorResponse](#alliance.alliance.QueryAllianceValidatorResponse) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  |  |






<a name="alliance.alliance.QueryAlliancesDelegationByValidatorRequest"></a>

### QueryAlliancesDelegationByValidatorRequest
AlliancesDelegationByValidator


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_addr` | [string](#string) |  |  |
| `validator_addr` | [string](#string) |  |  |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  |  |






<a name="alliance.alliance.QueryAlliancesDelegationsRequest"></a>

### QueryAlliancesDelegationsRequest
AlliancesDelegation


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_addr` | [string](#string) |  |  |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  |  |






<a name="alliance.alliance.QueryAlliancesDelegationsResponse"></a>

### QueryAlliancesDelegationsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegations` | [DelegationResponse](#alliance.alliance.DelegationResponse) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  |  |






<a name="alliance.alliance.QueryAlliancesRequest"></a>

### QueryAlliancesRequest
Alliances


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  |  |






<a name="alliance.alliance.QueryAlliancesResponse"></a>

### QueryAlliancesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `alliances` | [AllianceAsset](#alliance.alliance.AllianceAsset) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  |  |






<a name="alliance.alliance.QueryIBCAllianceDelegationRequest"></a>

### QueryIBCAllianceDelegationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_addr` | [string](#string) |  |  |
| `validator_addr` | [string](#string) |  |  |
| `hash` | [string](#string) |  |  |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  |  |






<a name="alliance.alliance.QueryIBCAllianceDelegationRewardsRequest"></a>

### QueryIBCAllianceDelegationRewardsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_addr` | [string](#string) |  |  |
| `validator_addr` | [string](#string) |  |  |
| `hash` | [string](#string) |  |  |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  |  |






<a name="alliance.alliance.QueryIBCAllianceRequest"></a>

### QueryIBCAllianceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `hash` | [string](#string) |  |  |






<a name="alliance.alliance.QueryParamsRequest"></a>

### QueryParamsRequest
Params






<a name="alliance.alliance.QueryParamsResponse"></a>

### QueryParamsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#alliance.alliance.Params) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="alliance.alliance.Query"></a>

### Query


| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Params` | [QueryParamsRequest](#alliance.alliance.QueryParamsRequest) | [QueryParamsResponse](#alliance.alliance.QueryParamsResponse) |  | GET|/terra/alliances/params|
| `Alliances` | [QueryAlliancesRequest](#alliance.alliance.QueryAlliancesRequest) | [QueryAlliancesResponse](#alliance.alliance.QueryAlliancesResponse) | Query paginated alliances | GET|/terra/alliances|
| `IBCAlliance` | [QueryIBCAllianceRequest](#alliance.alliance.QueryIBCAllianceRequest) | [QueryAllianceResponse](#alliance.alliance.QueryAllianceResponse) | Query a specific alliance by ibc hash @deprecated: this endpoint will be replaced for by the encoded version of the denom e.g.: GET:/terra/alliances/ibc%2Falliance | GET|/terra/alliances/ibc/{hash}|
| `AllAlliancesDelegations` | [QueryAllAlliancesDelegationsRequest](#alliance.alliance.QueryAllAlliancesDelegationsRequest) | [QueryAlliancesDelegationsResponse](#alliance.alliance.QueryAlliancesDelegationsResponse) | Query all paginated alliance delegations | GET|/terra/alliances/delegations|
| `AllianceValidator` | [QueryAllianceValidatorRequest](#alliance.alliance.QueryAllianceValidatorRequest) | [QueryAllianceValidatorResponse](#alliance.alliance.QueryAllianceValidatorResponse) | Query alliance validator | GET|/terra/alliances/validators/{validator_addr}|
| `AllAllianceValidators` | [QueryAllAllianceValidatorsRequest](#alliance.alliance.QueryAllAllianceValidatorsRequest) | [QueryAllianceValidatorsResponse](#alliance.alliance.QueryAllianceValidatorsResponse) | Query all paginated alliance validators | GET|/terra/alliances/validators|
| `AlliancesDelegation` | [QueryAlliancesDelegationsRequest](#alliance.alliance.QueryAlliancesDelegationsRequest) | [QueryAlliancesDelegationsResponse](#alliance.alliance.QueryAlliancesDelegationsResponse) | Query all paginated alliance delegations for a delegator addr | GET|/terra/alliances/delegations/{delegator_addr}|
| `AlliancesDelegationByValidator` | [QueryAlliancesDelegationByValidatorRequest](#alliance.alliance.QueryAlliancesDelegationByValidatorRequest) | [QueryAlliancesDelegationsResponse](#alliance.alliance.QueryAlliancesDelegationsResponse) | Query all paginated alliance delegations for a delegator addr and validator_addr | GET|/terra/alliances/delegations/{delegator_addr}/{validator_addr}|
| `AllianceDelegation` | [QueryAllianceDelegationRequest](#alliance.alliance.QueryAllianceDelegationRequest) | [QueryAllianceDelegationResponse](#alliance.alliance.QueryAllianceDelegationResponse) | Query a delegation to an alliance by delegator addr, validator_addr and denom | GET|/terra/alliances/delegations/{delegator_addr}/{validator_addr}/{denom}|
| `IBCAllianceDelegation` | [QueryIBCAllianceDelegationRequest](#alliance.alliance.QueryIBCAllianceDelegationRequest) | [QueryAllianceDelegationResponse](#alliance.alliance.QueryAllianceDelegationResponse) | Query a delegation to an alliance by delegator addr, validator_addr and denom @deprecated: this endpoint will be replaced for by the encoded version of the denom e.g.: GET:/terra/alliances/terradr1231/terravaloper41234/ibc%2Falliance | GET|/terra/alliances/delegations/{delegator_addr}/{validator_addr}/ibc/{hash}|
| `AllianceDelegationRewards` | [QueryAllianceDelegationRewardsRequest](#alliance.alliance.QueryAllianceDelegationRewardsRequest) | [QueryAllianceDelegationRewardsResponse](#alliance.alliance.QueryAllianceDelegationRewardsResponse) | Query for rewards by delegator addr, validator_addr and denom | GET|/terra/alliances/rewards/{delegator_addr}/{validator_addr}/{denom}|
| `IBCAllianceDelegationRewards` | [QueryIBCAllianceDelegationRewardsRequest](#alliance.alliance.QueryIBCAllianceDelegationRewardsRequest) | [QueryAllianceDelegationRewardsResponse](#alliance.alliance.QueryAllianceDelegationRewardsResponse) | Query for rewards by delegator addr, validator_addr and denom @deprecated: this endpoint will be replaced for by the encoded version of the denom e.g.: GET:/terra/alliances/terradr1231/terravaloper41234/ibc%2Falliance | GET|/terra/alliances/rewards/{delegator_addr}/{validator_addr}/ibc/{hash}|
| `AllianceUnbondingsByDenomAndDelegator` | [QueryAllianceUnbondingsByDenomAndDelegatorRequest](#alliance.alliance.QueryAllianceUnbondingsByDenomAndDelegatorRequest) | [QueryAllianceUnbondingsByDenomAndDelegatorResponse](#alliance.alliance.QueryAllianceUnbondingsByDenomAndDelegatorResponse) | Query for rewards by delegator addr, validator_addr and denom | GET|/terra/alliances/unbondings/{denom}/{delegator_addr}|
| `AllianceUnbondings` | [QueryAllianceUnbondingsRequest](#alliance.alliance.QueryAllianceUnbondingsRequest) | [QueryAllianceUnbondingsResponse](#alliance.alliance.QueryAllianceUnbondingsResponse) | Query for rewards by delegator addr, validator_addr and denom | GET|/terra/alliances/unbondings/{denom}/{delegator_addr}/{validator_addr}|
| `AllianceRedelegations` | [QueryAllianceRedelegationsRequest](#alliance.alliance.QueryAllianceRedelegationsRequest) | [QueryAllianceRedelegationsResponse](#alliance.alliance.QueryAllianceRedelegationsResponse) | Query redelegations by denom and delegator address | GET|/terra/alliances/redelegations/{denom}/{delegator_addr}|
| `Alliance` | [QueryAllianceRequest](#alliance.alliance.QueryAllianceRequest) | [QueryAllianceResponse](#alliance.alliance.QueryAllianceResponse) | Query a specific alliance by denom | GET|/terra/alliances/{denom}|

 <!-- end services -->



<a name="alliance/alliance/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## alliance/alliance/tx.proto



<a name="alliance.alliance.MsgClaimDelegationRewards"></a>

### MsgClaimDelegationRewards



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  |  |
| `validator_address` | [string](#string) |  |  |
| `denom` | [string](#string) |  |  |






<a name="alliance.alliance.MsgClaimDelegationRewardsResponse"></a>

### MsgClaimDelegationRewardsResponse







<a name="alliance.alliance.MsgCreateAlliance"></a>

### MsgCreateAlliance



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `authority` | [string](#string) |  |  |
| `denom` | [string](#string) |  | Denom of the asset. It could either be a native token or an IBC token |
| `reward_weight` | [string](#string) |  | The reward weight specifies the ratio of rewards that will be given to each alliance asset It does not need to sum to 1. rate = weight / total_weight Native asset is always assumed to have a weight of 1. |
| `take_rate` | [string](#string) |  | A positive take rate is used for liquid staking derivatives. It defines an annualized reward rate that will be redirected to the distribution rewards pool |
| `reward_change_rate` | [string](#string) |  |  |
| `reward_change_interval` | [google.protobuf.Duration](#google.protobuf.Duration) |  |  |
| `reward_weight_range` | [RewardWeightRange](#alliance.alliance.RewardWeightRange) |  | set a bound of weight range to limit how much reward weights can scale. |






<a name="alliance.alliance.MsgCreateAllianceResponse"></a>

### MsgCreateAllianceResponse







<a name="alliance.alliance.MsgDelegate"></a>

### MsgDelegate



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  |  |
| `validator_address` | [string](#string) |  |  |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |






<a name="alliance.alliance.MsgDelegateResponse"></a>

### MsgDelegateResponse







<a name="alliance.alliance.MsgDeleteAlliance"></a>

### MsgDeleteAlliance



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `authority` | [string](#string) |  |  |
| `denom` | [string](#string) |  |  |






<a name="alliance.alliance.MsgDeleteAllianceResponse"></a>

### MsgDeleteAllianceResponse







<a name="alliance.alliance.MsgRedelegate"></a>

### MsgRedelegate



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  |  |
| `validator_src_address` | [string](#string) |  |  |
| `validator_dst_address` | [string](#string) |  |  |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |






<a name="alliance.alliance.MsgRedelegateResponse"></a>

### MsgRedelegateResponse







<a name="alliance.alliance.MsgUndelegate"></a>

### MsgUndelegate



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  |  |
| `validator_address` | [string](#string) |  |  |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |






<a name="alliance.alliance.MsgUndelegateResponse"></a>

### MsgUndelegateResponse







<a name="alliance.alliance.MsgUpdateAlliance"></a>

### MsgUpdateAlliance



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `authority` | [string](#string) |  |  |
| `denom` | [string](#string) |  | Denom of the asset. It could either be a native token or an IBC token |
| `reward_weight` | [string](#string) |  | The reward weight specifies the ratio of rewards that will be given to each alliance asset It does not need to sum to 1. rate = weight / total_weight Native asset is always assumed to have a weight of 1. |
| `take_rate` | [string](#string) |  |  |
| `reward_change_rate` | [string](#string) |  |  |
| `reward_change_interval` | [google.protobuf.Duration](#google.protobuf.Duration) |  |  |






<a name="alliance.alliance.MsgUpdateAllianceResponse"></a>

### MsgUpdateAllianceResponse







<a name="alliance.alliance.MsgUpdateParams"></a>

### MsgUpdateParams



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  | the title of the update proposal |
| `description` | [string](#string) |  | the description of the proposal |
| `authority` | [string](#string) |  |  |
| `params` | [Params](#alliance.alliance.Params) |  |  |






<a name="alliance.alliance.MsgUpdateParamsResponse"></a>

### MsgUpdateParamsResponse






 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="alliance.alliance.Msg"></a>

### Msg


| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Delegate` | [MsgDelegate](#alliance.alliance.MsgDelegate) | [MsgDelegateResponse](#alliance.alliance.MsgDelegateResponse) |  | |
| `Redelegate` | [MsgRedelegate](#alliance.alliance.MsgRedelegate) | [MsgRedelegateResponse](#alliance.alliance.MsgRedelegateResponse) |  | |
| `Undelegate` | [MsgUndelegate](#alliance.alliance.MsgUndelegate) | [MsgUndelegateResponse](#alliance.alliance.MsgUndelegateResponse) |  | |
| `ClaimDelegationRewards` | [MsgClaimDelegationRewards](#alliance.alliance.MsgClaimDelegationRewards) | [MsgClaimDelegationRewardsResponse](#alliance.alliance.MsgClaimDelegationRewardsResponse) |  | |
| `UpdateParams` | [MsgUpdateParams](#alliance.alliance.MsgUpdateParams) | [MsgUpdateParamsResponse](#alliance.alliance.MsgUpdateParamsResponse) |  | |
| `CreateAlliance` | [MsgCreateAlliance](#alliance.alliance.MsgCreateAlliance) | [MsgCreateAllianceResponse](#alliance.alliance.MsgCreateAllianceResponse) |  | |
| `UpdateAlliance` | [MsgUpdateAlliance](#alliance.alliance.MsgUpdateAlliance) | [MsgUpdateAllianceResponse](#alliance.alliance.MsgUpdateAllianceResponse) |  | |
| `DeleteAlliance` | [MsgDeleteAlliance](#alliance.alliance.MsgDeleteAlliance) | [MsgDeleteAllianceResponse](#alliance.alliance.MsgDeleteAllianceResponse) |  | |

 <!-- end services -->



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |
