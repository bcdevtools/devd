package security_check

//goland:noinspection SpellCheckingInspection
import (
	"fmt"
	libutils "github.com/EscanBE/go-lib/utils"
	sectypes "github.com/bcdevtools/devd/cmd/security_check/types"
	"github.com/bcdevtools/devd/cmd/types"
	"github.com/bcdevtools/devd/cmd/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
	"path"
	"runtime"
	"runtime/debug"
)

var securityReport = &sectypes.SecurityReport{}

//goland:noinspection SpellCheckingInspection
const (
	secureCheckModuleDotNetrc           = ".NETRC"
	secureCheckModuleDotSecrets         = ".SECRETS"
	secureCheckModuleDotSsh             = ".SSH"
	secureCheckModuleUfw                = "UFW"
	secureCheckModuleSystemdServiceFile = "SDF"
	secureCheckModuleSshd               = "SSHD"
	secureCheckModuleUserPassword       = "UPWD"
	secureCheckModuleDefaultPortOpen    = "DPORT"
)

// Command registers a sub-tree of commands
func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "security-check",
		Aliases: []string{"secure-check", "sc"},
		Short:   "Perform several security check.",
		PreRun: func(cmd *cobra.Command, args []string) {
			ctx := types.UnwrapAppContext(cmd.Context())

			operationUserInfo := ctx.GetOperationUserInfo()

			if !operationUserInfo.OperatingAsSuperUser {
				libutils.PrintlnStdErr("ERR: must be ran as super user")
				libutils.PrintlnStdErr("*** Hint: Try again with 'sudo' ***")
				os.Exit(1)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			defer func() {
				if len(securityReport.SecurityRecords) < 1 {
					fmt.Println("CONGRATS: Everything seems good")
					return
				}

				libutils.PrintlnStdErr("\nFound", len(securityReport.SecurityRecords), "issues:")
				fatalIssuesCount := securityReport.CountFatal()
				if fatalIssuesCount < 1 {
					libutils.PrintlnStdErr("+ No fatal issue found.")
				} else if fatalIssuesCount == 1 {
					libutils.PrintlnStdErr("+ 1 fatal issue found!")
				} else {
					libutils.PrintlnStdErr("+", fatalIssuesCount, "fatal issues found!")
				}
				warningIssuesCount := securityReport.CountWarning()
				if warningIssuesCount < 1 {
					libutils.PrintlnStdErr("+ No warning issue found.")
				} else if warningIssuesCount == 1 {
					libutils.PrintlnStdErr("+ 1 warning issue found!")
				} else {
					libutils.PrintlnStdErr("+", warningIssuesCount, "warning issues found!")
				}

				libutils.PrintlnStdErr("\nSecurity Report:")
				securityReport.Sort()
				for i, record := range securityReport.SecurityRecords {
					libutils.PrintfStdErr("%d. %s\n", i+1, record)
				}

				if r := recover(); r != nil {
					libutils.PrintlnStdErr("\nFATAL: application exited with error:", r)
					libutils.PrintlnStdErr(string(debug.Stack()))
				}

				if fatalIssuesCount > 0 {
					os.Exit(2)
				} else {
					os.Exit(1)
				}
			}()

			homeDirs := getHomeDirs()

			if utils.IsLinux() {
				checkSshdConfig_LinuxOnly()
			} else {
				securityReport.Add(
					sectypes.
						NewSecurityRecord(secureCheckModuleSshd, false, fmt.Sprintf("Not yet implemented SSHD check for %s", runtime.GOOS)),
				)
			}

			for _, homeDir := range homeDirs {
				checkSshDir(homeDir)
				checkDotNetrc(homeDir)
				checkSecretDir(homeDir)
			}
			if utils.IsLinux() {
				checkUfwFirewall_LinuxOnly()
			} else {
				securityReport.Add(
					sectypes.
						NewSecurityRecord(secureCheckModuleUfw, false, fmt.Sprintf("Not yet implemented firewall check for %s", runtime.GOOS)),
				)
			}
			if utils.IsLinux() {
				checkLinuxUsersPasswdDisabled_LinuxOnly()
			} else {
				securityReport.Add(
					sectypes.
						NewSecurityRecord(secureCheckModuleUserPassword, false, fmt.Sprintf("Not yet implemented users password disabled check for %s", runtime.GOOS)),
				)
			}
			if utils.IsLinux() {
				checkServiceFiles_LinuxOnly()
			} else {
				securityReport.Add(
					sectypes.
						NewSecurityRecord(secureCheckModuleSystemdServiceFile, false, fmt.Sprintf("Not yet implemented checking systemd service files for %s", runtime.GOOS)),
				)
			}

			checkDefaultPortOpen()
		},
	}

	return cmd
}

func getHomeDirs() []string {
	var homeDirs []string

	if utils.IsLinux() {
		homeDirs = append(homeDirs, "/root")
	}

	usersRootDir := utils.GetUserGeneralRootDirLocation()
	var dirsInHome []os.DirEntry
	dirsInHome, err := os.ReadDir(usersRootDir)
	if err != nil {
		panic(errors.Wrap(err, fmt.Sprintf("failed to read users root dir in %s", usersRootDir)))
	}
	for _, dir := range dirsInHome {
		if dir == nil || !dir.IsDir() {
			continue
		}
		dirName := dir.Name()
		if utils.IsDarwin() {
			if dirName == "Shared" || dirName == "Guest" {
				continue
			}
		}
		homeDirs = append(homeDirs, path.Join(usersRootDir, dirName))
	}

	return homeDirs
}
