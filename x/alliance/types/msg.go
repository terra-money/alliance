package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	_ sdk.Msg = &MsgDelegate{}
	_ sdk.Msg = &MsgRedelegate{}
	_ sdk.Msg = &MsgUndelegate{}
	_ sdk.Msg = &MsgClaimDelegationRewards{}
	_ sdk.Msg = &MsgUpdateParams{}
	_ sdk.Msg = &MsgCreateAlliance{}
	_ sdk.Msg = &MsgUpdateAlliance{}
	_ sdk.Msg = &MsgDeleteAlliance{}

	_ legacytx.LegacyMsg = &MsgDelegate{}
	_ legacytx.LegacyMsg = &MsgRedelegate{}
	_ legacytx.LegacyMsg = &MsgUndelegate{}
	_ legacytx.LegacyMsg = &MsgClaimDelegationRewards{}
	_ legacytx.LegacyMsg = &MsgCreateAlliance{}
	_ legacytx.LegacyMsg = &MsgUpdateAlliance{}
	_ legacytx.LegacyMsg = &MsgDeleteAlliance{}
)

var (
	MsgDelegateType               = "msg_delegate"
	MsgUndelegateType             = "msg_undelegate"
	MsgRedelegateType             = "msg_redelegate"
	MsgClaimDelegationRewardsType = "claim_delegation_rewards"
	MsgCreateAllianceType         = "create_alliance"
	MsgUpdateAllianceType         = "update_alliance"
	MsgDeleteAllianceType         = "delete_alliance"
)

func NewMsgDelegate(delegatorAddress, validatorAddress string, amount sdk.Coin) *MsgDelegate {
	return &MsgDelegate{
		DelegatorAddress: delegatorAddress,
		ValidatorAddress: validatorAddress,
		Amount:           amount,
	}
}

func (msg MsgDelegate) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgDelegate) Route() string {
	return sdk.MsgTypeURL(&msg)
}

func (msg MsgDelegate) ValidateBasic() error {
	if !msg.Amount.Amount.GT(sdk.ZeroInt()) {
		return status.Errorf(codes.InvalidArgument, "Alliance delegation amount must be more than zero")
	}
	return nil
}

func (msg MsgDelegate) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		panic("DelegatorAddress signer from MsgDelegate is not valid")
	}
	return []sdk.AccAddress{signer}
}

func (msg MsgDelegate) Type() string { return MsgDelegateType }

func NewMsgRedelegate(delegatorAddress, validatorSrcAddress, validatorDstAddress string, amount sdk.Coin) *MsgRedelegate {
	return &MsgRedelegate{
		DelegatorAddress:    delegatorAddress,
		ValidatorSrcAddress: validatorSrcAddress,
		ValidatorDstAddress: validatorDstAddress,
		Amount:              amount,
	}
}

func (msg MsgRedelegate) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgRedelegate) Route() string {
	return sdk.MsgTypeURL(&msg)
}

func (msg MsgRedelegate) ValidateBasic() error {
	if msg.Amount.Amount.LTE(sdk.ZeroInt()) {
		return status.Errorf(codes.InvalidArgument, "Alliance redelegation amount must be more than zero")
	}
	return nil
}

func (msg MsgRedelegate) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		panic("DelegatorAddress signer from MsgRedelegate is not valid")
	}
	return []sdk.AccAddress{signer}
}

func (msg MsgRedelegate) Type() string { return MsgRedelegateType }

func NewMsgUndelegate(delegatorAddress, validatorAddress string, amount sdk.Coin) *MsgUndelegate {
	return &MsgUndelegate{
		DelegatorAddress: delegatorAddress,
		ValidatorAddress: validatorAddress,
		Amount:           amount,
	}
}

func (msg MsgUndelegate) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgUndelegate) Route() string {
	return sdk.MsgTypeURL(&msg)
}

func (msg MsgUndelegate) ValidateBasic() error {
	if msg.Amount.Amount.LTE(sdk.ZeroInt()) {
		return status.Errorf(codes.InvalidArgument, "Alliance undelegate amount must be more than zero")
	}
	return nil
}

func (msg MsgUndelegate) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		panic("DelegatorAddress signer from MsgUndelegate is not valid")
	}
	return []sdk.AccAddress{signer}
}

func (msg MsgUndelegate) Type() string { return MsgUndelegateType }

func NewMsgClaimDelegationRewards(delegatorAddress, validatorAddress, denom string) *MsgClaimDelegationRewards {
	return &MsgClaimDelegationRewards{
		DelegatorAddress: delegatorAddress,
		ValidatorAddress: validatorAddress,
		Denom:            denom,
	}
}

func (msg *MsgClaimDelegationRewards) ValidateBasic() error {
	if msg.Denom == "" {
		return status.Errorf(codes.InvalidArgument, "Alliance denom must have a value")
	}
	return nil
}

func (msg MsgClaimDelegationRewards) Route() string {
	return sdk.MsgTypeURL(&msg)
}

