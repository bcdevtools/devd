package convert

import (
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/bcdevtools/devd/v3/cmd/utils"
	"github.com/spf13/cobra"
)

func GetConvertDecimalToHexadecimalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "dec_2_hex [dec]",
		Aliases: []string{"d2h"},
		Short:   "Convert decimal to hexadecimal.",
		Long: `Convert decimal to hexadecimal.
Support pipe.`,
		Run: func(cmd *cobra.Command, args []string) {
			args, err := utils.ProvidedArgsOrFromPipe(args)
			utils.ExitOnErr(err, "failed to get args from pipe")
			utils.RequireExactArgsCount(args, 1, cmd)

			input := strings.ToLower(args[0])

			bi, err := utils.ReadCustomInteger(input)
			if err != nil {
				var ok bool
				bi, ok = new(big.Int).SetString(input, 10)
				if !ok {
					utils.PrintlnStdErr("ERR: failed to convert decimal to hexadecimal")
					os.Exit(1)
				}
			}

			fmt.Printf("0x%x\n", bi)
		},
	}

	return cmd
}
