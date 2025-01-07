package query

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/bcdevtools/devd/v2/cmd/utils"
	acbitypes "github.com/cometbft/cometbft/abci/types"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
	"strings"
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

			events := utils.ResolveBase64Events(resultTx.TxResult.Events)

			filters, _ := cmd.Flags().GetStringSlice(flagFilter)
			filters = func() []string {
				uniqueFilters := make(map[string]struct{}, len(filters))
				for _, filter := range filters {
					if filter == "" {
						continue
					}
					uniqueFilters[filter] = struct{}{}
				}
				return maps.Keys(uniqueFilters)
			}()
			containsFilterPattern := func(str string) bool {
				for _, filter := range filters {
					if strings.Contains(str, filter) {
						return true
					}
				}
				return false
			}
			if len(filters) > 0 {
				var filteredEvents []acbitypes.Event
				for _, event := range events {
					var contains bool
					if containsFilterPattern(event.Type) {
						contains = true
					} else {
						for _, attr := range event.Attributes {
							if containsFilterPattern(attr.Key) || containsFilterPattern(attr.Value) {
								contains = true
								break
							}
						}
					}

					if contains {
						filteredEvents = append(filteredEvents, event)
					}
				}

				if len(filteredEvents) < 1 {
					utils.PrintlnStdErr("ERR: no events found after filtered")
					return
				}

				events = filteredEvents
			}

			bz, err := json.MarshalIndent(events, "", "  ")
			utils.ExitOnErr(err, "failed to marshal events to json")
			fmt.Println(string(bz))
		},
	}

	cmd.Flags().String(flagTmRpc, "", flagTmRpcDesc)
	cmd.Flags().StringSlice(flagFilter, []string{}, "filter events, output only events which contains the filter string. If multiple filters are provided, events that contain one of the filters will be output.")

	return cmd
}
