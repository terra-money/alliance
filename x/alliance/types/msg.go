package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	_ sdk.Msg = &MsgCreateAlliance{}
	_ sdk.Msg = &MsgUpdateAlliance{}
	_ sdk.Msg = &MsgDeleteAlliance{}

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

// Execution allowed only from Governance Module
func (m MsgCreateAlliance) ValidateBasic() error {
	if m.Alliance.Denom == "" {
		return status.Errorf(codes.InvalidArgument, "Alliance denom must have a value")
	}

	if m.Alliance.RewardWeight.IsNil() || m.Alliance.RewardWeight.LTE(sdk.ZeroDec()) {
		return status.Errorf(codes.InvalidArgument, "Alliance rewardWeight must be a positive number")
	}

	if m.Alliance.TakeRate.IsNil() || m.Alliance.TakeRate.LTE(sdk.ZeroDec()) {
		return status.Errorf(codes.InvalidArgument, "Alliance takeRate must be a positive number")
	}

	return nil
}

func (m MsgCreateAlliance) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(m.Authority)
	if err != nil {
		panic("Authority signer from MsgCreateAlliance is not valid")
	}
	return []sdk.AccAddress{signer}
}

func (m MsgUpdateAlliance) ValidateBasic() error {
	if m.Denom == "" {
		return status.Errorf(codes.InvalidArgument, "Alliance denom must have a value")
	}

	if m.RewardWeight.IsNil() || m.RewardWeight.LTE(sdk.ZeroDec()) {
		return status.Errorf(codes.InvalidArgument, "Alliance rewardWeight must be a positive number")
	}

	if m.TakeRate.IsNil() || m.TakeRate.LTE(sdk.ZeroDec()) {
		return status.Errorf(codes.InvalidArgument, "Alliance takeRate must be a positive number")
	}

	return nil
}

func (m MsgUpdateAlliance) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(m.Authority)
	if err != nil {
		panic("Authority signer from MsgUpdateAlliance is not valid")
	}
	return []sdk.AccAddress{signer}
}

func (m MsgDeleteAlliance) ValidateBasic() error {
	if m.Denom == "" {
		return status.Errorf(codes.InvalidArgument, "Alliance denom must have a value")
	}

	return nil
}

func (m MsgDeleteAlliance) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(m.Authority)
	if err != nil {
		panic("Authority signer from MsgDeleteAlliance is not valid")
	}
	return []sdk.AccAddress{signer}
}

// Execution allowed from any account
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

func (msg MsgDelegate) Type() string { return MsgDelegateType }

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

func (msg MsgRedelegate) Type() string { return MsgRedelegateType }

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

func (msg MsgUndelegate) Type() string { return MsgUndelegateType }

func (m *MsgClaimDelegationRewards) ValidateBasic() error {
	if m.Denom == "" {
		return status.Errorf(codes.InvalidArgument, "Alliance denom must have a value")
	}
	return nil
}

func (m *MsgClaimDelegationRewards) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(m.DelegatorAddress)
	if err != nil {
		panic("DelegatorAddress signer from MsgClaimDelegationRewards is not valid")
	}
	return []sdk.AccAddress{signer}
}

func (msg MsgClaimDelegationRewards) Type() string { return MsgClaimDelegationRewardsType }
