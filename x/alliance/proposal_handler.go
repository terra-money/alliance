package alliance

import (
	"github.com/terra-money/alliance/x/alliance/keeper"
	"github.com/terra-money/alliance/x/alliance/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func NewAllianceProposalHandler(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.MsgCreateAllianceProposal:
			return k.CreateAlliance(sdk.WrapSDKContext(ctx), c)
		case *types.MsgUpdateAllianceProposal:
			return k.UpdateAlliance(sdk.WrapSDKContext(ctx), c)
		case *types.MsgDeleteAllianceProposal:
			return k.DeleteAlliance(sdk.WrapSDKContext(ctx), c)

		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized alliance proposal content type: %T", c)
		}
	}
}
