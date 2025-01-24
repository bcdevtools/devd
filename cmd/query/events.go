package query

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/bcdevtools/devd/v3/cmd/flags"
	"regexp"
	"strconv"
	"strings"

	"github.com/bcdevtools/devd/v3/cmd/utils"
	acbitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
)

const (
	flagFilter = "filter"
)

func GetQueryTxEventsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "events [height/tx hash]",
		Short: "Query block/tx events",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			txHashType := utils.DetectTxHashType(args[0])
			tendermintRpcHttpClient, _ := flags.MustGetTmRpc(cmd)

			var events []acbitypes.Event

			switch txHashType {
			case utils.TxHashTypeEvm:
				txHash := args[0]

				resTx, err := tendermintRpcHttpClient.TxSearch(context.Background(), "ethereum_tx.ethereumTxHash='"+txHash+"'", false, nil, nil, "")
				utils.ExitOnErr(err, "failed to find Cosmos tx by EVM tx hash")
				if len(resTx.Txs) == 0 {
					utils.PrintlnStdErr("ERR: no Cosmos tx found by EVM tx hash")
					return
				}

				events = resTx.Txs[0].TxResult.Events
			case utils.TxHashTypeCosmos:
				txHash := args[0]

				txHashBz, err := hex.DecodeString(txHash)
				utils.ExitOnErr(err, "failed to decode tx hash")

				resultTx, err := tendermintRpcHttpClient.Tx(context.Background(), txHashBz, false)
				utils.ExitOnErr(err, "failed to get tx")

				events = resultTx.TxResult.Events
			default:
				if !regexp.MustCompile(`^\d+$`).MatchString(args[0]) {
					utils.PrintlnStdErr("ERR: input is neither a height nor a tx hash")
					return
				}

				height, err := strconv.ParseInt(args[0], 10, 64)
				utils.ExitOnErr(err, "failed to parse height")

				resultBlockResult, err := tendermintRpcHttpClient.BlockResults(context.Background(), &height)
				utils.ExitOnErr(err, "failed to get block results")

				events = append(events, resultBlockResult.BeginBlockEvents...)
				for _, txResult := range resultBlockResult.TxsResults {
					events = append(events, txResult.Events...)
				}
				events = append(events, resultBlockResult.EndBlockEvents...)
			}

			events = utils.ResolveBase64Events(events)

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

	cmd.Flags().String(flags.FlagTendermintRpc, "", flags.FlagTmRpcDesc)
	cmd.Flags().StringSliceP(flagFilter, "f", []string{}, "filter events, output only events which contains the filter string. If multiple filters are provided, events that contain one of the filters will be output.")

	return cmd
}
