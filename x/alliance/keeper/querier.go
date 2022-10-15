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

func (k Keeper) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	// Define a variable that will store the params
	var params types.Params

	// Get context with the information about the environment
	ctx := sdk.UnwrapSDKContext(c)

	k.paramstore.GetParamSet(ctx, &params)

	return &types.QueryParamsResponse{
		Params: params,
	}, nil
}

func (k Keeper) Alliances(c context.Context, req *types.QueryAlliancesRequest) (*types.QueryAlliancesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// Define a variable that will store a list of assets
	var assets []*types.AllianceAsset

	// Get context with the information about the environment
	ctx := sdk.UnwrapSDKContext(c)

	// Get the key-value module store using the store key (in our case store key is "chain")
	store := ctx.KVStore(k.storeKey)

	// Get the part of the store that keeps assets (using asset key, which is "asset-value-")
	assetsStore := prefix.NewStore(store, []byte(types.AssetKey))

	// Paginate the assets store based on PageRequest
	pageRes, err := query.Paginate(assetsStore, req.Pagination, func(key []byte, value []byte) error {
		var asset types.AllianceAsset
		if err := k.cdc.Unmarshal(value, &asset); err != nil {
			return err
		}

		assets = append(assets, &asset)

		return nil
	})

	// Throw an error if pagination failed
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Return a struct containing a list of assets and pagination info
	return &types.QueryAlliancesResponse{
		Assets:     assets,
		Pagination: pageRes,
	}, nil
}

func (k Keeper) Alliance(c context.Context, req *types.QueryAllianceRequest) (*types.QueryAllianceResponse, error) {
	// Define a variable that will store a list of assets
	var asset types.AllianceAsset

	// Get context with the information about the environment
	ctx := sdk.UnwrapSDKContext(c)

	// Get the part of the store that keeps assets (using asset key, which is "asset-value-")
	asset = k.GetAssetByDenom(ctx, req.Denom)

	// Return parsed asset, true since the asset exists
	return &types.QueryAllianceResponse{
		Alliance: &asset,
	}, nil
}

func (k Keeper) AllianceDelegations(c context.Context, req *types.QueryAllianceDelegationsRequest) (*types.QueryAllianceDelegationsResponse, error) {
	var delegations []*types.Delegation

	// ctx := sdk.UnwrapSDKContext(c)
	// asset, error = k.GetDelegation(ctx, sdk.AccAddress(req.Denom))

	return &types.QueryAllianceDelegationsResponse{
		Delegations: delegations,
	}, nil
}

func (k Keeper) AllianceDelegation(c context.Context, req *types.QueryAllianceDelegationRequest) (*types.QueryAllianceDelegationResponse, error) {
	var delegations *types.Delegation

	// ctx := sdk.UnwrapSDKContext(c)
	// asset, error = k.GetDelegation(ctx, sdk.AccAddress(req.Denom))

	return &types.QueryAllianceDelegationResponse{
		Delegations: delegations,
	}, nil
}
