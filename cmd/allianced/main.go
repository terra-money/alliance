package main

import (
	"os"

	"github.com/terra-money/alliance/app"
	"github.com/terra-money/alliance/cmd/allianced/cmd"

	"github.com/cosmos/cosmos-sdk/server"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
)

func main() {
	app.SetPrefixes(app.AccountAddressPrefix)
	rootCmd, _ := cmd.NewRootCmd()

	if err := svrcmd.Execute(rootCmd, "ALLIANCED", app.DefaultNodeHome); err != nil {
		switch e := err.(type) {
		case server.ErrorCode:
			os.Exit(e.Code)

		default:
			os.Exit(1)
		}
	}
}
