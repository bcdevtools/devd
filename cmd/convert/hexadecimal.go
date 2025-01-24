package convert

import (
	"fmt"
	"math/big"
	"os"
	"regexp"
	"strings"

	"github.com/bcdevtools/devd/v3/cmd/utils"
	"github.com/spf13/cobra"
)

func GetConvertHexadecimalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "hexadecimal [0xHex or Dec]",
		Aliases: []string{"hex"},
		Short:   "Convert hexadecimal <> decimal depends on input.",
		Long: `Convert hexadecimal <> decimal depends on input.
Input hexadecimal must be prefixed with 0x.
Support pipe.`,
		Run: func(cmd *cobra.Command, args []string) {
			args, err := utils.ProvidedArgsOrFromPipe(args)
			utils.ExitOnErr(err, "failed to get args from pipe")
			utils.RequireExactArgsCount(args, 1, cmd)

			input := strings.ToLower(args[0])

			if regexp.MustCompile(`^0x[a-f\d]+$`).MatchString(input) { // input is hexadecimal with 0x prefix
				input = strings.TrimPrefix(input, "0x")

				bi, ok := new(big.Int).SetString(input, 16)
				if !ok {
					utils.PrintlnStdErr("ERR: failed to convert hexadecimal to decimal")
					os.Exit(1)
				}

				utils.PrintlnStdErr("INF: converting hexadecimal to decimal")
				fmt.Println(bi)
				return
			}

			if regexp.MustCompile(`^[a-f\d]+$`).MatchString(input) { // input is hexadecimal without 0x prefix or decimal
				if regexp.MustCompile(`[a-f]`).MatchString(input) {
					utils.PrintlnStdErr("ERR: hexadecimal must have 0x prefix")
					os.Exit(1)
				}

				utils.PrintlnStdErr("INF: converting decimal to hexadecimal")
				bi, ok := new(big.Int).SetString(input, 10)
				if !ok {
					utils.PrintlnStdErr("ERR: failed to convert decimal to hexadecimal")
					os.Exit(1)
				}

				fmt.Println("0x" + bi.Text(16))
				return
			}

			utils.PrintlnStdErr("ERR: unrecognized hexadecimal or decimal")
			os.Exit(1)
		},
	}

	return cmd
}
