package keeper

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdkerrors "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/terra-money/alliance/x/alliance/types"
)

type MsgServer struct {
	Keeper
}

var _ types.MsgServer = MsgServer{}

func (m MsgServer) Delegate(ctx context.Context, msg *types.MsgDelegate) (*types.MsgDelegateResponse, error) {
	if !msg.Amount.Amount.GT(math.ZeroInt()) {
		return nil, status.Errorf(codes.InvalidArgument, "Alliance delegation amount must be more than zero")
	}

	delAddr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil, err
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	valAddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	validator, err := m.Keeper.GetAllianceValidator(sdkCtx, valAddr)
	if err != nil {
		return nil, err
	}

	_, err = m.Keeper.Delegate(sdkCtx, delAddr, validator, msg.Amount)
	if err != nil {
		return nil, err
	}

	return &types.MsgDelegateResponse{}, nil
}

func (m MsgServer) Redelegate(ctx context.Context, msg *types.MsgRedelegate) (*types.MsgRedelegateResponse, error) {
	if msg.Amount.Amount.LTE(math.ZeroInt()) {
		return nil, status.Errorf(codes.InvalidArgument, "Alliance redelegation amount must be more than zero")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	delAddr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	srcValAddr, err := sdk.ValAddressFromBech32(msg.ValidatorSrcAddress)
	if err != nil {
		return nil, err
	}
	dstValAddr, err := sdk.ValAddressFromBech32(msg.ValidatorDstAddress)
	if err != nil {
		return nil, err
	}

	srcValidator, err := m.Keeper.GetAllianceValidator(sdkCtx, srcValAddr)
	if err != nil {
		return nil, err
	}

	dstValidator, err := m.Keeper.GetAllianceValidator(sdkCtx, dstValAddr)
	if err != nil {
		return nil, err
	}

	_, err = m.Keeper.Redelegate(sdkCtx, delAddr, srcValidator, dstValidator, msg.Amount)
	if err != nil {
		return nil, err
	}

	return &types.MsgRedelegateResponse{}, nil
}

func (m MsgServer) Undelegate(ctx context.Context, msg *types.MsgUndelegate) (*types.MsgUndelegateResponse, error) {
	if msg.Amount.Amount.LTE(math.ZeroInt()) {
		return nil, status.Errorf(codes.InvalidArgument, "Alliance undelegate amount must be more than zero")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	delAddr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil, err
	}
	valAddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	validator, err := m.Keeper.GetAllianceValidator(sdkCtx, valAddr)
	if err != nil {
		return nil, err
	}

	_, err = m.Keeper.Undelegate(sdkCtx, delAddr, validator, msg.Amount)
	if err != nil {
		return nil, err
	}

	return &types.MsgUndelegateResponse{}, nil
}

func (m MsgServer) ClaimDelegationRewards(ctx context.Context, msg *types.MsgClaimDelegationRewards) (*types.MsgClaimDelegationRewardsResponse, error) {
	if msg.Denom == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Alliance denom must have a value")
	}

	delAddr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil, err
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	valAddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	validator, err := m.Keeper.GetAllianceValidator(sdkCtx, valAddr)
	if err != nil {
		return nil, err
	}

	_, err = m.Keeper.ClaimDelegationRewards(sdkCtx, delAddr, validator, msg.Denom)

	return &types.MsgClaimDelegationRewardsResponse{}, err
}

func (m MsgServer) UpdateParams(ctx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return nil, sdkerrors.Wrap(err, "invalid authority address")
	}
	if err := types.ValidatePositiveDuration(msg.Params.RewardDelayTime); err != nil {
		return nil, err
	}

	if m.GetAuthority() != msg.Authority {
		return nil, sdkerrors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", m.GetAuthority(), msg.Authority)
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if err := m.SetParams(sdkCtx, msg.Params); err != nil {
		return nil, err
	}
	return &types.MsgUpdateParamsResponse{}, nil
}

