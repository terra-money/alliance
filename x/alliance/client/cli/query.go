package cli

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-money/alliance/x/alliance/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
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

	cmd.AddCommand(CmdQueryValidator())
	cmd.AddCommand(CmdQueryValidators())

	cmd.AddCommand(CmdQueryAllAlliancesDelegations())
	cmd.AddCommand(CmdQueryAlliancesDelegation())
	cmd.AddCommand(CmdQueryAlliancesDelegationByValidator())
	cmd.AddCommand(CmdQueryAllianceDelegation())
	cmd.AddCommand(CmdQueryRewards())

	return cmd
}

func CmdQueryAlliances() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alliances",
		Short: "Query paginated alliances",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAlliancesRequest{}

			res, err := queryClient.Alliances(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryAlliance() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alliance denom",
		Short: "Query a specific alliance by denom",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			denom := args[0]

			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			query := types.NewQueryClient(ctx)

			params := &types.QueryAllianceRequest{Denom: denom}

			res, err := query.Alliance(cmd.Context(), params)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryValidator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator validator-addr",
		Short: "Query a specific alliance validator by addr",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			valAddr, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			query := types.NewQueryClient(ctx)

			req := &types.QueryAllianceValidatorRequest{ValidatorAddr: valAddr.String()}

			res, err := query.AllianceValidator(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryValidators() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validators",
		Short: "Query all alliance validators",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			query := types.NewQueryClient(ctx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			req := &types.QueryAllAllianceValidatorsRequest{
				Pagination: pageReq,
			}

			res, err := query.AllAllianceValidators(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryAllAlliancesDelegations() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegations",
		Short: "Query all paginated alliances delegations",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			ctx := client.GetClientContextFromCmd(cmd)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			query := types.NewQueryClient(ctx)

			params := &types.QueryAllAlliancesDelegationsRequest{
				Pagination: pageReq,
			}

			res, err := query.AllAlliancesDelegations(context.Background(), params)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryAlliancesDelegation() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegations-by-delegator delegator_addr",
		Short: "Query all paginated alliances delegations for a delegator_addr",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			delegatorAddr := args[0]
			ctx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			query := types.NewQueryClient(ctx)

			params := &types.QueryAlliancesDelegationsRequest{
				DelegatorAddr: delegatorAddr,
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
		Use:   "delegations-by-delegator-and-validator delegator_addr validator_addr",
		Short: "Query all paginated alliance delegations for a delegator_addr and validator_addr",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			delegatorAddr := args[0]
			validatorAddr := args[1]
			ctx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			query := types.NewQueryClient(ctx)

			params := &types.QueryAlliancesDelegationByValidatorRequest{
				Pagination:    pageReq,
				DelegatorAddr: delegatorAddr,
				ValidatorAddr: validatorAddr,
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
		Use:   "delegation delegator_addr validator_addr denom",
		Short: "Query a delegation to an alliance by delegator_addr, validator_addr and denom",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			delegatorAddr := args[0]
			validatorAddr := args[1]
			denom := args[2]
			ctx := client.GetClientContextFromCmd(cmd)

			if err != nil {
				return err
			}
			query := types.NewQueryClient(ctx)

			params := &types.QueryAllianceDelegationRequest{
				DelegatorAddr: delegatorAddr,
				ValidatorAddr: validatorAddr,
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

func CmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query module parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Params(context.Background(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryRewards() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rewards delegator_addr validator_addr denom",
		Short: "Query module parameters",
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
