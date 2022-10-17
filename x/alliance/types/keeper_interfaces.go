package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	"time"
)

type StakingKeeper interface {
	Delegate(ctx sdk.Context, delAddr sdk.AccAddress, bondAmt math.Int, tokenSrc types.BondStatus,
		validator types.Validator, subtractAccount bool) (newShares sdk.Dec, err error)
	BeginRedelegation(
		ctx sdk.Context, delAddr sdk.AccAddress, valSrcAddr, valDstAddr sdk.ValAddress, sharesAmount sdk.Dec,
	) (completionTime time.Time, err error)
	GetValidator(ctx sdk.Context, addr sdk.ValAddress) (validator types.Validator, found bool)
	BondDenom(ctx sdk.Context) (res string)
	ValidateUnbondAmount(
		ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, amt math.Int,
	) (shares sdk.Dec, err error)
	RemoveRedelegation(ctx sdk.Context, red types.Redelegation)
}

type BankKeeper interface {
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
}

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	GetModuleAddress(name string) sdk.AccAddress
	GetModuleAccount(ctx sdk.Context, name string) authtypes.ModuleAccountI
}
