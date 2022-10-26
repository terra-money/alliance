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

	delAddr, err := sdk.AccAddressFromBech32(req.DelegatorAddr)
	if err != nil {
		return nil, err
	}

	// Get the key-value module store using the store key
	store := ctx.KVStore(k.storeKey)

	// Get the specific delegations key
	key := types.GetDelegationsKey(delAddr)

	// Get the part of the store that keeps assets
	delegationsStore := prefix.NewStore(store, key)

	// Paginate the assets store based on PageRequest
	pageRes, err := query.Paginate(delegationsStore, req.Pagination, func(key []byte, value []byte) error {
		var delegation types.Delegation
		if err := k.cdc.Unmarshal(value, &delegation); err != nil {
			return err
		}

		asset, found := k.GetAssetByDenom(ctx, delegation.Denom)
		if !found {
			return types.ErrUnknownAsset
		}

		valAddr, _ := sdk.ValAddressFromBech32(delegation.ValidatorAddress)
		validator, err := k.GetAllianceValidator(ctx, valAddr)
		if err != nil {
			return err
		}
		balance := types.GetDelegationTokens(delegation, validator, asset)

		delegationRes := types.DelegationResponse{
			Delegation: delegation,
			Balance:    balance,
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

	delAddr, err := sdk.AccAddressFromBech32(req.DelegatorAddr)
	if err != nil {
		return nil, err
	}

	valAddr, err := sdk.ValAddressFromBech32(req.ValidatorAddr)
	if err != nil {
		return nil, err
	}

	_, found := k.stakingKeeper.GetValidator(ctx, valAddr)
	if !found {
		return nil, status.Errorf(codes.NotFound, "Validator not found by address %s", req.ValidatorAddr)
	}

	// Get the key-value module store using the store key
	store := ctx.KVStore(k.storeKey)

	// Get the specific delegations key
	key := types.GetDelegationsKeyForAllDenoms(delAddr, valAddr)

	// Get the part of the store that keeps assets
	delegationStore := prefix.NewStore(store, key)

	// Paginate the assets store based on PageRequest
	pageRes, err := query.Paginate(delegationStore, req.Pagination, func(key []byte, value []byte) error {
		var delegation types.Delegation
		if err := k.cdc.Unmarshal(value, &delegation); err != nil {
			return err
		}

		asset, found := k.GetAssetByDenom(ctx, delegation.Denom)
		if !found {
			return types.ErrUnknownAsset
		}

		valAddr, _ := sdk.ValAddressFromBech32(delegation.ValidatorAddress)
		validator, err := k.GetAllianceValidator(ctx, valAddr)
		if err != nil {
			return err
		}
		balance := types.GetDelegationTokens(delegation, validator, asset)

		delegationRes := types.DelegationResponse{
			Delegation: delegation,
			Balance:    balance,
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

	delAddr, err := sdk.AccAddressFromBech32(req.DelegatorAddr)
	if err != nil {
		return nil, err
	}

	valAddr, err := sdk.ValAddressFromBech32(req.ValidatorAddr)
	if err != nil {
		return nil, err
	}

	validator, err := k.GetAllianceValidator(ctx, valAddr)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Validator not found by address %s", req.ValidatorAddr)
	}

	asset, found := k.GetAssetByDenom(ctx, req.Denom)

	if !found {
		return nil, status.Errorf(codes.NotFound, "AllianceAsset not found by denom %s", req.Denom)
	}

	delegation, found := k.GetDelegation(ctx, delAddr, validator, req.Denom)
	if !found {
		return &types.QueryAllianceDelegationResponse{
			Delegation: types.DelegationResponse{
				Delegation: types.NewDelegation(delAddr, valAddr, req.Denom, sdk.ZeroDec(), []types.RewardHistory{}),
				Balance:    sdk.NewCoin(req.Denom, sdk.ZeroInt()),
			}}, nil
	}

	balance := types.GetDelegationTokens(delegation, validator, asset)
	return &types.QueryAllianceDelegationResponse{
		Delegation: types.DelegationResponse{
			Delegation: delegation,
			Balance:    balance,
		},
	}, nil
}
