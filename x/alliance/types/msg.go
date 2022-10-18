package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ sdk.Msg = &MsgDelegate{}
	_ sdk.Msg = &MsgRedelegate{}
	_ sdk.Msg = &MsgUndelegate{}
	_ sdk.Msg = &MsgCreateAlliance{}
	_ sdk.Msg = &MsgUpdateAlliance{}
	_ sdk.Msg = &MsgDeleteAlliance{}
	_ sdk.Msg = &MsgClaimDelegationRewards{}
)

func (m MsgDelegate) ValidateBasic() error {
	if !m.Amount.Amount.GT(sdk.ZeroInt()) {
		return fmt.Errorf("delegation amount must be more than zero")
	}
	return nil
}

func (m MsgDelegate) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(m.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}
func (m MsgRedelegate) ValidateBasic() error {
	if !m.Amount.Amount.GTE(sdk.ZeroInt()) {
		return fmt.Errorf("redelegation amount must be more than zero")
	}
	return nil
}

func (m MsgRedelegate) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(m.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}

func (m MsgUndelegate) ValidateBasic() error {
	if !m.Amount.Amount.GTE(sdk.ZeroInt()) {
		return fmt.Errorf("redelegation amount must be more than zero")
	}
	return nil
}

func (m MsgUndelegate) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(m.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}

func (m MsgCreateAlliance) ValidateBasic() error {
	if m.Alliance.Denom != "" {
		return fmt.Errorf("denom must not be empty")
	}

	if m.Alliance.RewardWeight.GTE(sdk.ZeroDec()) {
		return fmt.Errorf("rewardWeight must be positive")
	}

	if m.Alliance.TakeRate.GTE(sdk.ZeroDec()) {
		return fmt.Errorf("takeRate must be positive")
	}

	return nil
}

func (m MsgCreateAlliance) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(m.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}

func (m MsgUpdateAlliance) ValidateBasic() error {
	if m.Denom != "" {
		return fmt.Errorf("denom must not be empty")
	}

	if m.RewardWeight.GTE(sdk.ZeroDec()) {
		return fmt.Errorf("rewardWeight must be positive")
	}

	if m.TakeRate.GTE(sdk.ZeroDec()) {
		return fmt.Errorf("takeRate must be positive")
	}

	return nil
}

func (m MsgUpdateAlliance) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(m.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}

func (m MsgDeleteAlliance) ValidateBasic() error {
	if m.Denom != "" {
		return fmt.Errorf("denom must not be empty")
	}

	return nil
}

func (m MsgDeleteAlliance) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(m.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}

func (m *MsgClaimDelegationRewards) ValidateBasic() error {
	if m.Denom != "" {
		return fmt.Errorf("denom must not be empty")
	}
	return nil
}

func (m *MsgClaimDelegationRewards) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(m.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}
