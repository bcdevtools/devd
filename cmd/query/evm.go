package query

import (
	"bytes"
	"encoding/hex"
	"fmt"
	libutils "github.com/EscanBE/go-lib/utils"
	"github.com/bcdevtools/devd/cmd/types"
	"github.com/bcdevtools/devd/cmd/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"io"
	"math/big"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

func doQuery(host string, qb types.JsonRpcQueryBuilder, optionalTimeout time.Duration) ([]byte, error) {
	var timeout = optionalTimeout
	if optionalTimeout == 0 {
		timeout = 5 * time.Second
	}
	if timeout < time.Second {
		timeout = time.Second
	}

	httpClient := http.Client{
		Timeout: timeout * time.Millisecond,
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

func decodeResponseToString(bz []byte, desc string) (string, error) {
	var resAny any
	resAny, err := types.ParseJsonRpcResponse[string](bz)
	if err != nil {
		return "", err
	}

	res := resAny.(*string)

	resStr := *res

	fmt.Println("Decoding response to string...", resStr)

	bz, err = hex.DecodeString(resStr[2:])
	if err != nil {
		libutils.PrintfStdErr("ERR: failed to decode hex response %s: %v\n", desc, err)
		os.Exit(1)
	}

	result, err := utils.AbiDecodeString(bz)
	if err != nil {
		libutils.PrintfStdErr("ERR: failed to decode ABI response %s: %v\n", desc, err)
		os.Exit(1)
	}

	return result, nil
}

func decodeResponseToBigInt(bz []byte, desc string) (*big.Int, error) {
	var resAny any
	resAny, err := types.ParseJsonRpcResponse[string](bz)
	if err != nil {
		return nil, err
	}

	res := resAny.(*string)

	resStr := *res

	bz, err = hex.DecodeString(resStr[2:])
	if err != nil {
		libutils.PrintfStdErr("ERR: failed to decode hex response %s: %v\n", desc, err)
		os.Exit(1)
	}

	return new(big.Int).SetBytes(bz), nil
}
