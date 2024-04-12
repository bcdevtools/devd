package check

import (
	"github.com/spf13/cobra"
)

// Commands registers a sub-tree of commands
func Commands() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "check",
		Aliases: []string{"chk"},
		Short:   "Check commands",
	}

	cmd.AddCommand(
		GetCheckPortCommand(),
	)

	return cmd
}
