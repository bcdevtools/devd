package query

import (
	"fmt"
	libutils "github.com/EscanBE/go-lib/utils"
	"github.com/bcdevtools/devd/cmd/types"
	"github.com/spf13/cobra"
	"os"
	"regexp"
	"strings"
)

// GetQueryTraceTxCommand registers a sub-tree of commands
func GetQueryTraceTxCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "debug_traceTransaction [0xhash]",
		Aliases: []string{"trace", "trace_tx"},
		Short:   "debug_traceTransaction",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			_, rpc := mustGetEthClient(cmd)

			input := strings.ToLower(args[0])

			if !regexp.MustCompile(`^0x[a-f\d]{64}$`).MatchString(input) {
				libutils.PrintlnStdErr("ERR: invalid EVM transaction hash format")
				os.Exit(1)
			}

			var params []types.JsonRpcQueryParam

			paramTransactionHash, err := types.NewJsonRpcStringQueryParam(input)
			libutils.ExitIfErr(err, "failed to create json rpc query param")
			params = append(params, paramTransactionHash)

			if tracer := cmd.Flag(flagTracer).Value.String(); tracer != "" {
				if !regexp.MustCompile(`^\w+$`).MatchString(tracer) {
					libutils.PrintlnStdErr("ERR: invalid tracer name:", tracer)
					os.Exit(1)
				}
				params = append(params, types.NewJsonRpcRawQueryParam(fmt.Sprintf(`{"tracer":"%s"}`, tracer)))
			}

			bz, err := doQuery(
				rpc,
				types.NewJsonRpcQueryBuilder(
					"debug_traceTransaction",
					params...,
				),
				0,
			)
			libutils.ExitIfErr(err, "failed to trace transaction")

			tryPrintBeautyJson(bz)
		},
	}

	cmd.Flags().String(flagRpc, "http://localhost:8545", "EVM Json-RPC url")
	cmd.Flags().String(flagTracer, "callTracer", "EVM tracer")

	return cmd
}
