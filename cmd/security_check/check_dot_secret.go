package security_check

//goland:noinspection SpellCheckingInspection
import (
	"fmt"
	libutils "github.com/EscanBE/go-lib/utils"
	sectypes "github.com/bcdevtools/devd/cmd/security_check/types"
	"github.com/bcdevtools/devd/cmd/utils"
	"github.com/bcdevtools/devd/constants"
	"os"
	"path"
	"path/filepath"
)

func checkSecretDir(homeDir string) {
	const module = secureCheckModuleDotSecrets

	fmt.Println("Checking .secrets directory in", homeDir)

	secretsDir := path.Join(homeDir, constants.SECRETS_DIR_NAME)
	exists, perm, err := utils.IsDirAndExists(secretsDir)
	if err != nil {
		libutils.PrintlnStdErr("ERR: failed to check if .secrets directory exists:", err)
		securityReport.Add(
			sectypes.NewSecurityRecord(module, true, fmt.Sprintf("Failed to check permission of %s", secretsDir)),
		)
		return
	}

	if !exists {
		return
	}

	if perm != constants.REQUIRE_SECRET_DIR_PERMISSION {
		securityReport.Add(
			sectypes.NewSecurityRecord(
				module, true, fmt.Sprintf(
					"Directory %s has BAD permission %o, must change to %o",
					secretsDir,
					perm,
					constants.REQUIRE_SECRET_DIR_PERMISSION,
				),
			).WithGuide(fmt.Sprintf("sudo chmod %o %s", constants.REQUIRE_SECRET_DIR_PERMISSION, secretsDir)),
		)
	}

	// inner files
	_ = filepath.Walk(secretsDir, func(path string, _ os.FileInfo, _ error) error {
		fi, err := os.Stat(path)
		if err != nil {
			libutils.PrintlnStdErr("ERR: failed to stat file", path, "with error:", err)
			sectypes.NewSecurityRecord(
				module, false, fmt.Sprintf("Failed to stat %s", path),
			)
			return err
		}

		if fi.IsDir() {
			return nil
		}

		perm = fi.Mode().Perm()
		if perm != constants.REQUIRE_SECRET_FILE_PERMISSION {
			sectypes.NewSecurityRecord(
				module, false, fmt.Sprintf(
					"Secret file at %s has BAD permission %o, must change to %o",
					path,
					perm,
					constants.REQUIRE_SECRET_FILE_PERMISSION,
				),
			).WithGuide(fmt.Sprintf("sudo chmod %o %s", constants.REQUIRE_SECRET_FILE_PERMISSION, path))
		}

		return nil
	})
}
