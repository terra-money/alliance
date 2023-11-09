package keeper

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/terra-money/alliance/x/alliance/types"
)

func (k Keeper) CreateAlliance(ctx context.Context, req *types.MsgCreateAllianceProposal) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	_, found := k.GetAssetByDenom(sdkCtx, req.Denom)

	if found {
		return status.Errorf(codes.AlreadyExists, "Asset with denom: %s already exists", req.Denom)
	}
	rewardStartTime := sdkCtx.BlockTime().Add(k.RewardDelayTime(sdkCtx))
	asset := types.AllianceAsset{
		Denom:                req.Denom,
		RewardWeight:         req.RewardWeight,
		RewardWeightRange:    req.RewardWeightRange,
		TakeRate:             req.TakeRate,
		TotalTokens:          sdk.ZeroInt(),
		TotalValidatorShares: sdk.ZeroDec(),
		RewardStartTime:      rewardStartTime,
		RewardChangeRate:     req.RewardChangeRate,
		RewardChangeInterval: req.RewardChangeInterval,
		LastRewardChangeTime: rewardStartTime,
	}
	k.SetAsset(sdkCtx, asset)
	return nil
}

func (k Keeper) UpdateAlliance(ctx context.Context, req *types.MsgUpdateAllianceProposal) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	asset, found := k.GetAssetByDenom(sdkCtx, req.Denom)

	if !found {
		return status.Errorf(codes.NotFound, "Asset with denom: %s does not exist", req.Denom)
	}
	asset.RewardWeightRange = req.RewardWeightRange
	if asset.RewardWeightRange.Min.GT(req.RewardWeight) || asset.RewardWeightRange.Max.LT(req.RewardWeight) {
		return types.ErrRewardWeightOutOfBound
	}
	asset.RewardWeight = req.RewardWeight
	asset.TakeRate = req.TakeRate
	asset.RewardChangeRate = req.RewardChangeRate
	asset.RewardChangeInterval = req.RewardChangeInterval

	err := k.UpdateAllianceAsset(sdkCtx, asset)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) DeleteAlliance(ctx context.Context, req *types.MsgDeleteAllianceProposal) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	asset, found := k.GetAssetByDenom(sdkCtx, req.Denom)

	if !found {
		return status.Errorf(codes.NotFound, "Asset with denom: %s does not exist", req.Denom)
	}

	if asset.TotalTokens.GT(math.ZeroInt()) {
		return status.Errorf(codes.Internal, "Asset cannot be deleted because there are still %s delegations associated with it", asset.TotalTokens)
	}

	err := k.DeleteAsset(sdkCtx, asset)
	if err != nil {
		return err
	}

	return nil
}
