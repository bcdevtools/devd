package convert

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

// GetConvertToLowerCaseCmd creates a helper command that convert input into lower case
func GetConvertToLowerCaseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "to_lower_case [text]",
		Aliases: []string{"lowercase"},
		Short:   "Convert input into lower case",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(strings.ToLower(strings.Join(args, " ")))
		},
	}

	return cmd
}
