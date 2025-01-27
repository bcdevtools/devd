package query

import (
	"github.com/bcdevtools/devd/v3/cmd/utils"
	"github.com/bcdevtools/devd/v3/constants"
	"github.com/spf13/cobra"
	"os"
)

func GetDeprecatedAliasBlockAsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "block [height]",
		Hidden: true,
		Short:  "Deprecated and removed alias of `eth_getBlockByNumber`",
		Args:   cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			utils.PrintfStdErr("ERR: `block` is a deprecated alias of `eth_getBlockByNumber`. Please use `%s q eth_getBlockByNumber` instead.\n", constants.BINARY_NAME)
			os.Exit(1)
		},
	}

	return cmd
}

func GetDeprecatedAliasTxAsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "tx [hash]",
		Hidden: true,
		Short:  "Deprecated and removed alias of `eth_getTransactionByHash`",
		Args:   cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			utils.PrintfStdErr("ERR: `tx` is a deprecated alias of `eth_getTransactionByHash`. Please use `%s q eth_getTransactionByHash` instead.\n", constants.BINARY_NAME)
			os.Exit(1)
		},
	}

	return cmd
}
