package convert

import (
	"fmt"
	"github.com/bcdevtools/devd/v2/cmd/utils"
	"github.com/spf13/cobra"
	"strings"
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
			var err error
			args, err = utils.ProvidedArgsOrFromPipe(args)
			utils.ExitOnErr(err, "failed to get args from pipe")
			utils.RequireArgs(args, cmd)

			fmt.Println(strings.ToUpper(strings.Join(args, " ")))
		},
	}

	return cmd
}
