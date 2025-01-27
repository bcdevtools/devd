package query

import (
	"github.com/spf13/cobra"
)

// Commands registers a sub-tree of commands
func Commands() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "Query commands",
	}

	cmd.AddCommand(
		GetQueryErc20Command(),
		GetQueryBalanceCommand(),
		GetQueryTxsInBlockCommand(),
		GetQueryTxEventsCommand(),
		GetQueryEvmRpcEthGetTransactionByHashCommand(),
		GetQueryEvmRpcEthGetTransactionReceiptCommand(),
		GetQueryEvmRpcEthGetBlockByNumberCommand(),
		GetQueryEvmRpcEthChainIdCommand(),
		GetQueryEvmRpcEthCallCommand(),
		GetQueryEvmRpcDebugTraceTransactionCommand(),
		// fake command for deprecated alias
		GetDeprecatedAliasBlockAsCommand(),
		GetDeprecatedAliasTxAsCommand(),
	)

	return cmd
}
