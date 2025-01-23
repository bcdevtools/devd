package tx

import (
	"fmt"
	"github.com/bcdevtools/devd/v2/cmd/utils"
	"github.com/bcdevtools/devd/v2/constants"
	"github.com/spf13/cobra"
	"math/big"
	"regexp"
	"strings"
)

const (
	flagRpc       = "rpc"
	flagSecretKey = "secret-key"
	flagGasLimit  = "gas"
	flagGasPrices = "gas-prices"
)

const (
	flagEvmRpcDesc    = "EVM Json-RPC endpoint, default is " + constants.DEFAULT_EVM_RPC + ", can be set by environment variable " + constants.ENV_EVM_RPC
	flagSecretKeyDesc = "Secret private key or mnemonic of the account, can be set by environment variable " + constants.ENV_SECRET_KEY
	flagGasLimitDesc  = "Gas limit for the transaction, support custom unit (eg: 1m equals to one million, 21k equals to thousand)"
	flagGasPricesDesc = "Gas prices for the transaction, support custom unit (eg: both 20b and 20g(wei) equals to twenty billion)"
)

// Commands registers a sub-tree of commands
func Commands() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tx",
		Short: "Tx commands",
	}

	cmd.AddCommand(
		GetSendEvmTxCommand(),
		GetDeployContractEvmTxCommand(),
		GetDeployErc20EvmTxCommand(),
	)

	return cmd
}

func readGasPrices(cmd *cobra.Command) (*big.Int, error) {
	gasPrices, _ := cmd.Flags().GetString(flagGasPrices)
	if gasPrices == "" {
		gasPrices = "20b"
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

	bi, err := utils.ReadCustomInteger(gasPrices)
	if err != nil {
		return nil, err
	}

	return bi, nil
}

func readGasLimit(cmd *cobra.Command) (uint64, error) {
	gasLimit, _ := cmd.Flags().GetString(flagGasLimit)
	if gasLimit == "" {
		gasLimit = "500k"
	}

	bi, err := utils.ReadCustomInteger(gasLimit)
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
