package query

import (
	"context"
	"fmt"
	"math/big"

	"github.com/bcdevtools/devd/v3/cmd/flags"

	"github.com/bcdevtools/devd/v3/cmd/utils"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
)

const (
	flagErc20 = "erc20"
)

func GetQueryBalanceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "balance [account] [?optional ERC-20 addrs..]",
		Aliases: []string{"b"},
		Short:   "Get ERC-20 token information. If account address is provided, it will query the balance of the account (bech32 is accepted).",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			evmAddrs, err := utils.GetEvmAddressFromAnyFormatAddress(args...)
			if err != nil {
				utils.PrintlnStdErr("ERR:", err)
				return
			}

			ethClient8545, _ := flags.MustGetEthClient(cmd)
			var restApiEndpoint string

			fetchErc20ModuleAndVfbc := cmd.Flags().Changed(flagErc20)
			if fetchErc20ModuleAndVfbc {
				restApiEndpoint = flags.MustGetCosmosRest(cmd)
			}

			contextHeight := readContextHeightFromFlag(cmd)

			accountAddr := evmAddrs[0]
			fmt.Println("Account", accountAddr)

			printRow := func(colType, colContract, colSymbol, colBalance, colRaw, colDecimals, colHigh, colLow, extra string) {
				fmt.Printf("%-7s | %42s | %10s | %28s | %27s | %8s | %9s | %18s | %-1s\n", colType, colContract, colSymbol, colBalance, colRaw, colDecimals, colHigh, colLow, extra)
			}

			printRow("Type", "Contract", "Symbol", "Balance", "Raw", "Decimals", "High", "Low", "Extra")

			nativeBalance, err := ethClient8545.BalanceAt(context.Background(), accountAddr, contextHeight)
			utils.ExitOnErr(err, "failed to get account balance")

			display, high, low, err := utils.ConvertNumberIntoDisplayWithExponent(nativeBalance, 18)
			utils.ExitOnErr(err, "failed to convert number into display with exponent")

			printRow("native", "-", "(native)", display, nativeBalance.String(), "18", high.String(), low.String(), "")

			for i := 1; i < len(evmAddrs); i++ {
				contractAddr := evmAddrs[i]

				tokenBalance, tokenBalanceDisplay, contractSymbol, contractDecimals, balancePartHigh, balancePartLow, err := fetchBalanceForErc20Contract(contractAddr, contextHeight, ethClient8545, accountAddr, "contract")
				if err != nil {
					continue
				}

				printRow("Input", contractAddr.String(), contractSymbol, tokenBalanceDisplay, tokenBalance.String(), contractDecimals.String(), balancePartHigh.String(), balancePartLow.String(), "")
			}

			if fetchErc20ModuleAndVfbc && restApiEndpoint != "" {
				erc20TokenPairs, statusCode, err := fetchErc20ModuleTokenPairsFromRest(restApiEndpoint)
				if err != nil {
					if statusCode == 501 {
						utils.PrintlnStdErr("WARN: `x/erc20` module is not available on the chain")
					} else {
						utils.PrintlnStdErr("ERR:", err)
					}
				} else {
					for _, erc20TokenPair := range erc20TokenPairs {
						if !erc20TokenPair.Enabled {
							continue
						}

						contractAddr := common.HexToAddress(erc20TokenPair.Erc20Address)

						tokenBalance, tokenBalanceDisplay, contractSymbol, contractDecimals, balancePartHigh, balancePartLow, err := fetchBalanceForErc20Contract(contractAddr, contextHeight, ethClient8545, accountAddr, "x/erc20 contract")
						if err != nil {
							continue
						}

						if tokenBalance.Sign() == 0 {
							continue
						}

						printRow("x/erc20", contractAddr.String(), contractSymbol, tokenBalanceDisplay, tokenBalance.String(), contractDecimals.String(), balancePartHigh.String(), balancePartLow.String(), erc20TokenPair.Denom)
					}
				}

				vfbcPairs, statusCode, err := fetchVirtualFrontierBankContractPairsFromRest(restApiEndpoint)
				if err != nil {
					if statusCode == 501 {
						utils.PrintlnStdErr("WARN: virtual frontier contract feature is not available on the chain")
					} else {
						utils.PrintlnStdErr("ERR:", err)
					}
				} else {
					for _, vfbcPair := range vfbcPairs {
						if !vfbcPair.Enabled {
							continue
						}

						contractAddr := common.HexToAddress(vfbcPair.ContractAddress)

						tokenBalance, tokenBalanceDisplay, contractSymbol, contractDecimals, balancePartHigh, balancePartLow, err := fetchBalanceForErc20Contract(contractAddr, contextHeight, ethClient8545, accountAddr, "VFBC")
						if err != nil {
							continue
						}

						if tokenBalance.Sign() == 0 {
							continue
						}

						printRow("vfbc", contractAddr.String(), contractSymbol, tokenBalanceDisplay, tokenBalance.String(), contractDecimals.String(), balancePartHigh.String(), balancePartLow.String(), vfbcPair.MinDenom)
					}
				}
			}
		},
	}

	cmd.Flags().String(flags.FlagEvmRpc, "", flags.FlagEvmRpcDesc)
	cmd.Flags().String(flags.FlagCosmosRest, "", flags.FlagCosmosRestDesc)
	cmd.Flags().Int64(flagHeight, 0, "query balance at specific height")
	cmd.Flags().Bool(flagErc20, false, "query balance of ERC-20 contracts of `x/erc20` module and virtual frontier bank contracts")

	return cmd
}

