package query

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"regexp"
	"strings"

	"github.com/bcdevtools/devd/v3/cmd/flags"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/bcdevtools/devd/v3/cmd/types"
	"github.com/bcdevtools/devd/v3/cmd/utils"
	"github.com/spf13/cobra"
)

const (
	flagFull = "full"
)

func GetQueryEvmRpcEthGetBlockByNumberCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "eth_getBlockByNumber [height dec or 0xHex]",
		Aliases: []string{"evm-block"},
		Short:   "Query `eth_getBlockByNumber` from EVM RPC",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			_, evmRpc := flags.MustGetEthClient(cmd)

			input := strings.ToLower(args[0])

			if regexp.MustCompile(`[a-f]`).MatchString(input) && !strings.HasPrefix(input, "0x") {
				utils.PrintlnStdErr("ERR: hexadecimal block number must have 0x prefix.")
				os.Exit(1)
			}

			var blockNumber *big.Int
			var ok bool
			if strings.HasPrefix(input, "0x") {
				blockNumber, ok = new(big.Int).SetString(input[2:], 16)
				if !ok {
					utils.PrintlnStdErr("ERR: invalid EVM hexadecimal block number")
					os.Exit(1)
				}
			} else {
				blockNumber, ok = new(big.Int).SetString(input, 10)
				if !ok {
					utils.PrintlnStdErr("ERR: invalid EVM decimal block number")
					os.Exit(1)
				}
			}

			if blockNumber.Sign() == 0 {
				blockNumber = nil
			}

			var paramBlockNumber types.JsonRpcQueryParam
			var err error

			if blockNumber == nil || blockNumber.Sign() == 0 {
				paramBlockNumber, err = types.NewJsonRpcStringQueryParam("latest")
				utils.ExitOnErr(err, "failed to create json rpc query param")
			} else {
				paramBlockNumber = types.NewJsonRpcInt64QueryParam(blockNumber.Int64())
			}

			bz, err := types.DoEvmRpcQuery(
				evmRpc,
				types.NewJsonRpcQueryBuilder(
					"eth_getBlockByNumber",
					paramBlockNumber,
					types.NewJsonRpcBoolQueryParam(cmd.Flag(flagFull).Changed),
				),
				0,
			)
			utils.ExitOnErr(err, "failed to get block by number")

			blockInfoAsMap, err := getResultObjectFromEvmRpcResponse(bz)
			utils.ExitOnErr(err, "failed to get result object from response")

			utils.TryInjectTranslatedFieldForEvmRpcObjects(&ethtypes.Block{}, blockInfoAsMap, "baseFeePerGas")
			utils.TryInjectTranslatedFieldForEvmRpcObjects(&ethtypes.Block{}, blockInfoAsMap, "gasLimit")
			utils.TryInjectTranslatedFieldForEvmRpcObjects(&ethtypes.Block{}, blockInfoAsMap, "gasUsed")
			utils.TryInjectTranslatedFieldForEvmRpcObjects(&ethtypes.Block{}, blockInfoAsMap, "number")
			utils.TryInjectTranslatedFieldForEvmRpcObjects(&ethtypes.Block{}, blockInfoAsMap, "size")
			utils.TryInjectTranslatedFieldForEvmRpcObjects(&ethtypes.Block{}, blockInfoAsMap, "timestamp")

			bz, err = json.Marshal(blockInfoAsMap)
			utils.ExitOnErr(err, "failed to marshal response block by number")

			utils.TryPrintBeautyJson(bz)
		},
	}

	cmd.Flags().String(flags.FlagEvmRpc, "", flags.FlagEvmRpcDesc)
	cmd.Flags().Bool(flagFull, false, "should returns the full transaction objects when this value is true otherwise, it returns only the hashes of the transactions")

	return cmd
}

func getResultObjectFromEvmRpcResponse(bz []byte) (map[string]interface{}, error) {
	var _map map[string]interface{}
	err := json.Unmarshal(bz, &_map)
	if err != nil {
		return nil, err
	}
	errObj, found := _map["error"]
	if found && errObj != nil {
		errMap, ok := errObj.(map[string]interface{})
		if ok {
			code, foundCode := errMap["code"]
			msg, foundMsg := errMap["message"]
			if foundCode && foundMsg {
				return nil, fmt.Errorf("RPC response error: code=%v, message=%v", code, msg)
			}
		}
		return nil, fmt.Errorf("RPC response error: %v", errObj)
	}
	result, found := _map["result"]
	if !found {
		return nil, nil
	}
	resultAsMap, ok := result.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to convert into map")
	}
	return resultAsMap, nil
}
