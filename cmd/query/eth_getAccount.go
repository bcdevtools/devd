package query

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/bcdevtools/devd/v3/cmd/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/bcdevtools/devd/v3/cmd/flags"
	"github.com/bcdevtools/devd/v3/cmd/utils"
	"github.com/spf13/cobra"
)

func GetQueryEvmRpcEthGetAccountCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "eth_getAccount [0xAddress/Bech32]",
		Aliases: []string{"evm-account"},
		Short:   "Query account using `eth_getAccount` via EVM RPC",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ethClient8545, evmRpc := flags.MustGetEthClient(cmd)

			evmAddrs, err := utils.GetEvmAddressFromAnyFormatAddress(args...)
			utils.ExitOnErr(err, "failed to get evm address from input")
			evmAddr := evmAddrs[0]

			contextHeight, err := flags.ReadFlagBlockNumberOrNil(cmd, flags.FlagHeight)
			utils.ExitOnErr(err, "failed to parse block number")

			var params []types.JsonRpcQueryParam

			paramsAddr, err := types.NewJsonRpcStringQueryParam(evmAddr.String())
			utils.ExitOnErr(err, "failed to create json rpc query param")
			params = append(params, paramsAddr)

			paramsContext, err := types.NewJsonRpcStringQueryParam(func() string {
				if contextHeight == nil {
					return "latest"
				}
				return contextHeight.Text(16)
			}())
			utils.ExitOnErr(err, "failed to create json rpc query param")
			params = append(params, paramsContext)

			bz, err := types.DoEvmRpcQuery(
				evmRpc,
				types.NewJsonRpcQueryBuilder(
					"eth_getAccount",
					params...,
				),
				0,
			)

			utils.ExitOnErr(err, "failed to query account")

			codeHashOfEmpty := "0x" + hex.EncodeToString(crypto.Keccak256(nil))

			accountInfoAsMap, err := getResultObjectFromEvmRpcResponse(bz)
			if err != nil && strings.Contains(err.Error(), "method eth_getAccount does not exist") {
				// fallback
				err = nil

				code, err := ethClient8545.CodeAt(context.Background(), evmAddr, contextHeight)
				utils.ExitOnErr(err, "failed to get code")

				balance, err := ethClient8545.BalanceAt(context.Background(), evmAddr, contextHeight)
				utils.ExitOnErr(err, "failed to get balance")

				nonce, err := ethClient8545.NonceAt(context.Background(), evmAddr, contextHeight)
				utils.ExitOnErr(err, "failed to get nonce")

				accountInfoAsMap = map[string]interface{}{
					"codeHash": func() string {
						if len(code) == 0 {
							return codeHashOfEmpty
						}

						return "0x" + hex.EncodeToString(crypto.Keccak256(code))
					}(),
					"balance": balance.String(),
					"nonce":   nonce,
				}
			}
			utils.ExitOnErr(err, "failed to get result object from response")

			if codeHashRaw, found := accountInfoAsMap["codeHash"]; found {
				// normalize
				if codeHashStr, ok := codeHashRaw.(string); ok && codeHashStr == "0x" {
					accountInfoAsMap["codeHash"] = codeHashOfEmpty
				}
			}
			if codeHashRaw, found := accountInfoAsMap["codeHash"]; found {
				if codeHashStr, ok := codeHashRaw.(string); ok {
					isContract := codeHashStr != codeHashOfEmpty
					accountInfoAsMap["_isContract"] = isContract

					if !isContract {
						if nonceRaw, found := accountInfoAsMap["nonce"]; found {
							nonceStr := fmt.Sprintf("%v", nonceRaw)
							nonce, ok := new(big.Int).SetString(nonceStr, 10)
							if !ok {
								utils.PrintlnStdErr("ERR: failed to parse nonce:", nonceStr)
							} else {
								txSent := nonce
								if txSent.Sign() > 0 && txSent.Cmp(big.NewInt(1_000_000)) > 0 {
									txSent = new(big.Int).Mod(txSent, big.NewInt(1_000_000_000)) // Dymension RollApps increases nonce at fraud happened
								}
								accountInfoAsMap["_txSent"] = txSent.String()
							}
						}
					}
				}
			}

			bz, err = json.Marshal(accountInfoAsMap)
			utils.ExitOnErr(err, "failed to marshal account info")

			utils.TryPrintBeautyJson(bz)
		},
	}

	cmd.Flags().String(flags.FlagEvmRpc, "", flags.FlagEvmRpcDesc)
	cmd.Flags().String(flags.FlagHeight, "", "query account info at specific height")

	return cmd
}
