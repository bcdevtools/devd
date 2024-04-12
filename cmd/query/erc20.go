package query

import (
	"context"
	"fmt"
	"github.com/bcdevtools/devd/v2/cmd/utils"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"math/big"
)

func GetQueryErc20Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "erc20 [contract address] [?account address]",
		Short: "Get ERC-20 token information. If account address is provided, it will query the balance of the account (bech32 is accepted).",
		Args:  cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			ethClient8545, _ := mustGetEthClient(cmd, true)

			evmAddrs, err := getEvmAddressFromAnyFormatAddress(args...)
			utils.ExitOnErr(err, "failed to get evm address from input")

			contextHeight := readContextHeightFromFlag(cmd)

			var contractAddr, accountAddr common.Address

			contractAddr = evmAddrs[0]
			if len(evmAddrs) > 1 {
				accountAddr = evmAddrs[1]
			}

			bz, err := ethClient8545.CallContract(context.Background(), ethereum.CallMsg{
				To:   &contractAddr,
				Data: []byte{0x95, 0xd8, 0x9b, 0x41}, // symbol()
			}, contextHeight)
			utils.ExitOnErr(err, "failed to get contract symbol")

			contractSymbol, err := utils.AbiDecodeString(bz)
			utils.ExitOnErr(err, "failed to decode contract symbol")

			bz, err = ethClient8545.CallContract(context.Background(), ethereum.CallMsg{
				To:   &contractAddr,
				Data: []byte{0x31, 0x3c, 0xe5, 0x67}, // decimals()
			}, contextHeight)
			utils.ExitOnErr(err, "failed to get contract decimals")

			contractDecimals := new(big.Int).SetBytes(bz)

			var accountBalance *big.Int
			if accountAddr != (common.Address{}) {
				bz, err = ethClient8545.CallContract(context.Background(), ethereum.CallMsg{
					To:   &contractAddr,
					Data: append([]byte{0x70, 0xa0, 0x82, 0x31}, common.BytesToHash(accountAddr.Bytes()).Bytes()...), // balanceOf(address)
				}, contextHeight)
				utils.ExitOnErr(err, "failed to get account token balance")

				accountBalance = new(big.Int).SetBytes(bz)
			}

			fmt.Println("Contract Symbol:", contractSymbol)
			fmt.Println("Contract Decimals:", contractDecimals.Uint64())
			if accountBalance != nil {
				decimals := contractDecimals.Uint64()
				if decimals == 0 {
					fmt.Println("Account token balance:", accountBalance, contractSymbol)
				} else {
					display, high, low, err := utils.ConvertNumberIntoDisplayWithExponent(accountBalance, int(decimals))
					utils.ExitOnErr(err, "failed to convert number into display with exponent")
					fmt.Println("Account token balance:")
					fmt.Println(" - Raw:", accountBalance)
					fmt.Println(" - Display:", display, contractSymbol)
					fmt.Println("  + High:", high)
					fmt.Println("  + Low:", low)
				}
			}
		},
	}

	cmd.Flags().String(flagRpc, "", flagEvmRpcDesc)
	cmd.Flags().String("host", "", fmt.Sprintf("deprecated flag, use '--%s' instead", flagRpc))
	cmd.Flags().Int64(flagHeight, 0, "query balance at specific height")

	return cmd
}
