package cmd

import (
	"fmt"
	libutils "github.com/EscanBE/go-lib/utils"
	"github.com/bcdevtools/devd/cmd/utils"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

// verifyToolsCmd represents the verify-tools command, it checks required tools exists
var verifyToolsCmd = &cobra.Command{
	Use:     "verify-tools",
	Short:   "Checks required tools exists",
	Aliases: []string{"vt"},
	Run: func(cmd *cobra.Command, args []string) {
		var anyMandatoryToolsError bool

		defer func() {
			if anyMandatoryToolsError {
				os.Exit(1)
			}
		}()

		fmt.Println("Mandatory tools checking...")

		if !utils.HasBinaryName("rsync") {
			libutils.PrintfStdErr("- \"rsync\" might not exists (sudo apt install -y rsync)\n")
			anyMandatoryToolsError = true
		}

		if !utils.HasBinaryName("sshpass") {
			libutils.PrintfStdErr("- \"sshpass\" might not exists (sudo apt install -y sshpass)\n")
			anyMandatoryToolsError = true
		}

		if !utils.HasBinaryName("aria2c") {
			libutils.PrintfStdErr("- \"aria2c\" might not exists (sudo apt install -y aria2)\n")
			anyMandatoryToolsError = true

			if !utils.HasBinaryName("wget") {
				libutils.PrintlnStdErr("- \"wget\" might not exists")
			}
		}

		if !utils.HasBinaryName("ssh-keygen") {
			libutils.PrintlnStdErr("- \"ssh-keygen\" might not exists")
			anyMandatoryToolsError = true
		}

		if !anyMandatoryToolsError {
			fmt.Println("Successfully checking, all mandatory tools were installed")
		}

		possiblyMissingOptionalTools := make(map[string]string)
		defer func() {
			if len(possiblyMissingOptionalTools) > 0 {
				fmt.Println("(Additionally) The following optional tools might not be installed yet:")
				for toolName, installCommand := range possiblyMissingOptionalTools {
					if len(installCommand) > 0 {
						fmt.Printf(" - %s (%s)\n", toolName, installCommand)
					} else {
						fmt.Println(" -", toolName)
					}
				}
			}
		}()

		checkToolPossiblyExists("telnet", "sudo apt install telnet -y", possiblyMissingOptionalTools)
		checkToolPossiblyExists("htop", "sudo apt install htop -y", possiblyMissingOptionalTools)
		checkToolPossiblyExists("screen", "sudo apt install screen -y", possiblyMissingOptionalTools)
		checkToolPossiblyExists("wget", "sudo apt install wget -y", possiblyMissingOptionalTools)
		checkToolPossiblyExists("jq", "sudo apt install jq -y", possiblyMissingOptionalTools)
		//goland:noinspection SpellCheckingInspection
		checkToolPossiblyExists("lz4", "sudo apt install snapd -y && sudo snap install lz4", possiblyMissingOptionalTools)
		if utils.IsLinux() {
			checkToolPossiblyExists("shred", "sudo apt install shred -y", possiblyMissingOptionalTools)
			checkToolPossiblyExists("srm", "sudo apt install secure-delete -y", possiblyMissingOptionalTools)
		}
	},
}

func checkToolPossiblyExists(tool string, installCommand string, tracker map[string]string) {
	_, err := exec.LookPath(tool)
	if err != nil {
		tracker[tool] = installCommand
	}
}

func init() {
	rootCmd.AddCommand(verifyToolsCmd)
}
