package files

import (
	"fmt"
	"github.com/EscanBE/go-ienumerable/goe"
	libutils "github.com/EscanBE/go-lib/utils"
	"github.com/bcdevtools/devd/cmd/utils"
	"github.com/bcdevtools/devd/constants"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

const (
	flagRemoteToLocal         = "remote-to-local"
	flagLocalToRemote         = "local-to-remote"
	flagLocalToLocal          = "local-to-local"
	flagNoPassword            = "no-password"
	flagSshPassPassphraseMode = "passphrase"
	flagLogFile               = "log-file"
	flagPasswordFile          = "password-file"
	flagToolOptions           = "tool-options"
)

const rsyncOptCopyDir = "--recursive"

var defaultRsyncOptions = []string{"--human-readable", "--compress", "--stats"}

// RsyncCommands registers a sub-tree of commands
func RsyncCommands() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rsync [src] [dest]",
		Short: fmt.Sprintf("Remotely/Locally transfer file: %s rsync src dest", constants.BINARY_NAME),
		Long: fmt.Sprintf(`Remotely/Locally transfer file using rsync.
- Send file:
> %s rsync '/var/log/*.log' 'backup@192.168.0.2:/mnt/md0/backup/logs'
- Receive file:
> %s rsync 'load-balancer:/var/log/*.log' '/mnt/md0/backup/logs'

Note:
- This is just a wrapper of rsync, you must know how to use rsync and got rsync installed in order to use this.
  Actual translated rsync command would look similar to:
  > /usr/bin/rsync --human-readable --compress --stats -e ssh "server:/var/logs/*.log" "/mnt/md0/backup/logs"
- In case copy directory from local, the argument '%s' will be passed to rsync to indicate coping directory.
- When transfer from/to remote server, you must connect to that remote server at least one time before to perform host key verification (one time action) because the transfer will be performed via ssh.
`, constants.BINARY_NAME, constants.BINARY_NAME, rsyncOptCopyDir),
		Args: cobra.ExactArgs(2),
		Run:  remoteTransferFile,
	}

	cmd.Flags().Bool(
		flagRemoteToLocal,
		false,
		"ensure the transfer direction is from remote server to local",
	)

	cmd.Flags().Bool(
		flagLocalToRemote,
		false,
		"ensure the transfer direction is from local to remote server",
	)

	cmd.Flags().Bool(
		flagLocalToLocal,
		false,
		"ensure the transfer direction is from local to local",
	)

	cmd.Flags().StringSlice(
		flagToolOptions,
		defaultRsyncOptions,
		"supply options passes to rsync, comma separated",
	)

	cmd.Flags().String(
		flagToolFile,
		"",
		"custom rsync file path (absolute)",
	)

	cmd.Flags().String(
		flagPasswordFile,
		"",
		"file path which store password to access remote server",
	)

	cmd.Flags().String(
		flagLogFile,
		"",
		"log what we're doing to the specified file",
	)

	cmd.Flags().Bool(
		flagNoPassword,
		false,
		"force connect remote server without password (when remote user does not have password or identity key does not protected by password)",
	)

	cmd.Flags().Bool(
		flagSshPassPassphraseMode,
		false,
		"by default sshpass (if sshpass exists) passes password. If you are authenticating using passphrase, program will be hang (search phrase not found), supply this flag to indicate and would fix it",
	)

	return cmd
}

