package utils

import "regexp"

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
