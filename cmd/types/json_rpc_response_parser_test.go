package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseJsonRpcResponse(t *testing.T) {
	type sampleResult struct {
		Field1 string `json:"field1"`
	}

	t.Run("normal", func(t *testing.T) {
		r, err := ParseJsonRpcResponse[sampleResult]([]byte(`{
    "jsonrpc": "2.0",
    "id": 1,
    "result": {
        "field1": "f1"
    }
}`))
		require.NoError(t, err)
		require.NotNil(t, r)

		rsr := r.(*sampleResult)
		require.Equal(t, "f1", rsr.Field1)
	})

	t.Run("normal", func(t *testing.T) {
		r, err := ParseJsonRpcResponse[sampleResult]([]byte(`{
    "jsonrpc": "2.0",
    "id": 1,
    "result": {
    }
}`))
		require.NoError(t, err)
		require.NotNil(t, r)

		rsr := r.(*sampleResult)
		require.Empty(t, rsr.Field1)
	})

	t.Run("no result", func(t *testing.T) {
		r, err := ParseJsonRpcResponse[sampleResult]([]byte(`{
    "jsonrpc": "2.0",
    "id": 1
}`))
		require.Error(t, err)
		require.Nil(t, r)
		require.ErrorContains(t, err, "missing response")
		require.False(t, IsErrUpstreamRpcReturnedError(err))
	})

	t.Run("result with error", func(t *testing.T) {
		r, err := ParseJsonRpcResponse[sampleResult]([]byte(`{
    "jsonrpc": "2.0",
    "id": 1,
	"error": {
        "code": -32601,
        "message": "the method abc does not exist/is not available"
    }
}`))
		require.Error(t, err)
		require.Nil(t, r)
		require.ErrorContains(t, err, "the method abc does not exist/is not available")
		require.ErrorContains(t, err, "-32601")
		require.True(t, IsErrUpstreamRpcReturnedError(err))
	})

	type complexResultL4 struct {
		Message string `json:"message"`
		Number  int64  `json:"number"`
	}
	type complexResultL3 struct {
		Level41 complexResultL4  `json:"level41"`
		Level42 *complexResultL4 `json:"level42"`
		Level43 complexResultL4  `json:"level43"`
		Message string           `json:"message"`
	}
	type complexResultL2 struct {
		Level3  complexResultL3 `json:"level3"`
		Message string          `json:"message"`
	}
	type complexResultL1 struct {
		Level2  complexResultL2 `json:"level2"`
		Message string          `json:"message"`
	}

	t.Run("multi-nested", func(t *testing.T) {
		r, err := ParseJsonRpcResponse[complexResultL1]([]byte(`{
    "jsonrpc": "2.0",
    "id": 1,
    "result": {
        "message": "l1",
		"level2": {
			"message": "l2",
			"level3": {
				"message": "l3",
				"level41": {
					"message": "l4",
					"number": 77,
					"float": 1.6
				}
			}
		}
    }
}`))
		require.NoError(t, err)
		require.NotNil(t, r)

		cr := r.(*complexResultL1)
		require.Equal(t, "l1", cr.Message)
		require.Equal(t, "l2", cr.Level2.Message)
		require.Equal(t, "l3", cr.Level2.Level3.Message)
		require.Equal(t, "l4", cr.Level2.Level3.Level41.Message)
		require.Equal(t, int64(77), cr.Level2.Level3.Level41.Number)
		require.Nil(t, cr.Level2.Level3.Level42)
		require.Empty(t, cr.Level2.Level3.Level43.Message)
	})
}
