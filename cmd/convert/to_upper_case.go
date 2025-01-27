package convert

import (
	"fmt"
	"github.com/bcdevtools/devd/v2/constants"
	"strings"

	"github.com/bcdevtools/devd/v2/cmd/utils"
	"github.com/spf13/cobra"
)

// GetConvertToUpperCaseCmd creates a helper command that convert input into upper case
func GetConvertToUpperCaseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "to_upper_case [text]",
		Aliases: []string{"uppercase"},
		Short:   "Convert input text into upper case.",
		Long: `Convert input text into upper case.
Support pipe.`,
		Run: func(cmd *cobra.Command, args []string) {
			utils.PrintfStdErr("WARN: from v3, this command will be renamed to `%s convert case ... --upper`\n", constants.BINARY_NAME)

			var err error
			args, err = utils.ProvidedArgsOrFromPipe(args)
			utils.ExitOnErr(err, "failed to get args from pipe")
			utils.RequireArgs(args, cmd)

			fmt.Println(strings.ToUpper(strings.Join(args, " ")))
		},
	}

	return cmd
}
