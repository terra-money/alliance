package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govcli "github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/spf13/cobra"
	"github.com/terra-money/alliance/x/alliance/types"
)

func CreateAlliance() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-alliance denom rewards-weight take-rate",
		Args:  cobra.ExactArgs(3),
		Short: "Create an alliance with the specified parameters",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			title, err := cmd.Flags().GetString(govcli.FlagTitle)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(govcli.FlagDescription)
			if err != nil {
				return err
			}

			rewardWeight, err := sdk.NewDecFromStr(args[1])
			if err != nil {
				return err
			}

			takeRate, err := sdk.NewDecFromStr(args[2])
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			depositStr, err := cmd.Flags().GetString(govcli.FlagDeposit)
			if err != nil {
				return err
			}

			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			content := types.NewMsgCreateAllianceProposal(
				title,
				description,
				args[0],
				rewardWeight,
				takeRate,
			)

			err = content.ValidateBasic()

			if err != nil {
				return err
			}

			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)

			if err != nil {
				return err
			}

			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(govcli.FlagTitle, "", "title of proposal")
	cmd.Flags().String(govcli.FlagDescription, "", "description of proposal")
	cmd.Flags().String(govcli.FlagDeposit, "", "deposit of proposal")
	return cmd
}

func UpdateAlliance() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-alliance denom rewards-weight take-rate",
		Args:  cobra.ExactArgs(3),
		Short: "Update an alliance with the specified parameters",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			title, err := cmd.Flags().GetString(govcli.FlagTitle)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(govcli.FlagDescription)
			if err != nil {
				return err
			}

			rewardWeight, err := sdk.NewDecFromStr(args[1])
			if err != nil {
				return err
			}

			takeRate, err := sdk.NewDecFromStr(args[2])
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			depositStr, err := cmd.Flags().GetString(govcli.FlagDeposit)
			if err != nil {
				return err
			}

			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			content := types.NewMsgCreateAllianceProposal(
				title,
				description,
				args[0],
				rewardWeight,
				takeRate,
			)

			err = content.ValidateBasic()

			if err != nil {
				return err
			}

			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)

			if err != nil {
				return err
			}

			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(govcli.FlagTitle, "", "title of proposal")
	cmd.Flags().String(govcli.FlagDescription, "", "description of proposal")
	cmd.Flags().String(govcli.FlagDeposit, "", "deposit of proposal")
	return cmd
}

func DeleteAlliance() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-alliance denom",
		Args:  cobra.ExactArgs(1),
		Short: "Delete an alliance with the specified denom",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			title, err := cmd.Flags().GetString(govcli.FlagTitle)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(govcli.FlagDescription)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			depositStr, err := cmd.Flags().GetString(govcli.FlagDeposit)
			if err != nil {
				return err
			}

			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			content := types.NewMsgDeleteAllianceProposal(
				title,
				description,
				args[0],
			)

			err = content.ValidateBasic()

			if err != nil {
				return err
			}

			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)

			if err != nil {
				return err
			}

			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(govcli.FlagTitle, "", "title of proposal")
	cmd.Flags().String(govcli.FlagDescription, "", "description of proposal")
	cmd.Flags().String(govcli.FlagDeposit, "", "deposit of proposal")
	return cmd
}
