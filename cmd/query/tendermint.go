package query

import (
	"fmt"
	"os"
	"strings"

	"github.com/bcdevtools/devd/v2/cmd/utils"
	"github.com/bcdevtools/devd/v2/constants"
	httpclient "github.com/cometbft/cometbft/rpc/client/http"
	jsonrpcclient "github.com/cometbft/cometbft/rpc/jsonrpc/client"
	"github.com/spf13/cobra"
)

func mustGetTmRpc(cmd *cobra.Command) (tendermintRpcHttpClient *httpclient.HTTP, tmRpc string) {
	var inputSource string

	if tmRpcFromFlag, _ := cmd.Flags().GetString(flagTmRpc); len(tmRpcFromFlag) > 0 {
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
		utils.PrintfStdErr("ERR: if you are using a custom TM RPC endpoint, please provide it via flag '--%s <your_custom>' or setting environment variable 'export %s=<your_custom>'.\n", flagTmRpc, constants.ENV_TM_RPC)
		os.Exit(1)
	}

	return
}
