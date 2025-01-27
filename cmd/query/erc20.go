package query

import (
	"context"
	"fmt"
	"math/big"

	"github.com/bcdevtools/devd/v3/cmd/flags"

	"github.com/bcdevtools/devd/v3/cmd/utils"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

func GetQueryErc20Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "erc20 [contract address] [?account address]",
		Short: "Get ERC-20 token information. Optionally query the balance of an account.",
		Long: `Get ERC-20 token information. If account address is provided, it will query the balance of the account.
Support bech32 address format`,
		Args: cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			ethClient8545, _ := flags.MustGetEthClient(cmd)

			evmAddrs, err := utils.GetEvmAddressFromAnyFormatAddress(args...)
			utils.ExitOnErr(err, "failed to get evm address from input")

			contextHeight, err := flags.ReadFlagBlockNumberOrNil(cmd, flags.FlagHeight)
			utils.ExitOnErr(err, "failed to parse block number")

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

			contractDecimals := new(big.Int).SetBytes(bz).Uint64()

			var totalSupply *big.Int
			bz, err = ethClient8545.CallContract(context.Background(), ethereum.CallMsg{
				To:   &contractAddr,
				Data: []byte{0x18, 0x16, 0x0d, 0xdd}, // totalSupply()
			}, contextHeight)
			if err == nil && len(bz) > 0 {
				totalSupply = new(big.Int).SetBytes(bz)
				if totalSupply.Sign() != 1 {
					totalSupply = nil
				}
			}

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
			fmt.Println("Contract Decimals:", contractDecimals)
			if totalSupply != nil {
				display, _, _, err := utils.ConvertNumberIntoDisplayWithExponent(totalSupply, int(contractDecimals))
				utils.ExitOnErr(err, "failed to convert number into display with exponent")
				fmt.Println("Total Supply:", display, contractSymbol)
			}
			if accountBalance != nil {
				display, high, low, err := utils.ConvertNumberIntoDisplayWithExponent(accountBalance, int(contractDecimals))
				utils.ExitOnErr(err, "failed to convert number into display with exponent")
				fmt.Println("Account token balance:")
				fmt.Println(" - Raw:", accountBalance)
				fmt.Println(" - Display:", display, contractSymbol)
				fmt.Println("  + High:", high)
				fmt.Println("  + Low:", low)
			}
		},
	}

	cmd.Flags().String(flags.FlagEvmRpc, "", flags.FlagEvmRpcDesc)
	cmd.Flags().Int64(flags.FlagHeight, 0, "query balance at specific height")

	return cmd
}
