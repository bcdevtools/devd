package convert

import (
	"encoding/hex"
	"fmt"
	"github.com/bcdevtools/devd/v2/cmd/utils"
	"github.com/spf13/cobra"
	"regexp"
	"strings"
)

func GetConvertAbiStringCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "abi_string [hex_or_text]",
		Short: `Convert ABI encoded hex to string or vice versa.`,
		Long: `Convert ABI encoded hex to string or vice versa.
Support pipe.`,
		Run: func(cmd *cobra.Command, args []string) {
			args, err := utils.ProvidedArgsOrFromPipe(args)
			utils.ExitOnErr(err, "failed to get args from pipe")
			utils.RequireArgs(args, cmd)

			if len(args) == 1 && len(args[0]) >= 192 {
				abiToString := func(input string) (nonEmpty bool, decodeErr error) {
					bz, decodeErr := hex.DecodeString(input)
					if decodeErr != nil {
						return false, decodeErr
					}
					str, decodeErr := utils.AbiDecodeString(bz)
					if decodeErr != nil {
						return false, decodeErr
					}
					if len(str) == 0 {
						return false, nil
					}
					fmt.Println(str)
					return true, nil
				}
				if regexp.MustCompile(`^0x[a-f\d]+$`).MatchString(args[0]) && (len(args[0])-2 /*exclude 0x*/)%64 == 0 {
					if nonEmpty, err := abiToString(args[0][2:]); err == nil && nonEmpty {
						return
					}
				} else if regexp.MustCompile(`^[a-f\d]+$`).MatchString(args[0]) && (len(args[0]))%64 == 0 {
					if nonEmpty, err := abiToString(args[0]); err == nil && nonEmpty {
						return
					}
				} else if regexp.MustCompile(`^0x08c379a0[a-f\d]+$`).MatchString(args[0]) && (len(args[0])-10 /*exclude sig of error*/)%64 == 0 {
					if nonEmpty, err := abiToString(args[0][10:]); err == nil && nonEmpty {
						return
					}
				}
			}

			bz, err := utils.AbiEncodeString(strings.Join(args, " "))
			utils.ExitOnErr(err, "failed to encode string to ABI hex")
			fmt.Println("0x" + hex.EncodeToString(bz))
		},
	}

	return cmd
}
