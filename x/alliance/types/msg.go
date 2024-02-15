package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
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
)

func NewMsgDelegate(delegatorAddress, validatorAddress string, amount sdk.Coin) *MsgDelegate {
	return &MsgDelegate{
		DelegatorAddress: delegatorAddress,
		ValidatorAddress: validatorAddress,
		Amount:           amount,
	}
}
