package flags

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/bcdevtools/devd/v3/cmd/utils"
	"github.com/bcdevtools/devd/v3/constants"
	httpclient "github.com/cometbft/cometbft/rpc/client/http"
	jsonrpcclient "github.com/cometbft/cometbft/rpc/jsonrpc/client"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
)

const (
	FlagEvmRpc        = "evm-rpc"
	FlagTendermintRpc = "tm-rpc"
	FlagCosmosRest    = "rest"
)

const (
	FlagEvmRpcDesc     = "EVM Json-RPC endpoint, default is " + constants.DEFAULT_EVM_RPC + ", can be set by environment variable " + constants.ENV_EVM_RPC
	FlagCosmosRestDesc = "Cosmos Rest API endpoint, default is " + constants.DEFAULT_COSMOS_REST + ", can be set by environment variable " + constants.ENV_COSMOS_REST
	FlagTmRpcDesc      = "Tendermint RPC endpoint, default is " + constants.DEFAULT_TM_RPC + ", can be set by environment variable " + constants.ENV_TM_RPC
)

func MustGetEthClient(cmd *cobra.Command) (ethClient8545 *ethclient.Client, evmRpc string) {
	var inputSource string
	var err error

	if evmRpcFromFlagRpc, _ := cmd.Flags().GetString(FlagEvmRpc); len(evmRpcFromFlagRpc) > 0 {
		evmRpc = evmRpcFromFlagRpc
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
		utils.PrintfStdErr("ERR: if you are using a custom EVM Json-RPC, please provide it via flag '--%s <your_custom>' or setting environment variable 'export %s=<your_custom>'.\n", FlagEvmRpc, constants.ENV_EVM_RPC)
		os.Exit(1)
	}

	return
}

func MustGetTmRpc(cmd *cobra.Command) (tendermintRpcHttpClient *httpclient.HTTP, tmRpc string) {
	var inputSource string

	if tmRpcFromFlag, _ := cmd.Flags().GetString(FlagTendermintRpc); len(tmRpcFromFlag) > 0 {
		tmRpc = tmRpcFromFlag
		inputSource = "flag"
	} else if restFromEnv := os.Getenv(constants.ENV_TM_RPC); len(restFromEnv) > 0 {
		tmRpc = restFromEnv
		inputSource = "environment variable"
	} else {
		tmRpc = constants.DEFAULT_TM_RPC
		inputSource = "default"
	}

	tmRpc = strings.TrimSuffix(tmRpc, "/")
	utils.PrintlnStdErr("INF: Connecting to Tendermint RPC", tmRpc, fmt.Sprintf("(from %s)", inputSource))

	httpClient26657, err := jsonrpcclient.DefaultHTTPClient(tmRpc)
	if err == nil {
		tendermintRpcHttpClient, err = httpclient.NewWithClient(tmRpc, "/websocket", httpClient26657)
	}

	if err != nil {
		utils.PrintlnStdErr("ERR:", err)
		utils.PrintlnStdErr("ERR: failed to connect to TM RPC, please check the connection and try again.")
		utils.PrintfStdErr("ERR: if you are using a custom TM RPC endpoint, please provide it via flag '--%s <your_custom>' or setting environment variable 'export %s=<your_custom>'.\n", FlagTendermintRpc, constants.ENV_TM_RPC)
		os.Exit(1)
	}

	return
}

func MustGetCosmosRest(cmd *cobra.Command) (rest string) {
	var inputSource string

	if restFromFlagRest, _ := cmd.Flags().GetString(FlagCosmosRest); len(restFromFlagRest) > 0 {
		rest = restFromFlagRest
		inputSource = "flag"
	} else if restFromEnv := os.Getenv(constants.ENV_COSMOS_REST); len(restFromEnv) > 0 {
		rest = restFromEnv
		inputSource = "environment variable"
	} else {
		rest = constants.DEFAULT_COSMOS_REST
		inputSource = "default"
	}

	rest = strings.TrimSuffix(rest, "/")

	utils.PrintlnStdErr("INF: Connecting to Cosmos Rest-API", rest, fmt.Sprintf("(from %s)", inputSource))

	// pre-flight check to ensure the connection is working
	_, err := http.Get(rest)
	if err != nil && strings.Contains(err.Error(), "connection refused") {
		utils.PrintlnStdErr("ERR: failed to connect to Rest API, please check the connection and try again.")
		utils.PrintfStdErr("ERR: if you are using a custom Rest API endpoint, please provide it via flag '--%s <your_custom>' or setting environment variable 'export %s=<your_custom>'.\n", FlagCosmosRest, constants.ENV_COSMOS_REST)
		os.Exit(1)
	}

	return
}
