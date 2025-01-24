package query

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/bcdevtools/devd/v3/cmd/utils"
	"github.com/bcdevtools/devd/v3/constants"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
)

func mustGetEthClient(cmd *cobra.Command, fallbackDeprecatedFlagHost bool) (ethClient8545 *ethclient.Client, evmRpc string) {
	var inputSource string
	var err error

	if evmRpcFromFlagRpc, _ := cmd.Flags().GetString(flagEvmRpc); len(evmRpcFromFlagRpc) > 0 {
		evmRpc = evmRpcFromFlagRpc
		inputSource = "flag"
	} else if evmRpcFromEnv := os.Getenv(constants.ENV_EVM_RPC); len(evmRpcFromEnv) > 0 {
		evmRpc = evmRpcFromEnv
		inputSource = "environment variable"
	} else if evmRpcFromFlagHost, _ := cmd.Flags().GetString("host"); fallbackDeprecatedFlagHost && len(evmRpcFromFlagHost) > 0 {
		utils.PrintfStdErr("WARN: flag '--host' is deprecated, use '--%s' instead\n", flagEvmRpc)
		evmRpc = evmRpcFromFlagHost
		inputSource = "flag"
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

func readContextHeightFromFlag(cmd *cobra.Command) *big.Int {
	height, _ := cmd.Flags().GetInt64(flagHeight)
	if height > 0 {
		return big.NewInt(height)
	}

	return nil
}
