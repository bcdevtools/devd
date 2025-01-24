package query

import (
	"math/big"

	"github.com/spf13/cobra"
)

const (
	flagHeight = "height"
)

// Commands registers a sub-tree of commands
func Commands() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "Query commands",
	}

	cmd.AddCommand(
		GetQueryErc20Command(),
		GetQueryTxCommand(),
		GetQueryTxReceiptCommand(),
		GetQueryBlockCommand(),
		GetQueryTraceTxCommand(),
		GetQueryBalanceCommand(),
		GetQueryTxsInBlockCommand(),
		GetQueryTxEventsCommand(),
	)

	return cmd
}

func readContextHeightFromFlag(cmd *cobra.Command) *big.Int {
	height, _ := cmd.Flags().GetInt64(flagHeight)
	if height > 0 {
		return big.NewInt(height)
	}

	return nil
}
