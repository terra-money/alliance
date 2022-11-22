package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type StakingKeeper interface {
	BondDenom(ctx sdk.Context) (res string)
}
