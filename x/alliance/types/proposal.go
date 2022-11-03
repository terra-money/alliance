package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	ProposalTypeCreateAlliance = "CreateAlliance"
	ProposalTypeUpdateAlliance = "UpdateAlliance"
	ProposalTypeDeleteAlliance = "DeleteAlliance"
)

var (
	_ govtypes.Content = &CreateAllianceProposal{}
	_ govtypes.Content = &UpdateAllianceProposal{}
	_ govtypes.Content = &DeleteAllianceProposal{}
)

func init() {
	govtypes.RegisterProposalType(ProposalTypeCreateAlliance)
	govtypes.RegisterProposalType(ProposalTypeUpdateAlliance)
	govtypes.RegisterProposalType(ProposalTypeDeleteAlliance)
}
func NewCreateAllianceProposal(title, description, denom string, rewardWeight, takeRate sdk.Dec) govtypes.Content {
	return &CreateAllianceProposal{
		Title:        title,
		Description:  description,
		Denom:        denom,
		RewardWeight: rewardWeight,
		TakeRate:     takeRate,
	}
}
func (m *CreateAllianceProposal) GetTitle() string       { return m.Title }
func (m *CreateAllianceProposal) GetDescription() string { return m.Description }
func (m *CreateAllianceProposal) ProposalRoute() string  { return RouterKey }
func (m *CreateAllianceProposal) ProposalType() string   { return ProposalTypeCreateAlliance }
func (m *CreateAllianceProposal) ValidateBasic() error {

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

func NewUpdateAllianceProposal(title, description, denom string, rewardWeight, takeRate sdk.Dec) govtypes.Content {
	return &UpdateAllianceProposal{
		Title:        title,
		Description:  description,
		Denom:        denom,
		RewardWeight: rewardWeight,
		TakeRate:     takeRate,
	}
}
func (m *UpdateAllianceProposal) GetTitle() string       { return m.Title }
func (m *UpdateAllianceProposal) GetDescription() string { return m.Description }
func (m *UpdateAllianceProposal) ProposalRoute() string  { return RouterKey }
func (m *UpdateAllianceProposal) ProposalType() string   { return ProposalTypeUpdateAlliance }
func (m *UpdateAllianceProposal) ValidateBasic() error {
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

func NewDeleteAllianceProposal(title, description, denom string) govtypes.Content {
	return &DeleteAllianceProposal{
		Title:       title,
		Description: description,
		Denom:       denom,
	}
}
func (m *DeleteAllianceProposal) GetTitle() string       { return m.Title }
func (m *DeleteAllianceProposal) GetDescription() string { return m.Description }
func (m *DeleteAllianceProposal) ProposalRoute() string  { return RouterKey }
func (m *DeleteAllianceProposal) ProposalType() string   { return ProposalTypeDeleteAlliance }
func (m *DeleteAllianceProposal) ValidateBasic() error {
	if m.Denom == "" {
		return status.Errorf(codes.InvalidArgument, "Alliance denom must have a value")
	}
	return nil
}
