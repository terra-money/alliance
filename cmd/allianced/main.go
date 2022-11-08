package main

import (
	"os"

	"github.com/terra-money/alliance/app"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ignite/cli/ignite/pkg/cosmoscmd"
	"github.com/ignite/cli/ignite/pkg/xstrings"
)

func main() {
	rootCmd, _ := cosmoscmd.NewRootCmd(
		app.Name,
		app.AccountAddressPrefix,
		app.DefaultNodeHome,
		xstrings.NoDash(app.Name),
		app.ModuleBasics,
		app.New,
		// this line is used by starport scaffolding # root/arguments
	)
	rootCmd.AddCommand(NewTestnetCmd(app.ModuleBasics, banktypes.GenesisBalancesIterator{}))
	if err := svrcmd.Execute(rootCmd, "", app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
