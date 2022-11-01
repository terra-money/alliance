package cli

import (
	"alliance/x/alliance/types"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
)

func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Alliance module subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(NewDelegateCmd(), NewRedelegateCmd(), NewUndelegateCmd(), NewClaimDelegationRewardsCmd())
	return txCmd
}

func NewDelegateCmd() *cobra.Command {
	bech32PrefixValAddr := sdk.GetConfig().GetBech32ValidatorAddrPrefix()

	cmd := &cobra.Command{
		Use:   "delegate validator-addr amount",
		Args:  cobra.ExactArgs(2),
		Short: "Delegate alliance enabled tokens to a validator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Delegate an amount of liquid alliance enabled coins to a validator from your wallet.

Example:
$ %s tx alliance delegate %s1l2rsakp388kuv9k8qzq6lrm9taddae7fpx59wm 1000stake --from mykey
`,
				version.AppName, bech32PrefixValAddr,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			amount, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			delAddr := clientCtx.GetFromAddress()
			valAddr, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			msg := &types.MsgDelegate{
				DelegatorAddress: delAddr.String(),
				ValidatorAddress: valAddr.String(),
				Amount:           amount,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewRedelegateCmd() *cobra.Command {
	bech32PrefixValAddr := sdk.GetConfig().GetBech32ValidatorAddrPrefix()

	cmd := &cobra.Command{
		Use:   "redelegate src-validator-addr dst-validator-addr amount",
		Args:  cobra.ExactArgs(3),
		Short: "Re-delegate alliance enabled tokens from a validator to another",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Re-delegate an amount of liquid alliance enabled coins from a validator to another from your wallet.

Example:
$ %s tx alliance redelegate %s1l2rsakp388kuv9k8qzq6lrm9taddae7fpx59wm %ss1l2rsakp388kuv9k8qzq6lrm9taddae7fpx59wm 1000stake --from mykey
`,
				version.AppName, bech32PrefixValAddr, bech32PrefixValAddr,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			delAddr := clientCtx.GetFromAddress()
			srcValAddr, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			dstValAddr, err := sdk.ValAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoinNormalized(args[2])
			if err != nil {
				return err
			}

			msg := &types.MsgRedelegate{
				DelegatorAddress:    delAddr.String(),
				ValidatorSrcAddress: srcValAddr.String(),
				ValidatorDstAddress: dstValAddr.String(),
				Amount:              amount,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewUndelegateCmd() *cobra.Command {
	bech32PrefixValAddr := sdk.GetConfig().GetBech32ValidatorAddrPrefix()

	cmd := &cobra.Command{
		Use:   "undelegate validator-addr amount",
		Args:  cobra.ExactArgs(2),
		Short: "Undelegate alliance enabled tokens to a validator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Undelegate an amount of liquid alliance enabled coins from a validator to your wallet (after the unbonding period has finished).

Example:
$ %s tx alliance undelegate %s1l2rsakp388kuv9k8qzq6lrm9taddae7fpx59wm 1000stake --from mykey
`,
				version.AppName, bech32PrefixValAddr,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			amount, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			delAddr := clientCtx.GetFromAddress()
			valAddr, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			msg := &types.MsgUndelegate{
				DelegatorAddress: delAddr.String(),
				ValidatorAddress: valAddr.String(),
				Amount:           amount,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewClaimDelegationRewardsCmd() *cobra.Command {
	bech32PrefixValAddr := sdk.GetConfig().GetBech32ValidatorAddrPrefix()
	cmd := &cobra.Command{
		Use:   "claim-rewards validator-addr denom",
		Args:  cobra.ExactArgs(2),
		Short: "Claim rewards from a delegation by specifying a validator address and denom",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Claim all rewards from a delegation
Example:
$ %s tx alliance claim-rewards %s1l2rsakp388kuv9k8qzq6lrm9taddae7fpx59wm stake --from mykey
`,
				version.AppName, bech32PrefixValAddr,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			delAddr := clientCtx.GetFromAddress()
			valAddr, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			denom := args[1]
			if denom == "" {
				return fmt.Errorf("denom is required")
			}
			msg := &types.MsgClaimDelegationRewards{
				DelegatorAddress: delAddr.String(),
				ValidatorAddress: valAddr.String(),
				Denom:            denom,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