func remoteTransferFile(cmd *cobra.Command, args []string) {
	src := strings.TrimSpace(args[0])
	if len(src) < 1 {
		panic("source file/dir is empty")
	}

	dest := strings.TrimSpace(args[1])
	if len(dest) < 1 {
		panic("destination file/dir is empty")
	}

	isSrcRemote := strings.Contains(src, ":")
	isDestRemote := strings.Contains(dest, ":")

	if isSrcRemote && isDestRemote {
		panic("not support transfer direction from remote to remote")
	} else if isSrcRemote && !isDestRemote {
		confirm, _ := cmd.Flags().GetBool(flagRemoteToLocal)
		if !confirm {
			panic(fmt.Errorf("detected transfer direction is from remote to local so flag --%s is required to confirm", flagRemoteToLocal))
		}
	} else if !isSrcRemote && isDestRemote {
		confirm, _ := cmd.Flags().GetBool(flagLocalToRemote)
		if !confirm {
			panic(fmt.Errorf("detected transfer direction is from local to remote so flag --%s is required to confirm", flagLocalToRemote))
		}
	} else if !isSrcRemote && !isDestRemote {
		confirm, _ := cmd.Flags().GetBool(flagLocalToLocal)
		if !confirm {
			panic(fmt.Errorf("detected local transfer so flag --%s is required to confirm", flagLocalToLocal))
		}
	}

	var isSrcLocalDir bool

	if !isSrcRemote {
		file, err := os.Stat(src)
		if os.IsNotExist(err) {
			panic(fmt.Errorf("local source file/dir does not exists: %s", src))
		}
		if err != nil {
			panic(errors.Wrap(err, fmt.Sprintf("problem while checking local source file %s", src)))
		}
		isSrcLocalDir = file.IsDir()
	}

	if !isDestRemote {
		fi, err := os.Stat(dest)
		if err == nil {
			if fi.IsDir() {
				// ok
			} else {
				panic(fmt.Errorf("local destination file/dir already exists: %s", dest))
			}
		} else {
			if os.IsNotExist(err) {
				// ok
			} else {
				panic(errors.Wrap(err, fmt.Sprintf("problem while checking local destination file %s", dest)))
			}
		}
	}

	toolName := "rsync"
	customToolName, _ := cmd.Flags().GetString(flagToolFile)
	customToolName = strings.TrimSpace(customToolName)
	if len(customToolName) > 0 {
		_, err := os.Stat(customToolName)
		if os.IsNotExist(err) {
			panic(fmt.Errorf("custom tool file does not exists: %s", customToolName))
		}

		toolName = customToolName
	}

	options, _ := cmd.Flags().GetStringSlice(flagToolOptions)
	if len(options) < 1 {
		options = defaultRsyncOptions
	}

	if isSrcLocalDir {
		ieOptions := goe.NewIEnumerable(options...)

		if !isDestRemote {
			// local to local transfer => remove compress flag
			ieOptions = ieOptions.Where(func(option string) bool {
				return !strings.EqualFold(option, "--compress")
			})
		}

		if !ieOptions.AnyBy(isOrContainsRsyncRecursiveFlag) {
			// in case copy from local dir, supply flag '--recursive'
			ieOptions = ieOptions.Append(rsyncOptCopyDir)
		}

		options = ieOptions.ToArray()
	}

	logFile, _ := cmd.Flags().GetString(flagLogFile)
	if len(logFile) > 0 {
		duplicated := goe.NewIEnumerable[string](options...).AnyBy(func(flag string) bool {
			return strings.HasPrefix(flag, "--log-file ") || strings.HasPrefix(flag, "--log-file=")
		})
		if duplicated {
			panic(fmt.Sprintf("duplicated flags --%s", flagLogFile))
		}
		options = append(options, "--log-file", logFile)
	}

	if !isSrcRemote && !isDestRemote {
		launchApp(toolName, append(options, src, dest), nil)
		return
	}

	noPassword, _ := cmd.Flags().GetBool(flagNoPassword)
	if noPassword {
		launchApp(toolName, append(options, "--rsh", "ssh", src, dest), nil)
		return
	}

	sshPassPhrase, _ := cmd.Flags().GetBool(flagSshPassPassphraseMode)

	passwordFile, _ := cmd.Flags().GetString(flagPasswordFile)
	if len(passwordFile) > 0 {
		password := utils.MustReadPasswordFromEnsuredSecurityFile(passwordFile) // placed it here to validate the file

		if utils.HasToolSshPass() {
			fmt.Println("Using sshpass to passing password file")

			var cmdArgs []string
			if sshPassPhrase {
				//goland:noinspection SpellCheckingInspection
				cmdArgs = []string{"-P", "assphrase", "-f", passwordFile}
			} else {
				cmdArgs = []string{"-f", passwordFile}
			}
			cmdArgs = append(cmdArgs, toolName)
			cmdArgs = append(cmdArgs, options...)
			cmdArgs = append(cmdArgs, "--rsh", "ssh", src, dest)

			launchApp("sshpass", cmdArgs, nil)
			return
		}

		fmt.Println("Using environment variable", constants.ENV_RSYNC_PASSWORD, "to passing password from password file to rsync")
		fmt.Println("**WARNING: if remote machine does not have rsync service running, password prompt still appears")
		launchApp(toolName, append(options, "--rsh", "ssh", src, dest), []string{fmt.Sprintf("%s=%s", constants.ENV_RSYNC_PASSWORD, password)})
		return
	}

	rsyncPassword := strings.TrimSpace(os.Getenv(constants.ENV_RSYNC_PASSWORD))
	sshPassword := strings.TrimSpace(os.Getenv(constants.ENV_SSHPASS))

	var password string
	if len(rsyncPassword) > 0 && len(sshPassword) > 0 {
		if rsyncPassword != sshPassword {
			panic(fmt.Errorf("both environment variables %s and %s are set but mis-match, consider remove one to take the rest", constants.ENV_RSYNC_PASSWORD, constants.ENV_SSHPASS))
		}

		password = rsyncPassword
	} else if len(rsyncPassword) > 0 {
		password = rsyncPassword
	} else if len(sshPassword) > 0 {
		password = sshPassword
	} else {
		panic(fmt.Errorf("missing password for remote server, either environment variable %s or %s or flag --%s is required", constants.ENV_RSYNC_PASSWORD, constants.ENV_SSHPASS, flagPasswordFile))
	}

	if utils.HasToolSshPass() {
		if len(rsyncPassword) > 0 && len(sshPassword) == 0 {
			fmt.Println("Copied environment variable value from", constants.ENV_RSYNC_PASSWORD, "to", constants.ENV_SSHPASS)
		}
		fmt.Println("Using sshpass to passing password via environment variable", constants.ENV_SSHPASS)

		var cmdArgs []string
		if sshPassPhrase {
			//goland:noinspection SpellCheckingInspection
			cmdArgs = []string{"-P", "assphrase", "-e", toolName}
		} else {
			cmdArgs = []string{"-e", toolName}
		}
		cmdArgs = append(cmdArgs, options...)
		cmdArgs = append(cmdArgs, "--rsh", "ssh", src, dest)

		launchApp("sshpass", cmdArgs, []string{fmt.Sprintf("%s=%s", constants.ENV_SSHPASS, password)})
		return
	}

	if len(sshPassword) > 0 && len(rsyncPassword) == 0 {
		fmt.Println("Copied environment variable value from", constants.ENV_SSHPASS, "to", constants.ENV_RSYNC_PASSWORD)
	}
	fmt.Println("Using environment variable", constants.ENV_RSYNC_PASSWORD, "to passing password to rsync")
	fmt.Println("**WARNING: if remote machine does not have rsync service running, password prompt still appears")
	launchApp(toolName, append(options, "--rsh", "ssh", src, dest), []string{fmt.Sprintf("%s=%s", constants.ENV_RSYNC_PASSWORD, password)})
}

func launchApp(toolName string, args []string, additionalEnvVars []string) {
	fmt.Println("Rsync arguments:\n", toolName, strings.Join(args, " "))
	fmt.Println("Begin rsync at", utils.NowStr())

	envVars := os.Environ()
	if len(additionalEnvVars) > 0 {
		envVars = append(envVars, additionalEnvVars...)
	}

	defer func() {
		fmt.Println("Finished rsync at", utils.NowStr())
	}()

	ec := utils.LaunchAppWithDirectStd(toolName, args, envVars)

	if ec != 0 {
		libutils.PrintlnStdErr("ERR: Failed to rsync at", utils.NowStr())
		os.Exit(ec)
	}
}

func isOrContainsRsyncRecursiveFlag(option string) bool {
	if strings.HasPrefix(option, "--") {
		return strings.EqualFold(option, "--recursive") || strings.EqualFold(option, "--archive")
	} else if strings.HasPrefix(option, "-") {
		return strings.Contains(option, "r") || strings.Contains(option, "a")
	} else {
		return false
	}
}
