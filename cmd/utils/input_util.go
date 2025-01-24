package utils

import (
	"fmt"
	"math/big"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

// ProvidedArgsOrFromPipe will prioritize provided args, if not provided, it will try to read from pipe.
func ProvidedArgsOrFromPipe(providedArgs []string) (outputArgs []string, err error) {
	if len(providedArgs) > 0 {
		outputArgs = providedArgs
	} else {
		var inputFromPipe string
		inputFromPipe, err = tryReadPipe()
		if err != nil {
			return
		}

		outputArgs = []string{inputFromPipe}
	}

	return
}

// RequireArgs will exit program if no arg provided.
func RequireArgs(args []string, cmd *cobra.Command) {
	if len(args) == 0 {
		PrintlnStdErr("ERR: require arg(s)\n")
		_ = cmd.Help()
		os.Exit(1)
	}
}

// RequireExactArgsCount will exit program if number of arguments is not equal to want.
func RequireExactArgsCount(args []string, want int, cmd *cobra.Command) {
	if len(args) != want {
		if want == 0 {
			PrintlnStdErr("ERR: require no arg\n")
		} else if len(args) == 0 {
			PrintlnStdErr(fmt.Sprintf("ERR: require %d arg(s)\n", want))
		} else {
			PrintlnStdErr(fmt.Sprintf("ERR: require %d arg(s), got %d\n", want, len(args)))
		}
		_ = cmd.Help()
		os.Exit(1)
	}
}

func tryReadPipe() (dataFromPipe string, err error) {
	fi, _ := os.Stdin.Stat()

	if (fi.Mode() & os.ModeCharDevice) == 0 {
		// data from pipe

		sb := strings.Builder{}
		buffer := make([]byte, 1024)

		for {
			n, _ := os.Stdin.Read(buffer)
			if n == 0 {
				break
			}

			sb.Write(buffer[:n])
		}

		if sb.Len() > 0 {
			dataFromPipe = sb.String()
			if dataFromPipe[len(dataFromPipe)-1] == '\n' {
				dataFromPipe = dataFromPipe[:len(dataFromPipe)-1]
			}
		}
	}

	return
}

func ReadShortInt(input string) (out *big.Int, err error) {
	normalizedInput := strings.ToLower(strings.TrimSpace(input))

	negative := strings.HasPrefix(normalizedInput, "-")
	defer func() {
		if negative {
			out = out.Neg(out)
		}
	}()

	positiveInput := strings.TrimPrefix(normalizedInput, "-")

	if regexp.MustCompile(`^\d+$`).MatchString(positiveInput) { // general format
		bi, ok := new(big.Int).SetString(positiveInput, 10)
		if !ok {
			err = fmt.Errorf("unexpected error, cannot read integer from %s", normalizedInput)
			return
		}
		out = bi
		return
	}

	if regexp.MustCompile(`^\d+e\d+$`).MatchString(positiveInput) { // scientific notation
		parts := strings.Split(positiveInput, "e")

		baseMultiplier, ok := new(big.Int).SetString(parts[0], 10)
		if !ok {
			err = fmt.Errorf("unexpected error, cannot read integer %s", parts[0])
			return
		}

		baseExp, ok := new(big.Int).SetString(parts[1], 10)
		if !ok {
			err = fmt.Errorf("unexpected error, cannot read integer %s", parts[1])
			return
		}
		exp := new(big.Int).Exp(big.NewInt(10), baseExp, nil)

		out = new(big.Int).Mul(baseMultiplier, exp)
		return
	}

	if regexp.MustCompile(`^\d+(\.\d+)?[kmb]+$`).MatchString(positiveInput) {
		var base *big.Float
		suffixMultiplier := big.NewInt(1)
		for true {
			if strings.HasSuffix(positiveInput, "k") {
				suffixMultiplier = new(big.Int).Mul(suffixMultiplier, big.NewInt(1_000))
				positiveInput = strings.TrimSuffix(positiveInput, "k")
				continue
			}

			if strings.HasSuffix(positiveInput, "m") {
				suffixMultiplier = new(big.Int).Mul(suffixMultiplier, big.NewInt(1_000_000))
				positiveInput = strings.TrimSuffix(positiveInput, "m")
				continue
			}

			if strings.HasSuffix(positiveInput, "b") {
				suffixMultiplier = new(big.Int).Mul(suffixMultiplier, big.NewInt(1_000_000_000))
				positiveInput = strings.TrimSuffix(positiveInput, "b")
				continue
			}

			var ok bool
			base, ok = new(big.Float).SetString(positiveInput)
			if !ok {
				err = fmt.Errorf("unexpected error, cannot read number from %s", positiveInput)
				return
			}

			break
		}

		result := new(big.Float).Mul(base, new(big.Float).SetInt(suffixMultiplier))
		out, _ = result.Int(nil)
		return
	}

	err = fmt.Errorf("unexpected error, unrecorgnized format: %s", normalizedInput)
	return
}
