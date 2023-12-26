package keeper

import (
	"fmt"

	"github.com/terra-money/alliance/x/alliance/types"

	"github.com/cometbft/cometbft/libs/log"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper struct {
	storeKey           storetypes.StoreKey
	cdc                codec.BinaryCodec
	accountKeeper      types.AccountKeeper
	bankKeeper         types.BankKeeper
	stakingKeeper      types.StakingKeeper
	distributionKeeper types.DistributionKeeper
	feeCollectorName   string // name of the FeeCollector ModuleAccount
	authorityAddr      string // name of the Gov ModuleAccount for permissioned messages
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	stakingKeeper types.StakingKeeper,
	distributionKeeper types.DistributionKeeper,
	feeCollectorName string,
	authorityAddr string,
) Keeper {
	// make sure the fee collector module account exists
	if accountKeeper.GetModuleAddress(feeCollectorName) == nil {
		panic(fmt.Sprintf("%s module account has not been set", feeCollectorName))
	}

	return Keeper{
		storeKey:           storeKey,
		cdc:                cdc,
		accountKeeper:      accountKeeper,
		bankKeeper:         bankKeeper,
		stakingKeeper:      stakingKeeper,
		distributionKeeper: distributionKeeper,
		feeCollectorName:   feeCollectorName,
		authorityAddr:      authorityAddr,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) StakingHooks() Hooks {
	return Hooks{
		k: k,
	}
}

func (k Keeper) StoreKey() storetypes.StoreKey {
	return k.storeKey
}

func (k Keeper) GetAuthority() string {
	return k.authorityAddr
}
