package types

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"math"
	"testing"
)

func Benchmark_nextRequestId(b *testing.B) {
	// Apr 03 2024: 13.87 ns/op
	for i := 0; i < b.N; i++ {
		nextRequestId()
	}
}

func Test_NewJsonRpcQueryBuilder(t *testing.T) {
	t.Run("request ID must be unique each time", func(t *testing.T) {
		const ccr = 6000
		resChan := make(chan uint64, ccr)
		defer func() {
			close(resChan)
		}()
		for i := 0; i < ccr; i++ {
			go func() {
				qb := NewJsonRpcQueryBuilder("x")
				resChan <- qb.(*jsonRpcQueryBuilder).requestId
			}()
		}
		unique := make(map[uint64]bool, ccr)
		for i := 0; i < ccr; i++ {
			reqId := <-resChan
			if _, found := unique[reqId]; found {
				t.Errorf("duplicated request ID found")
				t.FailNow()
			}
			unique[reqId] = true
		}
		require.Lenf(t, unique, ccr, "test executed wrongly")
	})
	t.Run("auto rotate request id", func(t *testing.T) {
		curReqId = math.MaxUint64 - 1
		require.Equal(t, uint64(math.MaxUint64), NewJsonRpcQueryBuilder("x").(*jsonRpcQueryBuilder).requestId)
		require.Equal(t, uint64(1), NewJsonRpcQueryBuilder("x").(*jsonRpcQueryBuilder).requestId)
	})
	t.Run("output string must match method & id", func(t *testing.T) {
		qb := NewJsonRpcQueryBuilder("x_req")
		require.Equal(t, fmt.Sprintf(`{
    "method": "x_req",
    "params": [],
    "id": %d,
    "jsonrpc": "2.0"
}`, qb.(*jsonRpcQueryBuilder).requestId), qb.String())
	})
	t.Run("output string must when param int", func(t *testing.T) {
		qb := NewJsonRpcQueryBuilder("x_reqInt", NewJsonRpcIntQueryParam(16))
		require.Equal(t, fmt.Sprintf(`{
    "method": "x_reqInt",
    "params": [16],
    "id": %d,
    "jsonrpc": "2.0"
}`, qb.(*jsonRpcQueryBuilder).requestId), qb.String())
	})
	t.Run("output string must match when param string", func(t *testing.T) {
		p, _ := NewJsonRpcStringQueryParam("aa")
		qb := NewJsonRpcQueryBuilder("x_reqStr", p)
		require.Equal(t, fmt.Sprintf(`{
    "method": "x_reqStr",
    "params": ["aa"],
    "id": %d,
    "jsonrpc": "2.0"
}`, qb.(*jsonRpcQueryBuilder).requestId), qb.String())
	})
	t.Run("output string must match when param string array", func(t *testing.T) {
		p, _ := NewJsonRpcStringArrayQueryParam("aa", "bb")
		qb := NewJsonRpcQueryBuilder("x_reqStrArr", p)
		require.Equal(t, fmt.Sprintf(`{
    "method": "x_reqStrArr",
    "params": [["aa","bb"]],
    "id": %d,
    "jsonrpc": "2.0"
}`, qb.(*jsonRpcQueryBuilder).requestId), qb.String())
	})
	t.Run("output string must match when params are mixed types", func(t *testing.T) {
		pStr, _ := NewJsonRpcStringQueryParam("aa")
		pStrArr, _ := NewJsonRpcStringArrayQueryParam("aa", "bb")
		qb := NewJsonRpcQueryBuilder(
			"x_reqMixed",
			pStr,
			pStrArr,
			NewJsonRpcInt64QueryParam(16),
		)
		require.Equal(t, fmt.Sprintf(`{
    "method": "x_reqMixed",
    "params": ["aa",["aa","bb"],16],
    "id": %d,
    "jsonrpc": "2.0"
}`, qb.(*jsonRpcQueryBuilder).requestId), qb.String())
	})
	t.Run("request id on json should maintains full-sized when big", func(t *testing.T) {
		curReqId = math.MaxUint64 - 10
		qb := NewJsonRpcQueryBuilder("x_req")
		require.Equal(t, `{
    "method": "x_req",
    "params": [],
    "id": 18446744073709551606,
    "jsonrpc": "2.0"
}`, qb.String())
	})
}

func Test_JsonRpcIntegerQueryParam(t *testing.T) {
	var p JsonRpcQueryParam

	p = NewJsonRpcInt64QueryParam(16)
	require.False(t, p.IsArray())
	require.Equal(t, "16", p.String())

	p = NewJsonRpcIntQueryParam(-16)
	require.False(t, p.IsArray())
	require.Equal(t, "-16", p.String())
}

func Test_JsonRpcStringQueryParam(t *testing.T) {
	var p JsonRpcQueryParam
	var err error

	p, err = NewJsonRpcStringQueryParam("\"16\"")
	require.Error(t, err)
	require.Nil(t, p)

	p, err = NewJsonRpcStringQueryParam("16")
	require.NoError(t, err)
	require.False(t, p.IsArray())
	require.Equal(t, `"16"`, p.String())
}

func Test_JsonRpcStringArrayQueryParam(t *testing.T) {
	var p JsonRpcQueryParam
	var err error

	p, err = NewJsonRpcStringArrayQueryParam("\"16\"")
	require.Error(t, err)
	require.Nil(t, p)

	p, err = NewJsonRpcStringArrayQueryParam("15", "\"16\"")
	require.Error(t, err)
	require.Nil(t, p)

	p, err = NewJsonRpcStringArrayQueryParam("15", "16")
	require.NoError(t, err)
	require.True(t, p.IsArray())
	require.Equal(t, `["15","16"]`, p.String())
}
