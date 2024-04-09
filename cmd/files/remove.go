package files

import (
	"fmt"
	libutils "github.com/EscanBE/go-lib/utils"
	"github.com/bcdevtools/devd/cmd/utils"
	"github.com/spf13/cobra"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

const (
	flagFast          = "fast"
	flagConfirmDelete = "delete"
	flagYes           = "yes"
)

// RemoveCommands registers a sub-tree of commands
func RemoveCommands() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rm [files]",
		Short: "Secure remove files and directory",
		Args:  cobra.ExactArgs(1),
		Run:   removeFiles,
	}

	if utils.IsLinux() {
		cmd.Flags().Bool(flagFast, false, "Fast mode, just run a single round")
	}
	cmd.Flags().Bool(flagConfirmDelete, false, "REQUIRED flag")
	cmd.Flags().BoolP(flagYes, "y", false, "Confirmed action and run immediately without having to wait the safe time")

	return cmd
}

func removeFiles(cmd *cobra.Command, args []string) {
	if !cmd.Flag(flagConfirmDelete).Changed {
		libutils.PrintfStdErr("ERR: For the purpose of preventing incident from incorrect typing/usage, flag '--%s' is required, indicate that you know what this is.\n", flagConfirmDelete)
		os.Exit(1)
	}

	target := args[0]

	fi, err := os.Stat(target)
	if err != nil {
		if os.IsNotExist(err) {
			libutils.PrintlnStdErr("ERR: target does not exists:", target)
		} else {
			libutils.PrintlnStdErr("ERR: failed to stat", target)
			libutils.PrintlnStdErr("ERR:", err)
		}
		os.Exit(1)
	}

	isDir := fi.IsDir()
	var execArgs []string
	if utils.IsLinux() {
		if isDir {
			execArgs = []string{"srm"}

			if cmd.Flag(flagFast).Changed {
				//goland:noinspection SpellCheckingInspection
				execArgs = append(execArgs, "-rfll")
			} else {
				execArgs = append(execArgs, "-r")
			}
		} else {
			execArgs = []string{"shred"}

			if cmd.Flag(flagFast).Changed {
				execArgs = append(execArgs, "--iterations=0")
			} else {
				execArgs = append(execArgs, "--iterations=3")
			}

			execArgs = append(execArgs, "--zero", "--verbose")
		}
	} else if utils.IsDarwin() {
		execArgs = []string{"rm", "-P"}

		if isDir {
			execArgs = append(execArgs, "-r")
		}
	}
	execArgs = append(execArgs, target)

	if !utils.HasBinaryName(execArgs[0]) {
		libutils.PrintlnStdErr("ERR: Missing required tool", execArgs[0])
		os.Exit(1)
	}

	if isDir {
		utils.WarnIfNotRunningUnderScreenSession("secure-delete")
	}

	if cmd.Flag(flagYes).Changed {
		fmt.Println("WARN: The following command is going to be used:")
		fmt.Println(strings.Join(execArgs, " "))
	} else {
		timeSleep := 15 * time.Second
		fmt.Printf("WARN: The following command is going to be executed after %v\n", timeSleep)
		fmt.Println(strings.Join(execArgs, " "))
		time.Sleep(timeSleep)
	}

	fmt.Println("Begin execution at", time.Now().Format("2006-01-02 15:04:05"))

	ec := utils.LaunchAppWithDirectStd(
		execArgs[0],
		execArgs[1:],
		nil,
	)

	fmt.Println("Finished execution at", time.Now().Format("2006-01-02 15:04:05"))

	if ec == 0 && isDir {
		fmt.Println("HINT: Paranoid? Do this:")
		tmpFile := path.Join(target, "tmp_file")
		tmpFileAbsPath, err := filepath.Abs(tmpFile)
		if err != nil {
			tmpFileAbsPath = tmpFile
		}
		fmt.Printf("> dd if=/dev/zero of=%s bs=1M count=999999999999 ; rm -f %s", tmpFileAbsPath, tmpFileAbsPath)
		fmt.Println()
		fmt.Printf("(create file %s, fill it up with zero bytes until out-of-diskspace, then delete the file)", tmpFileAbsPath)
		fmt.Println()
	}

	if ec != 0 {
		libutils.PrintlnStdErr("ERR: Is user missing write privileges?")
		os.Exit(ec)
	}
}
