package keeper

import (
	"alliance/x/alliance/types"
	"context"
	"time"

	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gov "github.com/cosmos/cosmos-sdk/x/gov/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	newShares, err := m.Keeper.Delegate(sdkCtx, delAddr, validator, msg.Amount)
	if err != nil {
		return nil, err
	}

	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDelegate,
			sdk.NewAttribute(types.AttributeKeyValidator, msg.ValidatorAddress),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyNewShares, newShares.String()),
		),
	})

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

	completionTime, err := m.Keeper.Redelegate(sdkCtx, delAddr, srcValidator, dstValidator, msg.Amount)
	if err != nil {
		return nil, err
	}

	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRedelegate,
			sdk.NewAttribute(types.AttributeKeySrcValidator, msg.ValidatorSrcAddress),
			sdk.NewAttribute(types.AttributeKeyDstValidator, msg.ValidatorDstAddress),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyCompletionTime, completionTime.Format(time.RFC3339)),
		),
	})
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

	completionTime, err := m.Keeper.Undelegate(sdkCtx, delAddr, validator, msg.Amount)
	if err != nil {
		return nil, err
	}

	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUndelegate,
			sdk.NewAttribute(types.AttributeKeyValidator, msg.ValidatorAddress),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyCompletionTime, completionTime.Format(time.RFC3339)),
		),
	})
	return &types.MsgUndelegateResponse{}, nil
}

func (m MsgServer) CreateAlliance(ctx context.Context, req *types.MsgCreateAlliance) (*types.MsgCreateAllianceResponse, error) {
	err := req.ValidateBasic()
	if err != nil {
		return nil, err
	}

	if m.Keeper.authority != req.Authority {
		return nil, errors.Wrapf(gov.ErrInvalidSigner, "expected %s got %s", m.Keeper.authority, req.Authority)
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	_, found := m.Keeper.GetAssetByDenom(sdkCtx, req.Alliance.Denom)

	if found {
		return nil, status.Errorf(codes.AlreadyExists, "Asset with denom: %s already exists", req.Alliance.Denom)
	}

	rewardStartTime := sdkCtx.BlockTime().Add(m.Keeper.RewardDelayTime(sdkCtx))
	asset := types.NewAllianceAsset(req.Alliance.Denom, req.Alliance.RewardWeight, req.Alliance.TakeRate, rewardStartTime)
	m.Keeper.SetAsset(sdkCtx, asset)

	return &types.MsgCreateAllianceResponse{}, nil
}

func (m MsgServer) UpdateAlliance(ctx context.Context, req *types.MsgUpdateAlliance) (*types.MsgUpdateAllianceResponse, error) {
	err := req.ValidateBasic()
	if err != nil {
		return nil, err
	}

	if m.Keeper.authority != req.Authority {
		return nil, errors.Wrapf(gov.ErrInvalidSigner, "expected %s got %s", m.Keeper.authority, req.Authority)
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	asset, found := m.Keeper.GetAssetByDenom(sdkCtx, req.Denom)

	if !found {
		return nil, status.Errorf(codes.NotFound, "Asset with denom: %s does not exist", req.Denom)
	}

	asset.RewardWeight = req.RewardWeight
	asset.TakeRate = req.TakeRate

	err = m.Keeper.UpdateAllianceAsset(sdkCtx, asset)
	if err != nil {
		return nil, err
	}
	return &types.MsgUpdateAllianceResponse{}, nil
}

func (m MsgServer) DeleteAlliance(ctx context.Context, req *types.MsgDeleteAlliance) (*types.MsgDeleteAllianceResponse, error) {
	err := req.ValidateBasic()
	if err != nil {
		return nil, err
	}

	if m.Keeper.authority != req.Authority {
		return nil, errors.Wrapf(gov.ErrInvalidSigner, "expected %s got %s", m.Keeper.authority, req.Authority)
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	asset, found := m.Keeper.GetAssetByDenom(sdkCtx, req.Denom)

	if !found {
		return nil, status.Errorf(codes.NotFound, "Asset with denom: %s does not exist", req.Denom)
	}

	if asset.TotalTokens.GT(math.ZeroInt()) {
		return nil, status.Errorf(codes.Internal, "Asset cannot be deleted because there are still %s delegations associated with it", asset.TotalTokens)
	}

	m.Keeper.DeleteAsset(sdkCtx, req.Denom)

	return &types.MsgDeleteAllianceResponse{}, nil
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

	coins, err := m.Keeper.ClaimDelegationRewards(sdkCtx, delAddr, validator, msg.Denom)

	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeClaimDelegationRewards,
			sdk.NewAttribute(types.AttributeKeyValidator, msg.ValidatorAddress),
			sdk.NewAttribute(sdk.AttributeKeyAmount, coins.String()),
		),
	})
	return &types.MsgClaimDelegationRewardsResponse{}, err
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &MsgServer{Keeper: keeper}
}
