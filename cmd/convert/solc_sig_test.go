package convert

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_normalizeEvmEventOrMethodInterface(t *testing.T) {
	tests := []struct {
		name       string
		_interface string
		want       string
	}{
		{
			name:       "trim space",
			_interface: " balanceOf()\n\t",
			want:       "balanceOf()",
		},
		{
			name:       "trim suffix ';'",
			_interface: " balanceOf() ; ",
			want:       "balanceOf()",
		},
		{
			name:       "trim suffix '{'",
			_interface: " balanceOf() { ",
			want:       "balanceOf()",
		},
		{
			name:       "trim suffix ';' & '{'",
			_interface: " balanceOf() { ; ",
			want:       "balanceOf()",
		},
		{
			name:       "drop after ')'",
			_interface: "function burnFrom(address account, uint256 amount) public virtual",
			want:       "burnFrom(address account,uint256 amount)",
		},
		{
			name:       "drop after ')'",
			_interface: "approve(address,(string,string,(string,uint256)[],string[])[]) public virtual",
			want:       "approve(address,(string,string,(string,uint256)[],string[])[])",
		},
		{
			name:       "remove duplicated spaces",
			_interface: "balanceOf(  \n\t)",
			want:       "balanceOf()",
		},
		{
			name:       "remove duplicated spaces",
			_interface: "balanceOf(  \t\n  )",
			want:       "balanceOf()",
		},
		{
			name:       "replace surrounding spaces",
			_interface: "balanceOf ()",
			want:       "balanceOf()",
		},
		{
			name:       "replace surrounding spaces",
			_interface: "balanceOf( address)",
			want:       "balanceOf(address)",
		},
		{
			name:       "replace surrounding spaces",
			_interface: "balanceOf(address )",
			want:       "balanceOf(address)",
		},
		{
			name:       "replace surrounding spaces",
			_interface: "approve(address,(string,string,(string,uint256)[],string[]) )",
			want:       "approve(address,(string,string,(string,uint256)[],string[]))",
		},
		{
			name:       "replace surrounding spaces",
			_interface: "approve( address a, (string,string,(string,uint256)[],string[]) b)",
			want:       "approve(address a,(string,string,(string,uint256)[],string[]) b)",
		},
		{
			name:       "remove keyword 'function'",
			_interface: "function  balanceOf(address)",
			want:       "balanceOf(address)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, normalizeEvmEventOrMethodInterface(tt._interface))
		})
	}
}

func Test_removeExtraSpaces(t *testing.T) {
	tests := []struct {
		name string
		str  string
		want string
	}{
		{
			name: "none",
			str:  "balanceOf()",
			want: "balanceOf()",
		},
		{
			name: "remove duplicated spaces",
			str:  "balanceOf( \n\t)",
			want: "balanceOf( )",
		},
		{
			name: "remove duplicated spaces",
			str:  "balanceOf(  \n\t)",
			want: "balanceOf( )",
		},
		{
			name: "remove duplicated spaces",
			str:  "balanceOf(  \t\n  )",
			want: "balanceOf( )",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, removeExtraSpaces(tt.str))
		})
	}
}

func Test_prepareInterfaceToHash(t *testing.T) {
	tests := []struct {
		name       string
		_interface string
		want       string
	}{
		{
			name:       "remove 'indexed' keyword",
			_interface: "balanceOf(address indexed account)",
			want:       "balanceOf(address)",
		},
		{
			name:       "remove 'indexed' keyword",
			_interface: "balanceOf(address   indexed\t\naccount)",
			want:       "balanceOf(address)",
		},
		{
			name:       "remove argument name",
			_interface: "approve(address a , address b,uint64 indexed c,(int64,(string, uint256 )[])[],(int64,(string,uint256)[])[]  d)",
			want:       "approve(address,address,uint64,(int64,(string,uint256)[])[],(int64,(string,uint256)[])[])",
		},
		{
			name:       "remove argument name",
			_interface: "approve((int64,(string,uint256)[])[],(int64,(string ,uint256)[])[] indexed d, address a ,address b,uint64 indexed c)",
			want:       "approve((int64,(string,uint256)[])[],(int64,(string,uint256)[])[],address,address,uint64)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_interface, err := prepareInterfaceToHash(tt._interface)
			require.NoError(t, err)
			require.Equal(t, tt.want, _interface)
		})
	}
}

func Test_getSignatureFromInterface(t *testing.T) {
	tests := []struct {
		_interface    string
		wantSignature string
		wantErr       bool
	}{
		{
			_interface:    "function transfer(address recipient, uint256 amount) public virtual override returns (bool) {",
			wantSignature: "0xa9059cbb",
			wantErr:       false,
		},
		{
			_interface:    "function transfer(address recipient, uint256 amount) public virtual override returns (bool);",
			wantSignature: "0xa9059cbb",
			wantErr:       false,
		},
		{
			_interface:    "transfer(address recipient, uint256 amount)",
			wantSignature: "0xa9059cbb",
			wantErr:       false,
		},
		{
			_interface:    "transfer(address, uint256)",
			wantSignature: "0xa9059cbb",
			wantErr:       false,
		},
		{
			_interface:    "_updateList(address[],address,address[])",
			wantSignature: "0x00199b79",
			wantErr:       false,
		},
		{
			_interface:    "proposeRepeated((address,bytes)[],uint256)",
			wantSignature: "0x013a652d",
			wantErr:       false,
		},
		{
			_interface:    "addSsTokensToSwap((address,address,bool,int128,int128)[])",
			wantSignature: "0x015d04c9",
			wantErr:       false,
		},
		{
			_interface:    "removeLiquidityWithPermit(address,address,uint256,uint256,uint256,address,uint256,bool,uint8,bytes32,bytes32,((address,address,address,uint256,uint256,address,uint256,uint8,bytes32,bytes32),uint256,address[])[])",
			wantSignature: "0x068694c6",
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt._interface, func(t *testing.T) {
			gotSignature, gotHash, _, err := getSignatureFromInterface(tt._interface)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantSignature, gotSignature)
			require.Truef(t, strings.HasPrefix(gotHash.Hex(), tt.wantSignature), "want hash %s has prefix %s", gotHash.Hex(), tt.wantSignature)
		})
	}
}
