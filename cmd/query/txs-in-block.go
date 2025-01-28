package query

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/bcdevtools/devd/v3/cmd/flags"

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

			tendermintRpcHttpClient, _ := flags.MustGetTmRpc(cmd)

			resBlock, err := tendermintRpcHttpClient.Block(context.Background(), &blockNumber)
			utils.ExitOnErr(err, "failed to get block")

			resBlockResult, err := tendermintRpcHttpClient.BlockResults(context.Background(), &blockNumber)
			utils.ExitOnErr(err, "failed to get block results")

			orTextEmpty := func(input any) any {
				if input == nil {
					return "(empty)"
				}
				if str, ok := input.(string); ok && str == "" {
					return "(empty)"
				}
				return input
			}

			tryExtractMsgTypeFromRawTx := func(input []byte) []string {
				matches := regexp.MustCompile(`/[a-z\d]+(\.[a-z\d]+)+\.Msg[a-zA-Z\d]+`).FindAll(input, -1)
				var result []string
				for _, match := range matches {
					result = append(result, string(match))
				}
				return result
			}

			for i, tx := range resBlock.Block.Txs {
				if i > 0 {
					utils.PrintlnStdErr("====================================")
				}
				txResult := resBlockResult.TxsResults[i]
				txResult.Events = utils.ResolveBase64Events(txResult.Events)
				fmt.Println("Index:", i)
				fmt.Println("Hash:", strings.ToUpper(hex.EncodeToString(tx.Hash())))
				if strings.Contains(string(tx), "/ethermint.evm.v1.MsgEthereumTx") {
				L1:
					for _, event := range txResult.Events {
						if event.Type != "ethereum_tx" {
							continue
						}

						for _, attr := range event.Attributes {
							if string(attr.Key) == "ethereumTxHash" {
								fmt.Println("EvmHash:", attr.Value)
								break L1
							}
						}
					}
				}
				fmt.Println("Type:", func() string {
					if msgsType := tryExtractMsgTypeFromRawTx(tx); len(msgsType) > 0 {
						return strings.Join(msgsType, ", ")
					}
					return "(failed to extract msg type)"
				}())
				fmt.Println("Code:", txResult.Code)
				fmt.Println("Gas:", txResult.GasUsed, "/", txResult.GasWanted)
				fmt.Println("Data:", strings.ReplaceAll(string(txResult.Data), "\n", ""))
				fmt.Println("Log:", orTextEmpty(txResult.Log))
				fmt.Println("Events:", func() string {
					if len(txResult.Events) == 0 {
						return orTextEmpty("").(string)
					}
					bz, err := json.Marshal(txResult.Events)
					utils.ExitOnErr(err, "failed to marshal events")
					return string(bz)
				}())
				fmt.Println("Info:", orTextEmpty(txResult.Info))
				fmt.Println("Code-space:", orTextEmpty(txResult.Codespace))
			}
		},
	}

	cmd.Flags().String(flags.FlagTendermintRpc, "", flags.FlagTmRpcDesc)

	return cmd
}
