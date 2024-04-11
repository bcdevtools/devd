package convert

import (
	"encoding/hex"
	"fmt"
	"github.com/bcdevtools/devd/cmd/utils"
	"github.com/spf13/cobra"
	"regexp"
	"strings"
)

func GetConvertAbiStringCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "abi_string [hex_or_text]",
		Short: `Convert ABI encoded hex to string or vice versa.
Support pipe.`,
		Run: func(cmd *cobra.Command, args []string) {
			args, err := utils.ProvidedArgsOrFromPipe(args)
			utils.ExitOnErr(err, "failed to get args from pipe")
			utils.RequireArgs(args, cmd)

			if len(args) == 1 && len(args[0]) >= 192 {
				abiToString := func(input string) {
					bz, decodeErr := hex.DecodeString(input)
					utils.ExitOnErr(decodeErr, "failed to decode hex string: "+input)
					str, decodeErr := utils.AbiDecodeString(bz)
					utils.ExitOnErr(decodeErr, "failed to decode ABI string: "+input)
					fmt.Println("\n# ABI encoded hex string to string:")
					fmt.Println(str)
				}
				if regexp.MustCompile(`^0x[a-fA-F\d]+$`).MatchString(args[0]) && (len(args[0])-2 /*exclude 0x*/)%64 == 0 {
					abiToString(args[0][2:])
				} else if regexp.MustCompile(`^[a-fA-F\d]+$`).MatchString(args[0]) && (len(args[0]))%64 == 0 {
					abiToString(args[0])
				}
			}

			fmt.Println("\n# String to ABI encoded hex:")
			bz, err := utils.AbiEncodeString(strings.Join(args, " "))
			utils.ExitOnErr(err, "failed to encode string to ABI hex")
			fmt.Println("0x" + hex.EncodeToString(bz))
		},
	}

	return cmd
}
