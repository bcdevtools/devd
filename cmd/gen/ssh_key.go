package gen

import (
	"fmt"
	libutils "github.com/EscanBE/go-lib/utils"
	"github.com/bcdevtools/devd/cmd/types"
	"github.com/bcdevtools/devd/cmd/utils"
	"github.com/spf13/cobra"
	"net/mail"
	"os"
	"path"
	"regexp"
	"runtime"
	"strings"
)

const defaultSshAlgorithm = "ed25519"
const flagNoPass = "no-pass"

// GenerateSshKeypairCommand registers a sub-tree of commands
func GenerateSshKeypairCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ssh-key [key_name] [email]",
		Short: fmt.Sprintf("Generate SSH key pair, using %s algorithm", defaultSshAlgorithm),
		Args:  cobra.ExactArgs(2),
		Run:   generateSshKeyPair,
	}

	cmd.Flags().Bool(flagNoPass, false, "do not prompt for passphrase, empty passphrase will be used")

	return cmd
}

func generateSshKeyPair(cmd *cobra.Command, args []string) {
	utils.EnsureBinaryNameExists("ssh-keygen")

	ctx := types.UnwrapAppContext(cmd.Context())

	operationUserInfo := ctx.GetOperationUserInfo()

	userInfo := ctx.GetWorkingUserInfo()

	keyName := strings.TrimSpace(args[0])
	email := strings.TrimSpace(args[1])

	if !regexp.MustCompile("^[a-z\\d\\-_]+$").MatchString(keyName) {
		libutils.PrintlnStdErr("ERR: Invalid key name", keyName)
		os.Exit(1)
	}
	_, err := mail.ParseAddress(email)
	if err != nil {
		libutils.PrintlnStdErr("ERR: Invalid email address", email)
		os.Exit(1)
	}

	sshDir := path.Join(userInfo.HomeDir, ".ssh")
	if exists, _, _ := utils.IsDirAndExists(sshDir); !exists {
		libutils.PrintlnStdErr("ERR: .ssh directory does not exists at required path", sshDir)
		os.Exit(1)
	}

	keyFilePath := path.Join(sshDir, keyName)
	if exists, _, _ := utils.IsFileAndExists(keyFilePath); exists {
		libutils.PrintlnStdErr("ERR: key file already exists at required path", keyFilePath)
		os.Exit(1)
	}

	fmt.Println("Generating SSH key pair...")
	fmt.Println("File:", keyFilePath)
	fmt.Println("Email:", email)

	sshKeyGenArgs := []string{"-t", defaultSshAlgorithm, "-C", email, "-f", keyFilePath}
	emptyPass, _ := cmd.Flags().GetBool(flagNoPass)
	if emptyPass {
		sshKeyGenArgs = append(sshKeyGenArgs, "-N", "")
	}

	ec := utils.LaunchAppWithDirectStd("ssh-keygen", sshKeyGenArgs, nil)
	if ec != 0 {
		libutils.PrintlnStdErr("ERR: Exited with code", ec)
		os.Exit(ec)
	}

	if operationUserInfo.IsSameUser {
		return
	}

	if !operationUserInfo.OperatingAsSuperUser {
		return
	}

	if !utils.IsLinux() && !utils.IsDarwin() {
		fmt.Println("Skipping updating owner of generated key file on this platform", runtime.GOOS)
		return
	}

	group := utils.GetPseudoDefaultGroupForUser(userInfo.Username)
	newOwner := fmt.Sprintf("%s:%s", userInfo.Username, group)
	fmt.Printf("Updating owner of generated key file to %s\n", newOwner)

	updateOwner := func(filePath string) {
		err = utils.UpdateOwner(filePath, userInfo.Username, group, false)
		if err != nil {
			libutils.PrintlnStdErr("ERR: failed to update owner of", filePath, "to", newOwner, "with error:", err)
			os.Exit(1)
		}
	}

	updateOwner(keyFilePath)
	updateOwner(keyFilePath + ".pub")
}
