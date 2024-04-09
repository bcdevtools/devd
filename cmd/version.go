package cmd

import (
	"fmt"
	"github.com/bcdevtools/devd/constants"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command, it prints the current version of the binary
var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "Show binary version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(constants.VERSION)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
