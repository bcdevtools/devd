package flags

import (
	"fmt"
	"math/big"
	"regexp"
	"strings"

	"github.com/bcdevtools/devd/v3/cmd/utils"
	"github.com/spf13/cobra"
)

const (
	FlagGasLimit  = "gas"
	FlagGasPrices = "gas-prices"
)

func ReadFlagGasLimit(cmd *cobra.Command, flag string, _default uint64) (uint64, error) {
	gasLimit, _ := cmd.Flags().GetString(flag)
	if gasLimit == "" {
		gasLimit = fmt.Sprintf("%d", _default)
	}

	bi, err := utils.ReadShortInt(gasLimit)
	if err != nil {
		return 0, err
	}

	if !bi.IsUint64() {
		return 0, fmt.Errorf("invalid gas limit %s", gasLimit)
	}

	num := bi.Uint64()
	if num < 21_000 {
		return 0, fmt.Errorf("minimum gas limit is 21k, too low: %s", gasLimit)
	}
	if num > 35_000_000 {
		return 0, fmt.Errorf("gas limit is too high: %s", gasLimit)
	}

	return num, nil
}

func ReadFlagGasPrices(cmd *cobra.Command, flag string, _default uint64) (*big.Int, error) {
	gasPrices, _ := cmd.Flags().GetString(flag)
	if gasPrices == "" {
		gasPrices = fmt.Sprintf("%d", _default)
	}

	if regexp.MustCompile(`^\d+g$`).MatchString(gasPrices) {
		gasPrices = strings.TrimSuffix(gasPrices, "g")
		bi, ok := new(big.Int).SetString(gasPrices, 10)
		if !ok {
			panic("failed to parse gas prices")
		}
		bi = new(big.Int).Mul(bi, big.NewInt(1e9))
		return bi, nil
	}

	if regexp.MustCompile(`^\d+gwei$`).MatchString(gasPrices) {
		gasPrices = strings.TrimSuffix(gasPrices, "gwei")
		bi, ok := new(big.Int).SetString(gasPrices, 10)
		if !ok {
			panic("failed to parse gas prices")
		}
		bi = new(big.Int).Mul(bi, big.NewInt(1e9))
		return bi, nil
	}

	bi, err := utils.ReadShortInt(gasPrices)
	if err != nil {
		return nil, err
	}

	return bi, nil
}
