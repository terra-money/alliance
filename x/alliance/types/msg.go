package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	_ sdk.Msg = &MsgDelegate{}
	_ sdk.Msg = &MsgRedelegate{}
	_ sdk.Msg = &MsgUndelegate{}
	_ sdk.Msg = &MsgClaimDelegationRewards{}
)

var (
	MsgDelegateType               = "msg_delegate"
	MsgUndelegateType             = "msg_undelegate"
	MsgRedelegateType             = "msg_redelegate"
	MsgClaimDelegationRewardsType = "claim_delegation_rewards"
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

func (m MsgDelegate) ValidateBasic() error {
	if !m.Amount.Amount.GT(sdk.ZeroInt()) {
		return status.Errorf(codes.InvalidArgument, "Alliance delegation amount must be more than zero")
	}
	return nil
}

func (m MsgDelegate) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(m.DelegatorAddress)
	if err != nil {
		panic("DelegatorAddress signer from MsgDelegate is not valid")
	}
	return []sdk.AccAddress{signer}
}

func (msg MsgDelegate) Type() string { return MsgDelegateType } //nolint:revive // TODO: we should figure out how to differentiate this from MsgDelegate above.

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

func (m MsgRedelegate) ValidateBasic() error {
	if m.Amount.Amount.LTE(sdk.ZeroInt()) {
		return status.Errorf(codes.InvalidArgument, "Alliance redelegation amount must be more than zero")
	}
	return nil
}

func (m MsgRedelegate) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(m.DelegatorAddress)
	if err != nil {
		panic("DelegatorAddress signer from MsgRedelegate is not valid")
	}
	return []sdk.AccAddress{signer}
}

func (msg MsgRedelegate) Type() string { return MsgRedelegateType } //nolint:revive // should make receivers consistent

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

func (m MsgUndelegate) ValidateBasic() error {
	if m.Amount.Amount.LTE(sdk.ZeroInt()) {
		return status.Errorf(codes.InvalidArgument, "Alliance undelegate amount must be more than zero")
	}
	return nil
}

func (m MsgUndelegate) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(m.DelegatorAddress)
	if err != nil {
		panic("DelegatorAddress signer from MsgUndelegate is not valid")
	}
	return []sdk.AccAddress{signer}
}

func (msg MsgUndelegate) Type() string { return MsgUndelegateType } //nolint:revive // should make receivers consistent

func (m *MsgClaimDelegationRewards) ValidateBasic() error {
	if m.Denom == "" {
		return status.Errorf(codes.InvalidArgument, "Alliance denom must have a value")
	}
	return nil
}

func (msg MsgClaimDelegationRewards) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (m *MsgClaimDelegationRewards) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(m.DelegatorAddress)
	if err != nil {
		panic("DelegatorAddress signer from MsgClaimDelegationRewards is not valid")
	}
	return []sdk.AccAddress{signer}
}

func (msg MsgClaimDelegationRewards) Type() string { return MsgClaimDelegationRewardsType } //nolint:revive // should make receivers consistent
