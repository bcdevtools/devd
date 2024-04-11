package query

import (
	"github.com/bcdevtools/devd/v2/constants"
	"github.com/spf13/cobra"
)

const (
	flagRpc         = "rpc"
	flagFull        = "full"
	flagTracer      = "tracer"
	flagHeight      = "height"
	flagNoTranslate = "no-translate"
)

const (
	flagEvmRpcDesc = "EVM Json-RPC endpoint, default is " + constants.DEFAULT_EVM_RPC + ", can be set by environment variable " + constants.ENV_EVM_RPC
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
		GetQueryTxCommand(),
		GetQueryTxReceiptCommand(),
		GetQueryBlockCommand(),
		GetQueryTraceTxCommand(),
		GetQueryBalanceCommand(),
	)

	return cmd
}
