package hash

import (
	"github.com/spf13/cobra"
)

// Commands registers a sub-tree of commands
func Commands() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hash",
		Short: "Hashing commands",
	}

	cmd.AddCommand(
		GetMd5Command(),
		GetKeccak256Command(),
		GetKeccak512Command(),
	)

	return cmd
}
