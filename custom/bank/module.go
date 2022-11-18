package bank

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	bankmodule "github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	custombankkeeper "github.com/terra-money/alliance/custom/bank/keeper"
)

// AppModule wraps around the bank module and the bank keeper to return the right total supply ignoring bonded tokens
// that the alliance module minted to rebalance the voting power
// It modifies the TotalSupply and SupplyOf GRPC queries
type AppModule struct {
	bankmodule.AppModule
	keeper custombankkeeper.Keeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(cdc codec.Codec, keeper custombankkeeper.Keeper, accountKeeper types.AccountKeeper) AppModule {
	bankModule := bankmodule.NewAppModule(cdc, keeper, accountKeeper)
	return AppModule{
		AppModule: bankModule,
		keeper:    keeper,
	}
}

// RegisterServices registers module services.
// NOTE: Overriding this method as not doing so will cause a panic
// when trying to force this custom keeper into a bankkeeper.BaseKeeper
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), bankkeeper.NewMsgServerImpl(am.keeper))
	types.RegisterQueryServer(cfg.QueryServer(), am.keeper)

	m := bankkeeper.NewMigrator(am.keeper.BaseKeeper)
	if err := cfg.RegisterMigration(types.ModuleName, 1, m.Migrate1to2); err != nil {
		panic(fmt.Sprintf("failed to migrate x/bank from version 1 to 2: %v", err))
	}

	if err := cfg.RegisterMigration(types.ModuleName, 2, m.Migrate2to3); err != nil {
		panic(fmt.Sprintf("failed to migrate x/bank from version 2 to 3: %v", err))
	}
}
