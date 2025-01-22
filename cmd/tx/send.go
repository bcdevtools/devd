package tx

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"

	"github.com/bcdevtools/devd/v2/cmd/utils"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/cobra"
)

const flagErc20 = "erc20"

func GetSendEvmTxCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send [to] [amount]",
		Short: "Send some native coin or ERC-20 token to an address",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			ethClient8545, _ := mustGetEthClient(cmd)

			evmAddrs, err := utils.GetEvmAddressFromAnyFormatAddress(args[0])
			utils.ExitOnErr(err, "failed to get evm address from input")

			var pErc20ContractAddress *common.Address
			erc20ContractAddress, _ := cmd.Flags().GetString(flagErc20)
			if erc20ContractAddress != "" {
				if !common.IsHexAddress(erc20ContractAddress) {
					utils.ExitOnErr(fmt.Errorf("invalid format"), "failed to parse ERC-20 contract address")
				}
				contractAddr := common.HexToAddress(erc20ContractAddress)
				pErc20ContractAddress = &contractAddr
			}

			receiverAddr := evmAddrs[0]
			amount, ok := new(big.Int).SetString(args[1], 10)
			if !ok {
				utils.ExitOnErr(fmt.Errorf("invalid amount %s", args[1]), "failed to parse amount")
			}

			var exponent int
			if pErc20ContractAddress != nil {
				bz, err := ethClient8545.CallContract(context.Background(), ethereum.CallMsg{
					To:   pErc20ContractAddress,
					Data: []byte{0x31, 0x3c, 0xe5, 0x67}, // decimals()
				}, nil)
				utils.ExitOnErr(err, "failed to get contract decimals")

				contractDecimals := new(big.Int).SetBytes(bz)
				exponent = int(contractDecimals.Int64())
			} else {
				exponent = 18
			}
			display, _, _, err := utils.ConvertNumberIntoDisplayWithExponent(amount, exponent)
			utils.ExitOnErr(err, "failed to convert amount into display with exponent")

			ecdsaPrivateKey, _, from := mustSecretEvmAccount(cmd)

			nonce, err := ethClient8545.NonceAt(context.Background(), *from, nil)
			utils.ExitOnErr(err, "failed to get nonce of sender")

			chainId, err := ethClient8545.ChainID(context.Background())
			utils.ExitOnErr(err, "failed to get chain ID")

			var txData ethtypes.LegacyTx
			if pErc20ContractAddress != nil {
				data := []byte{0xa9, 0x05, 0x9c, 0xbb}
				data = append(data, common.LeftPadBytes(receiverAddr.Bytes(), 32)...)
				data = append(data, common.LeftPadBytes(amount.Bytes(), 32)...)

				txData = ethtypes.LegacyTx{
					Nonce:    nonce,
					GasPrice: big.NewInt(20_000_000_000),
					Gas:      500_000,
					To:       pErc20ContractAddress,
					Value:    big.NewInt(0),
					Data:     data,
				}
			} else {
				txData = ethtypes.LegacyTx{
					Nonce:    nonce,
					GasPrice: big.NewInt(20_000_000_000),
					Gas:      21000,
					To:       &receiverAddr,
					Value:    amount,
				}
			}
			tx := ethtypes.NewTx(&txData)

			fmt.Println("Send", display, "from", from.Hex(), "to", receiverAddr.Hex())
			fmt.Println("EIP155 Chain ID:", chainId.String(), "and nonce", txData.Nonce)

			signedTx, err := ethtypes.SignTx(tx, ethtypes.LatestSignerForChainID(chainId), ecdsaPrivateKey)
			utils.ExitOnErr(err, "failed to sign tx")

			var buf bytes.Buffer
			err = signedTx.EncodeRLP(&buf)
			utils.ExitOnErr(err, "failed to encode tx")

			rawTxRLPHex := hex.EncodeToString(buf.Bytes())
			fmt.Printf("RawTx: 0x%s\n", rawTxRLPHex)

			fmt.Println("Tx hash", signedTx.Hash())

			err = ethClient8545.SendTransaction(context.Background(), signedTx)
			utils.ExitOnErr(err, "failed to send tx")

			if tx := waitForEthTx(ethClient8545, signedTx.Hash()); tx != nil {
				fmt.Println("Tx executed successfully")
			} else {
				fmt.Println("Timed out waiting for tx to be mined")
			}
		},
	}

	cmd.Flags().String(flagRpc, "", flagEvmRpcDesc)
	cmd.Flags().String(flagSecretKey, "", flagSecretKeyDesc)
	cmd.Flags().String(flagErc20, "", "contract address if you want to send ERC-20 token instead of native coin")

	return cmd
}
