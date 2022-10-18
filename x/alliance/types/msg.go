package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ sdk.Msg = &MsgDelegate{}
	_ sdk.Msg = &MsgRedelegate{}
	_ sdk.Msg = &MsgUndelegate{}
)

func (m MsgDelegate) ValidateBasic() error {
	if !m.Amount.Amount.GT(sdk.NewInt(0)) {
		return fmt.Errorf("Delegation amount must be more than zero")
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
	//TODO implement me
	panic("implement me")
}

func (m MsgRedelegate) GetSigners() []sdk.AccAddress {
	//TODO implement me
	panic("implement me")
}
func (m MsgUndelegate) ValidateBasic() error {
	//TODO implement me
	panic("implement me")
}

func (m MsgUndelegate) GetSigners() []sdk.AccAddress {
	//TODO implement me
	panic("implement me")
}