func (msg MsgClaimDelegationRewards) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg *MsgClaimDelegationRewards) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		panic("DelegatorAddress signer from MsgClaimDelegationRewards is not valid")
	}
	return []sdk.AccAddress{signer}
}

func (msg MsgClaimDelegationRewards) Type() string { return MsgClaimDelegationRewardsType }

func (msg *MsgCreateAlliance) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return sdkerrors.Wrap(err, "invalid authority address")
	}

	if msg.Denom == "" {
		return status.Errorf(codes.InvalidArgument, "Alliance denom must have a value")
	}

	if err := sdk.ValidateDenom(msg.Denom); err != nil {
		return err
	}

	if msg.RewardWeight.IsNil() || msg.RewardWeight.LT(sdk.ZeroDec()) {
		return status.Errorf(codes.InvalidArgument, "Alliance rewardWeight must be zero or a positive number")
	}

	if msg.RewardWeightRange.Min.IsNil() || msg.RewardWeightRange.Min.LT(sdk.ZeroDec()) ||
		msg.RewardWeightRange.Max.IsNil() || msg.RewardWeightRange.Max.LT(sdk.ZeroDec()) {
		return status.Errorf(codes.InvalidArgument, "Alliance rewardWeight min and max must be zero or a positive number")
	}

	if msg.RewardWeightRange.Min.GT(msg.RewardWeightRange.Max) {
		return status.Errorf(codes.InvalidArgument, "Alliance rewardWeight min must be less or equal to rewardWeight max")
	}

	if msg.RewardWeight.LT(msg.RewardWeightRange.Min) || msg.RewardWeight.GT(msg.RewardWeightRange.Max) {
		return status.Errorf(codes.InvalidArgument, "Alliance rewardWeight must be bounded in RewardWeightRange")
	}

	if msg.TakeRate.IsNil() || msg.TakeRate.IsNegative() || msg.TakeRate.GTE(sdk.OneDec()) {
		return status.Errorf(codes.InvalidArgument, "Alliance takeRate must be more or equals to 0 but strictly less than 1")
	}

	if msg.RewardChangeRate.IsZero() || msg.RewardChangeRate.IsNegative() {
		return status.Errorf(codes.InvalidArgument, "Alliance rewardChangeRate must be strictly a positive number")
	}

	if msg.RewardChangeInterval < 0 {
		return status.Errorf(codes.InvalidArgument, "Alliance rewardChangeInterval must be strictly a positive number")
	}

	return nil
}

func (msg MsgCreateAlliance) Route() string {
	return sdk.MsgTypeURL(&msg)
}

func (msg MsgCreateAlliance) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg *MsgCreateAlliance) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic("Authority is not valid")
	}
	return []sdk.AccAddress{signer}
}

func (msg MsgCreateAlliance) Type() string { return MsgCreateAllianceType }

func (msg *MsgUpdateAlliance) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return sdkerrors.Wrap(err, "invalid authority address")
	}

	if msg.Denom == "" {
		return status.Errorf(codes.InvalidArgument, "Alliance denom must have a value")
	}

	if msg.RewardWeight.IsNil() || msg.RewardWeight.LT(sdk.ZeroDec()) {
		return status.Errorf(codes.InvalidArgument, "Alliance rewardWeight must be zero or a positive number")
	}

	if msg.TakeRate.IsNil() || msg.TakeRate.IsNegative() || msg.TakeRate.GTE(sdk.OneDec()) {
		return status.Errorf(codes.InvalidArgument, "Alliance takeRate must be more or equals to 0 but strictly less than 1")
	}

	if msg.RewardChangeRate.IsZero() || msg.RewardChangeRate.IsNegative() {
		return status.Errorf(codes.InvalidArgument, "Alliance rewardChangeRate must be strictly a positive number")
	}

	if msg.RewardChangeInterval < 0 {
		return status.Errorf(codes.InvalidArgument, "Alliance rewardChangeInterval must be strictly a positive number")
	}

	return nil
}

func (msg MsgUpdateAlliance) Route() string {
	return sdk.MsgTypeURL(&msg)
}

func (msg MsgUpdateAlliance) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg *MsgUpdateAlliance) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic("Authority is not valid")
	}
	return []sdk.AccAddress{signer}
}

func (msg MsgUpdateAlliance) Type() string { return MsgUpdateAllianceType }

func (msg *MsgDeleteAlliance) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return sdkerrors.Wrap(err, "invalid authority address")
	}

	if msg.Denom == "" {
		return status.Errorf(codes.InvalidArgument, "Alliance denom must have a value")
	}
	return nil
}

func (msg MsgDeleteAlliance) Route() string {
	return sdk.MsgTypeURL(&msg)
}

func (msg MsgDeleteAlliance) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg *MsgDeleteAlliance) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic("Authority is not valid")
	}
	return []sdk.AccAddress{signer}
}

func (msg MsgDeleteAlliance) Type() string { return MsgDeleteAllianceType }
