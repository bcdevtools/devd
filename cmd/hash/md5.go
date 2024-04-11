package hash

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	libutils "github.com/EscanBE/go-lib/utils"
	"github.com/bcdevtools/devd/cmd/utils"
	"github.com/spf13/cobra"
	"strings"
)

func GetMd5Command() *cobra.Command {
	cmd := &cobra.Command{
		Use: "md5 [input]",
		Short: `md5 hashing input.
Support pipe.`,
		Run: func(cmd *cobra.Command, args []string) {
			args, err := utils.ProvidedArgsOrFromPipe(args)
			libutils.ExitIfErr(err, "failed to get args from pipe")
			utils.RequireArgs(args, cmd)

			input := strings.Join(args, " ")
			hash := md5.Sum([]byte(input))
			fmt.Println(hex.EncodeToString(hash[:]))
		},
	}

	return cmd
}
