package cli

import (
	"alliance/x/alliance/types"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdQueryParams())

	cmd.AddCommand(CmdQueryAlliances())
	cmd.AddCommand(CmdQueryAlliance())

	cmd.AddCommand(CmdQueryAlliancesDelegation())
	cmd.AddCommand(CmdQueryAlliancesDelegationByValidator())
	cmd.AddCommand(CmdQueryAllianceDelegation())
	cmd.AddCommand(CmdQueryRewards())

	return cmd
}
