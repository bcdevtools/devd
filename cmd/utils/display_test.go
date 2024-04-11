package utils

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestConvertNumberIntoDisplayWithExponentAndViceVersa(t *testing.T) {
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

			// reverse

			gotNumber, gotHighNumber, gotLowNumber, err := ConvertDisplayWithExponentIntoRaw(gotDisplay, tt.exponent, '.')
			require.NoError(t, err)
			require.Equal(t, tt.number, gotNumber)
			require.Equal(t, tt.wantHigh, gotHighNumber)
			require.Equal(t, tt.wantLow, gotLowNumber)
		})
	}
}

func TestConvertDisplayWithExponentIntoRawAndViceVersa(t *testing.T) {
	tests := []struct {
		name                string
		display             string
		exponent            int
		customDecimalsPoint rune
		wantNumber          *big.Int
		wantHighNumber      *big.Int
		wantLowNumber       *big.Int
		wantErr             bool
		wantReverseDisplay  string
	}{
		{
			name:           "normal",
			display:        "1.0",
			exponent:       6,
			wantNumber:     big.NewInt(1_000_000),
			wantHighNumber: common.Big1,
			wantLowNumber:  common.Big0,
			wantErr:        false,
		},
		{
			name:               "normal, multiple zero tails",
			display:            "1.0000",
			exponent:           6,
			wantNumber:         big.NewInt(1_000_000),
			wantHighNumber:     common.Big1,
			wantLowNumber:      common.Big0,
			wantErr:            false,
			wantReverseDisplay: "1.0",
		},
		{
			name:           "normal, padding right",
			display:        "1.1",
			exponent:       6,
			wantNumber:     big.NewInt(1_100_000),
			wantHighNumber: common.Big1,
			wantLowNumber:  big.NewInt(100_000),
			wantErr:        false,
		},
		{
			name:           "normal, 6 exponent, lower than 1",
			display:        "0.000101",
			exponent:       6,
			wantNumber:     big.NewInt(101),
			wantHighNumber: common.Big0,
			wantLowNumber:  big.NewInt(101),
			wantErr:        false,
		},
		{
			name:           "normal, 6 exponent, greater than 1",
			display:        "2.000101",
			exponent:       6,
			wantNumber:     big.NewInt(2_000_101),
			wantHighNumber: big.NewInt(2),
			wantLowNumber:  big.NewInt(101),
			wantErr:        false,
		},
		{
			name:           "normal, 18 exponent, lower than 1",
			display:        "0.000000000000000101",
			exponent:       18,
			wantNumber:     big.NewInt(101),
			wantHighNumber: common.Big0,
			wantLowNumber:  big.NewInt(101),
			wantErr:        false,
		},
		{
			name:           "normal, 18 exponent, greater than 1",
			display:        "2.000000000000000101",
			exponent:       18,
			wantNumber:     big.NewInt(2_000_000_000_000_000_101),
			wantHighNumber: big.NewInt(2),
			wantLowNumber:  big.NewInt(101),
			wantErr:        false,
		},
		{
			name:                "normal, 6 exponent, lower than 1, custom decimals point",
			display:             "0,000101",
			exponent:            6,
			customDecimalsPoint: ',',
			wantNumber:          big.NewInt(101),
			wantHighNumber:      common.Big0,
			wantLowNumber:       big.NewInt(101),
			wantErr:             false,
			wantReverseDisplay:  "0.000101",
		},
		{
			name:                "normal, 6 exponent, greater than 1, custom decimals point",
			display:             "2,000101",
			exponent:            6,
			customDecimalsPoint: ',',
			wantNumber:          big.NewInt(2_000_101),
			wantHighNumber:      big.NewInt(2),
			wantLowNumber:       big.NewInt(101),
			wantErr:             false,
			wantReverseDisplay:  "2.000101",
		},
		{
			name:     "negative exponent",
			display:  "0.000101",
			exponent: -1,
			wantErr:  true,
		},
		{
			name:     "over ranged exponent",
			display:  "0.000101",
			exponent: 19,
			wantErr:  true,
		},
		{
			name:     "multiple decimals points",
			display:  "0.0.00101",
			exponent: 6,
			wantErr:  true,
		},
		{
			name:                "multiple decimals points",
			display:             "0,0,00101",
			customDecimalsPoint: ',',
			exponent:            6,
			wantErr:             true,
		},
		{
			name:     "left part is not a number",
			display:  "a.000101",
			exponent: 6,
			wantErr:  true,
		},
		{
			name:     "right part is not a number",
			display:  "0.000101a",
			exponent: 6,
			wantErr:  true,
		},
		{
			name:     "left part is not a number, custom decimals point",
			display:  "a,000101",
			exponent: 6,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var decimalsPoint rune
			if tt.customDecimalsPoint != 0 {
				decimalsPoint = tt.customDecimalsPoint
			} else {
				decimalsPoint = '.'
			}

			gotNumber, gotHighNumber, gotLowNumber, err := ConvertDisplayWithExponentIntoRaw(tt.display, tt.exponent, decimalsPoint)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantNumber, gotNumber)
			require.Equal(t, tt.wantHighNumber, gotHighNumber)
			require.Equal(t, tt.wantLowNumber, gotLowNumber)

			// reverse

			var reverseDisplay string
			if tt.wantReverseDisplay != "" {
				reverseDisplay = tt.wantReverseDisplay
			} else {
				reverseDisplay = tt.display
			}

			gotDisplay, gotHigh, gotLow, err := ConvertNumberIntoDisplayWithExponent(gotNumber, tt.exponent)

			require.NoError(t, err)
			require.Equal(t, reverseDisplay, gotDisplay)
			require.Equal(t, tt.wantHighNumber, gotHigh)
			require.Equal(t, tt.wantLowNumber, gotLow)
		})
	}
}
