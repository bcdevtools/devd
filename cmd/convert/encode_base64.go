package convert

import (
	"encoding/base64"
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

// GetEncodeBase64CaseCmd creates a helper command that encode input into base64
func GetEncodeBase64CaseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "encode_base64 [text]",
		Aliases: []string{"base64"},
		Short:   "Encode input into base64",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(base64.StdEncoding.EncodeToString([]byte(strings.Join(args, " "))))
		},
	}

	return cmd
}
