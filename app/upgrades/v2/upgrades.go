package v2

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

var (
	UpgradeName      = "v2"
	FaucetAddressHex = "85cbde5c7dc29a8600f25aad2d6966c1f7d8975f"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	bankkeeper minttypes.BankKeeper,
	stakingkeeper stakingkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		// Mint more tokens for the faucet
		bondDenom := stakingkeeper.BondDenom(ctx)
		coins := sdk.NewCoins(sdk.NewCoin(bondDenom, sdk.NewInt(1_000_000_000_000_000)))
		bankkeeper.MintCoins(ctx, minttypes.ModuleName, coins)
		faucetAddress, err := sdk.AccAddressFromHexUnsafe(FaucetAddressHex)
		if err != nil {
			return nil, err
		}
		bankkeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, faucetAddress, coins)

		// Continue with the migration
		return mm.RunMigrations(ctx, configurator, vm)
	}
}
