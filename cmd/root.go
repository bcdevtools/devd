package cmd

//goland:noinspection GoSnakeCaseUsage
import (
	"os"

	"github.com/bcdevtools/devd/v3/cmd/check"
	"github.com/bcdevtools/devd/v3/cmd/convert"
	"github.com/bcdevtools/devd/v3/cmd/debug"
	"github.com/bcdevtools/devd/v3/cmd/hash"
	"github.com/bcdevtools/devd/v3/cmd/query"
	"github.com/bcdevtools/devd/v3/cmd/tx"
	"github.com/bcdevtools/devd/v3/cmd/types"
	"github.com/bcdevtools/devd/v3/constants"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   constants.BINARY_NAME,
	Short: constants.BINARY_NAME,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.CompletionOptions.HiddenDefaultCmd = true    // hide the 'completion' subcommand
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true}) // hide the 'help' subcommand

	ctx := types.NewContext()

	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(convert.Commands())
	rootCmd.AddCommand(debug.Commands())
	rootCmd.AddCommand(query.Commands())
	rootCmd.AddCommand(hash.Commands())
	rootCmd.AddCommand(check.Commands())
	rootCmd.AddCommand(tx.Commands())

	rootCmd.PersistentFlags().Bool("help", false, "show help")
}
