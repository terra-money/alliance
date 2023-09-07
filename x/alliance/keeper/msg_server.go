package keeper

import (
	"context"

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
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
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
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
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
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
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
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
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
	if m.GetAuthority() != msg.Authority {
		return nil, sdkerrors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", m.GetAuthority(), msg.Authority)
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if err := m.SetParams(sdkCtx, msg.Plan); err != nil {
		return nil, err
	}
	return &types.MsgUpdateParamsResponse{}, nil
}

func (m MsgServer) CreateAlliance(ctx context.Context, req *types.MsgCreateAlliance) (*types.MsgCreateAllianceResponse, error) {
	if m.GetAuthority() != req.Authority {
		return nil, sdkerrors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", m.GetAuthority(), req.Authority)
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	_, found := m.GetAssetByDenom(sdkCtx, req.Plan.Denom)

	if found {
		return nil, types.ErrAlreadyExists
	}
	rewardStartTime := sdkCtx.BlockTime().Add(m.RewardDelayTime(sdkCtx))
	asset := types.AllianceAsset{
		Denom:                req.Plan.Denom,
		RewardWeight:         req.Plan.RewardWeight,
		RewardWeightRange:    req.Plan.RewardWeightRange,
		TakeRate:             req.Plan.TakeRate,
		TotalTokens:          sdk.ZeroInt(),
		TotalValidatorShares: sdk.ZeroDec(),
		RewardStartTime:      rewardStartTime,
		RewardChangeRate:     req.Plan.RewardChangeRate,
		RewardChangeInterval: req.Plan.RewardChangeInterval,
		LastRewardChangeTime: rewardStartTime,
	}
	m.SetAsset(sdkCtx, asset)
	return &types.MsgCreateAllianceResponse{}, nil
}

func (m MsgServer) UpdateAlliance(ctx context.Context, req *types.MsgUpdateAlliance) (*types.MsgUpdateAllianceResponse, error) {
	if m.GetAuthority() != req.Authority {
		return nil, sdkerrors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", m.GetAuthority(), req.Authority)
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	asset, found := m.GetAssetByDenom(sdkCtx, req.Plan.Denom)

	if !found {
		return nil, types.ErrUnknownAsset
	}
	if asset.RewardWeightRange.Min.GT(req.Plan.RewardWeight) || asset.RewardWeightRange.Max.LT(req.Plan.RewardWeight) {
		return nil, types.ErrRewardWeightOutOfBound
	}
	asset.RewardWeight = req.Plan.RewardWeight
	asset.TakeRate = req.Plan.TakeRate
	asset.RewardChangeRate = req.Plan.RewardChangeRate
	asset.RewardChangeInterval = req.Plan.RewardChangeInterval

	err := m.UpdateAllianceAsset(sdkCtx, asset)
	if err != nil {
		return nil, err
	}

	return &types.MsgUpdateAllianceResponse{}, nil
}

func (m MsgServer) DeleteAlliance(ctx context.Context, req *types.MsgDeleteAlliance) (*types.MsgDeleteAllianceResponse, error) {
	if m.GetAuthority() != req.Authority {
		return nil, sdkerrors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", m.GetAuthority(), req.Authority)
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	asset, found := m.GetAssetByDenom(sdkCtx, req.Plan.Denom)

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
