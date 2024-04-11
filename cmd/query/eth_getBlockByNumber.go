package query

import (
	"github.com/bcdevtools/devd/cmd/types"
	"github.com/bcdevtools/devd/cmd/utils"
	"github.com/spf13/cobra"
	"math/big"
	"os"
	"regexp"
	"strings"
)

func GetQueryBlockCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "eth_getBlockByNumber [height dec or 0xHex]",
		Aliases: []string{"block"},
		Short:   "eth_getBlockByNumber",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
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

			bz, err := doQuery(
				rpc,
				types.NewJsonRpcQueryBuilder(
					"eth_getBlockByNumber",
					paramBlockNumber,
					types.NewJsonRpcBoolQueryParam(cmd.Flag(flagFull).Changed),
				),
				0,
			)
			utils.ExitOnErr(err, "failed to get block by number")

			utils.TryPrintBeautyJson(bz)
		},
	}

	cmd.Flags().String(flagRpc, "", flagEvmRpcDesc)
	cmd.Flags().Bool(flagFull, false, "should returns the full transaction objects when this value is true otherwise, it returns only the hashes of the transactions")

	return cmd
}
