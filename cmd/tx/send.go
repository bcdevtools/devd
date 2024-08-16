package tx

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/bcdevtools/devd/v2/cmd/utils"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/cobra"
)

func GetSendEvmTxCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send [to] [amount]",
		Short: "Send some token to an address via EVM transfer",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			ethClient8545, _ := mustGetEthClient(cmd)

			evmAddrs, err := utils.GetEvmAddressFromAnyFormatAddress(args[0])
			utils.ExitOnErr(err, "failed to get evm address from input")

			receiverAddr := evmAddrs[0]
			amount, ok := new(big.Int).SetString(args[1], 10)
			if !ok {
				utils.ExitOnErr(fmt.Errorf("invalid amount %s", args[1]), "failed to parse amount")
			}
			display, _, _, err := utils.ConvertNumberIntoDisplayWithExponent(amount, 18)
			utils.ExitOnErr(err, "failed to convert amount into display with exponent")

			_, ecdsaPrivateKey, _, from := mustSecretEvmAccount(cmd)

			nonce, err := ethClient8545.NonceAt(context.Background(), *from, nil)
			utils.ExitOnErr(err, "failed to get nonce of sender")

			chainId, err := ethClient8545.ChainID(context.Background())
			utils.ExitOnErr(err, "failed to get chain ID")

			txData := ethtypes.LegacyTx{
				Nonce:    nonce,
				GasPrice: big.NewInt(20_000_000_000),
				Gas:      21000,
				To:       &receiverAddr,
				Value:    amount,
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

			err = ethClient8545.SendTransaction(context.Background(), signedTx)
			utils.ExitOnErr(err, "failed to send tx")
		},
	}

	cmd.Flags().String(flagRpc, "", flagEvmRpcDesc)
	cmd.Flags().String(flagSecretKey, "", flagSecretKeyDesc)

	return cmd
}
