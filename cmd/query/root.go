package query

import (
	"math/big"

	"github.com/spf13/cobra"
)

const (
	flagHeight = "height"
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
		GetQueryEvmRpcDebugTraceTransactionCommand(),
		// fake command for deprecated alias
		GetDeprecatedAliasBlockAsCommand(),
		GetDeprecatedAliasTxAsCommand(),
	)

	return cmd
}

func readContextHeightFromFlag(cmd *cobra.Command) *big.Int {
	height, _ := cmd.Flags().GetInt64(flagHeight)
	if height > 0 {
		return big.NewInt(height)
	}

	return nil
}
