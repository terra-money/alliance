package cli

import (
	"cosmossdk.io/math"
	"time"

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
		Use:   "create-alliance denom reward-weight reward-weight-min reward-weight-max take-rate reward-change-rate reward-change-interval",
		Args:  cobra.ExactArgs(7),
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

			description, err := cmd.Flags().GetString(govcli.FlagDescription) //nolint:staticcheck // SA1019: govcli.FlagDescription is deprecated
			if err != nil {
				return err
			}

			denom := args[0]

			rewardWeight, err := math.LegacyNewDecFromStr(args[1])
			if err != nil {
				return err
			}

			rewardWeightMin, err := math.LegacyNewDecFromStr(args[2])
			if err != nil {
				return err
			}

			rewardWeightMax, err := math.LegacyNewDecFromStr(args[3])
			if err != nil {
				return err
			}

			takeRate, err := math.LegacyNewDecFromStr(args[4])
			if err != nil {
				return err
			}

			rewardChangeRate, err := math.LegacyNewDecFromStr(args[5])
			if err != nil {
				return err
			}

			rewardChangeInterval, err := time.ParseDuration(args[6])
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
				denom,
				rewardWeight,
				types.RewardWeightRange{
					Min: rewardWeightMin,
					Max: rewardWeightMax,
				},
				takeRate,
				rewardChangeRate,
				rewardChangeInterval,
			)

			err = content.ValidateBasic()

			if err != nil {
				return err
			}

			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(govcli.FlagTitle, "", "title of proposal")
	cmd.Flags().String(govcli.FlagDescription, "", "description of proposal") //nolint:staticcheck
	cmd.Flags().String(govcli.FlagDeposit, "", "deposit of proposal")
	return cmd
}

func UpdateAlliance() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-alliance denom reward-weight reward-weight-min reward-weight-max take-rate reward-change-rate reward-change-interval",
		Args:  cobra.ExactArgs(5),
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

			description, err := cmd.Flags().GetString(govcli.FlagDescription) //nolint:staticcheck
			if err != nil {
				return err
			}

			denom := args[0]

			rewardWeight, err := math.LegacyNewDecFromStr(args[1])
			if err != nil {
				return err
			}

			rewardWeightMin, err := math.LegacyNewDecFromStr(args[2])
			if err != nil {
				return err
			}

			rewardWeightMax, err := math.LegacyNewDecFromStr(args[3])
			if err != nil {
				return err
			}

			takeRate, err := math.LegacyNewDecFromStr(args[2])
			if err != nil {
				return err
			}

			rewardChangeRate, err := math.LegacyNewDecFromStr(args[3])
			if err != nil {
				return err
			}

			rewardChangeInterval, err := time.ParseDuration(args[4])
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

			content := types.NewMsgUpdateAllianceProposal(
				title,
				description,
				denom,
				rewardWeight,
				types.RewardWeightRange{
					Min: rewardWeightMin,
					Max: rewardWeightMax,
				},
				takeRate,
				rewardChangeRate,
				rewardChangeInterval,
			)

			err = content.ValidateBasic()

			if err != nil {
				return err
			}

			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(govcli.FlagTitle, "", "title of proposal")
	cmd.Flags().String(govcli.FlagDescription, "", "description of proposal") //nolint:staticcheck // SA1019: govcli.FlagDescription is deprecated
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

			description, err := cmd.Flags().GetString(govcli.FlagDescription) //nolint:staticcheck // SA1019: govcli.FlagDescription is deprecated
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			denom := args[0]

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
				denom,
			)

			err = content.ValidateBasic()

			if err != nil {
				return err
			}

			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(govcli.FlagTitle, "", "title of proposal")
	cmd.Flags().String(govcli.FlagDescription, "", "description of proposal") //nolint:staticcheck // SA1019: govcli.FlagDescription is deprecated: use FlagDescription instead
	cmd.Flags().String(govcli.FlagDeposit, "", "deposit of proposal")
	return cmd
}
