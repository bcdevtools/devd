package gen

import (
	"github.com/spf13/cobra"
)

// Commands registers a sub-tree of commands
func Commands() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "generate [template]",
		Aliases: []string{"gen"},
		Short:   "Generate template",
	}

	cmd.AddCommand(
		GenerateVisudoCommand(),
		GenerateSshKeypairCommand(),
		GenerateUfwAllowCommand(),
	)

	return cmd
}