func fetchBalanceForErc20Contract(contractAddr common.Address, contextHeight *big.Int, ethClient8545 *ethclient.Client, accountAddr common.Address, contractType string) (
	tokenBalance *big.Int, tokenBalanceDisplay, contractSymbol string,
	contractDecimals, balancePartHigh, balancePartLow *big.Int,
	err error,
) {
	bz, err := ethClient8545.CallContract(context.Background(), ethereum.CallMsg{
		To:   &contractAddr,
		Data: []byte{0x95, 0xd8, 0x9b, 0x41}, // symbol()
	}, contextHeight)
	if err != nil {
		utils.PrintlnStdErr("ERR: failed to get symbol for", contractType, contractAddr, ":", err)
		return
	}

	contractSymbol, err = utils.AbiDecodeString(bz)
	if err != nil {
		utils.PrintlnStdErr("ERR: failed to decode symbol for", contractType, contractAddr, ":", err)
		return
	}

	bz, err = ethClient8545.CallContract(context.Background(), ethereum.CallMsg{
		To:   &contractAddr,
		Data: []byte{0x31, 0x3c, 0xe5, 0x67}, // decimals()
	}, contextHeight)
	if err != nil {
		utils.PrintlnStdErr("ERR: failed to get decimals for", contractType, contractAddr, ":", err)
		return
	}

	contractDecimals = new(big.Int).SetBytes(bz)

	bz, err = ethClient8545.CallContract(context.Background(), ethereum.CallMsg{
		To:   &contractAddr,
		Data: append([]byte{0x70, 0xa0, 0x82, 0x31}, common.BytesToHash(accountAddr.Bytes()).Bytes()...), // balanceOf(address)
	}, contextHeight)
	if err != nil {
		utils.PrintlnStdErr("ERR: failed to get", contractType, "token", contractAddr, "balance for", accountAddr, ":", err)
		return
	}

	tokenBalance = new(big.Int).SetBytes(bz)

	tokenBalanceDisplay, balancePartHigh, balancePartLow, err = utils.ConvertNumberIntoDisplayWithExponent(tokenBalance, int(contractDecimals.Int64()))
	if err != nil {
		utils.PrintlnStdErr("ERR: failed to convert number", tokenBalance.String(), "decimals", contractDecimals.String(), "into display with exponent for", contractType, "token balance:", err)
		return
	}

	return
}
