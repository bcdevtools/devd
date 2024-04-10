package hash

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
	"strings"
)

func GetKeccak256Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "keccak256 [input]",
		Short: "keccak256 hashing input",
		Args:  cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			input := strings.Join(args, " ")
			hash := crypto.Keccak256Hash([]byte(input))
			fmt.Println(hex.EncodeToString(hash.Bytes()))
		},
	}

	return cmd
}
