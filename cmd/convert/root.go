package convert

import (
	"github.com/spf13/cobra"
)

// Commands registers a sub-tree of commands
func Commands() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "convert",
		Aliases: []string{"c"},
		Short:   "Convert commands",
	}

	cmd.AddCommand(
		GetConvertAddressCmd(),
		GetConvertAbiStringCmd(),
		GetConvertDecimalToHexadecimalCmd(),
		GetConvertHexadecimalToDecimalCmd(),
		GetConvertSolcSignatureCmd(),
		GetConvertToLowerCaseCmd(),
		GetConvertToUpperCaseCmd(),
		GetDecodeBase64CaseCmd(),
		GetEncodeBase64CaseCmd(),
		GetDisplayBalanceCmd(),
		GetRawBalanceCmd(),
	)

	return cmd
}
