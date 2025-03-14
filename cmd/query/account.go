package query

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256r1"
	"io"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bcdevtools/devd/v3/third_party/ethsecp256k1"
	sdkcodec "github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting/exported"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"

	sdkcodectypes "github.com/cosmos/cosmos-sdk/codec/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"

	"github.com/bcdevtools/devd/v3/cmd/flags"

	"github.com/bcdevtools/devd/v3/cmd/utils"
	"github.com/spf13/cobra"
)

func GetQueryAccountCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "account [0xAccount/Bech32]",
		Aliases: []string{"acc"},
		Short:   "Get account details based on address",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			evmAddrs, err := utils.GetEvmAddressFromAnyFormatAddress(args...)
			if err != nil {
				utils.PrintlnStdErr("ERR:", err)
				return
			}

			restApiEndpoint := flags.MustGetCosmosRest(cmd)

			utils.PrintlnStdErr("INF: querying bech32 prefix")
			bech32Prefix, statusCode, err := fetchBech32PrefixFromRest(restApiEndpoint)
			if err != nil {
				if statusCode == 501 {
					utils.PrintlnStdErr("ERR: REST API does not support query bech32 prefix info")
				} else {
					utils.PrintlnStdErr("ERR: failed to fetch bech32 prefix:", err)
				}
				os.Exit(1)
			}
			bech32, err := sdk.Bech32ifyAddressBytes(bech32Prefix, evmAddrs[0].Bytes())
			if err == nil && bech32 == "" {
				err = errors.New("output bech32 address is empty")
			}
			utils.ExitOnErr(err, "failed to convert address to bech32")
			utils.PrintlnStdErr("INF: querying account", bech32)
			response, statusCode, err := fetchAccountDetailsFromRest(restApiEndpoint, bech32)
			if err != nil {
				if statusCode == 501 {
					utils.PrintlnStdErr("ERR: REST API does not support this feature")
				} else {
					utils.PrintlnStdErr("ERR: failed to fetch account details:", err)
				}
				os.Exit(1)
			}

			var accountInfoAsMap map[string]interface{}
			err = json.Unmarshal([]byte(response), &accountInfoAsMap)
			utils.ExitOnErr(err, "failed to unmarshal account details")

			if accountRaw, found := accountInfoAsMap["account"]; found {
				if accountMap, ok := accountRaw.(map[string]interface{}); ok && accountMap != nil {
					if typeRaw, found := accountMap["@type"]; found {
						if typeString, ok := typeRaw.(string); ok {
							if typeString == "/ethermint.types.v1.EthAccount" {
								codeHashOfEmpty := "0x" + hex.EncodeToString(crypto.Keccak256(nil))

								if codeHashRaw, found := accountMap["code_hash"]; found {
									if codeHashStr, ok := codeHashRaw.(string); ok {
										isContract := codeHashStr != codeHashOfEmpty
										accountInfoAsMap["_isContract"] = isContract

										if !isContract {
											if baseAccountRaw, found := accountMap["base_account"]; found {
												if baseAccountMap, ok := baseAccountRaw.(map[string]interface{}); ok {
													if accountNumberRaw, found := baseAccountMap["account_number"]; found {
														nonceStr := fmt.Sprintf("%v", accountNumberRaw)
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
									}
								}
							} else if strings.HasPrefix(typeString, "/cosmos.vesting.v1beta1.") {
								tendermintRpcHttpClient, _ := flags.MustGetTmRpc(cmd)

								req := authtypes.QueryAccountRequest{
									Address: bech32,
								}
								bz, err := req.Marshal()
								utils.ExitOnErr(err, "failed to marshal ABCI query account details")

								resultABCIQuery, err := tendermintRpcHttpClient.ABCIQuery(context.Background(), "/cosmos.auth.v1beta1.Query/Account", bz)
								utils.ExitOnErr(err, "failed to query account details")

								var queryAccountResponse authtypes.QueryAccountResponse
								err = queryAccountResponse.Unmarshal(resultABCIQuery.Response.Value)
								utils.ExitOnErr(err, "failed to unmarshal ABCI query account details")

								ir := sdkcodectypes.NewInterfaceRegistry()
								vestingtypes.RegisterInterfaces(ir)
								ir.RegisterInterface(
									"ethermint.crypto.v1.ethsecp256k1.PubKey",
									(*cryptotypes.PubKey)(nil),
									&ethsecp256k1.PubKey{},
								)
								ir.RegisterInterface(
									"cosmos.crypto.secp256k1.PubKey",
									(*cryptotypes.PubKey)(nil),
									&secp256k1.PubKey{},
								)
								ir.RegisterInterface(
									"cosmos.crypto.secp256r1.PubKey",
									(*cryptotypes.PubKey)(nil),
									&secp256r1.PubKey{},
								)
								cdc := sdkcodec.NewProtoCodec(ir)

								var account exported.VestingAccount
								err = cdc.UnpackAny(queryAccountResponse.Account, &account)
								utils.ExitOnErr(err, "failed to unpack account details")

								now := time.Now()
								accountInfoAsMap["_vesting"] = map[string]interface{}{
									"vested":  account.GetVestedCoins(now),
									"vesting": account.GetVestingCoins(now),
									"locked":  account.LockedCoins(now),
								}
							}
						}
					}
				}
			}

			bz, err := json.Marshal(accountInfoAsMap)
			utils.TryPrintBeautyJson(bz)
		},
	}

	cmd.Flags().String(flags.FlagCosmosRest, "", flags.FlagCosmosRestDesc)
	cmd.Flags().String(flags.FlagTendermintRpc, "", fmt.Sprintf("(conditionally used, not always) %s", flags.FlagTmRpcDesc))

	return cmd
}

func fetchAccountDetailsFromRest(rest, bech32Address string) (response string, statusCode int, err error) {
	var resp *http.Response
	resp, err = http.Get(rest + "/cosmos/auth/v1beta1/accounts/" + bech32Address)
	if err != nil {
		err = errors.Wrap(err, "failed to fetch account details")
		return
	}

	statusCode = resp.StatusCode

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("failed to fetch account details! Status code: %d", resp.StatusCode)
		return
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	bz, err := io.ReadAll(resp.Body)
	if err != nil {
		err = errors.Wrap(err, "failed to read response body of account details")
		return
	}

	response = string(bz)
	return
}

func fetchBech32PrefixFromRest(rest string) (bech32Prefix string, statusCode int, err error) {
	var resp *http.Response
	resp, err = http.Get(rest + "/cosmos/auth/v1beta1/bech32")
	if err != nil {
		err = errors.Wrap(err, "failed to fetch bech32 prefix info")
		return
	}

	statusCode = resp.StatusCode

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("failed to fetch bech32 prefix info! Status code: %d", resp.StatusCode)
		return
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	bz, err := io.ReadAll(resp.Body)
	if err != nil {
		err = errors.Wrap(err, "failed to read response body of bech32 prefix info")
		return
	}

	type bech32PrefixResponse struct {
		Bech32Prefix string `json:"bech32_prefix"`
	}

	var bech32PrefixResp bech32PrefixResponse
	err = json.Unmarshal(bz, &bech32PrefixResp)
	if err != nil {
		err = errors.Wrap(err, "failed to unmarshal bech32 prefix info")
		return
	}
	if bech32PrefixResp.Bech32Prefix == "" {
		err = errors.New("bech32 prefix is empty")
		return
	}

	bech32Prefix = bech32PrefixResp.Bech32Prefix
	return
}
