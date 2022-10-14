package types

import sdk "github.com/cosmos/cosmos-sdk/types"

var (
	_ sdk.Msg = &MsgDelegate{}
	_ sdk.Msg = &MsgRedelegate{}
	_ sdk.Msg = &MsgUndelegate{}
)

func (m MsgDelegate) ValidateBasic() error {
	//TODO implement me
	panic("implement me")
}

func (m MsgDelegate) GetSigners() []sdk.AccAddress {
	//TODO implement me
	panic("implement me")
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
