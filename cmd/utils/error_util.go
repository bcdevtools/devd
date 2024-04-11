package utils

import (
	libutils "github.com/EscanBE/go-lib/utils"
	"os"
)

// ExitOnErr exit the app with code 1 if error. Or does nothing if no error.
func ExitOnErr(err error, msg string) {
	if err == nil {
		return
	}
	libutils.PrintfStdErr("ERR: %s: %v\n", msg, err)
	os.Exit(1)
}
