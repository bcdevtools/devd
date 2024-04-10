package debug

import (
	"encoding/hex"
	"fmt"
	libutils "github.com/EscanBE/go-lib/utils"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/spf13/cobra"
	"os"
	"regexp"
	"strings"
)

// GetIntrinsicCommand registers a sub-tree of commands
func GetIntrinsicCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "intrinsic_gas [0xdata]",
		Short: "Get intrinsic gas used by the given EVM transaction input data",
		Args:  cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			input := strings.ToLower(args[0])
			if !regexp.MustCompile(`^(0x)?[a-f\d]+$`).MatchString(input) {
				libutils.PrintlnStdErr("ERR: invalid EVM transaction input data format")
				os.Exit(1)
			}
			if strings.HasPrefix(input, "0x") {
				input = input[2:]
			}
			if len(input)%2 != 0 {
				libutils.PrintlnStdErr("ERR: invalid EVM transaction input string length", len(input), ", must be an even number of characters")
				os.Exit(1)
			}

			bz, err := hex.DecodeString(input)
			libutils.ExitIfErr(err, "failed to decode input hex data")

			var zeroByteCount, nonZeroByteCount int
			for _, b := range bz {
				if b == 0 {
					zeroByteCount++
				} else {
					nonZeroByteCount++
				}
			}
			intrinsicGas := getIntrinsicGasFromInputData(bz)

			fmt.Println("Zero byte count:", zeroByteCount)
			fmt.Println("Non-zero byte count:", nonZeroByteCount)
			fmt.Println("Intrinsic gas:", intrinsicGas)

			recompute := params.TxGas + params.TxDataNonZeroGasEIP2028*uint64(nonZeroByteCount) + params.TxDataZeroGas*uint64(zeroByteCount)
			if recompute == intrinsicGas {
				fmt.Println("=", "tx gas", params.TxGas, "+", "non-zero byte gas", params.TxDataNonZeroGasEIP2028, "x", nonZeroByteCount, "+", "zero byte gas", params.TxDataZeroGas, "x", zeroByteCount)
			}
		},
	}

	return cmd
}

func getIntrinsicGasFromInputData(bz []byte) uint64 {
	intrinsicGas, err := core.IntrinsicGas(bz, ethtypes.AccessList{}, false, true, true)
	libutils.ExitIfErr(err, "failed to get intrinsic gas")
	return intrinsicGas
}
