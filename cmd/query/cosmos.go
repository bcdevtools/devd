package query

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/bcdevtools/devd/v3/cmd/utils"
	"github.com/bcdevtools/devd/v3/constants"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func mustGetRest(cmd *cobra.Command) (rest string) {
	var inputSource string

	if restFromFlagRest, _ := cmd.Flags().GetString(flagRest); len(restFromFlagRest) > 0 {
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
		utils.PrintfStdErr("ERR: if you are using a custom Rest API endpoint, please provide it via flag '--%s <your_custom>' or setting environment variable 'export %s=<your_custom>'.\n", flagRest, constants.ENV_COSMOS_REST)
		os.Exit(1)
	}

	return
}

func fetchErc20ModuleTokenPairsFromRest(rest string) (erc20ModuleTokenPairs []Erc20ModuleTokenPair, statusCode int, err error) {
	var resp *http.Response
	resp, err = http.Get(rest + "/evmos/erc20/v1/token_pairs")
	if err != nil {
		err = errors.Wrap(err, "failed to fetch ERC-20 module token pairs")
		return
	}

	statusCode = resp.StatusCode

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("failed to fetch ERC-20 module token pairs! Status code: %d", resp.StatusCode)
		return
	}

	type responseStruct struct {
		TokenPairs []Erc20ModuleTokenPair `json:"token_pairs"`
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	bz, err := io.ReadAll(resp.Body)
	if err != nil {
		err = errors.Wrap(err, "failed to read response body of ERC-20 module token pairs")
		return
	}

	var response responseStruct
	err = json.Unmarshal(bz, &response)
	if err != nil {
		err = errors.Wrap(err, "failed to unmarshal response body of ERC-20 module token pairs")
		return
	}

	erc20ModuleTokenPairs = response.TokenPairs
	return
}

type Erc20ModuleTokenPair struct {
	Erc20Address string `json:"erc20_address"`
	Denom        string `json:"denom"`
	Enabled      bool   `json:"enabled"`
}

func fetchVirtualFrontierBankContractPairsFromRest(rest string) (vfbcPairs []VfbcTokenPair, statusCode int, err error) {
	var resp *http.Response
	resp, err = http.Get(rest + "/ethermint/evm/v1/virtual_frontier_bank_contracts")
	if err != nil {
		err = errors.Wrap(err, "failed to fetch VFBC pairs")
		return
	}

	statusCode = resp.StatusCode

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("failed to fetch VFBC pairs! Status code: %d", resp.StatusCode)
		return
	}

	type responseStruct struct {
		Pairs []VfbcTokenPair `json:"pairs"`
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	bz, err := io.ReadAll(resp.Body)
	if err != nil {
		err = errors.Wrap(err, "failed to read response body of VFBC pairs")
		return
	}

	var response responseStruct
	err = json.Unmarshal(bz, &response)
	if err != nil {
		err = errors.Wrap(err, "failed to unmarshal response body of VFBC pairs")
		return
	}

	vfbcPairs = response.Pairs
	return
}

type VfbcTokenPair struct {
	ContractAddress string `json:"contract_address"`
	MinDenom        string `json:"min_denom"`
	Enabled         bool   `json:"enabled"`
}
