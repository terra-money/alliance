package types

import (
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

	_ legacytx.LegacyMsg = &MsgDelegate{}
	_ legacytx.LegacyMsg = &MsgRedelegate{}
	_ legacytx.LegacyMsg = &MsgUndelegate{}
	_ legacytx.LegacyMsg = &MsgClaimDelegationRewards{}
	_ legacytx.LegacyMsg = &MsgUpdateParams{}
)

var (
	MsgDelegateType               = "msg_delegate"
	MsgUndelegateType             = "msg_undelegate"
	MsgRedelegateType             = "msg_redelegate"
	MsgClaimDelegationRewardsType = "claim_delegation_rewards"
	MsgUpdateParamsType           = "update_params"
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

func NewMsgUpdateParams(authority string, params Params) *MsgUpdateParams {
	return &MsgUpdateParams{
		authority,
		params,
	}
}


func (msg *MsgUpdateParams) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return sdkerrors.Wrap(err, "invalid authority address")
	}
	if err := ValidatePositiveDuration(msg.Params.RewardDelayTime); err != nil {
		return err
	}
	return ValidatePositiveDuration(msg.Params.TakeRateClaimInterval)
}

func (msg MsgUpdateParams) Route() string {
	return sdk.MsgTypeURL(&msg)
}

func (msg MsgUpdateParams) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg *MsgUpdateParams) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic("Authority is not valid")
	}
	return []sdk.AccAddress{signer}
}

func (msg MsgUpdateParams) Type() string { return MsgUpdateParamsType }
