package convert

import (
	"fmt"
	"github.com/bcdevtools/devd/v2/constants"
	"math/big"
	"os"
	"regexp"
	"strings"

	"github.com/bcdevtools/devd/v2/cmd/utils"
	"github.com/spf13/cobra"
)

func GetConvertHexadecimalToDecimalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "hex_2_dec [hex]",
		Aliases: []string{"h2d"},
		Short:   "Convert hexadecimal to decimal.",
		Long: `Convert hexadecimal to decimal.
Support pipe.`,
		Run: func(cmd *cobra.Command, args []string) {
			utils.PrintfStdErr("WARN: from v3, this command will be renamed to `%s convert hex [0xHex]`\n", constants.BINARY_NAME)

			args, err := utils.ProvidedArgsOrFromPipe(args)
			utils.ExitOnErr(err, "failed to get args from pipe")
			utils.RequireExactArgsCount(args, 1, cmd)

			input := strings.ToLower(args[0])

			if regexp.MustCompile(`^0x[a-f\d]+$`).MatchString(input) {
				// hex with 0x

				bi, ok := new(big.Int).SetString(input[2:], 16)
				if !ok {
					utils.PrintlnStdErr("ERR: failed to convert hexadecimal to decimal")
					os.Exit(1)
				}

				fmt.Println(bi)
				return
			}

			if regexp.MustCompile(`^[a-f\d]+$`).MatchString(input) {
				// hex without 0x

				bi, ok := new(big.Int).SetString(input, 16)
				if !ok {
					utils.PrintlnStdErr("ERR: failed to convert hexadecimal to decimal")
					os.Exit(1)
				}

				fmt.Println(bi)
				return
			}

			utils.PrintlnStdErr("ERR: unrecognized hexadecimal")
			os.Exit(1)
		},
	}

	return cmd
}
