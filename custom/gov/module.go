package gov

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	govmodule "github.com/cosmos/cosmos-sdk/x/gov"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	"github.com/cosmos/cosmos-sdk/x/gov/types"

	customgovkeeper "github.com/terra-money/alliance/custom/gov/keeper"
)

type AppModule struct {
	govmodule.AppModule
	keeper customgovkeeper.Keeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(cdc codec.Codec, keeper customgovkeeper.Keeper, accountKeeper types.AccountKeeper, bankKeeper types.BankKeeper) AppModule {
	govmodule := govmodule.NewAppModule(cdc, keeper.Keeper, accountKeeper, bankKeeper)
	return AppModule{
		AppModule: govmodule,
		keeper:    keeper,
	}
}

// RegisterServices registers module services.
// NOTE: Overriding this method as not doing so will cause a panic
// when trying to force this custom keeper into a govkeeper
func (am AppModule) RegisterServices(cfg module.Configurator) {
	m := govkeeper.NewMigrator(am.keeper.Keeper)
	if err := cfg.RegisterMigration(types.ModuleName, 1, m.Migrate1to2); err != nil {
		panic(fmt.Sprintf("failed to migrate x/gov from version 1 to 2: %v", err))
	}

	if err := cfg.RegisterMigration(types.ModuleName, 2, m.Migrate2to3); err != nil {
		panic(fmt.Sprintf("failed to migrate x/gov from version 2 to 3: %v", err))
	}
}
