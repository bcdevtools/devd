package utils

import (
	"fmt"
	"os"
)

// PrintlnStdErr does println to StdErr
func PrintlnStdErr(a ...any) {
	_, _ = fmt.Fprintln(os.Stderr, a...)
}

// PrintfStdErr does printf to StdErr
func PrintfStdErr(format string, a ...any) {
	_, _ = fmt.Fprintf(os.Stderr, format, a...)
}
