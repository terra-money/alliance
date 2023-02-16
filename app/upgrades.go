package app

import v2 "github.com/terra-money/alliance/app/upgrades/v2"

func (app *App) setupUpgradeHandlers() {
	// v2 upgrade handler
	app.UpgradeKeeper.SetUpgradeHandler(
		v2.UpgradeName,
		v2.CreateUpgradeHandler(app.mm, app.configurator),
	)
}
