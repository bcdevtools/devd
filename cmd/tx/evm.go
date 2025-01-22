package tx

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/bcdevtools/devd/v2/cmd/utils"
	"github.com/bcdevtools/devd/v2/constants"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
)

func mustGetEthClient(cmd *cobra.Command) (ethClient8545 *ethclient.Client, rpc string) {
	var inputSource string
	var err error

	if rpcFromFlagRpc, _ := cmd.Flags().GetString(flagRpc); len(rpcFromFlagRpc) > 0 {
		rpc = rpcFromFlagRpc
		inputSource = "flag"
	} else if rpcFromEnv := os.Getenv(constants.ENV_EVM_RPC); len(rpcFromEnv) > 0 {
		rpc = rpcFromEnv
		inputSource = "environment variable"
	} else {
		rpc = constants.DEFAULT_EVM_RPC
		inputSource = "default"
	}

	utils.PrintlnStdErr("INF: Connecting to EVM Json-RPC", rpc, fmt.Sprintf("(from %s)", inputSource))

	ethClient8545, err = ethclient.Dial(rpc)
	utils.ExitOnErr(err, "failed to connect to EVM Json-RPC")

	// pre-flight check to ensure the connection is working
	_, err = ethClient8545.BlockNumber(context.Background())
	if err != nil && strings.Contains(err.Error(), "connection refused") {
		utils.PrintlnStdErr("ERR: failed to connect to EVM Json-RPC, please check the connection and try again.")
		utils.PrintfStdErr("ERR: if you are using a custom EVM Json-RPC, please provide it via flag '--%s <your_custom>' or setting environment variable 'export %s=<your_custom>'.\n", flagRpc, constants.ENV_EVM_RPC)
		os.Exit(1)
	}

	return
}

func mustSecretEvmAccount(cmd *cobra.Command) (privKey string, ecdsaPrivateKey *ecdsa.PrivateKey, ecdsaPubKey *ecdsa.PublicKey, account *common.Address) {
	var inputSource string
	var err error
	var ok bool

	if secretFromFlag, _ := cmd.Flags().GetString(flagSecretKey); len(secretFromFlag) > 0 {
		privKey = secretFromFlag
		inputSource = "flag"
	} else if secretFromEnv := os.Getenv(constants.ENV_SECRET_KEY); len(secretFromEnv) > 0 {
		privKey = secretFromEnv
		inputSource = "environment variable"
	} else {
		utils.PrintlnStdErr("ERR: secret key is required")
		utils.PrintfStdErr("ERR: secret key can be set by flag '--%s <your_secret_key>' or environment variable 'export %s=<your_secret_key>'.\n", flagSecretKey, constants.ENV_SECRET_KEY)
		os.Exit(1)
	}

	if regexp.MustCompile(`^(0x)?[a-fA-F\d]{64}$`).MatchString(privKey) {
		// private key
		privKey = strings.TrimPrefix(privKey, "0x")

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
	} else if mnemonicCount := len(strings.Split(privKey, " ")); mnemonicCount == 12 || mnemonicCount == 24 {
		// is mnemonic
		mnemonic := privKey
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
