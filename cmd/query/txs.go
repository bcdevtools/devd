package query

import (
	"context"
	"fmt"
	"strconv"

	"github.com/bcdevtools/devd/v3/cmd/utils"
	"github.com/spf13/cobra"
)

func GetQueryTxsInBlockCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "txs-in-block [block-number]",
		Short: "Query txs in a block",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			blockNumber, _ := strconv.ParseInt(args[0], 10, 64)
			if blockNumber < 1 {
				utils.PrintlnStdErr("ERR: invalid block number")
				return
			}

			tendermintRpcHttpClient, _ := mustGetTmRpc(cmd)

			resBlock, err := tendermintRpcHttpClient.Block(context.Background(), &blockNumber)
			utils.ExitOnErr(err, "failed to get block")

			for _, tx := range resBlock.Block.Txs {
				fmt.Println(tx.Hash())
			}

			resBlockResult, err := tendermintRpcHttpClient.BlockResults(context.Background(), &blockNumber)
			utils.ExitOnErr(err, "failed to get block results")

			for _, txResult := range resBlockResult.TxsResults {
				fmt.Println(txResult.Log)
			}
		},
	}

	cmd.Flags().String(flagTmRpc, "", flagTmRpcDesc)

	return cmd
}
