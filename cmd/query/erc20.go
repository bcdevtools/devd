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
)

// GetQueryErc20Command registers a sub-tree of commands
func GetQueryErc20Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "erc20 [contract_address] [?account_address]",
		Short: "Get ERC-20 token information. If account address is provided, it will query the balance of the account (bech32 is accepted).",
		Args:  cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			var rpc string
			if rpc, _ = cmd.Flags().GetString("host"); len(rpc) > 0 {
				// accepted deprecated flag
				libutils.PrintfStdErr("WARN: flag '--host' is deprecated, use '--%s' instead\n", flagRpc)
			} else if rpc, _ = cmd.Flags().GetString(flagRpc); len(rpc) > 0 {
				// accepted new flag
			} else {
				libutils.PrintlnStdErr("ERR: missing RPC to query")
				os.Exit(1)
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
				rpc,
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
				0,
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
				rpc,
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
				0,
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
					rpc,
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
					0,
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

	cmd.Flags().StringP(flagRpc, "p", "http://localhost:8545", "EVM Json-RPC url")
	cmd.Flags().String("host", "", fmt.Sprintf("deprecated flag, use '--%s' instead", flagRpc))

	return cmd
}
