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

func (k Keeper) Alliances(c context.Context, req *types.QueryAlliancesRequest) (*types.QueryAlliancesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// Define a variable that will store a list of assets
	var alliances []types.AllianceAsset

	// Get context with the information about the environment
	ctx := sdk.UnwrapSDKContext(c)

	// Get the key-value module store using the store key
	store := ctx.KVStore(k.storeKey)

	// Get the part of the store that keeps assets
	assetsStore := prefix.NewStore(store, []byte(types.AssetKey))

	// Paginate the assets store based on PageRequest
	pageRes, err := query.Paginate(assetsStore, req.Pagination, func(key []byte, value []byte) error {
		var asset types.AllianceAsset
		if err := k.cdc.Unmarshal(value, &asset); err != nil {
			return err
		}

		alliances = append(alliances, asset)

		return nil
	})

	// Throw an error if pagination failed
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Return a struct containing a list of assets and pagination info
	return &types.QueryAlliancesResponse{
		Alliances:  alliances,
		Pagination: pageRes,
	}, nil
}

func (k Keeper) Alliance(c context.Context, req *types.QueryAllianceRequest) (*types.QueryAllianceResponse, error) {
	var asset types.AllianceAsset

	// Get context with the information about the environment
	ctx := sdk.UnwrapSDKContext(c)

	// Get the part of the store that keeps assets
	asset, _ = k.GetAssetByDenom(ctx, req.Denom)

	// Return parsed asset, true since the asset exists
	return &types.QueryAllianceResponse{
		Alliance: &asset,
	}, nil
}
