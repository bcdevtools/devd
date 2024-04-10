package hash

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

func GetMd5Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "md5 [input]",
		Short: "MD5 hashing input",
		Args:  cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			input := strings.Join(args, " ")
			hash := md5.Sum([]byte(input))
			fmt.Println(hex.EncodeToString(hash[:]))
		},
	}

	return cmd
}
