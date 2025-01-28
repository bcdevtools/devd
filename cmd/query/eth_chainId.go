package query

import (
	"context"
	"fmt"

	"github.com/bcdevtools/devd/v3/cmd/flags"
	"github.com/bcdevtools/devd/v3/cmd/utils"
	"github.com/spf13/cobra"
)

func GetQueryEvmRpcEthChainIdCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "eth_chainId",
		Short: "Query `eth_chainId` from EVM RPC",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			ethClient, _ := flags.MustGetEthClient(cmd)

			chainId, err := ethClient.ChainID(context.Background())
			utils.ExitOnErr(err, "failed to get chain id")

			fmt.Println(chainId.String())
		},
	}

	cmd.Flags().String(flags.FlagEvmRpc, "", flags.FlagEvmRpcDesc)

	return cmd
}
