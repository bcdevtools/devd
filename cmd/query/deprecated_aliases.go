package query

import (
	"os"

	"github.com/bcdevtools/devd/v3/cmd/utils"
	"github.com/bcdevtools/devd/v3/constants"
	"github.com/spf13/cobra"
)

func GetDeprecatedAliasBlockAsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "block [height]",
		Hidden: true,
		Short:  "Deprecated and removed alias of `eth_getBlockByNumber`, renamed to `evm-block`",
		Args:   cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			utils.PrintfStdErr("ERR: `block` is a deprecated alias of `eth_getBlockByNumber` and replaced by `evm-block`. Please use `%s q evm-block/eth_getBlockByNumber` instead.\n", constants.BINARY_NAME)
			os.Exit(1)
		},
	}

	return cmd
}

func GetDeprecatedAliasTxAsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "tx [hash]",
		Hidden: true,
		Short:  "Deprecated and removed alias of `eth_getTransactionByHash`, renamed to `evm-tx`",
		Args:   cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			utils.PrintfStdErr("ERR: `tx` is a deprecated alias of `eth_getTransactionByHash` and replaced by `evm-tx`. Please use `%s q evm-tx/eth_getTransactionByHash` instead.\n", constants.BINARY_NAME)
			os.Exit(1)
		},
	}

	return cmd
}

func GetDeprecatedAliasTraceAsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "trace [hash]",
		Hidden: true,
		Short:  "Deprecated and removed alias of `debug_traceTransaction`, renamed to `evm-trace`",
		Args:   cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			utils.PrintfStdErr("ERR: `trace` is a deprecated alias of `debug_traceTransaction` and replaced by `evm-trace`. Please use `%s q evm-trace/debug_traceTransaction` instead.\n", constants.BINARY_NAME)
			os.Exit(1)
		},
	}

	return cmd
}

func GetDeprecatedAliasReceiptAsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "receipt [hash]",
		Hidden: true,
		Short:  "Deprecated and removed alias of `eth_getTransactionReceipt`, renamed to `evm-receipt`",
		Args:   cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			utils.PrintfStdErr("ERR: `receipt` is a deprecated alias of `eth_getTransactionReceipt` and replaced by `evm-receipt`. Please use `%s q evm-receipt/eth_getTransactionReceipt` instead.\n", constants.BINARY_NAME)
			os.Exit(1)
		},
	}

	return cmd
}
