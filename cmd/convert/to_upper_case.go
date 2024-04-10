package convert

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

// GetConvertToUpperCaseCmd creates a helper command that convert input into upper case
func GetConvertToUpperCaseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "to_upper_case [text]",
		Aliases: []string{"uppercase"},
		Short:   "Convert input into upper case",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(strings.ToUpper(strings.Join(args, " ")))
		},
	}

	return cmd
}
