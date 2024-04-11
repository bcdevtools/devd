package convert

import (
	"encoding/base64"
	"fmt"
	libutils "github.com/EscanBE/go-lib/utils"
	"github.com/bcdevtools/devd/cmd/utils"
	"github.com/spf13/cobra"
)

// GetDecodeBase64CaseCmd creates a helper command that decode base64
func GetDecodeBase64CaseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "decode_base64 [base64]",
		Short: `Decode base64.
Support pipe.`,
		Run: func(cmd *cobra.Command, args []string) {
			args, err := utils.ProvidedArgsOrFromPipe(args)
			libutils.ExitIfErr(err, "failed to get args from pipe")
			utils.RequireExactArgsCount(args, 1, cmd)

			data, err := base64.StdEncoding.DecodeString(args[0])
			libutils.ExitIfErr(err, "failed to decode base64")

			fmt.Println(string(data))
		},
	}

	return cmd
}
