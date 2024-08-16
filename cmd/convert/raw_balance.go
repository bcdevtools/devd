package convert

import (
	"fmt"
	"os"
	"strconv"

	"github.com/bcdevtools/devd/v2/cmd/utils"
	"github.com/spf13/cobra"
)

const (
	flagCustomDecimalsPoint          = "decimals-point"
	flagCustomDecimalsPointShorthand = "d"
)

// GetRawBalanceCmd creates a helper command that convert display balance into raw balance
func GetRawBalanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "raw_balance [display balance] [decimals]",
		Aliases: []string{"rbal"},
		Short:   "Convert display balance into raw balance.",
		Long: `Convert display balance into raw balance.
Sample: 10.0111 with 6 exponent => 10011100`,
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			displayBalanceStr := args[0]
			decimalsStr := args[1]

			decimals, err := strconv.ParseInt(decimalsStr, 10, 64)
			utils.ExitOnErr(err, "failed to read, decimals is not a number")

			decimalsPoint, _ := cmd.Flags().GetString(flagCustomDecimalsPoint)
			if len(decimalsPoint) == 0 {
				decimalsPoint = "."
			} else if decimalsPoint != "." && decimalsPoint != "," {
				utils.PrintlnStdErr("ERR: decimals point must be either '.' or ','")
				os.Exit(1)
			}

			raw, _, _, err := utils.ConvertDisplayWithExponentIntoRaw(displayBalanceStr, int(decimals), rune(decimalsPoint[0]))
			utils.ExitOnErr(err, "failed to convert display balance into raw balance")

			fmt.Println(raw)
		},
	}

	cmd.Flags().StringP(flagCustomDecimalsPoint, flagCustomDecimalsPointShorthand, ".", "decimals point used to split parts of display balance")

	return cmd
}
