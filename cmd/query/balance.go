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
)

func GetQueryBalanceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "balance [account_address] [?optional_erc20_contracts...]",
		Aliases: []string{"b"},
		Short:   "Get ERC-20 token information. If account address is provided, it will query the balance of the account (bech32 is accepted).",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ethClient8545, _ := mustGetEthClient(cmd, false)
			var bz []byte

			evmAddrs, err := getEvmAddressFromAnyFormatAddress(args...)
			if err != nil {
				libutils.PrintlnStdErr("ERR:", err)
				return
			}

			contextHeight := readContextHeightFromFlag(cmd)

			accountAddr := evmAddrs[0]
			fmt.Println("Account", accountAddr)

			if len(evmAddrs) == 1 {
				nativeBalance, err := ethClient8545.BalanceAt(context.Background(), accountAddr, contextHeight)
				utils.ExitOnErr(err, "failed to get account balance")

				display, high, low, err := utils.ConvertNumberIntoDisplayWithExponent(nativeBalance, 18)
				utils.ExitOnErr(err, "failed to convert number into display with exponent")
				fmt.Println("> Native balance:")
				fmt.Println(" - Decimal:", 18)
				fmt.Println(" - Raw:", nativeBalance)
				fmt.Println(" - Display:", display)
				fmt.Println("  + High:", high)
				fmt.Println("  + Low:", low)
			}

			for i := 1; i < len(evmAddrs); i++ {
				contractAddr := evmAddrs[i]

				bz, err = ethClient8545.CallContract(context.Background(), ethereum.CallMsg{
					To:   &contractAddr,
					Data: []byte{0x95, 0xd8, 0x9b, 0x41}, // symbol()
				}, contextHeight)
				if err != nil {
					libutils.PrintlnStdErr("ERR: failed to get symbol for contract", contractAddr, ":", err)
					continue
				}

				contractSymbol, err := utils.AbiDecodeString(bz)
				if err != nil {
					libutils.PrintlnStdErr("ERR: failed to decode symbol for contract", contractAddr, ":", err)
					continue
				}

				bz, err = ethClient8545.CallContract(context.Background(), ethereum.CallMsg{
					To:   &contractAddr,
					Data: []byte{0x31, 0x3c, 0xe5, 0x67}, // decimals()
				}, contextHeight)
				if err != nil {
					libutils.PrintlnStdErr("ERR: failed to get decimals for contract", contractAddr, ":", err)
					continue
				}

				contractDecimals := new(big.Int).SetBytes(bz)

				var tokenBalance *big.Int
				bz, err = ethClient8545.CallContract(context.Background(), ethereum.CallMsg{
					To:   &contractAddr,
					Data: append([]byte{0x70, 0xa0, 0x82, 0x31}, common.BytesToHash(accountAddr.Bytes()).Bytes()...), // balanceOf(address)
				}, contextHeight)
				if err != nil {
					libutils.PrintlnStdErr("ERR: failed to get contract token", contractAddr, "balance for", accountAddr, ":", err)
					continue
				}

				tokenBalance = new(big.Int).SetBytes(bz)

				display, high, low, err := utils.ConvertNumberIntoDisplayWithExponent(tokenBalance, int(contractDecimals.Int64()))
				utils.ExitOnErr(err, "failed to convert number into display with exponent")

				fmt.Printf("> ERC-20 %s\n", contractAddr)
				fmt.Println(" - Symbol:", contractSymbol)
				fmt.Println(" - Decimals:", contractDecimals.Uint64())
				fmt.Println(" - Raw:", tokenBalance)
				fmt.Println(" - Display:", display, contractSymbol)
				fmt.Println("  + High:", high)
				fmt.Println("  + Low:", low)
			}
		},
	}

	cmd.Flags().String(flagRpc, "", flagEvmRpcDesc)
	cmd.Flags().Int64(flagHeight, 0, "query balance at specific height")

	return cmd
}
