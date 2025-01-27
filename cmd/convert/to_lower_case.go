package convert

import (
	"fmt"
	"github.com/bcdevtools/devd/v2/constants"
	"strings"

	"github.com/bcdevtools/devd/v2/cmd/utils"
	"github.com/spf13/cobra"
)

// GetConvertToLowerCaseCmd creates a helper command that convert input into lower case
func GetConvertToLowerCaseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "to_lower_case [text]",
		Aliases: []string{"lowercase"},
		Short:   "Convert input text into lower case.",
		Long: `Convert input text into lower case.
Support pipe.`,
		Run: func(cmd *cobra.Command, args []string) {
			utils.PrintfStdErr("WARN: from v3, this command will be renamed to `%s convert case ...`\n", constants.BINARY_NAME)

			var err error
			args, err = utils.ProvidedArgsOrFromPipe(args)
			utils.ExitOnErr(err, "failed to get args from pipe")
			utils.RequireArgs(args, cmd)

			fmt.Println(strings.ToLower(strings.Join(args, " ")))
		},
	}

	return cmd
}
