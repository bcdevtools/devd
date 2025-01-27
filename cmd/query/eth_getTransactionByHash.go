package query

import (
	"context"
	"fmt"
	"github.com/bcdevtools/devd/v2/constants"
	"os"
	"regexp"
	"strings"

	"github.com/bcdevtools/devd/v2/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

func GetQueryTxCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "eth_getTransactionByHash [0xhash]",
		Aliases: []string{"tx"},
		Short:   "eth_getTransactionByHash",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			utils.PrintfStdErr("WARN! Deprecation notice: from v3, command alias `tx` will be removed, please use `%s q eth_getTransactionByHash ...` instead of `%s q tx ...`\n", constants.BINARY_NAME, constants.BINARY_NAME)

			ethClient, _ := mustGetEthClient(cmd, false)

			input := strings.ToLower(args[0])

			if !regexp.MustCompile(`^0x[a-f\d]{64}$`).MatchString(input) {
				utils.PrintlnStdErr("ERR: invalid EVM transaction hash format")
				os.Exit(1)
			}

			tx, _, err := ethClient.TransactionByHash(context.Background(), common.HexToHash(input))
			utils.ExitOnErr(err, "failed to get transaction by hash")

			bz, err := utils.MarshalPrettyJsonEvmTx(tx, &utils.PrettyMarshalJsonEvmTxOption{
				InjectFrom:                true,
				InjectTranslateAbleFields: true,
			})
			utils.ExitOnErr(err, "failed to marshal transaction to json")

			fmt.Println(string(bz))
		},
	}

	cmd.Flags().String(flagRpc, "", flagEvmRpcDesc)

	return cmd
}
