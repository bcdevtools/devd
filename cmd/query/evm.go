package query

import (
	"context"
	"fmt"
	"github.com/bcdevtools/devd/v2/cmd/utils"
	"github.com/bcdevtools/devd/v2/constants"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
	"math/big"
	"os"
	"strings"
)

func mustGetEthClient(cmd *cobra.Command, fallbackDeprecatedFlagHost bool) (ethClient8545 *ethclient.Client, rpc string) {
	var inputSource string
	var err error

	if rpcFromFlagRpc, _ := cmd.Flags().GetString(flagRpc); len(rpcFromFlagRpc) > 0 {
		rpc = rpcFromFlagRpc
		inputSource = "flag"
	} else if rpcFromEnv := os.Getenv(constants.ENV_EVM_RPC); len(rpcFromEnv) > 0 {
		rpc = rpcFromEnv
		inputSource = "environment variable"
	} else if rpcFromFlagHost, _ := cmd.Flags().GetString("host"); fallbackDeprecatedFlagHost && len(rpcFromFlagHost) > 0 {
		utils.PrintfStdErr("WARN: flag '--host' is deprecated, use '--%s' instead\n", flagRpc)
		rpc = rpcFromFlagHost
		inputSource = "flag"
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

func readContextHeightFromFlag(cmd *cobra.Command) *big.Int {
	height, _ := cmd.Flags().GetInt64(flagHeight)
	if height > 0 {
		return big.NewInt(height)
	}

	return nil
}
