package convert

import (
	"fmt"
	libutils "github.com/EscanBE/go-lib/utils"
	"github.com/bcdevtools/devd/cmd/utils"
	"github.com/spf13/cobra"
	"strings"
)

// GetConvertToUpperCaseCmd creates a helper command that convert input into upper case
func GetConvertToUpperCaseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "to_upper_case [text]",
		Aliases: []string{"uppercase"},
		Short: `Convert input into upper case.
Support pipe.`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			args, err = utils.ProvidedArgsOrFromPipe(args)
			libutils.ExitIfErr(err, "failed to get args from pipe")
			utils.RequireArgs(args, cmd)

			fmt.Println(strings.ToUpper(strings.Join(args, " ")))
		},
	}

	return cmd
}
