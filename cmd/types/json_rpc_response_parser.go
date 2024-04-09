package types

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
)

func ParseJsonRpcResponse[T any](bz []byte) (any, error) {
	type genericResponseErrorStruct struct {
		Code    int64  `json:"code"`
		Message string `json:"message"`
	}
	type genericResponseStruct[T any] struct {
		JsonRpc string                     `json:"json-rpc"`
		Id      uint64                     `json:"id"`
		Result  *T                         `json:"result,omitempty"`
		Err     genericResponseErrorStruct `json:"error,omitempty"`
	}

	var response genericResponseStruct[T]
	err := json.Unmarshal(bz, &response)
	if err != nil {
		return nil, err
	}

	if response.Err.Code != 0 {
		return nil, errors.Wrapf(ErrUpstreamRpcReturnedError, "error code: %d, message: %s", response.Err.Code, response.Err.Message)
	}

	if response.Result == nil {
		return nil, fmt.Errorf("missing response")
	}

	return response.Result, nil
}
