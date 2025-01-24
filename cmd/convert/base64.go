package convert

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/bcdevtools/devd/v3/cmd/utils"
	"github.com/spf13/cobra"
)

const (
	flagDecode = "decode"
)

// GetBase64CaseCmd creates a helper command that encode/decode base64
func GetBase64CaseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "base64 [input]",
		Short: fmt.Sprintf("Encode base64, use --%s to decode", flagDecode),
		Long: fmt.Sprintf(`Encode/decode base64, use --%s to decode.
Support pipe.`, flagDecode),
		Run: func(cmd *cobra.Command, args []string) {
			args, err := utils.ProvidedArgsOrFromPipe(args)
			utils.ExitOnErr(err, "failed to get args from pipe")

			if cmd.Flag(flagDecode).Changed {
				utils.RequireExactArgsCount(args, 1, cmd)
				utils.PrintlnStdErr("INF: decoding base64")
				data, err := base64.StdEncoding.DecodeString(args[0])
				utils.ExitOnErr(err, "failed to decode base64")
				fmt.Println(string(data))
			} else {
				utils.RequireArgs(args, cmd)
				utils.PrintlnStdErr("INF: encoding base64")
				fmt.Println(base64.StdEncoding.EncodeToString([]byte(strings.Join(args, " "))))
			}
		},
	}

	cmd.Flags().BoolP(flagDecode, "d", false, "decode base64 instead of encode")

	return cmd
}
