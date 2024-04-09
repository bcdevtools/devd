package files

import (
	"fmt"
	libutils "github.com/EscanBE/go-lib/utils"
	"github.com/bcdevtools/devd/cmd/types"
	"github.com/bcdevtools/devd/cmd/utils"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path"
	"regexp"
	"sort"
	"strings"
	"time"
)

const (
	flagConfirmParanoid = "i-am-paranoid"
)

// RemoveParanoidCommands registers a sub-tree of commands
func RemoveParanoidCommands() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "paranoid-rm [/dev/sdX]",
		Short: "Secure remove files and directory (requires super user)",
		Args:  cobra.ExactArgs(1),
		Run:   removeFilesInParanoidMode,
	}

	cmd.Flags().Bool(flagConfirmDelete, false, "REQUIRED one of the two flags used to prevent incident")
	cmd.Flags().Bool(flagConfirmParanoid, false, "REQUIRED one of the two flags used to prevent incident")

	return cmd
}

func removeFilesInParanoidMode(cmd *cobra.Command, args []string) {
	confirmDelete := cmd.Flag(flagConfirmDelete).Changed
	confirmParanoid := cmd.Flag(flagConfirmParanoid).Changed
	if !confirmDelete || !confirmParanoid {
		libutils.PrintfStdErr("ERR: For the purpose of preventing incident from incorrect typing/usage, both flags '--%s' and '--%s' are required, indicate that you know what this is and you want to do this because of your mental illness (paranoid, psychotic).\n", flagConfirmDelete, flagConfirmParanoid)
		os.Exit(1)
	}

	ctx := types.UnwrapAppContext(cmd.Context())

	operationUserInfo := ctx.GetOperationUserInfo()
	operationUserInfo.RequireSuperUser()
	operationUserInfo.RequireOperatingAsSuperUser()

	targetDisk := args[0]
	targetDisk = strings.TrimSuffix(targetDisk, "/")

	const diskPattern = `(\/dev\/((sd[a-z\d]+)|(nvme[a-z\d]+)|(disk[a-z\d]+))\/?)|(\/mnt\/[a-z\d]+)`
	var regexDiskPattern = regexp.MustCompile(diskPattern)
	if !regexDiskPattern.MatchString(targetDisk) {
		libutils.PrintlnStdErr("ERR: Not accepted input", targetDisk)
		libutils.PrintlnStdErr(`ERR: Accepted patterns:
- /dev/sdX
- /dev/nvmeX
- /dev/diskX
- /mnt/X`)
		//goland:noinspection SpellCheckingInspection
		bz, err := exec.Command("/bin/bash", "-c", "df | awk '{ print $1; }' | grep -v -e 'Filesystem' -e map -e devfs -e tmpfs").Output()
		if err == nil {
			var candidates []string
			spl := strings.Split(string(bz), "\n")
			for _, line := range spl {
				normalized := strings.TrimSpace(line)
				if normalized == "" {
					continue
				}

				if !regexDiskPattern.MatchString(normalized) {
					continue
				}

				candidates = append(candidates, normalized)
			}

			if len(candidates) > 0 {
				sort.Strings(candidates)
				libutils.PrintlnStdErr("HINT: Can be one of the following disks:", strings.Join(candidates, ", "))
			} else {
				libutils.PrintlnStdErr("ERR: failed to hint disks")
			}
		}
		os.Exit(1)
	}

	utils.WarnIfNotRunningUnderScreenSession("paranoid-wipe-disk")

	timeSleep := 2 * time.Minute
	actionTime := time.Now().Add(timeSleep)
	for time.Now().Before(actionTime) {
		remainingTime := actionTime.Sub(time.Now()).Round(time.Second)
		fmt.Printf("WARN: The psychotic action will be executed after %v\n", remainingTime)
		if remainingTime >= time.Minute {
			time.Sleep(15 * time.Second)
		} else if remainingTime >= 30*time.Second {
			time.Sleep(5 * time.Second)
		} else {
			time.Sleep(time.Second)
		}
	}

	fmt.Println("Begin execution at", time.Now().Format("2006-01-02 15:04:05"))
	var exitCode int
	defer func() {
		fmt.Println("Finished execution at", time.Now().Format("2006-01-02 15:04:05"))
		if exitCode != 0 {
			libutils.PrintlnStdErr("ERR: Exited with error!")
			os.Exit(exitCode)
		}
	}()

	if strings.HasPrefix(targetDisk, "/dev/") {
		_ = utils.LaunchAppWithDirectStd(
			"sudo",
			[]string{"dd", "if=/dev/zero", fmt.Sprintf("of=%s", targetDisk), "bs=1M", "status=progress"},
			nil,
		)
	} else if strings.HasPrefix(targetDisk, "/mnt/") {
		tmpFile := path.Join(targetDisk, "tmp_file_paranoid")
		fmt.Println("INF: temp file for filling up data:", tmpFile)

		_ = utils.LaunchAppWithDirectStd(
			"sudo",
			[]string{"dd", "if=/dev/zero", fmt.Sprintf("of=%s", tmpFile), "bs=1M", "count=999999999999", "status=progress"},
			nil,
		)

		err := os.Remove(tmpFile)
		if err != nil {
			exitCode = 1
			libutils.PrintlnStdErr("ERR: failed to remove temp file", tmpFile)
		}
	} else {
		libutils.PrintlnStdErr("ERR: not implemented for pattern of", targetDisk)
		exitCode = 2
	}
}
