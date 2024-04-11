package debug

import (
	"github.com/spf13/cobra"
)

// Commands registers a sub-tree of commands
func Commands() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "debug",
		Aliases: []string{"d"},
		Short:   "Debug commands",
	}

	cmd.AddCommand(
		GetIntrinsicCommand(),
	)

	return cmd
}
