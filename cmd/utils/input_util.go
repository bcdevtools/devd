package utils

import (
	"fmt"
	libutils "github.com/EscanBE/go-lib/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io"
	"os"
	"strings"
)

// ProvidedArgsOrFromPipe will prioritize provided args, if not provided, it will try to read from pipe.
func ProvidedArgsOrFromPipe(providedArgs []string) (outputArgs []string, err error) {
	if len(providedArgs) > 0 {
		outputArgs = providedArgs
	} else {
		outputArgs, err = tryReadPipe()
	}

	return
}

// RequireArgs will exit program if no arg provided.
func RequireArgs(args []string, cmd *cobra.Command) {
	if len(args) == 0 {
		libutils.PrintlnStdErr("ERR: require arg(s)\n")
		_ = cmd.Help()
		os.Exit(1)
	}
}

// RequireExactArgsCount will exit program if number of arguments is not equal to want.
func RequireExactArgsCount(args []string, want int, cmd *cobra.Command) {
	if len(args) != want {
		if want == 0 {
			libutils.PrintlnStdErr("ERR: require no arg\n")
		} else if len(args) == 0 {
			libutils.PrintlnStdErr(fmt.Sprintf("ERR: require %d arg(s)\n", want))
		} else {
			libutils.PrintlnStdErr(fmt.Sprintf("ERR: require %d arg(s), got %d\n", want, len(args)))
		}
		_ = cmd.Help()
		os.Exit(1)
	}
}

func tryReadPipe() (args []string, err error) {
	fi, _ := os.Stdin.Stat()

	if (fi.Mode() & os.ModeCharDevice) == 0 {
		// data from pipe
		var input string
		for {
			n, errScan := fmt.Scan(&input)
			if errScan != nil {
				if errScan == io.EOF {
					break
				}
				if strings.Contains(errScan.Error(), "unexpected newline") {
					break
				}
				err = errors.Wrap(errScan, "failed to read input")
				return
			}
			if n < 1 {
				break
			}
			args = append(args, input)
		}
	}

	return
}
