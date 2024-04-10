package query

import (
	"context"
	"fmt"
	libutils "github.com/EscanBE/go-lib/utils"
	"github.com/bcdevtools/devd/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"os"
	"regexp"
	"strings"
)

// GetQueryTxReceiptCommand registers a sub-tree of commands
func GetQueryTxReceiptCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "eth_getTransactionReceipt [0xhash]",
		Aliases: []string{"receipt"},
		Short:   "eth_getTransactionReceipt",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ethClient, _ := mustGetEthClient(cmd)

			input := strings.ToLower(args[0])

			if !regexp.MustCompile(`^0x[a-f\d]{64}$`).MatchString(input) {
				libutils.PrintlnStdErr("ERR: invalid EVM transaction hash format")
				os.Exit(1)
			}

			receipt, err := ethClient.TransactionReceipt(context.Background(), common.HexToHash(input))
			libutils.ExitIfErr(err, "failed to get transaction by hash")

			bz, err := receipt.MarshalJSON()
			libutils.ExitIfErr(err, "failed to marshal receipt to json")

			beautifyBz, err := utils.BeautifyJson(bz)
			if err != nil {
				libutils.PrintlnStdErr("failed to beautify json:", err)
				fmt.Println(string(bz))
			} else {
				fmt.Println(string(beautifyBz))
			}
		},
	}

	cmd.Flags().StringP(flagRpc, "p", "http://localhost:8545", "EVM Json-RPC url")

	return cmd
}
