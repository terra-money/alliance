package alliance

import (
	cosmoserrors "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	"github.com/terra-money/alliance/x/alliance/keeper"
	"github.com/terra-money/alliance/x/alliance/types"
)

func NewAllianceProposalHandler(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.MsgCreateAllianceProposal:
			return k.CreateAlliance(ctx, c)
		case *types.MsgUpdateAllianceProposal:
			return k.UpdateAlliance(ctx, c)
		case *types.MsgDeleteAllianceProposal:
			return k.DeleteAlliance(ctx, c)

		default:
			return cosmoserrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized alliance proposal content type: %T", c)
		}
	}
}