func (m MsgServer) CreateAlliance(ctx context.Context, msg *types.MsgCreateAlliance) (*types.MsgCreateAllianceResponse, error) {

	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return nil, sdkerrors.Wrap(err, "invalid authority address")
	}

	if msg.Denom == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Alliance denom must have a value")
	}

	if err := sdk.ValidateDenom(msg.Denom); err != nil {
		return nil, err
	}

	if msg.RewardWeight.IsNil() || msg.RewardWeight.LT(math.LegacyZeroDec()) {
		return nil, status.Errorf(codes.InvalidArgument, "Alliance rewardWeight must be zero or a positive number")
	}

	if msg.RewardWeightRange.Min.IsNil() || msg.RewardWeightRange.Min.LT(math.LegacyZeroDec()) ||
		msg.RewardWeightRange.Max.IsNil() || msg.RewardWeightRange.Max.LT(math.LegacyZeroDec()) {
		return nil, status.Errorf(codes.InvalidArgument, "Alliance rewardWeight min and max must be zero or a positive number")
	}

	if msg.RewardWeightRange.Min.GT(msg.RewardWeightRange.Max) {
		return nil, status.Errorf(codes.InvalidArgument, "Alliance rewardWeight min must be less or equal to rewardWeight max")
	}

	if msg.RewardWeight.LT(msg.RewardWeightRange.Min) || msg.RewardWeight.GT(msg.RewardWeightRange.Max) {
		return nil, status.Errorf(codes.InvalidArgument, "Alliance rewardWeight must be bounded in RewardWeightRange")
	}

	if msg.TakeRate.IsNil() || msg.TakeRate.IsNegative() || msg.TakeRate.GTE(math.LegacyOneDec()) {
		return nil, status.Errorf(codes.InvalidArgument, "Alliance takeRate must be more or equals to 0 but strictly less than 1")
	}

	if msg.RewardChangeRate.IsZero() || msg.RewardChangeRate.IsNegative() {
		return nil, status.Errorf(codes.InvalidArgument, "Alliance rewardChangeRate must be strictly a positive number")
	}

	if msg.RewardChangeInterval < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Alliance rewardChangeInterval must be strictly a positive number")
	}

	if m.GetAuthority() != msg.Authority {
		return nil, sdkerrors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", m.GetAuthority(), msg.Authority)
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	_, found := m.GetAssetByDenom(sdkCtx, msg.Denom)

	if found {
		return nil, types.ErrAlreadyExists
	}
	rewardStartTime := sdkCtx.BlockTime().Add(m.RewardDelayTime(sdkCtx))
	asset := types.AllianceAsset{
		Denom:                msg.Denom,
		RewardWeight:         msg.RewardWeight,
		RewardWeightRange:    msg.RewardWeightRange,
		TakeRate:             msg.TakeRate,
		TotalTokens:          math.ZeroInt(),
		TotalValidatorShares: math.LegacyZeroDec(),
		RewardStartTime:      rewardStartTime,
		RewardChangeRate:     msg.RewardChangeRate,
		RewardChangeInterval: msg.RewardChangeInterval,
		LastRewardChangeTime: rewardStartTime,
	}
	m.SetAsset(sdkCtx, asset)
	return &types.MsgCreateAllianceResponse{}, nil
}

func (m MsgServer) UpdateAlliance(ctx context.Context, msg *types.MsgUpdateAlliance) (*types.MsgUpdateAllianceResponse, error) {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return nil, sdkerrors.Wrap(err, "invalid authority address")
	}

	if msg.Denom == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Alliance denom must have a value")
	}

	if msg.RewardWeight.IsNil() || msg.RewardWeight.LT(math.LegacyZeroDec()) {
		return nil, status.Errorf(codes.InvalidArgument, "Alliance rewardWeight must be zero or a positive number")
	}

	if msg.TakeRate.IsNil() || msg.TakeRate.IsNegative() || msg.TakeRate.GTE(math.LegacyOneDec()) {
		return nil, status.Errorf(codes.InvalidArgument, "Alliance takeRate must be more or equals to 0 but strictly less than 1")
	}

	if msg.RewardChangeRate.IsZero() || msg.RewardChangeRate.IsNegative() {
		return nil, status.Errorf(codes.InvalidArgument, "Alliance rewardChangeRate must be strictly a positive number")
	}

	if msg.RewardChangeInterval < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Alliance rewardChangeInterval must be strictly a positive number")
	}

	if m.GetAuthority() != msg.Authority {
		return nil, sdkerrors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", m.GetAuthority(), msg.Authority)
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	asset, found := m.GetAssetByDenom(sdkCtx, msg.Denom)

	if !found {
		return nil, types.ErrUnknownAsset
	}
	asset.RewardWeightRange = msg.RewardWeightRange
	if asset.RewardWeightRange.Min.GT(msg.RewardWeight) || asset.RewardWeightRange.Max.LT(msg.RewardWeight) {
		return nil, types.ErrRewardWeightOutOfBound
	}
	asset.RewardWeight = msg.RewardWeight
	asset.TakeRate = msg.TakeRate
	asset.RewardChangeRate = msg.RewardChangeRate
	asset.RewardChangeInterval = msg.RewardChangeInterval

	err := m.UpdateAllianceAsset(sdkCtx, asset)
	if err != nil {
		return nil, err
	}

	return &types.MsgUpdateAllianceResponse{}, nil
}

func (m MsgServer) DeleteAlliance(ctx context.Context, msg *types.MsgDeleteAlliance) (*types.MsgDeleteAllianceResponse, error) {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return nil, sdkerrors.Wrap(err, "invalid authority address")
	}

	if msg.Denom == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Alliance denom must have a value")
	}

	if m.GetAuthority() != msg.Authority {
		return nil, sdkerrors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", m.GetAuthority(), msg.Authority)
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	asset, found := m.GetAssetByDenom(sdkCtx, msg.Denom)

	if !found {
		return nil, types.ErrUnknownAsset
	}

	if asset.TotalTokens.GT(math.ZeroInt()) {
		return nil, types.ErrActiveDelegationsExists
	}

	err := m.DeleteAsset(sdkCtx, asset)
	if err != nil {
		return nil, err
	}

	return &types.MsgDeleteAllianceResponse{}, nil
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &MsgServer{Keeper: keeper}
}
