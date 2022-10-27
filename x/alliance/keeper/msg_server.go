package keeper

import (
	"alliance/x/alliance/types"
	"context"

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

func (m MsgServer) Delegate(ctx context.Context, delegate *types.MsgDelegate) (*types.MsgDelegateResponse, error) {
	err := delegate.ValidateBasic()
	if err != nil {
		return nil, err
	}

	delAddr, err := sdk.AccAddressFromBech32(delegate.DelegatorAddress)
	if err != nil {
		return nil, err
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	valAddr, err := sdk.ValAddressFromBech32(delegate.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	validator, err := m.Keeper.GetAllianceValidator(sdkCtx, valAddr)
	if err != nil {
		return nil, err
	}

	_, err = m.Keeper.Delegate(sdkCtx, delAddr, validator, delegate.Amount)
	if err != nil {
		return nil, err
	}
	return &types.MsgDelegateResponse{}, nil
}

func (m MsgServer) Redelegate(ctx context.Context, redelegate *types.MsgRedelegate) (*types.MsgRedelegateResponse, error) {
	err := redelegate.ValidateBasic()
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	delAddr, err := sdk.AccAddressFromBech32(redelegate.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	srcValAddr, err := sdk.ValAddressFromBech32(redelegate.ValidatorSrcAddress)
	if err != nil {
		return nil, err
	}
	dstValAddr, err := sdk.ValAddressFromBech32(redelegate.ValidatorDstAddress)
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

	_, err = m.Keeper.Redelegate(sdkCtx, delAddr, srcValidator, dstValidator, redelegate.Amount)
	if err != nil {
		return nil, err
	}
	return &types.MsgRedelegateResponse{}, nil
}

func (m MsgServer) Undelegate(ctx context.Context, undelegate *types.MsgUndelegate) (*types.MsgUndelegateResponse, error) {
	err := undelegate.ValidateBasic()
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	delAddr, err := sdk.AccAddressFromBech32(undelegate.DelegatorAddress)
	if err != nil {
		return nil, err
	}
	valAddr, err := sdk.ValAddressFromBech32(undelegate.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	validator, err := m.Keeper.GetAllianceValidator(sdkCtx, valAddr)
	if err != nil {
		return nil, err
	}

	err = m.Keeper.Undelegate(sdkCtx, delAddr, validator, undelegate.Amount)
	if err != nil {
		return nil, err
	}
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

	m.Keeper.SetAsset(sdkCtx, req.Alliance)

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

func (m MsgServer) ClaimDelegationRewards(ctx context.Context, request *types.MsgClaimDelegationRewards) (*types.MsgClaimDelegationRewardsResponse, error) {
	err := request.ValidateBasic()
	if err != nil {
		return nil, err
	}

	delAddr, err := sdk.AccAddressFromBech32(request.DelegatorAddress)
	if err != nil {
		return nil, err
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	valAddr, err := sdk.ValAddressFromBech32(request.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	validator, err := m.Keeper.GetAllianceValidator(sdkCtx, valAddr)
	if err != nil {
		return nil, err
	}

	_, err = m.Keeper.ClaimDelegationRewards(sdkCtx, delAddr, validator, request.Denom)
	return &types.MsgClaimDelegationRewardsResponse{}, err
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &MsgServer{Keeper: keeper}
}
