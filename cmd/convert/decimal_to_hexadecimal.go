package convert

import (
	"fmt"
	"github.com/bcdevtools/devd/v2/cmd/utils"
	"github.com/spf13/cobra"
	"math/big"
	"os"
	"strings"
)

func GetConvertDecimalToHexadecimalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "dec_2_hex [dec]",
		Aliases: []string{"d2h"},
		Short: `Convert decimal to hexadecimal.
Support pipe.`,
		Run: func(cmd *cobra.Command, args []string) {
			args, err := utils.ProvidedArgsOrFromPipe(args)
			utils.ExitOnErr(err, "failed to get args from pipe")
			utils.RequireExactArgsCount(args, 1, cmd)

			input := strings.ToLower(args[0])

			bi, ok := new(big.Int).SetString(input, 10)
			if !ok {
				utils.PrintlnStdErr("ERR: failed to convert decimal to hexadecimal")
				os.Exit(1)
			}
			fmt.Printf("0x%x\n", bi)
		},
	}

	return cmd
}
