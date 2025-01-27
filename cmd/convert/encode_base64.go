package convert

import (
	"encoding/base64"
	"fmt"
	"github.com/bcdevtools/devd/v2/constants"
	"strings"

	"github.com/bcdevtools/devd/v2/cmd/utils"
	"github.com/spf13/cobra"
)

// GetEncodeBase64CaseCmd creates a helper command that encode input into base64
func GetEncodeBase64CaseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "encode_base64 [text]",
		Aliases: []string{"base64"},
		Short:   "Encode input text into base64.",
		Long: `Encode input text into base64.
Support pipe.`,
		Run: func(cmd *cobra.Command, args []string) {
			utils.PrintfStdErr("WARN: from v3, this command will be renamed to `%s convert base64 ...`\n", constants.BINARY_NAME)

			args, err := utils.ProvidedArgsOrFromPipe(args)
			utils.ExitOnErr(err, "failed to get args from pipe")
			utils.RequireArgs(args, cmd)

			fmt.Println(base64.StdEncoding.EncodeToString([]byte(strings.Join(args, " "))))
		},
	}

	return cmd
}
