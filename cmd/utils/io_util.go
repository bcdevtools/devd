package utils

import (
	"fmt"
	"os"
)

// PrintlnStdErr does println to StdErr
func PrintlnStdErr(a ...any) {
	fmt.Fprintln(os.Stderr, a...)
}

// PrintfStdErr does printf to StdErr
func PrintfStdErr(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format, a...)
}
