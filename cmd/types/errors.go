package types

import (
	"github.com/pkg/errors"
	"strings"
)

// ErrUpstreamRpcReturnedError is the error when the upstream RPC returned error
var ErrUpstreamRpcReturnedError = errors.New("upstream RPC returned error")

// IsErrUpstreamRpcReturnedError returns true if the error is built from upstream RPC response error,
// when error code is not zero.
func IsErrUpstreamRpcReturnedError(err error) bool {
	return strings.Contains(err.Error(), ErrUpstreamRpcReturnedError.Error())
}
