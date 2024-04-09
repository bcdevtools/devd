package convert

import (
	"fmt"
	"github.com/spf13/cobra"
	"math/big"
	"regexp"
	"strings"
)

func GetConvertHexadecimalToDecimalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "hex_2_dec [hex or dec]",
		Aliases: []string{"h2d"},
		Short:   `Convert hex to dec or vice versa.`,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			input := strings.ToLower(args[0])
			if regexp.MustCompile(`^0x[a-f\d]+$`).MatchString(input) {
				fmt.Println("# Hex to Dec:")
				bi, ok := new(big.Int).SetString(input[2:], 16)
				if !ok {
					panic("failed to convert hexadecimal to decimal")
				}
				fmt.Println(bi)
			} else if regexp.MustCompile(`^\d+$`).MatchString(input) {
				// can be hex or dec

				fmt.Println("# Hex to Dec:")
				bi, ok := new(big.Int).SetString(input, 16)
				if !ok {
					panic("failed to convert hexadecimal to decimal")
				}
				fmt.Println(bi)

				fmt.Println("\n# Dec to Hex:")

				bi, ok = new(big.Int).SetString(input, 10)
				if !ok {
					panic("failed to convert string to decimal")
				}
				fmt.Println(fmt.Sprintf("0x%x", bi))
			} else {
				// is hex without 0x
				fmt.Println("# Hex to Dec:")
				bi, ok := new(big.Int).SetString(input, 16)
				if !ok {
					panic("failed to convert hexadecimal to decimal")
				}
				fmt.Println(bi)
			}
		},
	}

	return cmd
}
