package query

import (
	"encoding/json"
	"fmt"
	"github.com/bcdevtools/devd/v2/constants"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"os"
	"regexp"
	"strings"

	"github.com/bcdevtools/devd/v2/cmd/types"
	"github.com/bcdevtools/devd/v2/cmd/utils"
	"github.com/spf13/cobra"
)

func GetQueryBlockCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "eth_getBlockByNumber [height dec or 0xHex]",
		Aliases: []string{"block", "evm-block"},
		Short:   "eth_getBlockByNumber",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			utils.PrintfStdErr("WARN! Deprecation notice: from v3, command alias `block` will be replaced by `evm-block`, please use `%s q evm-block/eth_getBlockByNumber ...` instead of `%s q block ...`\n", constants.BINARY_NAME, constants.BINARY_NAME)

			_, rpc := mustGetEthClient(cmd, false)

			input := strings.ToLower(args[0])

			if regexp.MustCompile(`[a-f]`).MatchString(input) && !strings.HasPrefix(input, "0x") {
				utils.PrintlnStdErr("Hexadecimal block number must have 0x prefix.")
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

			bz, err := types.DoEvmQuery(
				rpc,
				types.NewJsonRpcQueryBuilder(
					"eth_getBlockByNumber",
					paramBlockNumber,
					types.NewJsonRpcBoolQueryParam(cmd.Flag(flagFull).Changed),
				),
				0,
			)
			utils.ExitOnErr(err, "failed to get block by number")

			blockInfoAsMap, err := getResultObjectFromEvmRpcResponse(bz)
			if err != nil {
				utils.TryPrintBeautyJson(bz)
				return
			}

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

	cmd.Flags().String(flagRpc, "", flagEvmRpcDesc)
	cmd.Flags().Bool(flagFull, false, "should returns the full transaction objects when this value is true otherwise, it returns only the hashes of the transactions")

	return cmd
}

func getResultObjectFromEvmRpcResponse(bz []byte) (map[string]interface{}, error) {
	var _map map[string]interface{}
	err := json.Unmarshal(bz, &_map)
	if err != nil {
		return nil, err
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
