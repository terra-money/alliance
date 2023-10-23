package types

import (
	"time"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	ProposalTypeUpdateAllianceParams = "msg_update_alliance_params"
	ProposalTypeCreateAlliance       = "msg_create_alliance_proposal"
	ProposalTypeUpdateAlliance       = "msg_update_alliance_proposal"
	ProposalTypeDeleteAlliance       = "msg_delete_alliance_proposal"
)

var (
	_ govtypes.Content = &MsgUpdateParams{}
	_ govtypes.Content = &MsgCreateAllianceProposal{}
	_ govtypes.Content = &MsgUpdateAllianceProposal{}
	_ govtypes.Content = &MsgDeleteAllianceProposal{}
)

func init() {
	govtypes.RegisterProposalType(ProposalTypeUpdateAllianceParams)
	govtypes.RegisterProposalType(ProposalTypeCreateAlliance)
	govtypes.RegisterProposalType(ProposalTypeUpdateAlliance)
	govtypes.RegisterProposalType(ProposalTypeDeleteAlliance)
}

func NewMsgUpdateParams(title, description string,
	rewardDelayTime, takeRateClaimInterval time.Duration,
	lastTakeRateClaimTime time.Time) govtypes.Content {
	return &MsgUpdateParams{
		Title:       title,
		Description: description,
		Params: Params{
			RewardDelayTime:       rewardDelayTime,
			TakeRateClaimInterval: takeRateClaimInterval,
			LastTakeRateClaimTime: lastTakeRateClaimTime,
		},
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
func (m *MsgUpdateParams) ProposalRoute() string { return RouterKey }
func (m *MsgUpdateParams) ProposalType() string  { return ProposalTypeCreateAlliance }

func NewMsgCreateAllianceProposal(title, description, denom string, rewardWeight sdk.Dec, rewardWeightRange RewardWeightRange, takeRate sdk.Dec, rewardChangeRate sdk.Dec, rewardChangeInterval time.Duration) govtypes.Content {
	return &MsgCreateAllianceProposal{
		Title:                title,
		Description:          description,
		Denom:                denom,
		RewardWeight:         rewardWeight,
		RewardWeightRange:    rewardWeightRange,
		TakeRate:             takeRate,
		RewardChangeRate:     rewardChangeRate,
		RewardChangeInterval: rewardChangeInterval,
	}
}
func (m *MsgCreateAllianceProposal) ProposalRoute() string { return RouterKey }
func (m *MsgCreateAllianceProposal) ProposalType() string  { return ProposalTypeCreateAlliance }
func (m *MsgCreateAllianceProposal) ValidateBasic() error {
	if m.Denom == "" {
		return status.Errorf(codes.InvalidArgument, "Alliance denom must have a value")
	}

	if err := sdk.ValidateDenom(m.Denom); err != nil {
		return err
	}

	if m.RewardWeight.IsNil() || m.RewardWeight.LT(sdk.ZeroDec()) {
		return status.Errorf(codes.InvalidArgument, "Alliance rewardWeight must be zero or a positive number")
	}

	if m.RewardWeightRange.Min.IsNil() || m.RewardWeightRange.Min.LT(sdk.ZeroDec()) ||
		m.RewardWeightRange.Max.IsNil() || m.RewardWeightRange.Max.LT(sdk.ZeroDec()) {
		return status.Errorf(codes.InvalidArgument, "Alliance rewardWeight min and max must be zero or a positive number")
	}

	if m.RewardWeightRange.Min.GT(m.RewardWeightRange.Max) {
		return status.Errorf(codes.InvalidArgument, "Alliance rewardWeight min must be less or equal to rewardWeight max")
	}

	if m.RewardWeight.LT(m.RewardWeightRange.Min) || m.RewardWeight.GT(m.RewardWeightRange.Max) {
		return status.Errorf(codes.InvalidArgument, "Alliance rewardWeight must be bounded in RewardWeightRange")
	}

	if m.TakeRate.IsNil() || m.TakeRate.IsNegative() || m.TakeRate.GTE(sdk.OneDec()) {
		return status.Errorf(codes.InvalidArgument, "Alliance takeRate must be more or equals to 0 but strictly less than 1")
	}

	if m.RewardChangeRate.IsZero() || m.RewardChangeRate.IsNegative() {
		return status.Errorf(codes.InvalidArgument, "Alliance rewardChangeRate must be strictly a positive number")
	}

	if m.RewardChangeInterval < 0 {
		return status.Errorf(codes.InvalidArgument, "Alliance rewardChangeInterval must be strictly a positive number")
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
func (m *MsgUpdateAllianceProposal) ProposalRoute() string { return RouterKey }
func (m *MsgUpdateAllianceProposal) ProposalType() string  { return ProposalTypeUpdateAlliance }
func (m *MsgUpdateAllianceProposal) ValidateBasic() error {
	if m.Denom == "" {
		return status.Errorf(codes.InvalidArgument, "Alliance denom must have a value")
	}

	if m.RewardWeight.IsNil() || m.RewardWeight.LT(sdk.ZeroDec()) {
		return status.Errorf(codes.InvalidArgument, "Alliance rewardWeight must be zero or a positive number")
	}

	if m.TakeRate.IsNil() || m.TakeRate.IsNegative() || m.TakeRate.GTE(sdk.OneDec()) {
		return status.Errorf(codes.InvalidArgument, "Alliance takeRate must be more or equals to 0 but strictly less than 1")
	}

	if m.RewardChangeRate.IsZero() || m.RewardChangeRate.IsNegative() {
		return status.Errorf(codes.InvalidArgument, "Alliance rewardChangeRate must be strictly a positive number")
	}

	if m.RewardChangeInterval < 0 {
		return status.Errorf(codes.InvalidArgument, "Alliance rewardChangeInterval must be strictly a positive number")
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
func (m *MsgDeleteAllianceProposal) ProposalRoute() string { return RouterKey }
func (m *MsgDeleteAllianceProposal) ProposalType() string  { return ProposalTypeDeleteAlliance }
func (m *MsgDeleteAllianceProposal) ValidateBasic() error {
	if m.Denom == "" {
		return status.Errorf(codes.InvalidArgument, "Alliance denom must have a value")
	}
	return nil
}
