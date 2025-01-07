package query

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/bcdevtools/devd/v2/cmd/utils"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/spf13/cobra"
)

func GetQueryTxEventsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "events [hash]",
		Short: "Query tx events by hash",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			txHash := args[0]

			txHashType := utils.DetectTxHashType(txHash)
			if txHashType == utils.TxHashTypeInvalid {
				utils.PrintlnStdErr("ERR: invalid tx hash format")
				return
			}

			isEvmTxHash := txHashType == utils.TxHashTypeEvm

			tendermintRpcHttpClient, _ := mustGetTmRpc(cmd)

			var resultTx *coretypes.ResultTx
			if isEvmTxHash {
				resTx, err := tendermintRpcHttpClient.TxSearch(context.Background(), "ethereum_tx.ethereumTxHash='"+txHash+"'", false, nil, nil, "")
				utils.ExitOnErr(err, "failed to find Cosmos tx by EVM tx hash")
				if len(resTx.Txs) == 0 {
					utils.PrintlnStdErr("ERR: no Cosmos tx found by EVM tx hash")
					return
				}

				resultTx = resTx.Txs[0]
			} else {
				txHashBz, err := hex.DecodeString(txHash)
				utils.ExitOnErr(err, "failed to decode tx hash")

				resultTx, err = tendermintRpcHttpClient.Tx(context.Background(), txHashBz, false)
				utils.ExitOnErr(err, "failed to get tx")
			}

			bz, err := json.MarshalIndent(utils.ResolveBase64Events(resultTx.TxResult.Events), "", "  ")
			utils.ExitOnErr(err, "failed to marshal events to json")
			fmt.Println(string(bz))
		},
	}

	cmd.Flags().String(flagTmRpc, "", flagTmRpcDesc)

	return cmd
}
