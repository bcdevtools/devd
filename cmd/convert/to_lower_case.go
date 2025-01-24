package convert

import (
	"fmt"
	"strings"

	"github.com/bcdevtools/devd/v3/cmd/utils"
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
			var err error
			args, err = utils.ProvidedArgsOrFromPipe(args)
			utils.ExitOnErr(err, "failed to get args from pipe")
			utils.RequireArgs(args, cmd)

			fmt.Println(strings.ToLower(strings.Join(args, " ")))
		},
	}

	return cmd
}
