package convert

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/bcdevtools/devd/v3/cmd/utils"
	"github.com/spf13/cobra"
)

func GetDecodeRawEvmTxCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "decode-raw-tx [raw RLP-encoded EVM tx hex]",
		Short: `Decode the raw RLP-encoded EVM tx to see inner details, additional information will be injected with prefix '_'`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			rawTx := strings.ToLower(args[0])
			if !regexp.MustCompile(`^(0x)?[a-fA-F\d]+$`).MatchString(rawTx) {
				utils.PrintlnStdErr("ERR: invalid raw EVM tx, must be valid hex-encoded string")
				os.Exit(1)
			}

			tx, err := utils.DecodeRawEvmTx(rawTx)
			utils.ExitOnErr(err, "failed to decode into EVM tx")

			bz, err := utils.MarshalPrettyJsonEvmTx(tx, &utils.PrettyMarshalJsonEvmTxOption{
				InjectFrom:                true,
				InjectTranslateAbleFields: true,
			})
			utils.ExitOnErr(err, "failed to marshal tx to json")

			fmt.Println(string(bz))
		},
	}

	return cmd
}
