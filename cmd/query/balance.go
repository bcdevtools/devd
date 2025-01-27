package query

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"slices"

	"github.com/pkg/errors"

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
		Short:   "Get native balance of account. Optionally query ERC-20 token balances.",
		Long: fmt.Sprintf(`Get native balance of account. Optionally query ERC-20 token balances of if 2nd arg is provided or flag --%s is used.
Bech32 account address is accepted.`, flagErc20),
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			evmAddrs, err := utils.GetEvmAddressFromAnyFormatAddress(args...)
			if err != nil {
				utils.PrintlnStdErr("ERR:", err)
				return
			}

			ethClient8545, _ := flags.MustGetEthClient(cmd)

			// construct contract list to query

			type contract struct {
				contractAddr    common.Address
				source          string
				skipZeroBalance bool
				extra           string
			}

			var erc20Contracts []contract

			for i := 1; i < len(evmAddrs); i++ {
				erc20Contracts = append(erc20Contracts, contract{
					contractAddr:    evmAddrs[i],
					source:          "input",
					skipZeroBalance: false,
					extra:           "",
				})
			}

			// if flag --erc20 is used, query x/erc20 module and virtual frontier bank contracts
			// (order them by denom to ensure balance of higher priority are fetched/displayed first)

			if cmd.Flags().Changed(flagErc20) {
				restApiEndpoint := flags.MustGetCosmosRest(cmd)

				existingErc20TokenPairs, statusCode, err := fetchErc20ModuleTokenPairsFromRest(restApiEndpoint)
				if err != nil {
					if statusCode == 501 {
						utils.PrintlnStdErr("WARN: `x/erc20` module is not available on the chain")
					} else {
						utils.PrintlnStdErr("ERR: failed to check x/erc20 contract list", err)
					}
				} else {
					slices.SortFunc(existingErc20TokenPairs, func(l, r Erc20ModuleTokenPair) int {
						ln := utils.OrderNumberForDenom(l.Denom)
						rn := utils.OrderNumberForDenom(r.Denom)
						if ln < rn {
							return -1
						} else if ln > rn {
							return 1
						} else {
							return 0
						}
					})

					for _, erc20TokenPair := range existingErc20TokenPairs {
						if !erc20TokenPair.Enabled {
							continue
						}

						erc20Contracts = append(erc20Contracts, contract{
							contractAddr:    common.HexToAddress(erc20TokenPair.Erc20Address),
							source:          "x/erc20",
							skipZeroBalance: true,
							extra:           erc20TokenPair.Denom,
						})
					}
				}

				existingVfbcPairs, statusCode, err := fetchVirtualFrontierBankContractPairsFromRest(restApiEndpoint)
				if err != nil {
					if statusCode == 501 {
						utils.PrintlnStdErr("WARN: virtual frontier contract feature is not available on the chain")
					} else {
						utils.PrintlnStdErr("ERR: failed to check virtual frontier contract list", err)
					}
				} else {
					slices.SortFunc(existingVfbcPairs, func(l, r VfbcTokenPair) int {
						ln := utils.OrderNumberForDenom(l.MinDenom)
						rn := utils.OrderNumberForDenom(r.MinDenom)
						if ln < rn {
							return -1
						} else if ln > rn {
							return 1
						} else {
							return 0
						}
					})

					for _, vfbcPair := range existingVfbcPairs {
						if !vfbcPair.Enabled {
							continue
						}

						erc20Contracts = append(erc20Contracts, contract{
							contractAddr:    common.HexToAddress(vfbcPair.ContractAddress),
							source:          "vfbc",
							skipZeroBalance: true,
							extra:           vfbcPair.MinDenom,
						})
					}
				}
			}

			// start fetching balances

			contextHeight, err := flags.ReadFlagBlockNumberOrNil(cmd, flags.FlagHeight)
			utils.ExitOnErr(err, "failed to parse block number")

			accountAddr := evmAddrs[0]
			utils.PrintlnStdErr("INF: Account", accountAddr)

			printRow := func(colType, colContract, colSymbol, colBalance, colRaw, colDecimals, extra string) {
				fmt.Printf("%-7s | %42s | %-10s | %28s | %27s | %8s | %-1s\n", colType, colContract, colSymbol, colBalance, colRaw, colDecimals, extra)
			}

			printRow("Type", "Contract", "Symbol", "Balance", "Raw", "Decimals", "Extra")

			nativeBalance, err := ethClient8545.BalanceAt(context.Background(), accountAddr, contextHeight)
			utils.ExitOnErr(err, "failed to get account balance")

			display, _, _, err := utils.ConvertNumberIntoDisplayWithExponent(nativeBalance, 18)
			utils.ExitOnErr(err, "failed to convert number into display with exponent")

			printRow("native", "-", "(native)", display, nativeBalance.String(), "18", "")

			for _, erc20Contract := range erc20Contracts {
				tokenBalance, tokenBalanceDisplay, contractSymbol, contractDecimals, _, _, err := fetchBalanceForErc20Contract(erc20Contract.contractAddr, contextHeight, ethClient8545, accountAddr, erc20Contract.source)
				if err != nil {
					continue
				}

				if tokenBalance.Sign() == 0 && erc20Contract.skipZeroBalance {
					continue
				}

				printRow(erc20Contract.source, erc20Contract.contractAddr.String(), contractSymbol, tokenBalanceDisplay, tokenBalance.String(), contractDecimals.String(), erc20Contract.extra)
			}
		},
	}

	cmd.Flags().String(flags.FlagEvmRpc, "", flags.FlagEvmRpcDesc)
	cmd.Flags().String(flags.FlagCosmosRest, "", flags.FlagCosmosRestDesc)
	cmd.Flags().Int64(flags.FlagHeight, 0, "query balance at specific height")
	cmd.Flags().Bool(flagErc20, false, "query balance of ERC-20 contracts of `x/erc20` module and virtual frontier bank contracts")

	return cmd
}

