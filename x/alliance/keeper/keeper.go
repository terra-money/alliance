package keeper

import (
	"fmt"

	"github.com/terra-money/alliance/x/alliance/types"

	"github.com/cometbft/cometbft/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

type Keeper struct {
	storeKey           storetypes.StoreKey
	paramstore         paramtypes.Subspace
	cdc                codec.BinaryCodec
	accountKeeper      types.AccountKeeper
	bankKeeper         types.BankKeeper
	stakingKeeper      types.StakingKeeper
	distributionKeeper types.DistributionKeeper
	feeCollectorName   string // name of the FeeCollector ModuleAccount
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	stakingKeeper types.StakingKeeper,
	distributionKeeper types.DistributionKeeper,
	feeCollectorName string,
) Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		kt := paramtypes.NewKeyTable().RegisterParamSet(&types.Params{})
		ps = ps.WithKeyTable(kt)
	}

	// make sure the fee collector module account exists
	if accountKeeper.GetModuleAddress(feeCollectorName) == nil {
		panic(fmt.Sprintf("%s module account has not been set", feeCollectorName))
	}

	return Keeper{
		storeKey:           storeKey,
		paramstore:         ps,
		cdc:                cdc,
		accountKeeper:      accountKeeper,
		bankKeeper:         bankKeeper,
		stakingKeeper:      stakingKeeper,
		distributionKeeper: distributionKeeper,
		feeCollectorName:   feeCollectorName,
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
