package cli

import (
	"context"

	"alliance/x/alliance/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func CmdQueryRewards() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rewards delegator-addr validator-addr denom",
		Short: "Query rewards generated for a specific alliance by delegator-addr validator-addr denom",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			delegatorAddr := args[0]
			validatorAddr := args[1]
			denom := args[2]
			ctx := client.GetClientContextFromCmd(cmd)
			query := types.NewQueryClient(ctx)
			params := &types.QueryAllianceDelegationRewardsRequest{
				DelegatorAddr: delegatorAddr,
				ValidatorAddr: validatorAddr,
				Denom:         denom,
			}

			res, err := query.AllianceDelegationRewards(context.Background(), params)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
