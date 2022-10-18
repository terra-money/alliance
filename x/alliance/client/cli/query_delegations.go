package cli

import (
	"context"
	"strconv"

	"alliance/x/alliance/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdQueryAlliancesDelegation() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegations [delegator_addr]",
		Short: "Query all paginated alliance delegations for a delegator_addr",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			delegator_addr := args[0]
			ctx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			query := types.NewQueryClient(ctx)

			params := &types.QueryAlliancesDelegationsRequest{
				DelegatorAddr: delegator_addr,
				Pagination:    pageReq,
			}

			res, err := query.AlliancesDelegation(context.Background(), params)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryAlliancesDelegationByValidator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegations_validator [delegator_addr] [validator_addr]",
		Short: "Query all paginated alliance delegations for a delegator_addr and validator_addr",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			delegator_addr := args[0]
			validator_addr := args[1]
			ctx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			query := types.NewQueryClient(ctx)

			params := &types.QueryAlliancesDelegationByValidatorRequest{
				Pagination:    pageReq,
				DelegatorAddr: delegator_addr,
				ValidatorAddr: validator_addr,
			}

			res, err := query.AlliancesDelegationByValidator(context.Background(), params)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryAllianceDelegation() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegation [delegator_addr] [validator_addr] [denom]",
		Short: "Query a delegation to an alliance by delegator_addr, validator_addr and denom",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			delegator_addr := args[0]
			validator_addr := args[1]
			denom := args[2]
			ctx := client.GetClientContextFromCmd(cmd)

			if err != nil {
				return err
			}
			query := types.NewQueryClient(ctx)

			params := &types.QueryAllianceDelegationRequest{
				DelegatorAddr: delegator_addr,
				ValidatorAddr: validator_addr,
				Denom:         denom,
			}

			res, err := query.AllianceDelegation(context.Background(), params)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
