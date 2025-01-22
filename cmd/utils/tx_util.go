package utils

import (
	"encoding/hex"
	"encoding/json"
	"math/big"
	"regexp"
	"strings"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

type TxHashType uint8

const (
	TxHashTypeInvalid TxHashType = iota
	TxHashTypeEvm
	TxHashTypeCosmos
)

func DetectTxHashType(txHash string) TxHashType {
	if match, _ := regexp.MatchString(`^0x[\da-f]{64}$`, txHash); match {
		return TxHashTypeEvm
	} else if match, _ := regexp.MatchString(`^[\dA-F]{64}$`, txHash); match {
		return TxHashTypeCosmos
	} else {
		return TxHashTypeInvalid
	}
}

func DecodeRawEvmTx(rawTx string) (*ethtypes.Transaction, error) {
	rawTx = strings.TrimPrefix(rawTx, "0x")
	decoded, err := hex.DecodeString(rawTx)
	if err != nil {
		return nil, err
	}

	tx := &ethtypes.Transaction{}
	err = tx.UnmarshalBinary(decoded)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

type PrettyMarshalJsonEvmTxOption struct {
	InjectFrom               bool
	InjectHexTranslatedField bool
}

func MarshalPrettyJsonEvmTx(tx *ethtypes.Transaction, option *PrettyMarshalJsonEvmTxOption) ([]byte, error) {
	bz, err := tx.MarshalJSON()
	if err != nil {
		return nil, err
	}

	var _map map[string]interface{}
	err = json.Unmarshal(bz, &_map)
	if err != nil {
		return nil, err
	}

	if option != nil {
		if option.InjectHexTranslatedField {
			tryInjectHexTranslatedField := func(name string) {
				defer func() {
					_ = recover() // omit any panic
				}()

				v, found := _map[name]
				if !found || v == nil {
					return
				}

				hexStr, ok := v.(string)
				if !ok || !strings.HasPrefix(hexStr, "0x") || len(hexStr) < 3 {
					return
				}
				hexStr = hexStr[2:]

				bi, ok := new(big.Int).SetString(hexStr, 16)
				if !ok {
					return
				}

				_map["_"+name] = bi.String()
			}
			tryInjectHexTranslatedField("chainId")
			tryInjectHexTranslatedField("gas")
			tryInjectHexTranslatedField("gasPrice")
			tryInjectHexTranslatedField("maxFeePerGas")
			tryInjectHexTranslatedField("maxPriorityFeePerGas")
			tryInjectHexTranslatedField("nonce")
			tryInjectHexTranslatedField("value")
		}

		if option.InjectFrom {
			signer := ethtypes.LatestSignerForChainID(tx.ChainId())
			from, err := ethtypes.Sender(signer, tx)
			if err != nil {
				return nil, err
			}
			_map["_from"] = from.String()
		}
	}

	return json.MarshalIndent(_map, "", "  ")
}
