package query

import (
	"encoding/hex"
	"fmt"
	libutils "github.com/EscanBE/go-lib/utils"
	"github.com/bcdevtools/devd/cmd/types"
	"github.com/bcdevtools/devd/cmd/utils"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/spf13/cobra"
	"os"
	"regexp"
	"strings"
)

func GetQueryTraceTxCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "debug_traceTransaction [0xhash]",
		Aliases: []string{"trace", "trace_tx"},
		Short:   "debug_traceTransaction",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			_, rpc := mustGetEthClient(cmd, false)

			input := strings.ToLower(args[0])

			if !regexp.MustCompile(`^0x[a-f\d]{64}$`).MatchString(input) {
				libutils.PrintlnStdErr("ERR: invalid EVM transaction hash format")
				os.Exit(1)
			}

			var params []types.JsonRpcQueryParam

			paramTransactionHash, err := types.NewJsonRpcStringQueryParam(input)
			utils.ExitOnErr(err, "failed to create json rpc query param")
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
			utils.ExitOnErr(err, "failed to trace transaction")

			utils.TryPrintBeautyJson(bz)

			if !cmd.Flag(flagNoTranslate).Changed {
				// try to decode error message if any

				type structOfError struct {
					Error  string `json:"error"`
					Output string `json:"output"`
				}

				res, err := types.ParseJsonRpcResponse[structOfError](bz)
				if err == nil {
					if resStruct, ok := res.(*structOfError); ok {
						if resStruct.Error != "" && resStruct.Output != "" && resStruct.Error == vm.ErrExecutionReverted.Error() {
							if regexp.MustCompile(`^0x08c379a0[a-f\d]{192,}$`).MatchString(resStruct.Output) && (len(resStruct.Output)-10 /*exclude sig of error*/)%64 == 0 {
								bz, err := hex.DecodeString(resStruct.Output[10:])
								errMsg, err := utils.AbiDecodeString(bz)
								if err == nil && errMsg != "" {
									libutils.PrintfStdErr(
										"ERR: EVM execution reverted with message [%s], this translation can be omitted by providing flag '--%s'\n",
										errMsg,
										flagNoTranslate,
									)
								}
							}
						}
					}
				}
			}
		},
	}

	cmd.Flags().String(flagRpc, "", flagEvmRpcDesc)
	cmd.Flags().String(flagTracer, "callTracer", "EVM tracer")
	cmd.Flags().Bool(flagNoTranslate, false, "do not translate and print EVM revert error message")

	return cmd
}
