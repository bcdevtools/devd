package convert

import (
	"fmt"
	"github.com/bcdevtools/devd/v2/cmd/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"os"
	"regexp"
	"strings"
)

// GetConvertAddressCmd creates a helper command that convert account bech32 address into hex address or vice versa
func GetConvertAddressCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "address [account] [bech32 hrp]",
		Aliases: []string{"a"},
		Short:   "Convert account bech32 address into hex address or vice versa.",
		Long: `Convert account bech32 address into hex address or vice versa.
- Case 1: if the input address is an EVM address, bech32 HRP as the second argument is required.
- Case 2: if the input address is a bech32 address, bech32 HRP as the second argument is present, convert.
- Case 3: if the input address is a bech32 address, will show the EVM address.
`,
		Args: cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			defer func() {
				utils.PrintlnStdErr("WARN: DO NOT use this command to convert address across chains with different HD-Path! (eg: Ethermint 60 and Cosmos 118)")
			}()

			normalizedInputAddress := strings.TrimSpace(strings.ToLower(args[0]))
			var nextConvertToBech32Hrp string
			if len(args) > 1 {
				nextConvertToBech32Hrp = strings.TrimSpace(strings.ToLower(args[1]))
				nextConvertToBech32Hrp = strings.TrimSuffix(nextConvertToBech32Hrp, "1")
			}

			// conditional processing

			if regexp.MustCompile(`^(0x)?[a-f\d]{40}$`).MatchString(normalizedInputAddress) {
				// case 1, is EVM address

				if len(nextConvertToBech32Hrp) < 1 {
					utils.PrintlnStdErr("ERR: missing Bech32 HRP as the second argument")
					os.Exit(1)
				}

				evmAddress := common.HexToAddress(normalizedInputAddress)
				bech32Addr, err := bech32.ConvertAndEncode(nextConvertToBech32Hrp, evmAddress.Bytes())
				utils.ExitOnErr(err, "failed to convert EVM address to bech32 address")

				fmt.Println(bech32Addr)

				return
			}

			if regexp.MustCompile(`^(0x)?[a-f\d]{64}$`).MatchString(normalizedInputAddress) {
				// case 1, but bytes addr (eg: interchain account)

				if len(nextConvertToBech32Hrp) < 1 {
					utils.PrintlnStdErr("ERR: missing Bech32 HRP as the second argument")
					os.Exit(1)
				}

				bytesAddress := common.HexToHash(normalizedInputAddress)
				bech32Addr, err := bech32.ConvertAndEncode(nextConvertToBech32Hrp, bytesAddress.Bytes())
				utils.ExitOnErr(err, "failed to convert bytes address to bech32 address")

				fmt.Println(bech32Addr)

				return
			}

			// case 2 + 3
			spl := strings.Split(normalizedInputAddress, "1")
			if len(spl) != 2 || len(spl[0]) < 1 || len(spl[1]) < 1 {
				utils.PrintlnStdErr("ERR: invalid bech32 address")
				os.Exit(1)
			}

			bz, err := sdk.GetFromBech32(normalizedInputAddress, spl[0])
			utils.ExitOnErr(err, "failed to decode bech32 address")

			if len(nextConvertToBech32Hrp) > 0 {
				// case 2
				bech32Addr, err := bech32.ConvertAndEncode(nextConvertToBech32Hrp, bz)
				utils.ExitOnErr(
					err,
					fmt.Sprintf(
						"failed to convert bech32 address %s into bech32 address with HRP %s",
						normalizedInputAddress,
						nextConvertToBech32Hrp,
					),
				)

				fmt.Println(bech32Addr)
				return
			}

			// case 3
			fmt.Printf("0x%x\n", bz)
		},
	}

	return cmd
}
