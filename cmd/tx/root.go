package tx

import (
	"github.com/bcdevtools/devd/v2/constants"
	"github.com/spf13/cobra"
)

const (
	flagRpc       = "rpc"
	flagSecretKey = "secret-key"
)

const (
	flagEvmRpcDesc    = "EVM Json-RPC endpoint, default is " + constants.DEFAULT_EVM_RPC + ", can be set by environment variable " + constants.ENV_EVM_RPC
	flagSecretKeyDesc = "Secret private key or mnemonic of the account, can be set by environment variable " + constants.ENV_SECRET_KEY
)

// Commands registers a sub-tree of commands
func Commands() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tx",
		Short: "Tx commands",
	}

	cmd.AddCommand(
		GetSendEvmTxCommand(),
		GetDeployContractEvmTxCommand(),
		GetDeployErc20EvmTxCommand(),
	)

	return cmd
}
