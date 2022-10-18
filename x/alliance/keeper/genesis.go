package keeper

import (
	"alliance/x/alliance/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, g *types.GenesisState) []abci.ValidatorUpdate {
	k.SetParams(ctx, g.Params)
	for _, asset := range g.Assets {
		k.SetAsset(ctx, asset)
	}
	return []abci.ValidatorUpdate{}
}

func (k Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	k.paramstore.SetParamSet(ctx, &params)
	return nil
}
