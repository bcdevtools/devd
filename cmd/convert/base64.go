package convert

import (
	"encoding/base64"
	"encoding/hex"
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
	const flagShowBuffer = "show-buffer"
	const flagFromBuffer = "from-buffer"

	cmd := &cobra.Command{
		Use:   "base64 [input]",
		Short: fmt.Sprintf("Encode base64, use --%s to decode", flagDecode),
		Long: fmt.Sprintf(`Encode/decode base64, use --%s to decode.
Support pipe.`, flagDecode),
		Run: func(cmd *cobra.Command, args []string) {
			args, err := utils.ProvidedArgsOrFromPipe(args)
			utils.ExitOnErr(err, "failed to get args from pipe")

			decode := cmd.Flag(flagDecode).Changed
			showBuffer := cmd.Flag(flagShowBuffer).Changed
			fromBuffer := cmd.Flag(flagFromBuffer).Changed
			if decode {
				if fromBuffer {
					panic(fmt.Sprintf("flag --%s is not supported on decoding mode", flagFromBuffer))
				}
				utils.RequireExactArgsCount(args, 1, cmd)
				utils.PrintlnStdErr("INF: decoding base64")
				data, err := base64.StdEncoding.DecodeString(args[0])
				utils.ExitOnErr(err, "failed to decode base64")
				if showBuffer {
					utils.PrintfStdErr("(buffer: %s)\n", hex.EncodeToString(data))
				}
				fmt.Println(string(data))
			} else {
				utils.RequireArgs(args, cmd)
				utils.PrintfStdErr("INF: encoding base64 (use --%s to decode)\n", flagDecode)
				data := []byte(strings.Join(args, " "))
				if fromBuffer {
					data, err = hex.DecodeString(string(data))
					utils.ExitOnErr(err, "failed to decode hex buffer")
				}
				if showBuffer {
					utils.PrintfStdErr("(buffer: %s)\n", hex.EncodeToString(data))
				}
				fmt.Println(base64.StdEncoding.EncodeToString(data))
			}
		},
	}

	cmd.Flags().BoolP(flagDecode, "d", false, "decode base64 instead of encode")
	cmd.Flags().Bool(flagShowBuffer, false, "show buffer in hex format")
	cmd.Flags().Bool(flagFromBuffer, false, "input is hex buffer, effective ONLY on encoding")

	return cmd
}
