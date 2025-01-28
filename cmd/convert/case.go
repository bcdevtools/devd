package convert

import (
	"fmt"
	"strings"

	"github.com/bcdevtools/devd/v3/cmd/utils"
	"github.com/spf13/cobra"
)

const (
	flagToUpperCase = "upper"
)

// GetConvertCaseCmd creates a helper command that convert input into lower/upper case
func GetConvertCaseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "case [text]",
		Short: fmt.Sprintf("Convert input text into lower case, use --%s to upper case", flagToUpperCase),
		Long: fmt.Sprintf(`Convert input text into lower case, use --%s to upper case.
Support pipe.`, flagToUpperCase),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			args, err = utils.ProvidedArgsOrFromPipe(args)
			utils.ExitOnErr(err, "failed to get args from pipe")
			utils.RequireArgs(args, cmd)

			input := strings.Join(args, " ")
			if cmd.Flag(flagToUpperCase).Changed {
				utils.PrintlnStdErr("INF: converting to upper case")
				fmt.Println(strings.ToUpper(input))
				return
			}

			utils.PrintfStdErr("INF: converting to lower case (use --%s to upper case)\n", flagToUpperCase)
			fmt.Println(strings.ToLower(input))
		},
	}

	cmd.Flags().BoolP(flagToUpperCase, "u", false, "convert input to upper case instead of lower case")

	return cmd
}
