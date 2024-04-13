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

func GetQueryBalanceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "balance [account] [?optional ERC-20 addrs..]",
		Aliases: []string{"b"},
		Short:   "Get ERC-20 token information. If account address is provided, it will query the balance of the account (bech32 is accepted).",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			evmAddrs, err := getEvmAddressFromAnyFormatAddress(args...)
			if err != nil {
				utils.PrintlnStdErr("ERR:", err)
				return
			}

			//fetchErc20ModuleAndVfbc := cmd.Flags().Changed(flagErc20)

			ethClient8545, _ := mustGetEthClient(cmd, false)
			var bz []byte

			contextHeight := readContextHeightFromFlag(cmd)

			accountAddr := evmAddrs[0]
			fmt.Println("Account", accountAddr)

			printRow := func(colType, colContract, colSymbol, colBalance, colRaw, colDecimals, colHigh, colLow, extra string) {
				fmt.Printf("%-6s | %42s | %10s | %28s | %27s | %8s | %9s | %18s | %-1s\n", colType, colContract, colSymbol, colBalance, colRaw, colDecimals, colHigh, colLow, extra)
			}

			printRow("Type", "Contract", "Symbol", "Balance", "Raw", "Decimals", "High", "Low", "Extra")

			nativeBalance, err := ethClient8545.BalanceAt(context.Background(), accountAddr, contextHeight)
			utils.ExitOnErr(err, "failed to get account balance")

			display, high, low, err := utils.ConvertNumberIntoDisplayWithExponent(nativeBalance, 18)
			utils.ExitOnErr(err, "failed to convert number into display with exponent")

			printRow("Native", "-", "(native)", display, nativeBalance.String(), "18", high.String(), low.String(), "")

			for i := 1; i < len(evmAddrs); i++ {
				contractAddr := evmAddrs[i]

				bz, err = ethClient8545.CallContract(context.Background(), ethereum.CallMsg{
					To:   &contractAddr,
					Data: []byte{0x95, 0xd8, 0x9b, 0x41}, // symbol()
				}, contextHeight)
				if err != nil {
					utils.PrintlnStdErr("ERR: failed to get symbol for contract", contractAddr, ":", err)
					continue
				}

				contractSymbol, err := utils.AbiDecodeString(bz)
				if err != nil {
					utils.PrintlnStdErr("ERR: failed to decode symbol for contract", contractAddr, ":", err)
					continue
				}

				bz, err = ethClient8545.CallContract(context.Background(), ethereum.CallMsg{
					To:   &contractAddr,
					Data: []byte{0x31, 0x3c, 0xe5, 0x67}, // decimals()
				}, contextHeight)
				if err != nil {
					utils.PrintlnStdErr("ERR: failed to get decimals for contract", contractAddr, ":", err)
					continue
				}

				contractDecimals := new(big.Int).SetBytes(bz)

				var tokenBalance *big.Int
				bz, err = ethClient8545.CallContract(context.Background(), ethereum.CallMsg{
					To:   &contractAddr,
					Data: append([]byte{0x70, 0xa0, 0x82, 0x31}, common.BytesToHash(accountAddr.Bytes()).Bytes()...), // balanceOf(address)
				}, contextHeight)
				if err != nil {
					utils.PrintlnStdErr("ERR: failed to get contract token", contractAddr, "balance for", accountAddr, ":", err)
					continue
				}

				tokenBalance = new(big.Int).SetBytes(bz)

				display, high, low, err := utils.ConvertNumberIntoDisplayWithExponent(tokenBalance, int(contractDecimals.Int64()))
				utils.ExitOnErr(err, "failed to convert number into display with exponent")

				printRow("Input", contractAddr.String(), contractSymbol, display, tokenBalance.String(), contractDecimals.String(), high.String(), low.String(), "")
			}
		},
	}

	cmd.Flags().String(flagRpc, "", flagEvmRpcDesc)
	cmd.Flags().Int64(flagHeight, 0, "query balance at specific height")
	cmd.Flags().String(flagErc20, "", "query balance of ERC-20 contracts of `x/erc20` module and virtual frontier bank contracts")

	return cmd
}
