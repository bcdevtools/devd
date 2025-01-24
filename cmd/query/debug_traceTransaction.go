package query

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/bcdevtools/devd/v3/cmd/flags"

	"github.com/bcdevtools/devd/v3/cmd/types"
	"github.com/bcdevtools/devd/v3/cmd/utils"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/spf13/cobra"
)

const (
	flagTracer      = "tracer"
	flagNoTranslate = "no-translate"
)

func GetQueryTraceTxCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "debug_traceTransaction [0xhash]",
		Aliases: []string{"trace"},
		Short:   "debug_traceTransaction",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			_, evmRpc := flags.MustGetEthClient(cmd)

			input := strings.ToLower(args[0])

			if !regexp.MustCompile(`^0x[a-f\d]{64}$`).MatchString(input) {
				utils.PrintlnStdErr("ERR: invalid EVM transaction hash format")
				os.Exit(1)
			}

			var params []types.JsonRpcQueryParam

			paramTransactionHash, err := types.NewJsonRpcStringQueryParam(input)
			utils.ExitOnErr(err, "failed to create json rpc query param")
			params = append(params, paramTransactionHash)

			if tracer := cmd.Flag(flagTracer).Value.String(); tracer != "" {
				if !regexp.MustCompile(`^\w+$`).MatchString(tracer) {
					utils.PrintlnStdErr("ERR: invalid tracer name:", tracer)
					os.Exit(1)
				}
				params = append(params, types.NewJsonRpcRawQueryParam(fmt.Sprintf(`{"tracer":"%s"}`, tracer)))
			}

			bz, err := types.DoEvmQuery(
				evmRpc,
				types.NewJsonRpcQueryBuilder(
					"debug_traceTransaction",
					params...,
				),
				0,
			)
			utils.ExitOnErr(err, "failed to trace transaction")

			traceContentAsMap, err := getResultObjectFromEvmRpcResponse(bz)
			if err == nil {
				recursivelyTranslateTraceFrames(traceContentAsMap)

				bz, err = json.Marshal(traceContentAsMap)
				utils.ExitOnErr(err, "failed to marshal response trace tx")

				utils.TryPrintBeautyJson(bz)
			} else {
				utils.TryPrintBeautyJson(bz)
			}

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
									utils.PrintfStdErr(
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

	cmd.Flags().String(flags.FlagEvmRpc, "", flags.FlagEvmRpcDesc)
	cmd.Flags().String(flagTracer, "callTracer", "EVM tracer")
	cmd.Flags().Bool(flagNoTranslate, false, "do not translate and print EVM revert error message")

	return cmd
}

func recursivelyTranslateTraceFrames(_map map[string]interface{}) {
	if _map == nil {
		return
	}

	utils.TryInjectTranslatedFieldForEvmRpcObjects(nil, _map, "gas")
	utils.TryInjectTranslatedFieldForEvmRpcObjects(nil, _map, "gasUsed")
	utils.TryInjectTranslatedFieldForEvmRpcObjects(nil, _map, "value")

	calls, found := _map["calls"]
	if !found {
		return
	}

	callsArray, ok := calls.([]interface{})
	if !ok {
		return
	}

	for _, call := range callsArray {
		callAsMap, ok := call.(map[string]interface{})
		if !ok {
			continue
		}
		recursivelyTranslateTraceFrames(callAsMap)
	}
}
