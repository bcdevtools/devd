package query

import (
	"context"
	"fmt"
	"github.com/bcdevtools/devd/v3/cmd/flags"
	"os"
	"regexp"
	"strings"

	"github.com/bcdevtools/devd/v3/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

func GetQueryTxReceiptCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "eth_getTransactionReceipt [0xhash]",
		Aliases: []string{"receipt"},
		Short:   "eth_getTransactionReceipt",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ethClient, _ := flags.MustGetEthClient(cmd)

			input := strings.ToLower(args[0])

			if !regexp.MustCompile(`^0x[a-f\d]{64}$`).MatchString(input) {
				utils.PrintlnStdErr("ERR: invalid EVM transaction hash format")
				os.Exit(1)
			}

			receipt, err := ethClient.TransactionReceipt(context.Background(), common.HexToHash(input))
			utils.ExitOnErr(err, "failed to get transaction by hash")

			bz, err := utils.MarshalPrettyJsonEvmTxReceipt(receipt, &utils.PrettyMarshalJsonEvmTxReceiptOption{
				InjectTranslateAbleFields: true,
			})
			utils.ExitOnErr(err, "failed to marshal receipt to json")

			fmt.Println(string(bz))
		},
	}

	cmd.Flags().String(flags.FlagEvmRpc, "", flags.FlagEvmRpcDesc)

	return cmd
}
