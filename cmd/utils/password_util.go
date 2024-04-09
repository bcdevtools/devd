package utils

import (
	libutils "github.com/EscanBE/go-lib/utils"
	"github.com/bcdevtools/devd/constants"
	"os"
	"strings"
)

// MustReadPasswordFromEnsuredSecurityFile returns password contained in the file.
//
// This function will panic if one of the following conditions is not met:
//
// 1. File exists, not empty
//
// 2. File permission is X-0-0 (only owner can read. Groups/Others have none permission). Suggested permission is 400
func MustReadPasswordFromEnsuredSecurityFile(passwordFile string) string {
	exists, passwordFPerm, err := IsFileAndExists(passwordFile)
	if err != nil {
		libutils.PrintlnStdErr("ERR: failed to check existence of password file:", err)
		os.Exit(1)
	}
	if !exists {
		libutils.PrintlnStdErr("ERR: password file does not exists:", passwordFile)
		os.Exit(1)
	}

	errPerm := ValidatePasswordFileMode(passwordFPerm)
	if errPerm != nil {
		libutils.PrintfStdErr("ERR: Incorrect permission '%o' of password file: %v\n", passwordFPerm, errPerm)
		libutils.PrintfStdErr("ERR: Suggest setting permission to '%o'\n", constants.REQUIRE_SECRET_FILE_PERMISSION)
		libutils.PrintfStdErr("> chmod %o %s\n", constants.REQUIRE_SECRET_FILE_PERMISSION, passwordFile)
		os.Exit(1)
	}

	bz, err := os.ReadFile(passwordFile)
	if err != nil {
		libutils.PrintfStdErr("ERR: failed to read password file %s: %v\n", passwordFile, err)
		os.Exit(1)
	}

	password := strings.TrimSpace(string(bz))
	if len(password) < 1 {
		libutils.PrintlnStdErr("ERR: password file is empty:", passwordFile)
		os.Exit(1)
	}

	return password
}
