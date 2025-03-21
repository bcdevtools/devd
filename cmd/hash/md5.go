package hash

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/bcdevtools/devd/v3/cmd/utils"
	"github.com/spf13/cobra"
)

func GetMd5Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "md5 [input]",
		Short: "md5 hashing input.",
		Long: `md5 hashing input.
Support pipe.`,
		Run: func(cmd *cobra.Command, args []string) {
			args, err := utils.ProvidedArgsOrFromPipe(args)
			utils.ExitOnErr(err, "failed to get args from pipe")
			utils.RequireArgs(args, cmd)

			input := strings.Join(args, " ")
			hash := md5.Sum([]byte(input))
			fmt.Println(hex.EncodeToString(hash[:]))
		},
	}

	return cmd
}
