package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgDelegate{}, "alliance/MsgDelegate", nil)
	cdc.RegisterConcrete(&MsgRedelegate{}, "alliance/MsgRedelegate", nil)
	cdc.RegisterConcrete(&MsgUndelegate{}, "alliance/MsgUndelegate", nil)
	cdc.RegisterConcrete(&MsgClaimDelegationRewards{}, "alliance/MsgClaimDelegationRewards", nil)

	cdc.RegisterConcrete(&MsgCreateAllianceProposal{}, "alliance/MsgCreateAllianceProposal", nil)
	cdc.RegisterConcrete(&MsgUpdateAllianceProposal{}, "alliance/MsgUpdateAllianceProposal", nil)
	cdc.RegisterConcrete(&MsgDeleteAllianceProposal{}, "alliance/MsgDeleteAllianceProposal", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgDelegate{},
		&MsgRedelegate{},
		&MsgUndelegate{},
		&MsgClaimDelegationRewards{},
	)

	registry.RegisterImplementations((*govtypes.Content)(nil),
		&MsgCreateAllianceProposal{},
		&MsgUpdateAllianceProposal{},
		&MsgDeleteAllianceProposal{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	sdk.RegisterLegacyAminoCodec(amino)

	// Register all Amino interfaces and concrete types on the authz Amino codec so that this can later be
	// used to properly serialize MsgGrant and MsgExec instances
	//RegisterLegacyAminoCodec(authzcodec.Amino)
}
