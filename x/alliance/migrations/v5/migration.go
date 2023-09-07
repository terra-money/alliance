package v5

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	alliancekeeper "github.com/terra-money/alliance/x/alliance/keeper"
	"github.com/terra-money/alliance/x/alliance/types"
)

func Migrate(k alliancekeeper.Keeper, subspace paramtypes.Subspace) func(ctx sdk.Context) error {
	return func(ctx sdk.Context) error {
		var params types.Params
		subspace.GetParamSet(ctx, &params)
		return k.SetParams(ctx, params)
	}
}
