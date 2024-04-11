package convert

import (
	"encoding/base64"
	"fmt"
	"github.com/bcdevtools/devd/cmd/utils"
	"github.com/spf13/cobra"
	"strings"
)

// GetEncodeBase64CaseCmd creates a helper command that encode input into base64
func GetEncodeBase64CaseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "encode_base64 [text]",
		Aliases: []string{"base64"},
		Short: `Encode input into base64.
Support pipe.`,
		Run: func(cmd *cobra.Command, args []string) {
			args, err := utils.ProvidedArgsOrFromPipe(args)
			utils.ExitOnErr(err, "failed to get args from pipe")
			utils.RequireArgs(args, cmd)

			fmt.Println(base64.StdEncoding.EncodeToString([]byte(strings.Join(args, " "))))
		},
	}

	return cmd
}
