package utils

import (
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

func GetEvmAddressFromAnyFormatAddress(addrs ...string) (evmAddrs []common.Address, err error) {
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
