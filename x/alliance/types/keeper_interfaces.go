package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	"time"
)

type StakingKeeper interface {
	UnbondingTime(ctx sdk.Context) (res time.Duration)
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
	Unbond(
		ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, shares sdk.Dec,
	) (amount math.Int, err error)
	GetDelegation(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (delegation types.Delegation, found bool)
	TotalBondedTokens(ctx sdk.Context) math.Int
	GetDelegatorBonded(ctx sdk.Context, delegator sdk.AccAddress) math.Int
	RemoveValidatorTokensAndShares(ctx sdk.Context, validator types.Validator,
		sharesToRemove sdk.Dec,
	) (valOut types.Validator, removedTokens math.Int)
	RemoveValidatorTokens(ctx sdk.Context,
		validator types.Validator, tokensToRemove math.Int,
	) types.Validator
	IterateDelegatorDelegations(ctx sdk.Context, delegator sdk.AccAddress, cb func(delegation types.Delegation) (stop bool))
	GetAllValidators(ctx sdk.Context) (validators []types.Validator)
}

type BankKeeper interface {
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
}

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	GetModuleAddress(name string) sdk.AccAddress
	GetModuleAccount(ctx sdk.Context, name string) authtypes.ModuleAccountI
}

type DistributionKeeper interface {
	WithdrawDelegationRewards(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (sdk.Coins, error)
}
