package debug

import (
	"encoding/hex"
	"github.com/stretchr/testify/require"
	"testing"
)

//goland:noinspection SpellCheckingInspection
func Test_getIntrinsicGasFromInputData(t *testing.T) {
	require.Equal(t, uint64(21000), getIntrinsicGasFromInputData(nil))
	input := "0xa9059cbb000000000000000000000000aabbccddeeff112233d56f43e1699bbbe466301d0000000000000000000000000000000000000000000000000000000000000011"
	bz, err := hex.DecodeString(input[2:])
	require.NoError(t, err)
	intrinsicGas := getIntrinsicGasFromInputData(bz)
	require.Equal(t, uint64(21572), intrinsicGas)
}
