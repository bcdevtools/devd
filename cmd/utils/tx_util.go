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
	InjectFrom                bool
	InjectTranslateAbleFields bool
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
		if option.InjectTranslateAbleFields {
			tryInjectHexTranslatedFieldForEvmTx(_map, "chainId")
			tryInjectHexTranslatedFieldForEvmTx(_map, "gas")
			tryInjectHexTranslatedFieldForEvmTx(_map, "gasPrice")
			tryInjectHexTranslatedFieldForEvmTx(_map, "maxFeePerGas")
			tryInjectHexTranslatedFieldForEvmTx(_map, "maxPriorityFeePerGas")
			tryInjectHexTranslatedFieldForEvmTx(_map, "nonce")
			tryInjectHexTranslatedFieldForEvmTx(_map, "value")
			tryInjectTranslatedTypeForEvmTx(_map, "type")
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

type PrettyMarshalJsonEvmTxReceiptOption struct {
	InjectTranslateAbleFields bool
}

func MarshalPrettyJsonEvmTxReceipt(receipt *ethtypes.Receipt, option *PrettyMarshalJsonEvmTxReceiptOption) ([]byte, error) {
	bz, err := receipt.MarshalJSON()
	if err != nil {
		return nil, err
	}

	var _map map[string]interface{}
	err = json.Unmarshal(bz, &_map)
	if err != nil {
		return nil, err
	}

	if option != nil {
		if option.InjectTranslateAbleFields {
			tryInjectTranslatedTypeForEvmTx(_map, "type")
			if receipt.Status == ethtypes.ReceiptStatusSuccessful {
				_map["_status"] = "Success"
			} else {
				_map["_status"] = "Failed"
			}
			tryInjectHexTranslatedFieldForEvmTx(_map, "cumulativeGasUsed")
			tryInjectHexTranslatedFieldForEvmTx(_map, "gasUsed")
			tryInjectHexTranslatedFieldForEvmTx(_map, "blockNumber")
			tryInjectHexTranslatedFieldForEvmTx(_map, "transactionIndex")
		}
	}

	return json.MarshalIndent(_map, "", "  ")
}

func tryInjectTranslatedTypeForEvmTx(_map map[string]interface{}, key string) {
	keyTransform := func(key string) string {
		return "_" + key
	}

	valueTransform := func(v interface{}) (interface{}, bool) {
		if v == nil {
			return nil, false
		}

		hexStr, ok := v.(string)
		if !ok || !strings.HasPrefix(hexStr, "0x") || len(hexStr) < 3 {
			return nil, false
		}

		switch hexStr {
		case "0x0":
			return "Legacy", true
		case "0x1":
			return "Access List", true
		case "0x2":
			return "Dynamic Fee (EIP-1559)", true
		default:
			return nil, false
		}
	}

	TryInjectTranslationValueByKey(_map, key, keyTransform, valueTransform)
}

func tryInjectHexTranslatedFieldForEvmTx(_map map[string]interface{}, key string) {
	keyTransform := func(key string) string {
		return "_" + key
	}

	valueTransform := func(v interface{}) (interface{}, bool) {
		if v == nil {
			return nil, false
		}

		hexStr, ok := v.(string)
		if !ok || !strings.HasPrefix(hexStr, "0x") || len(hexStr) < 3 {
			return nil, false
		}
		hexStr = hexStr[2:]

		bi, ok := new(big.Int).SetString(hexStr, 16)
		if !ok {
			return nil, false
		}

		return bi.String(), true
	}

	TryInjectTranslationValueByKey(_map, key, keyTransform, valueTransform)
}
