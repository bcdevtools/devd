package convert

import (
	"encoding/hex"
	"fmt"
	libutils "github.com/EscanBE/go-lib/utils"
	"github.com/bcdevtools/devd/cmd/utils"
	"github.com/spf13/cobra"
	"regexp"
	"strings"
)

func GetConvertAbiStringCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "abi_string [hex_or_text]",
		Short: `Convert ABI encoded hex to string or vice versa.`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 1 && len(args[0]) >= 192 {
				abiToString := func(input string) {
					bz, err := hex.DecodeString(input)
					libutils.ExitIfErr(err, "failed to decode hex string: "+input)
					str, err := utils.AbiDecodeString(bz)
					libutils.ExitIfErr(err, "failed to decode ABI string: "+input)
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
			libutils.ExitIfErr(err, "failed to encode string to ABI hex")
			fmt.Println("0x" + hex.EncodeToString(bz))
		},
	}

	return cmd
}
