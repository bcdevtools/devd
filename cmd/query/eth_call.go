package query

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/bcdevtools/devd/v3/cmd/flags"
	"github.com/bcdevtools/devd/v3/cmd/utils"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

func GetQueryEvmRpcEthCallCommand() *cobra.Command {
	const flagFrom = "from"
	const flagValue = "value"

	cmd := &cobra.Command{
		Use:   "eth_call [contract address] [call data]",
		Short: "Call `eth_call` of EVM RPC: executes a new EVM message call immediately without creating a transaction on the block chain",
		Long: `Call "eth_call" of EVM RPC: executes a new EVM message call immediately without creating a transaction on the block chain.
Bech32 account address is accepted.`,
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			ethClient, _ := flags.MustGetEthClient(cmd)

			evmAddrs, err := utils.GetEvmAddressFromAnyFormatAddress(args[0])
			utils.ExitOnErr(err, "failed to parse contract address")
			contractAddress := evmAddrs[0]
			callData, err := hex.DecodeString(strings.TrimPrefix(strings.ToLower(args[1]), "0x"))
			utils.ExitOnErr(err, "failed to decode call data")

			contextHeight := func() *big.Int {
				height, err := flags.ReadFlagBlockNumberOrNil(cmd, flags.FlagHeight)
				utils.ExitOnErr(err, "failed to parse block number")
				if height != nil && height.Sign() == 1 {
					utils.PrintfStdErr("INF: using block number: %s\n", height.String())
				}
				return height
			}()

			result, err := ethClient.CallContract(context.Background(), ethereum.CallMsg{
				From: func() common.Address {
					from, _ := cmd.Flags().GetString(flagFrom)
					if from == "" {
						return common.Address{}
					}
					evmAddrs, err := utils.GetEvmAddressFromAnyFormatAddress(from)
					utils.ExitOnErr(err, fmt.Sprintf("failed to parse address of --%s", flagFrom))
					fromAddr := evmAddrs[0]
					utils.PrintfStdErr("INF: using from: %s\n", fromAddr)
					return fromAddr
				}(),
				To: &contractAddress,
				Gas: func() uint64 {
					gas, err := flags.ReadFlagShortIntOrHexOrZero(cmd, flags.FlagGasLimit)
					utils.ExitOnErr(err, "failed to parse gas")
					if gas > 0 {
						utils.PrintfStdErr("INF: using gas: %d\n", gas)
						if gas < 21000 {
							utils.PrintlnStdErr("WARN: gas is less than 21000, it may not enough for the call")
						}
					}
					return gas
				}(),
				GasPrice: func() *big.Int {
					gasPrice, err := flags.ReadFlagShortIntOrHexOrNil(cmd, flags.FlagGasPrices)
					utils.ExitOnErr(err, "failed to parse gas prices")
					if gasPrice != nil && gasPrice.Sign() == 1 {
						utils.PrintfStdErr("INF: using gas-prices: %s\n", gasPrice)
					}
					return gasPrice
				}(),
				GasFeeCap: nil,
				GasTipCap: nil,
				Value: func() *big.Int {
					value, err := flags.ReadFlagShortIntOrHexOrNil(cmd, flagValue)
					utils.ExitOnErr(err, "failed to parse value")
					if value != nil && value.Sign() == 1 {
						utils.PrintfStdErr("INF: using value: %s\n", value)
					}
					return value
				}(),
				Data:       callData,
				AccessList: nil,
			}, contextHeight)
			utils.ExitOnErr(err, "failed to call contract")

			fmt.Println("0x" + hex.EncodeToString(result))
		},
	}

	cmd.Flags().String(flags.FlagEvmRpc, "", flags.FlagEvmRpcDesc)
	cmd.Flags().StringP(flagFrom, "f", "", "the address from which the transaction is sent, support 0x or bech32")
	cmd.Flags().StringP(flags.FlagGasLimit, "g", "", "gas provided for the transaction execution, support short int and hex")
	cmd.Flags().StringP(flags.FlagGasPrices, "p", "", "gas prices used for each paid gas, support short int and hex")
	cmd.Flags().StringP(flagValue, "v", "", "value sent with this transaction, support short int and hex")
	cmd.Flags().StringP(flags.FlagHeight, "h", "latest", "the context height of the block to exec, accept \"latest\"/short int/hex")

	return cmd
}
