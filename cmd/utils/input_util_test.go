package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReadCustomInteger(t *testing.T) {
	tests := []struct {
		input   string
		wantOut string
		wantErr bool
	}{
		{
			input:   "1",
			wantOut: "1",
		},
		{
			input:   "-1",
			wantOut: "-1",
		},
		{
			input:   "1234",
			wantOut: "1234",
		},
		{
			input:   "1e18",
			wantOut: "1000000000000000000",
		},
		{
			input:   "-1e18",
			wantOut: "-1000000000000000000",
		},
		{
			input:   "33e18",
			wantOut: "33000000000000000000",
		},
		{
			input:   "23k",
			wantOut: "23000",
		},
		{
			input:   "23.1k",
			wantOut: "23100",
		},
		{
			input:   "23m",
			wantOut: "23000000",
		},
		{
			input:   "23.001m",
			wantOut: "23001000",
		},
		{
			input:   "23b",
			wantOut: "23000000000",
		},
		{
			input:   "23.123456b",
			wantOut: "23123456000",
		},
		{
			input:   "-35.123456b",
			wantOut: "-35123456000",
		},
		{
			input:   "53kb",
			wantOut: "53000000000000",
		},
		{
			input:   "53mb",
			wantOut: "53000000000000000",
		},
		{
			input:   "53bb",
			wantOut: "53000000000000000000",
		},
		{
			input:   "53.001bb",
			wantOut: "53001000000000000000",
		},
		{
			input:   "-53.001bb",
			wantOut: "-53001000000000000000",
		},
		{
			input:   "53kbb",
			wantOut: "53000000000000000000000",
		},
		{
			input:   "-53bbk",
			wantOut: "-53000000000000000000000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.wantOut, func(t *testing.T) {
			gotOut, err := ReadCustomInteger(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantOut, gotOut.String())
		})
	}
}
