package query

import (
	"github.com/bcdevtools/devd/v2/constants"
	"github.com/spf13/cobra"
)

const (
	flagRpc         = "rpc"
	flagRest        = "rest"
	flagTmRpc       = "tm-rpc"
	flagFull        = "full"
	flagTracer      = "tracer"
	flagHeight      = "height"
	flagNoTranslate = "no-translate"
	flagErc20       = "erc20"
)

const (
	flagEvmRpcDesc     = "EVM Json-RPC endpoint, default is " + constants.DEFAULT_EVM_RPC + ", can be set by environment variable " + constants.ENV_EVM_RPC
	flagCosmosRestDesc = "Cosmos Rest API endpoint, default is " + constants.DEFAULT_COSMOS_REST + ", can be set by environment variable " + constants.ENV_COSMOS_REST
	flagTmRpcDesc      = "Tendermint RPC endpoint, default is " + constants.DEFAULT_TM_RPC + ", can be set by environment variable " + constants.ENV_TM_RPC
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
		GetQueryTxsInBlockCommand(),
		GetQueryTxEventsCommand(),
	)

	return cmd
}
