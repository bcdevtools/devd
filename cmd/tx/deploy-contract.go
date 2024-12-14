package tx

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/bcdevtools/devd/v2/cmd/utils"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/cobra"
)

func GetDeployContractEvmTxCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy-contract [bytecode]",
		Short: `Deploy an EVM contract using bytecode.`,
		Long: `Deploy an EVM contract.
Predefined bytecode: erc20`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			bytecode := strings.ToLower(args[0])
			if bytecode == "erc20" {
				bytecode = BytecodeErc20Contract
			}
			deployEvmContract(bytecode, cmd)
		},
	}

	cmd.Flags().String(flagRpc, "", flagEvmRpcDesc)
	cmd.Flags().String(flagSecretKey, "", flagSecretKeyDesc)

	return cmd
}

func deployEvmContract(bytecode string, cmd *cobra.Command) {
	ethClient8545, _ := mustGetEthClient(cmd)

	_, ecdsaPrivateKey, _, from := mustSecretEvmAccount(cmd)

	nonce, err := ethClient8545.NonceAt(context.Background(), *from, nil)
	utils.ExitOnErr(err, "failed to get nonce of sender")

	chainId, err := ethClient8545.ChainID(context.Background())
	utils.ExitOnErr(err, "failed to get chain ID")

	if strings.HasPrefix(bytecode, "0x") {
		bytecode = bytecode[2:]
	}
	deploymentBytes, err := hex.DecodeString(bytecode)
	utils.ExitOnErr(err, "failed to parse deployment bytecode")

	txData := ethtypes.LegacyTx{
		Nonce:    nonce,
		GasPrice: big.NewInt(20_000_000_000),
		Gas:      2_000_000,
		To:       nil,
		Data:     deploymentBytes,
		Value:    common.Big0,
	}
	tx := ethtypes.NewTx(&txData)

	newContractAddress := crypto.CreateAddress(*from, nonce)

	fmt.Println("Deploying new contract using account", from)

	signedTx, err := ethtypes.SignTx(tx, ethtypes.LatestSignerForChainID(chainId), ecdsaPrivateKey)
	utils.ExitOnErr(err, "failed to sign tx")

	var buf bytes.Buffer
	err = signedTx.EncodeRLP(&buf)
	utils.ExitOnErr(err, "failed to encode tx")

	err = ethClient8545.SendTransaction(context.Background(), signedTx)
	utils.ExitOnErr(err, "failed to send tx")

	fmt.Println("Tx hash", signedTx.Hash())

	var found bool
	for try := 1; try <= 6; try++ {
		txByHash, pending, err := ethClient8545.TransactionByHash(context.Background(), signedTx.Hash())
		if err == nil && !pending && txByHash != nil {
			found = true
			break
		}

		time.Sleep(time.Second)
	}

	if found {
		fmt.Println("New contract deployed at:")
	} else {
		fmt.Println("Timed-out waiting for tx to be mined, contract may have been deployed.")
		fmt.Println("Expected contract address:")
	}
	fmt.Println(newContractAddress)
}
