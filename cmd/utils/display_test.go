package utils

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestConvertNumberIntoDisplayWithExponent(t *testing.T) {
	tests := []struct {
		name        string
		number      *big.Int
		exponent    int
		wantDisplay string
		wantHigh    *big.Int
		wantLow     *big.Int
		wantErr     bool
	}{
		{
			name:        "normal, 6 exponent, lower than 1",
			number:      big.NewInt(101),
			exponent:    6,
			wantDisplay: "0.000101",
			wantHigh:    common.Big0,
			wantLow:     big.NewInt(101),
			wantErr:     false,
		},
		{
			name:        "normal, 6 exponent, greater than 1",
			number:      big.NewInt(2_000_101),
			exponent:    6,
			wantDisplay: "2.000101",
			wantHigh:    big.NewInt(2),
			wantLow:     big.NewInt(101),
			wantErr:     false,
		},
		{
			name:        "normal, 18 exponent, lower than 1",
			number:      big.NewInt(101),
			exponent:    18,
			wantDisplay: "0.000000000000000101",
			wantHigh:    common.Big0,
			wantLow:     big.NewInt(101),
			wantErr:     false,
		},
		{
			name:        "normal, 18 exponent, greater than 1",
			number:      big.NewInt(2_000_000_000_000_000_101),
			exponent:    18,
			wantDisplay: "2.000000000000000101",
			wantHigh:    big.NewInt(2),
			wantLow:     big.NewInt(101),
			wantErr:     false,
		},
		{
			name:     "negative exponent",
			number:   big.NewInt(100),
			exponent: -1,
			wantErr:  true,
		},
		{
			name:     "over ranged exponent",
			number:   big.NewInt(100),
			exponent: 19,
			wantErr:  true,
		},
		{
			name:    "negative number",
			number:  big.NewInt(-1),
			wantErr: true,
		},
		{
			name:        "zero number, positive exponent",
			number:      common.Big0,
			exponent:    6,
			wantDisplay: "0.0",
			wantHigh:    common.Big0,
			wantLow:     common.Big0,
			wantErr:     false,
		},
		{
			name:        "positive number, zero exponent",
			number:      big.NewInt(100),
			exponent:    0,
			wantDisplay: "100.0",
			wantHigh:    big.NewInt(100),
			wantLow:     common.Big0,
			wantErr:     false,
		},
		{
			name:        "zero low number, positive exponent",
			number:      big.NewInt(1_000_000),
			exponent:    6,
			wantDisplay: "1.0",
			wantHigh:    common.Big1,
			wantLow:     common.Big0,
			wantErr:     false,
		},
		{
			name:        "zero low number, positive exponent",
			number:      big.NewInt(10_000_000),
			exponent:    6,
			wantDisplay: "10.0",
			wantHigh:    big.NewInt(10),
			wantLow:     common.Big0,
			wantErr:     false,
		},
		{
			name:        "low number need padding left zero",
			number:      big.NewInt(10_000_111),
			exponent:    6,
			wantDisplay: "10.000111",
			wantHigh:    big.NewInt(10),
			wantLow:     big.NewInt(111),
			wantErr:     false,
		},
		{
			name:        "low number truncate right zero",
			number:      big.NewInt(10_111_000),
			exponent:    6,
			wantDisplay: "10.111",
			wantHigh:    big.NewInt(10),
			wantLow:     big.NewInt(111_000),
			wantErr:     false,
		},
		{
			name:        "low number padding left zero and truncate right zero",
			number:      big.NewInt(10_011_100),
			exponent:    6,
			wantDisplay: "10.0111",
			wantHigh:    big.NewInt(10),
			wantLow:     big.NewInt(11_100),
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDisplay, gotHigh, gotLow, err := ConvertNumberIntoDisplayWithExponent(tt.number, tt.exponent)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantDisplay, gotDisplay)
			require.Equal(t, tt.wantHigh, gotHigh)
			require.Equal(t, tt.wantLow, gotLow)
		})
	}
}
