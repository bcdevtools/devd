package utils

import (
	"os"
)

// ExitOnErr exit the app with code 1 if error. Or does nothing if no error.
func ExitOnErr(err error, msg string) {
	if err == nil {
		return
	}
	PrintfStdErr("ERR: %s: %v\n", msg, err)
	os.Exit(1)
}

// PanicIfErr raises a panic. If the `err` is nil, this method does nothing
func PanicIfErr(err error, msg string) {
	if err == nil {
		return
	}
	PrintlnStdErr("Exit with error:", msg)
	panic(err)
}
