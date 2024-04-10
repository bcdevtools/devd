package hash

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
	"strings"
)

func GetKeccak512Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "keccak512 [input]",
		Short: "keccak512 hashing input",
		Args:  cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			input := strings.Join(args, " ")
			hash := crypto.Keccak512([]byte(input))
			fmt.Println(hex.EncodeToString(hash))
		},
	}

	return cmd
}
