package flags

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/bcdevtools/devd/v3/cmd/utils"
	"github.com/spf13/cobra"
)

const (
	FlagHeight = "height"
)

// ReadFlagShortIntOrHexOrNil reads a flag value as decimal or hexadecimal and returns a big.Int.
// If the flag value is empty, it returns nil.
func ReadFlagShortIntOrHexOrNil(cmd *cobra.Command, flag string) (*big.Int, error) {
	value, _ := cmd.Flags().GetString(flag)
	if value == "" {
		return nil, nil
	}

	return utils.ReadShortIntOrHex(value)
}

// ReadFlagShortIntOrHexOrZero reads a flag value as decimal or hexadecimal and returns an uint64.
// If the flag value is empty, it returns zero.
func ReadFlagShortIntOrHexOrZero(cmd *cobra.Command, flag string) (finalValue uint64, err error) {
	var bi *big.Int
	defer func() {
		if err == nil && bi != nil {
			if !bi.IsUint64() {
				err = fmt.Errorf("value is out of range: %s", bi.String())
			} else {
				finalValue = bi.Uint64()
			}
		}
	}()

	value, _ := cmd.Flags().GetString(flag)
	if value == "" {
		return
	}

	bi, err = utils.ReadShortIntOrHex(value)
	return
}

// ReadFlagBlockNumberOrNil reads a flag value as block number and returns a big.Int.
// If the flag value is empty or zero or "latest", it returns nil.
func ReadFlagBlockNumberOrNil(cmd *cobra.Command, flag string) (*big.Int, error) {
	heightStr, _ := cmd.Flags().GetString(flag)
	heightStr = strings.TrimSpace(strings.ToLower(heightStr))
	if heightStr == "" || heightStr == "latest" {
		return nil, nil
	}

	height, err := ReadFlagShortIntOrHexOrNil(cmd, flag)
	if err != nil {
		return nil, err
	}

	if height.Sign() == 0 {
		return nil, nil
	}

	return height, nil
}
