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
		GetConvertHexadecimalCmd(),
		GetConvertSolcSignatureCmd(),
		GetConvertCaseCmd(),
		GetBase64CaseCmd(),
		GetDisplayBalanceCmd(),
		GetRawBalanceCmd(),
	)

	return cmd
}
