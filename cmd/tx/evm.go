package tx

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/bcdevtools/devd/v3/cmd/utils"
	"github.com/bcdevtools/devd/v3/constants"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
)

func mustGetEthClient(cmd *cobra.Command) (ethClient8545 *ethclient.Client, evmRpc string) {
	var inputSource string
	var err error

	if evmRpcFromFlag, _ := cmd.Flags().GetString(flagEvmRpc); len(evmRpcFromFlag) > 0 {
		evmRpc = evmRpcFromFlag
		inputSource = "flag"
	} else if evmRpcFromEnv := os.Getenv(constants.ENV_EVM_RPC); len(evmRpcFromEnv) > 0 {
		evmRpc = evmRpcFromEnv
		inputSource = "environment variable"
	} else {
		evmRpc = constants.DEFAULT_EVM_RPC
		inputSource = "default"
	}

	utils.PrintlnStdErr("INF: Connecting to EVM Json-RPC", evmRpc, fmt.Sprintf("(from %s)", inputSource))

	ethClient8545, err = ethclient.Dial(evmRpc)
	utils.ExitOnErr(err, "failed to connect to EVM Json-RPC")

	// pre-flight check to ensure the connection is working
	_, err = ethClient8545.BlockNumber(context.Background())
	if err != nil && strings.Contains(err.Error(), "connection refused") {
		utils.PrintlnStdErr("ERR: failed to connect to EVM Json-RPC, please check the connection and try again.")
		utils.PrintfStdErr("ERR: if you are using a custom EVM Json-RPC, please provide it via flag '--%s <your_custom>' or setting environment variable 'export %s=<your_custom>'.\n", flagEvmRpc, constants.ENV_EVM_RPC)
		os.Exit(1)
	}

	return
}

func mustSecretEvmAccount(cmd *cobra.Command) (ecdsaPrivateKey *ecdsa.PrivateKey, ecdsaPubKey *ecdsa.PublicKey, account *common.Address) {
	var inputSource, secret string
	var err error
	var ok bool

	if secretFromFlag, _ := cmd.Flags().GetString(flagSecretKey); len(secretFromFlag) > 0 {
		secret = secretFromFlag
		inputSource = "flag"
	} else if secretFromEnv := os.Getenv(constants.ENV_SECRET_KEY); len(secretFromEnv) > 0 {
		secret = secretFromEnv
		inputSource = "environment variable"
	} else {
		utils.PrintlnStdErr("ERR: secret key is required")
		utils.PrintfStdErr("ERR: secret key can be set by flag '--%s <your_secret_key>' or environment variable 'export %s=<your_secret_key>'.\n", flagSecretKey, constants.ENV_SECRET_KEY)
		os.Exit(1)
	}

	if regexp.MustCompile(`^(0x)?[a-fA-F\d]{64}$`).MatchString(secret) {
		// private key
		privKey := strings.TrimPrefix(secret, "0x")

		pKeyBytes, err := hexutil.Decode("0x" + privKey)
		if err != nil {
			utils.PrintlnStdErr("ERR: failed to decode private key")
			os.Exit(1)
		}

		ecdsaPrivateKey, err = crypto.ToECDSA(pKeyBytes)
		if err != nil {
			utils.PrintlnStdErr("ERR: failed to convert private key to ECDSA")
			os.Exit(1)
		}
	} else if mnemonicCount := len(strings.Split(secret, " ")); mnemonicCount == 12 || mnemonicCount == 24 {
		// is mnemonic
		mnemonic := secret
		ecdsaPrivateKey, err = utils.FromMnemonicToPrivateKey(mnemonic, "" /*no password protected*/)
		utils.ExitOnErr(err, "failed to convert mnemonic to private key")
	} else {
		utils.PrintlnStdErr("ERR: invalid secret key format")
		os.Exit(1)
	}

	publicKey := ecdsaPrivateKey.Public()
	ecdsaPubKey, ok = publicKey.(*ecdsa.PublicKey)
	if !ok {
		utils.PrintlnStdErr("ERR: failed to convert secret public key to ECDSA")
		os.Exit(1)
	}

	fromAddress := crypto.PubkeyToAddress(*ecdsaPubKey)
	account = &fromAddress

	fmt.Println("Account Address:", account.Hex(), "(from", inputSource, ")")

	return
}
