package cmd

//goland:noinspection GoSnakeCaseUsage
import (
	"fmt"
	libutils "github.com/EscanBE/go-lib/utils"
	"github.com/bcdevtools/devd/cmd/convert"
	"github.com/bcdevtools/devd/cmd/debug"
	"github.com/bcdevtools/devd/cmd/query"
	"github.com/bcdevtools/devd/cmd/types"
	"github.com/bcdevtools/devd/cmd/utils"
	"github.com/bcdevtools/devd/constants"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"time"
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

	operationUserInfo, err := utils.GetOperationUserInfo()
	if err != nil {
		libutils.PrintlnStdErr("ERR: failed to get operation user info:", err)
		os.Exit(1)
	}

	if operationUserInfo.OperatingAsSuperUser {
		fmt.Println("WARN: Running as super user")
	}

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		switch cmd.Name() {
		case "version":
			return
		}

		changeWorkingUser(cmd)
		ensureRequireWorkingUsername(cmd)
	}

	ctx := types.NewContext(operationUserInfo)

	err = rootCmd.ExecuteContext(ctx)
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(convert.Commands())
	rootCmd.AddCommand(debug.Commands())
	rootCmd.AddCommand(query.Commands())

	rootCmd.PersistentFlags().String(constants.FLAG_USE_WORKING_USERNAME, "", "Use the specified username as working context username, must be either effective user or real user, if not specified, will use default selected working user")
	rootCmd.PersistentFlags().String(constants.FLAG_REQUIRE_WORKING_USERNAME, "", "Ensure working user is the specified user: if working user selected by application has username different with the specified username, application will exit with error")
	rootCmd.PersistentFlags().Bool("help", false, "show help")
}

func changeWorkingUser(cmd *cobra.Command) {
	selectedWorkingUsername, err := cmd.Flags().GetString(constants.FLAG_USE_WORKING_USERNAME)
	libutils.PanicIfErr(err, "failed to read flag")

	selectedWorkingUsername = strings.TrimSpace(selectedWorkingUsername)

	if len(selectedWorkingUsername) < 1 {
		return
	}

	ctx := types.UnwrapAppContext(cmd.Context())
	if ctx.GetWorkingUserInfo().Username == selectedWorkingUsername {
		return
	}

	var newWorkingUser *types.UserInfo

	operationUserInfo := ctx.GetOperationUserInfo()

	if operationUserInfo.EffectiveUserInfo.Username == selectedWorkingUsername {
		newWorkingUser = operationUserInfo.EffectiveUserInfo
	} else if operationUserInfo.RealUserInfo.Username == selectedWorkingUsername {
		newWorkingUser = operationUserInfo.RealUserInfo
	} else {
		libutils.PrintfStdErr("ERR: selected working user %s is not either effective user %s or real user %s\n", selectedWorkingUsername, operationUserInfo.EffectiveUserInfo.Username, operationUserInfo.RealUserInfo.Username)
		os.Exit(1)
	}

	fmt.Println("WARN: changing working user to", newWorkingUser.Username, "instead of default", ctx.GetWorkingUserInfo().Username)
	time.Sleep(2 * time.Second)

	ctx = ctx.WithWorkingUserInfo(newWorkingUser)
	cmd.SetContext(ctx)
}

func ensureRequireWorkingUsername(cmd *cobra.Command) {
	requireWorkingUsername, err := cmd.Flags().GetString(constants.FLAG_REQUIRE_WORKING_USERNAME)
	libutils.PanicIfErr(err, "failed to read flag")

	requireWorkingUsername = strings.TrimSpace(requireWorkingUsername)

	if len(requireWorkingUsername) < 1 {
		return
	}

	workingUser := types.UnwrapAppContext(cmd.Context()).GetWorkingUserInfo()
	if workingUser.Username != requireWorkingUsername {
		libutils.PrintfStdErr("ERR: working user is %s, but required user is %s\n", workingUser.Username, requireWorkingUsername)
		os.Exit(1)
	}
}
