package keeper

import (
	"context"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/terra-money/alliance/x/alliance/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		TakeRate:             req.TakeRate,
		TotalTokens:          sdk.ZeroInt(),
		TotalValidatorShares: sdk.ZeroDec(),
		RewardStartTime:      rewardStartTime,
		RewardDecayRate:      req.RewardDecayRate,
		RewardDecayInterval:  req.RewardDecayInterval,
	}
	k.SetAsset(sdkCtx, asset)

	if asset.HasPositiveDecay() {
		k.QueueRewardWeightDecayEvent(sdkCtx, asset)
	}

	return nil
}

func (k Keeper) UpdateAlliance(ctx context.Context, req *types.MsgUpdateAllianceProposal) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	asset, found := k.GetAssetByDenom(sdkCtx, req.Denom)

	if !found {
		return status.Errorf(codes.NotFound, "Asset with denom: %s does not exist", req.Denom)
	}

	asset.RewardWeight = req.RewardWeight
	asset.TakeRate = req.TakeRate
	asset.RewardDecayRate = req.RewardDecayRate
	asset.RewardDecayInterval = req.RewardDecayInterval

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

	k.DeleteAsset(sdkCtx, req.Denom)

	return nil
}
