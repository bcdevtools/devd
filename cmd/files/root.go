package files

import (
	"github.com/spf13/cobra"
)

const (
	flagToolFile = "tool-file"
)

// Commands registers a sub-tree of commands
func Commands() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "files",
		Aliases: []string{"f", "file"},
		Short:   "Interacting with files",
	}

	cmd.AddCommand(
		RsyncCommands(),
		RemoveCommands(),
		RemoveParanoidCommands(),
		CompressLz4Command(),
		DecompressLz4Command(),
	)

	return cmd
}
