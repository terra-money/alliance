package types

import (
	"context"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

type StakingKeeper interface {
	UnbondingTime(ctx context.Context) (res time.Duration, err error)
	Delegate(ctx context.Context, delAddr sdk.AccAddress, bondAmt math.Int, tokenSrc types.BondStatus,
		validator types.Validator, subtractAccount bool) (newShares math.LegacyDec, err error)
	BeginRedelegation(
		ctx context.Context, delAddr sdk.AccAddress, valSrcAddr, valDstAddr sdk.ValAddress, sharesAmount math.LegacyDec,
	) (completionTime time.Time, err error)
	GetValidator(ctx context.Context, addr sdk.ValAddress) (validator types.Validator, err error)
	BondDenom(ctx context.Context) (res string, err error)
	ValidateUnbondAmount(
		ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, amt math.Int,
	) (shares math.LegacyDec, err error)
	RemoveRedelegation(ctx context.Context, red types.Redelegation) (err error)
	Unbond(
		ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, shares math.LegacyDec,
	) (amount math.Int, err error)
	GetDelegation(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (delegation types.Delegation, err error)
	TotalBondedTokens(ctx context.Context) (math.Int, error)
	GetDelegatorBonded(ctx context.Context, delegator sdk.AccAddress) (math.Int, error)
	RemoveValidatorTokensAndShares(ctx context.Context, validator types.Validator,
		sharesToRemove math.LegacyDec,
	) (valOut types.Validator, removedTokens math.Int, err error)
	RemoveValidatorTokens(ctx context.Context,
		validator types.Validator, tokensToRemove math.Int,
	) (types.Validator, error)
	IterateDelegatorDelegations(ctx context.Context, delegator sdk.AccAddress, cb func(delegation types.Delegation) (stop bool)) error
	GetAllValidators(ctx context.Context) ([]types.Validator, error)
}

type BankKeeper interface {
	MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx context.Context, senderModule, recipientModule string, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	GetAllBalances(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
	SpendableCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins
}

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
	GetModuleAddress(name string) sdk.AccAddress
	GetModuleAccount(ctx context.Context, name string) sdk.ModuleAccountI
}

type DistributionKeeper interface {
	WithdrawDelegationRewards(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (sdk.Coins, error)
}
