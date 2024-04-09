package query

import (
	"github.com/spf13/cobra"
)

const (
	flagHost = "host"
)

// Commands registers a sub-tree of commands
func Commands() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "Query commands",
	}

	cmd.AddCommand(
		GetQueryErc20Command(),
	)

	return cmd
}
