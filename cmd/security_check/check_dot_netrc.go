package security_check

//goland:noinspection SpellCheckingInspection
import (
	"fmt"
	libutils "github.com/EscanBE/go-lib/utils"
	sectypes "github.com/bcdevtools/devd/cmd/security_check/types"
	"github.com/bcdevtools/devd/cmd/utils"
	"github.com/bcdevtools/devd/constants"
	"path"
)

func checkDotNetrc(homeDir string) {
	const module = secureCheckModuleDotNetrc

	fmt.Println("Checking .netrc file in", homeDir)

	// authorized_keys
	netrcFilePath := path.Join(homeDir, ".netrc")
	exists, perm, err := utils.IsFileAndExists(netrcFilePath)
	if err != nil {
		libutils.PrintlnStdErr("ERR: failed to check if .netrc file exists:", err)
		securityReport.Add(
			sectypes.NewSecurityRecord(module, false, fmt.Sprintf("Failed to check existence of %s", netrcFilePath)),
		)
		return
	}

	if !exists {
		return
	}

	if perm == constants.REQUIRE_NETRC_FILE_PERMISSION {
		return
	}

	securityReport.Add(
		sectypes.NewSecurityRecord(
			module, true, fmt.Sprintf(
				"File %s has BAD permission %o, must change to %o",
				netrcFilePath,
				perm,
				constants.REQUIRE_NETRC_FILE_PERMISSION,
			),
		).WithGuide(fmt.Sprintf("sudo chmod %o %s", constants.REQUIRE_NETRC_FILE_PERMISSION, netrcFilePath)),
	)
}
