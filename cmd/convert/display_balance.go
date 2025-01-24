package convert

import (
	"fmt"
	"math/big"
	"os"
	"strconv"

	"github.com/bcdevtools/devd/v3/cmd/utils"
	"github.com/spf13/cobra"
)

// GetDisplayBalanceCmd creates a helper command that convert raw balance into display balance
func GetDisplayBalanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "display_balance [raw balance] [decimals]",
		Aliases: []string{"dbal"},
		Short:   "Convert raw balance into display balance.",
		Long: `Convert raw balance into display balance.
Sample: 10011100 with 6 exponent => 10.0111`,
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			rawBalanceStr := args[0]
			decimalsStr := args[1]

			balance, err := utils.ReadCustomInteger(rawBalanceStr)
			if err != nil {
				var ok bool
				balance, ok = new(big.Int).SetString(rawBalanceStr, 10)
				if !ok {
					utils.PrintlnStdErr("ERR: failed to read, raw balance is not a number")
					os.Exit(1)
				}
			}

			decimals, err := strconv.ParseInt(decimalsStr, 10, 64)
			utils.ExitOnErr(err, "failed to read, decimals is not a number")

			display, _, _, err := utils.ConvertNumberIntoDisplayWithExponent(balance, int(decimals))
			utils.ExitOnErr(err, "failed to convert raw balance into display balance")

			fmt.Println(display)
		},
	}

	return cmd
}
