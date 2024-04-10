package convert

import (
	"encoding/hex"
	"fmt"
	libutils "github.com/EscanBE/go-lib/utils"
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
		Use:     "address [address] [bech32_hrp]",
		Aliases: []string{"a"},
		Short:   "Convert account bech32 address into hex address or vice versa",
		Long: `Convert account bech32 address into hex address or vice versa.
- If the input address is an EVM address, bech32 HRP as the second argument is required.
- If the input address is a bech32 address, will show the EVM address.
- If the input address is a bech32 address, bech32 HRP as the second argument is present, convert.
`,
		Args: cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			defer func() {
				fmt.Println("WARN: DO NOT use this command to convert address across chains with different HD-Path! (eg: Ethermint 60 and Cosmos 118)")
			}()

			if len(args) < 1 {
				fmt.Println("not enough arguments")
				os.Exit(1)
			}

			normalizedInputAddress := strings.TrimSpace(strings.ToLower(args[0]))
			var nextConvertToBech32Hrp string
			if len(args) > 1 {
				nextConvertToBech32Hrp = strings.TrimSpace(strings.ToLower(args[1]))
				nextConvertToBech32Hrp = strings.TrimSuffix(nextConvertToBech32Hrp, "1")
			}

			// conditional processing

			if regexp.MustCompile(`^(0x)?[a-f\d]{40}$`).MatchString(normalizedInputAddress) { // case 1
				// is EVM address

				if len(nextConvertToBech32Hrp) < 1 {
					libutils.PrintlnStdErr("ERR: missing Bech32 HRP as the second argument")
					os.Exit(1)
				}

				evmAddress := common.HexToAddress(normalizedInputAddress)
				bech32Addr, err := bech32.ConvertAndEncode(nextConvertToBech32Hrp, evmAddress.Bytes())
				if err != nil {
					libutils.PrintlnStdErr("ERR: failed to convert EVM address to bech32 address:", err)
					os.Exit(1)
				}
				fmt.Printf("Bech32 (%s): %s\n", nextConvertToBech32Hrp, bech32Addr)
			} else if regexp.MustCompile(`^(0x)?[a-f\d]{64}$`).MatchString(normalizedInputAddress) { // case 1, but bytes addr (eg: interchain account)
				// is bytes address

				if len(nextConvertToBech32Hrp) < 1 {
					libutils.PrintlnStdErr("ERR: missing Bech32 HRP as the second argument")
					os.Exit(1)
				}

				bytesAddress := common.HexToHash(normalizedInputAddress)
				bech32Addr, err := bech32.ConvertAndEncode(nextConvertToBech32Hrp, bytesAddress.Bytes())
				if err != nil {
					libutils.PrintlnStdErr("ERR: failed to convert bytes address to bech32 address:", err)
					os.Exit(1)
				}
				fmt.Printf("Bech32 (%s): %s\n", nextConvertToBech32Hrp, bech32Addr)
			} else { // case 2 + 3
				spl := strings.Split(normalizedInputAddress, "1")
				if len(spl) != 2 || len(spl[0]) < 1 || len(spl[1]) < 1 {
					libutils.PrintlnStdErr("ERR: invalid bech32 address")
					os.Exit(1)
				}

				bz, err := sdk.GetFromBech32(normalizedInputAddress, spl[0])
				if err != nil {
					libutils.PrintlnStdErr("ERR: failed to decode bech32 address:", err)
					os.Exit(1)
				}

				fmt.Printf("EVM address: 0x%s\n", hex.EncodeToString(bz))

				if len(nextConvertToBech32Hrp) > 0 {
					bech32Addr, err := bech32.ConvertAndEncode(nextConvertToBech32Hrp, bz)
					if err != nil {
						libutils.PrintlnStdErr("ERR: failed to convert EVM address to bech32 address:", err)
						os.Exit(1)
					}
					fmt.Printf("Bech32 (%s): %s\n", nextConvertToBech32Hrp, bech32Addr)
				}
			}
		},
	}

	return cmd
}
