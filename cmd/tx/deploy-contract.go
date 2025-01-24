package tx

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/bcdevtools/devd/v3/cmd/flags"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/bcdevtools/devd/v3/cmd/utils"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/cobra"
)

func GetDeployContractEvmTxCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy-contract [bytecode]",
		Short: `Deploy an EVM contract using bytecode. Constructor calldata is needed if contract has constructor.`,
		Long: `Deploy an EVM contract. Constructor calldata is needed if contract has constructor.
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

	cmd.Flags().String(flags.FlagEvmRpc, "", flags.FlagEvmRpcDesc)
	cmd.Flags().String(flags.FlagSecretKey, "", flags.FlagSecretKeyDesc)
	cmd.Flags().String(flagGasLimit, "4m", flagGasLimitDesc)
	cmd.Flags().String(flagGasPrices, "20b", flagGasPricesDesc)

	return cmd
}

func deployEvmContract(bytecode string, cmd *cobra.Command) {
	ethClient8545, _ := flags.MustGetEthClient(cmd)

	ecdsaPrivateKey, _, from := flags.MustSecretEvmAccount(cmd)

	gasPrices, err := readGasPrices(cmd)
	utils.ExitOnErr(err, "failed to parse gas price")

	gasLimit, err := readGasLimit(cmd)
	utils.ExitOnErr(err, "failed to parse gas limit")

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
		GasPrice: gasPrices,
		Gas:      gasLimit,
		To:       nil,
		Data:     deploymentBytes,
		Value:    common.Big0,
	}
	tx := ethtypes.NewTx(&txData)

	newContractAddress := crypto.CreateAddress(*from, nonce)

	utils.PrintlnStdErr("INF: Deploying new contract using account", from)

	signedTx, err := ethtypes.SignTx(tx, ethtypes.LatestSignerForChainID(chainId), ecdsaPrivateKey)
	utils.ExitOnErr(err, "failed to sign tx")

	var buf bytes.Buffer
	err = signedTx.EncodeRLP(&buf)
	utils.ExitOnErr(err, "failed to encode tx")

	utils.PrintlnStdErr("INF: Tx hash", signedTx.Hash())

	err = ethClient8545.SendTransaction(context.Background(), signedTx)
	utils.ExitOnErr(err, "failed to send tx")

	if tx := waitForEthTx(ethClient8545, signedTx.Hash()); tx != nil {
		utils.PrintlnStdErr("INF: New contract deployed at:")
	} else {
		utils.PrintlnStdErr("WARN: Timed-out waiting for tx to be mined, contract may have been deployed.")
		utils.PrintlnStdErr("INF: Expected contract address:")
	}
	fmt.Println(newContractAddress)
}

func waitForEthTx(ethClient8545 *ethclient.Client, txHash common.Hash) *ethtypes.Transaction {
	for try := 1; try <= 6; try++ {
		tx, _, err := ethClient8545.TransactionByHash(context.Background(), txHash)
		if err == nil && tx != nil {
			return tx
		}

		time.Sleep(time.Second)
	}

	return nil
}
