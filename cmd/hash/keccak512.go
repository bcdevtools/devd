package hash

import (
	"encoding/hex"
	"fmt"
	"github.com/bcdevtools/devd/cmd/utils"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
	"strings"
)

func GetKeccak512Command() *cobra.Command {
	cmd := &cobra.Command{
		Use: "keccak512 [input]",
		Short: `keccak512 hashing input.
Support pipe.`,
		Run: func(cmd *cobra.Command, args []string) {
			args, err := utils.ProvidedArgsOrFromPipe(args)
			utils.ExitOnErr(err, "failed to get args from pipe")
			utils.RequireArgs(args, cmd)

			input := strings.Join(args, " ")
			hash := crypto.Keccak512([]byte(input))
			fmt.Println(hex.EncodeToString(hash))
		},
	}

	return cmd
}
