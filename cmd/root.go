package cmd

//goland:noinspection GoSnakeCaseUsage
import (
	"github.com/bcdevtools/devd/v2/cmd/convert"
	"github.com/bcdevtools/devd/v2/cmd/debug"
	"github.com/bcdevtools/devd/v2/cmd/hash"
	"github.com/bcdevtools/devd/v2/cmd/query"
	"github.com/bcdevtools/devd/v2/cmd/types"
	"github.com/bcdevtools/devd/v2/constants"
	"github.com/spf13/cobra"
	"os"
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

	rootCmd.PersistentFlags().Bool("help", false, "show help")
}
