package flags

import (
	"crypto/ecdsa"
	"github.com/bcdevtools/devd/v3/cmd/utils"
	"github.com/bcdevtools/devd/v3/constants"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
	"os"
	"regexp"
	"strings"
)

const (
	FlagSecretKey = "secret-key"
)

const (
	FlagSecretKeyDesc = "Secret private key or mnemonic of the account, can be set by environment variable " + constants.ENV_SECRET_KEY
)

func MustSecretEvmAccount(cmd *cobra.Command) (ecdsaPrivateKey *ecdsa.PrivateKey, ecdsaPubKey *ecdsa.PublicKey, account *common.Address) {
	var inputSource, secret string
	var err error
	var ok bool

	if secretFromFlag, _ := cmd.Flags().GetString(FlagSecretKey); len(secretFromFlag) > 0 {
		secret = secretFromFlag
		inputSource = "flag"
	} else if secretFromEnv := os.Getenv(constants.ENV_SECRET_KEY); len(secretFromEnv) > 0 {
		secret = secretFromEnv
		inputSource = "environment variable"
	} else {
		utils.PrintlnStdErr("ERR: secret key is required")
		utils.PrintfStdErr("ERR: secret key can be set by flag '--%s <your_secret_key>' or environment variable 'export %s=<your_secret_key>'.\n", FlagSecretKey, constants.ENV_SECRET_KEY)
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

	utils.PrintlnStdErr("INF: Account Address:", account.Hex(), "(from", inputSource, ")")

	return
}
