package convert

import (
	"encoding/base64"
	"fmt"

	"github.com/bcdevtools/devd/v3/cmd/utils"
	"github.com/spf13/cobra"
)

// GetDecodeBase64CaseCmd creates a helper command that decode base64
func GetDecodeBase64CaseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "decode_base64 [base64]",
		Short: "Decode base64.",
		Long: `Decode base64.
Support pipe.`,
		Run: func(cmd *cobra.Command, args []string) {
			args, err := utils.ProvidedArgsOrFromPipe(args)
			utils.ExitOnErr(err, "failed to get args from pipe")
			utils.RequireExactArgsCount(args, 1, cmd)

			data, err := base64.StdEncoding.DecodeString(args[0])
			utils.ExitOnErr(err, "failed to decode base64")

			fmt.Println(string(data))
		},
	}

	return cmd
}
