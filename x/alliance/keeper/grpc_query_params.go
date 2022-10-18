package keeper

import (
	"alliance/x/alliance/types"
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

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
