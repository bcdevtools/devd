package tx

import (
	"bytes"
	"encoding/hex"

	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/bcdevtools/devd/v3/cmd/utils"
	"github.com/spf13/cobra"
)

const (
	flagRawTx = "raw-tx"
)

const (
	flagGasLimitDesc  = "Gas limit for the transaction, support custom unit (eg: 1m equals to one million, 21k equals to thousand)"
	flagGasPricesDesc = "Gas prices for the transaction, support custom unit (eg: both 20b and 20g(wei) equals to twenty billion)"
	flagRawTxDesc     = "Print raw tx"
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
	)

	return cmd
}

func printRawEvmTx(signedTx *ethtypes.Transaction) {
	var buf bytes.Buffer
	err := signedTx.EncodeRLP(&buf)
	utils.ExitOnErr(err, "failed to encode tx")

	rawTxRlpEncoded := hex.EncodeToString(buf.Bytes())
	utils.PrintfStdErr("INF: Raw EVM tx: 0x%s\n", rawTxRlpEncoded)
}
