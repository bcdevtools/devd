package query

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/bcdevtools/devd/v2/cmd/types"
	"github.com/bcdevtools/devd/v2/cmd/utils"
	"github.com/bcdevtools/devd/v2/constants"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io"
	"math/big"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

const generalQueryTimeout = 3 * time.Second

func doQuery(host string, qb types.JsonRpcQueryBuilder, optionalTimeout time.Duration) ([]byte, error) {
	var timeout = optionalTimeout
	if optionalTimeout == 0 {
		timeout = generalQueryTimeout
	}
	if timeout < time.Second {
		timeout = time.Second
	}

	httpClient := http.Client{
		Timeout: timeout,
	}

	fmt.Println("Querying", host, strings.ReplaceAll(strings.ReplaceAll(qb.String(), "\n", " "), " ", ""))

	resp, err := httpClient.Post(host, "application/json", bytes.NewBuffer([]byte(qb.String())))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-OK status code: %d", resp.StatusCode)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	bz, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}

	return bz, nil
}

func getEvmAddressFromAnyFormatAddress(addrs ...string) (evmAddrs []common.Address, err error) {
	for _, addr := range addrs {
		normalizedAddr := strings.ToLower(addr)

		if regexp.MustCompile(`^(0x)?[a-f\d]{40}$`).MatchString(normalizedAddr) {
			evmAddrs = append(evmAddrs, common.HexToAddress(normalizedAddr))
		} else if regexp.MustCompile(`^(0x)?[a-f\d]{64}$`).MatchString(normalizedAddr) {
			err = fmt.Errorf("ERR: invalid address format: %s", normalizedAddr)
			return
		} else { // bech32
			spl := strings.Split(normalizedAddr, "1")
			if len(spl) != 2 || len(spl[0]) < 1 || len(spl[1]) < 1 {
				err = fmt.Errorf("ERR: invalid bech32 address: %s", normalizedAddr)
				return
			}

			var bz []byte
			bz, err = sdk.GetFromBech32(normalizedAddr, spl[0])
			if err != nil {
				err = fmt.Errorf("ERR: failed to decode bech32 address %s: %s", normalizedAddr, err)
				return
			}

			if len(bz) != 20 {
				err = fmt.Errorf("ERR: bech32 address %s has invalid length, must be 20 bytes, got %s %d bytes", normalizedAddr, hex.EncodeToString(bz), len(bz))
				return
			}

			evmAddrs = append(evmAddrs, common.BytesToAddress(bz))
		}
	}

	return
}

func mustGetEthClient(cmd *cobra.Command, fallbackDeprecatedFlagHost bool) (*ethclient.Client, string) {
	var rpc, inputSource string

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

	fmt.Println("Connecting to", rpc, fmt.Sprintf("(from %s)", inputSource))

	ethClient8545, err := ethclient.Dial(rpc)
	utils.ExitOnErr(err, "failed to connect to EVM Json-RPC")

	return ethClient8545, rpc
}

func readContextHeightFromFlag(cmd *cobra.Command) *big.Int {
	height, _ := cmd.Flags().GetInt64(flagHeight)
	if height > 0 {
		return big.NewInt(height)
	}

	return nil
}
