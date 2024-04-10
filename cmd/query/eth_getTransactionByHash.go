package query

import (
	"context"
	libutils "github.com/EscanBE/go-lib/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"os"
	"regexp"
	"strings"
)

// GetQueryTxCommand registers a sub-tree of commands
func GetQueryTxCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "eth_getTransactionByHash [0xhash]",
		Aliases: []string{"tx"},
		Short:   "eth_getTransactionByHash",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ethClient, _ := mustGetEthClient(cmd)

			input := strings.ToLower(args[0])

			if !regexp.MustCompile(`^0x[a-f\d]{64}$`).MatchString(input) {
				libutils.PrintlnStdErr("ERR: invalid EVM transaction hash format")
				os.Exit(1)
			}

			tx, _, err := ethClient.TransactionByHash(context.Background(), common.HexToHash(input))
			libutils.ExitIfErr(err, "failed to get transaction by hash")

			bz, err := tx.MarshalJSON()
			libutils.ExitIfErr(err, "failed to marshal transaction to json")

			tryPrintBeautyJson(bz)
		},
	}

	cmd.Flags().String(flagRpc, "http://localhost:8545", "EVM Json-RPC url")

	return cmd
}
