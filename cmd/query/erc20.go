package query

import (
	"encoding/hex"
	"fmt"
	libutils "github.com/EscanBE/go-lib/utils"
	"github.com/bcdevtools/devd/cmd/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"math/big"
	"os"
	"strings"
	"time"
)

// GetQueryErc20Command registers a sub-tree of commands
func GetQueryErc20Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "erc20 [contract_address] [?account_address]",
		Short: "Get ERC-20 token information. If account address is provided, it will query the balance of the account (bech32 is accepted).",
		Args:  cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			const queryTimeout = 3 * time.Second

			host, _ := cmd.Flags().GetString(flagHost)
			if len(host) == 0 {
				libutils.PrintlnStdErr("ERR: missing host")
				return
			}

			evmAddrs, err := getEvmAddressFromAnyFormatAddress(args...)
			if err != nil {
				libutils.PrintlnStdErr(err)
				return
			}

			var contractAddr, accountAddr common.Address

			contractAddr = evmAddrs[0]

			if len(evmAddrs) > 1 {
				accountAddr = evmAddrs[1]
			}

			paramLatest, err := types.NewJsonRpcStringQueryParam("latest")
			libutils.ExitIfErr(err, "failed to create json rpc query param")

			fmt.Println("Getting contract symbol...")
			bz, err := doQuery(
				host,
				types.NewJsonRpcQueryBuilder(
					"eth_call",
					types.NewJsonRpcRawQueryParam(
						fmt.Sprintf(
							`{"from":null,"to":"%s","data":"0x95d89b41"}`,
							strings.ToLower(contractAddr.String()),
						),
					),
					paramLatest,
				),
				queryTimeout,
			)
			if err != nil {
				libutils.PrintlnStdErr("ERR: failed to query contract symbol:", err)
				os.Exit(1)
			}

			contractSymbol, err := decodeResponseToString(bz, "contract symbol")
			if err != nil {
				libutils.PrintlnStdErr("ERR: failed to decode contract symbol:", err)
				os.Exit(1)
			}

			fmt.Println("Getting contract decimals...")

			bz, err = doQuery(
				host,
				types.NewJsonRpcQueryBuilder(
					"eth_call",
					types.NewJsonRpcRawQueryParam(
						fmt.Sprintf(
							`{"from":null,"to":"%s","data":"0x313ce567"}`,
							strings.ToLower(contractAddr.String()),
						),
					),
					paramLatest,
				),
				queryTimeout,
			)

			contractDecimals, err := decodeResponseToBigInt(bz, "contract decimals")
			if err != nil {
				libutils.PrintlnStdErr("ERR: failed to decode contract decimals:", err)
				os.Exit(1)
			}

			var accountBalance *big.Int
			if accountAddr != (common.Address{}) {
				fmt.Println("Getting account balance...")

				bz, err = doQuery(
					host,
					types.NewJsonRpcQueryBuilder(
						"eth_call",
						types.NewJsonRpcRawQueryParam(
							fmt.Sprintf(`{"from":null,"to":"%s","data":"0x70a08231000000000000000000000000%s"}`,
								strings.ToLower(contractAddr.String()),
								strings.ToLower(hex.EncodeToString(accountAddr.Bytes())),
							),
						),
						paramLatest,
					),
					queryTimeout,
				)

				accountBalance, err = decodeResponseToBigInt(bz, "account balance")
				if err != nil {
					libutils.PrintlnStdErr("ERR: failed to decode account balance:", err)
					os.Exit(1)
				}
			}

			fmt.Println("Contract Symbol:", contractSymbol)
			fmt.Println("Contract Decimals:", contractDecimals.Uint64())
			if accountBalance != nil {
				decimals := contractDecimals.Uint64()
				if decimals == 0 {
					fmt.Println("Account Balance:", accountBalance, contractSymbol)
				} else {
					fmt.Println("Account Balance:")
					fmt.Println(" - Raw:", accountBalance)
					fmt.Println(" - Display:")
					pow := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
					fmt.Println("  + High:", new(big.Int).Div(accountBalance, pow), contractSymbol)
					fmt.Println("  + Low:", new(big.Int).Mod(accountBalance, pow))
				}
			}
		},
	}

	cmd.Flags().String(flagHost, "http://localhost:8545", "EVM RPC host")

	return cmd
}
