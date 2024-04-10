package convert

import (
	"encoding/base64"
	"fmt"
	libutils "github.com/EscanBE/go-lib/utils"
	"github.com/spf13/cobra"
)

// GetDecodeBase64CaseCmd creates a helper command that decode base64
func GetDecodeBase64CaseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "decode_base64 [base64]",
		Short: "Decode base64",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			data, err := base64.StdEncoding.DecodeString(args[0])
			libutils.ExitIfErr(err, "failed to decode base64")

			fmt.Println(string(data))
		},
	}

	return cmd
}
