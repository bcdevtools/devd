package utils

import (
	"bytes"
	"encoding/json"
)

func BeautifyJson(bzJson []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, bzJson, "", "  ")
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}