func fetchBalanceForErc20Contract(contractAddr common.Address, contextHeight *big.Int, ethClient8545 *ethclient.Client, accountAddr common.Address, sourceInputContract string) (
	tokenBalance *big.Int, tokenBalanceDisplay, contractSymbol string,
	contractDecimals, balancePartHigh, balancePartLow *big.Int,
	err error,
) {
	bz, err := ethClient8545.CallContract(context.Background(), ethereum.CallMsg{
		To:   &contractAddr,
		Data: []byte{0x95, 0xd8, 0x9b, 0x41}, // symbol()
	}, contextHeight)
	if err != nil {
		utils.PrintlnStdErr("ERR: failed to get symbol for", sourceInputContract, "contract", contractAddr, ":", err)
		return
	}

	contractSymbol, err = utils.AbiDecodeString(bz)
	if err != nil {
		utils.PrintlnStdErr("ERR: failed to decode symbol for", sourceInputContract, "contract", contractAddr, ":", err)
		return
	}

	bz, err = ethClient8545.CallContract(context.Background(), ethereum.CallMsg{
		To:   &contractAddr,
		Data: []byte{0x31, 0x3c, 0xe5, 0x67}, // decimals()
	}, contextHeight)
	if err != nil {
		utils.PrintlnStdErr("ERR: failed to get decimals for", sourceInputContract, "contract", contractAddr, ":", err)
		return
	}

	contractDecimals = new(big.Int).SetBytes(bz)

	bz, err = ethClient8545.CallContract(context.Background(), ethereum.CallMsg{
		To:   &contractAddr,
		Data: append([]byte{0x70, 0xa0, 0x82, 0x31}, common.BytesToHash(accountAddr.Bytes()).Bytes()...), // balanceOf(address)
	}, contextHeight)
	if err != nil {
		utils.PrintlnStdErr("ERR: failed to get", sourceInputContract, "contract token", contractAddr, "balance for", accountAddr, ":", err)
		return
	}

	tokenBalance = new(big.Int).SetBytes(bz)

	tokenBalanceDisplay, balancePartHigh, balancePartLow, err = utils.ConvertNumberIntoDisplayWithExponent(tokenBalance, int(contractDecimals.Int64()))
	if err != nil {
		utils.PrintlnStdErr("ERR: failed to convert number", tokenBalance.String(), "decimals", contractDecimals.String(), "into display with exponent for", sourceInputContract, "contract token balance:", err)
		return
	}

	return
}

type Erc20ModuleTokenPair struct {
	Erc20Address string `json:"erc20_address"`
	Denom        string `json:"denom"`
	Enabled      bool   `json:"enabled"`
}

func fetchErc20ModuleTokenPairsFromRest(rest string) (erc20ModuleTokenPairs []Erc20ModuleTokenPair, statusCode int, err error) {
	var resp *http.Response
	resp, err = http.Get(rest + "/evmos/erc20/v1/token_pairs?pagination.limit=10000")
	if err != nil {
		err = errors.Wrap(err, "failed to fetch ERC-20 module token pairs")
		return
	}

	statusCode = resp.StatusCode

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("failed to fetch ERC-20 module token pairs! Status code: %d", resp.StatusCode)
		return
	}

	type responseStruct struct {
		TokenPairs []Erc20ModuleTokenPair `json:"token_pairs"`
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	bz, err := io.ReadAll(resp.Body)
	if err != nil {
		err = errors.Wrap(err, "failed to read response body of ERC-20 module token pairs")
		return
	}

	var response responseStruct
	err = json.Unmarshal(bz, &response)
	if err != nil {
		err = errors.Wrap(err, "failed to unmarshal response body of ERC-20 module token pairs")
		return
	}

	erc20ModuleTokenPairs = response.TokenPairs
	return
}

type VfbcTokenPair struct {
	ContractAddress string `json:"contract_address"`
	MinDenom        string `json:"min_denom"`
	Enabled         bool   `json:"enabled"`
}

func fetchVirtualFrontierBankContractPairsFromRest(rest string) (vfbcPairs []VfbcTokenPair, statusCode int, err error) {
	var resp *http.Response
	resp, err = http.Get(rest + "/ethermint/evm/v1/virtual_frontier_bank_contracts?pagination.limit=10000")
	if err != nil {
		err = errors.Wrap(err, "failed to fetch VFBC pairs")
		return
	}

	statusCode = resp.StatusCode

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("failed to fetch VFBC pairs! Status code: %d", resp.StatusCode)
		return
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	bz, err := io.ReadAll(resp.Body)
	if err != nil {
		err = errors.Wrap(err, "failed to read response body of VFBC pairs")
		return
	}

	type responseStruct struct {
		Pairs []VfbcTokenPair `json:"pairs"`
	}

	var response responseStruct
	err = json.Unmarshal(bz, &response)
	if err != nil {
		err = errors.Wrap(err, "failed to unmarshal response body of VFBC pairs")
		return
	}

	vfbcPairs = response.Pairs
	return
}
