package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// BeautifyJson beautifies the given json.
func BeautifyJson(bzJson []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, bzJson, "", "  ")
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

// TryPrintBeautyJson tries to beautify the given json and print it.
// If failed to beautify, it will print the original json.
func TryPrintBeautyJson(bz []byte) {
	beautifyBz, err := BeautifyJson(bz)
	if err != nil {
		PrintlnStdErr("ERR: Failed to beautify json:", err)
		fmt.Println(string(bz))
	} else {
		fmt.Println(string(beautifyBz))
	}
}
