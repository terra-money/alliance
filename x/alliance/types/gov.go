package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

const (
	ProposalTypeCreateAlliance = "msg_create_alliance_proposal"
	ProposalTypeUpdateAlliance = "msg_update_alliance_proposal"
	ProposalTypeDeleteAlliance = "msg_delete_alliance_proposal"
)

var (
	_ govtypes.Content = &MsgCreateAllianceProposal{}
	_ govtypes.Content = &MsgUpdateAllianceProposal{}
	_ govtypes.Content = &MsgDeleteAllianceProposal{}
)

func init() {
	govtypes.RegisterProposalType(ProposalTypeCreateAlliance)
	govtypes.RegisterProposalType(ProposalTypeUpdateAlliance)
	govtypes.RegisterProposalType(ProposalTypeDeleteAlliance)
}
func NewMsgCreateAllianceProposal(title, description, denom string, rewardWeight, takeRate sdk.Dec) govtypes.Content {
	return &MsgCreateAllianceProposal{
		Title:        title,
		Description:  description,
		Denom:        denom,
		RewardWeight: rewardWeight,
		TakeRate:     takeRate,
	}
}
func (m *MsgCreateAllianceProposal) GetTitle() string       { return m.Title }
func (m *MsgCreateAllianceProposal) GetDescription() string { return m.Description }
func (m *MsgCreateAllianceProposal) ProposalRoute() string  { return RouterKey }
func (m *MsgCreateAllianceProposal) ProposalType() string   { return ProposalTypeCreateAlliance }

func (m *MsgCreateAllianceProposal) ValidateBasic() error {

	if m.Denom == "" {
		return status.Errorf(codes.InvalidArgument, "Alliance denom must have a value")
	}

	if m.RewardWeight.IsNil() || m.RewardWeight.LTE(sdk.ZeroDec()) {
		return status.Errorf(codes.InvalidArgument, "Alliance rewardWeight must be a positive number")
	}

	if m.TakeRate.IsNil() || m.TakeRate.LTE(sdk.ZeroDec()) {
		return status.Errorf(codes.InvalidArgument, "Alliance takeRate must be zero or a positive number")
	}

	if m.RewardChangeRate.IsZero() || m.RewardChangeRate.IsNegative() {
		return status.Errorf(codes.InvalidArgument, "Alliance rewardChangeRate must be strictly a positive number")
	}

	return nil
}

func NewMsgUpdateAllianceProposal(title, description, denom string, rewardWeight, takeRate sdk.Dec, rewardChangeRate sdk.Dec, rewardChangeInterval time.Duration) govtypes.Content {
	return &MsgUpdateAllianceProposal{
		Title:                title,
		Description:          description,
		Denom:                denom,
		RewardWeight:         rewardWeight,
		TakeRate:             takeRate,
		RewardChangeRate:     rewardChangeRate,
		RewardChangeInterval: rewardChangeInterval,
	}
}
func (m *MsgUpdateAllianceProposal) GetTitle() string       { return m.Title }
func (m *MsgUpdateAllianceProposal) GetDescription() string { return m.Description }
func (m *MsgUpdateAllianceProposal) ProposalRoute() string  { return RouterKey }
func (m *MsgUpdateAllianceProposal) ProposalType() string   { return ProposalTypeUpdateAlliance }

func (m *MsgUpdateAllianceProposal) ValidateBasic() error {
	if m.Denom == "" {
		return status.Errorf(codes.InvalidArgument, "Alliance denom must have a value")
	}

	if m.RewardWeight.IsNil() || m.RewardWeight.LTE(sdk.ZeroDec()) {
		return status.Errorf(codes.InvalidArgument, "Alliance rewardWeight must be a positive number")
	}

	if m.TakeRate.IsNil() || m.TakeRate.LTE(sdk.ZeroDec()) {
		return status.Errorf(codes.InvalidArgument, "Alliance takeRate must be a positive number")
	}

	if m.RewardChangeRate.IsZero() || m.RewardChangeRate.IsNegative() {
		return status.Errorf(codes.InvalidArgument, "Alliance rewardChangeRate must be strictly a positive number")
	}

	return nil
}

func NewMsgDeleteAllianceProposal(title, description, denom string) govtypes.Content {
	return &MsgDeleteAllianceProposal{
		Title:       title,
		Description: description,
		Denom:       denom,
	}
}
func (m *MsgDeleteAllianceProposal) GetTitle() string       { return m.Title }
func (m *MsgDeleteAllianceProposal) GetDescription() string { return m.Description }
func (m *MsgDeleteAllianceProposal) ProposalRoute() string  { return RouterKey }
func (m *MsgDeleteAllianceProposal) ProposalType() string   { return ProposalTypeDeleteAlliance }

func (m *MsgDeleteAllianceProposal) ValidateBasic() error {
	if m.Denom == "" {
		return status.Errorf(codes.InvalidArgument, "Alliance denom must have a value")
	}
	return nil
}
