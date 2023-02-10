package gov

import (
	"github.com/cosmos/cosmos-sdk/codec"
	govmodule "github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/gov/types"

	customgovkeeper "github.com/terra-money/alliance/custom/gov/keeper"
)

type AppModule struct {
	govmodule.AppModule
	keeper customgovkeeper.Keeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(cdc codec.Codec, keeper customgovkeeper.Keeper, accountKeeper types.AccountKeeper, bankKeeper types.BankKeeper) AppModule {
	govmodule := govmodule.NewAppModule(cdc, keeper, accountKeeper, bankKeeper)
	return AppModule{
		AppModule: govmodule,
		keeper:    keeper,
	}
}
