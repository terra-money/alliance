package keeper

import (
	"alliance/x/alliance/types"
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Querier struct {
	Keeper
}

var _ types.QueryServer = Querier{}

func (k Keeper) AlliancesDelegation(c context.Context, req *types.QueryAlliancesDelegationsRequest) (*types.QueryAlliancesDelegationsResponse, error) {
	var delegationsRes []types.DelegationResponse

	// Get context with the information about the environment
	ctx := sdk.UnwrapSDKContext(c)

	// Get the key-value module store using the store key
	store := ctx.KVStore(k.storeKey)

	// Get the part of the store that keeps assets
	delegationsStore := prefix.NewStore(store, types.GetDelegationsKey(sdk.AccAddress(req.DelegatorAddr)))

	// Paginate the assets store based on PageRequest
	pageRes, err := query.Paginate(delegationsStore, req.Pagination, func(key []byte, value []byte) error {
		var delegation types.Delegation
		if err := k.cdc.Unmarshal(value, &delegation); err != nil {
			return err
		}

		asset, _ := k.GetAssetByDenom(ctx, delegation.Denom)
		valAddr, _ := sdk.ValAddressFromBech32(delegation.ValidatorAddress)
		aVal := k.GetOrCreateValidator(ctx, valAddr)

		delegationRes := types.DelegationResponse{
			Delegation: delegation,
			Balance: sdk.Coin{
				Denom:  delegation.Denom,
				Amount: convertNewShareToToken(asset.TotalTokens, aVal.TotalSharesWithDenom(delegation.Denom), delegation.Shares),
			},
		}

		delegationsRes = append(delegationsRes, delegationRes)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &types.QueryAlliancesDelegationsResponse{
		Delegations: delegationsRes,
		Pagination:  pageRes,
	}, nil
}

func (k Keeper) AlliancesDelegationByValidator(c context.Context, req *types.QueryAlliancesDelegationByValidatorRequest) (*types.QueryAlliancesDelegationsResponse, error) {
	var delegationsRes []types.DelegationResponse
	ctx := sdk.UnwrapSDKContext(c)

	valAddr, err := sdk.ValAddressFromBech32(req.ValidatorAddr)
	if err != nil {
		return nil, err
	}

	// Get the key-value module store using the store key
	store := ctx.KVStore(k.storeKey)

	// Get the part of the store that keeps assets
	delegationStore := prefix.NewStore(store, types.GetDelegationsKeyForAllDenoms(sdk.AccAddress(req.DelegatorAddr), valAddr))

	// Paginate the assets store based on PageRequest
	pageRes, err := query.Paginate(delegationStore, req.Pagination, func(key []byte, value []byte) error {
		var delegation types.Delegation
		if err := k.cdc.Unmarshal(value, &delegation); err != nil {
			return err
		}

		asset, _ := k.GetAssetByDenom(ctx, delegation.Denom)
		valAddr, _ := sdk.ValAddressFromBech32(delegation.ValidatorAddress)
		aVal := k.GetOrCreateValidator(ctx, valAddr)

		delegationRes := types.DelegationResponse{
			Delegation: delegation,
			Balance: sdk.Coin{
				Denom:  delegation.Denom,
				Amount: convertNewShareToToken(asset.TotalTokens, aVal.TotalSharesWithDenom(delegation.Denom), delegation.Shares),
			},
		}

		delegationsRes = append(delegationsRes, delegationRes)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &types.QueryAlliancesDelegationsResponse{
		Delegations: delegationsRes,
		Pagination:  pageRes,
	}, nil
}

func (k Keeper) AllianceDelegation(c context.Context, req *types.QueryAllianceDelegationRequest) (*types.QueryAllianceDelegationResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	valAddr, err := sdk.ValAddressFromBech32(req.ValidatorAddr)
	if err != nil {
		return nil, err
	}

	validator, ok := k.stakingKeeper.GetValidator(ctx, valAddr)
	if !ok {
		return nil, status.Errorf(codes.NotFound, "Cannot recover the validator %s", req.ValidatorAddr)
	}

	asset, _ := k.GetAssetByDenom(ctx, req.Denom)

	delegation, success := k.GetDelegation(ctx, sdk.AccAddress(req.DelegatorAddr), validator, req.Denom)
	if !success {
		return nil, status.Errorf(
			codes.Unknown,
			"Could not find delegation with combination %s %s %s",
			req.DelegatorAddr, req.ValidatorAddr, req.Denom,
		)
	}

	aVal := k.GetOrCreateValidator(ctx, valAddr)
	return &types.QueryAllianceDelegationResponse{
		Delegation: types.DelegationResponse{
			Delegation: delegation,
			Balance: sdk.Coin{
				Denom:  delegation.Denom,
				Amount: convertNewShareToToken(asset.TotalTokens, aVal.TotalSharesWithDenom(delegation.Denom), delegation.Shares),
			},
		},
	}, nil
}
