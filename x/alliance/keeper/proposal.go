package keeper

import (
	"context"

	"github.com/terra-money/alliance/x/alliance/types"
)

func (k Keeper) CreateAlliance(ctx context.Context, req *types.MsgCreateAllianceProposal) error {
	ms := MsgServer{k}
	_, err := ms.CreateAlliance(ctx, &types.MsgCreateAlliance{
		Authority:            k.GetAuthority(),
		Denom:                req.Denom,
		RewardWeight:         req.RewardWeight,
		TakeRate:             req.TakeRate,
		RewardChangeRate:     req.RewardChangeRate,
		RewardChangeInterval: req.RewardChangeInterval,
		RewardWeightRange:    req.RewardWeightRange,
	})
	return err
}

func (k Keeper) UpdateAlliance(ctx context.Context, req *types.MsgUpdateAllianceProposal) error {
	ms := MsgServer{k}
	_, err := ms.UpdateAlliance(ctx, &types.MsgUpdateAlliance{
		Authority:            k.GetAuthority(),
		Denom:                req.Denom,
		RewardWeight:         req.RewardWeight,
		TakeRate:             req.TakeRate,
		RewardChangeRate:     req.RewardChangeRate,
		RewardWeightRange:    req.RewardWeightRange,
		RewardChangeInterval: req.RewardChangeInterval,
	})
	return err
}

func (k Keeper) DeleteAlliance(ctx context.Context, req *types.MsgDeleteAllianceProposal) error {
	ms := MsgServer{k}
	_, err := ms.DeleteAlliance(ctx, &types.MsgDeleteAlliance{
		Authority: k.GetAuthority(),
		Denom:     req.Denom,
	})
	return err
}
