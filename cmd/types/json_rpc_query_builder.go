package types

import (
	"fmt"
	"github.com/bcdevtools/devd/v2/cmd/utils"
	"strings"
	"sync"
)

var syncId = sync.Mutex{}

// curReqId is the current request id, concurrent increment will be handled by syncId
var curReqId uint64 = 1

func nextRequestId() uint64 {
	syncId.Lock()
	defer syncId.Unlock()
	curReqId++
	if curReqId == 0 {
		// avoid request id = 0
		curReqId++
	}
	return curReqId
}

type JsonRpcQueryBuilder interface {
	String() string
}

var _ JsonRpcQueryBuilder = &jsonRpcQueryBuilder{}

type jsonRpcQueryBuilder struct {
	requestId   uint64
	method      string
	queryParams []JsonRpcQueryParam
}

func NewJsonRpcQueryBuilder(
	method string,
	params ...JsonRpcQueryParam,
) JsonRpcQueryBuilder {
	return &jsonRpcQueryBuilder{
		requestId:   nextRequestId(),
		method:      method,
		queryParams: params,
	}
}

func (j *jsonRpcQueryBuilder) String() string {
	return fmt.Sprintf(`{
    "method": "%s",
    "params": [%s],
    "id": %d,
    "jsonrpc": "2.0"
}`,
		j.method,
		func() string {
			if len(j.queryParams) == 0 {
				return ""
			}
			if len(j.queryParams) == 1 {
				return j.queryParams[0].String()
			}

			sb := strings.Builder{}
			sb.WriteString(j.queryParams[0].String())
			for i := 1; i < len(j.queryParams); i++ {
				sb.WriteRune(',')
				sb.WriteString(j.queryParams[i].String())
			}
			return sb.String()
		}(),
		j.requestId,
	)
}

type JsonRpcQueryParam interface {
	// IsArray returns true if the param is array formed
	IsArray() bool

	// String returns string of the param
	String() string
}

var _ JsonRpcQueryParam = jsonRpcIntegerQueryParam{}

type jsonRpcIntegerQueryParam struct {
	value string
}

func NewJsonRpcInt64QueryParam(num int64) JsonRpcQueryParam {
	return jsonRpcIntegerQueryParam{fmt.Sprintf("%d", num)}
}

func NewJsonRpcIntQueryParam(num int) JsonRpcQueryParam {
	return jsonRpcIntegerQueryParam{fmt.Sprintf("%d", num)}
}

func (j jsonRpcIntegerQueryParam) IsArray() bool {
	return false
}

func (j jsonRpcIntegerQueryParam) String() string {
	return j.value
}

var _ JsonRpcQueryParam = jsonRpcStringQueryParam{}

type jsonRpcStringQueryParam struct {
	value string
}

func NewJsonRpcStringQueryParam(str string) (JsonRpcQueryParam, error) {
	if strings.Contains(str, "\"") {
		return nil, fmt.Errorf("prohibited double quote on string param: %s", str)
	}
	return jsonRpcStringQueryParam{str}, nil
}

func (j jsonRpcStringQueryParam) IsArray() bool {
	return false
}

func (j jsonRpcStringQueryParam) String() string {
	return fmt.Sprintf("\"%s\"", j.value)
}

var _ JsonRpcQueryParam = jsonRpcStringArrayQueryParam{}

type jsonRpcStringArrayQueryParam struct {
	value []string
}

func NewJsonRpcStringArrayQueryParam(strArr ...string) (JsonRpcQueryParam, error) {
	for _, str := range strArr {
		if strings.Contains(str, "\"") {
			return nil, fmt.Errorf("prohibited double quote on string param: %s", str)
		}
	}
	return jsonRpcStringArrayQueryParam{strArr}, nil
}

func (j jsonRpcStringArrayQueryParam) IsArray() bool {
	return true
}

func (j jsonRpcStringArrayQueryParam) String() string {
	sb := strings.Builder{}
	var err error
	assertErr := func() {
		utils.PanicIfErr(err, "failed to build query params")
	}

	_, err = sb.WriteRune('[')
	assertErr()

	for i, str := range j.value {
		if i > 0 {
			_, err = sb.WriteRune(',')
			assertErr()
		}
		_, err = sb.WriteRune('"')
		assertErr()
		_, err = sb.WriteString(str)
		assertErr()
		_, err = sb.WriteRune('"')
		assertErr()
	}

	_, err = sb.WriteRune(']')
	assertErr()

	return sb.String()
}

var _ JsonRpcQueryParam = jsonRpcRawQueryParam{}

type jsonRpcRawQueryParam struct {
	value string
}

func NewJsonRpcRawQueryParam(rawStr string) JsonRpcQueryParam {
	return jsonRpcRawQueryParam{rawStr}
}

func (j jsonRpcRawQueryParam) IsArray() bool {
	return false
}

func (j jsonRpcRawQueryParam) String() string {
	return j.value
}

var _ JsonRpcQueryParam = jsonRpcBoolQueryParam{}

type jsonRpcBoolQueryParam struct {
	value string
}

func NewJsonRpcBoolQueryParam(value bool) JsonRpcQueryParam {
	return jsonRpcBoolQueryParam{fmt.Sprintf("%t", value)}
}

func (j jsonRpcBoolQueryParam) IsArray() bool {
	return false
}

func (j jsonRpcBoolQueryParam) String() string {
	return j.value
}
