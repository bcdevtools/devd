package query

import (
	"context"
	"fmt"
	libutils "github.com/EscanBE/go-lib/utils"
	"github.com/bcdevtools/devd/cmd/utils"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"math/big"
	"os"
)

func GetQueryErc20Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "erc20 [contract_address] [?account_address]",
		Short: "Get ERC-20 token information. If account address is provided, it will query the balance of the account (bech32 is accepted).",
		Args:  cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			ethClient8545, _ := mustGetEthClient(cmd, true)

			evmAddrs, err := getEvmAddressFromAnyFormatAddress(args...)
			if err != nil {
				libutils.PrintlnStdErr("ERR:", err)
				return
			}

			var contractAddr, accountAddr common.Address

			contractAddr = evmAddrs[0]
			if len(evmAddrs) > 1 {
				accountAddr = evmAddrs[1]
			}

			fmt.Println("Getting contract symbol...")

			bz, err := ethClient8545.CallContract(context.Background(), ethereum.CallMsg{
				To:   &contractAddr,
				Data: []byte{0x95, 0xd8, 0x9b, 0x41}, // symbol()
			}, nil)
			if err != nil {
				libutils.PrintlnStdErr("ERR: failed to get contract symbol:", err)
				os.Exit(1)
			}

			contractSymbol, err := utils.AbiDecodeString(bz)
			if err != nil {
				libutils.PrintlnStdErr("ERR: failed to decode contract symbol:", err)
				os.Exit(1)
			}

			fmt.Println("Getting contract decimals...")

			bz, err = ethClient8545.CallContract(context.Background(), ethereum.CallMsg{
				To:   &contractAddr,
				Data: []byte{0x31, 0x3c, 0xe5, 0x67}, // decimals()
			}, nil)
			if err != nil {
				libutils.PrintlnStdErr("ERR: failed to get contract decimals:", err)
				os.Exit(1)
			}

			contractDecimals := new(big.Int).SetBytes(bz)

			var accountBalance *big.Int
			if accountAddr != (common.Address{}) {
				fmt.Println("Getting account balance...")

				bz, err = ethClient8545.CallContract(context.Background(), ethereum.CallMsg{
					To:   &contractAddr,
					Data: append([]byte{0x70, 0xa0, 0x82, 0x31}, common.BytesToHash(accountAddr.Bytes()).Bytes()...), // balanceOf(address)
				}, nil)
				if err != nil {
					libutils.PrintlnStdErr("ERR: failed to get account token balance:", err)
					os.Exit(1)
				}

				accountBalance = new(big.Int).SetBytes(bz)
			}

			fmt.Println("Contract Symbol:", contractSymbol)
			fmt.Println("Contract Decimals:", contractDecimals.Uint64())
			if accountBalance != nil {
				decimals := contractDecimals.Uint64()
				if decimals == 0 {
					fmt.Println("Account token balance:", accountBalance, contractSymbol)
				} else {
					fmt.Println("Account token balance:")
					fmt.Println(" - Raw:", accountBalance)
					fmt.Println(" - Display:")
					pow := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
					fmt.Println("  + High:", new(big.Int).Div(accountBalance, pow), contractSymbol)
					fmt.Println("  + Low:", new(big.Int).Mod(accountBalance, pow))
				}
			}
		},
	}

	cmd.Flags().String(flagRpc, "", flagEvmRpcDesc)
	cmd.Flags().String("host", "", fmt.Sprintf("deprecated flag, use '--%s' instead", flagRpc))

	return cmd
}
