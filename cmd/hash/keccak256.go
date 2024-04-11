package hash

import (
	"encoding/hex"
	"fmt"
	libutils "github.com/EscanBE/go-lib/utils"
	"github.com/bcdevtools/devd/cmd/utils"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
	"strings"
)

func GetKeccak256Command() *cobra.Command {
	cmd := &cobra.Command{
		Use: "keccak256 [input]",
		Short: `keccak256 hashing input.
Support pipe.`,
		Run: func(cmd *cobra.Command, args []string) {
			args, err := utils.ProvidedArgsOrFromPipe(args)
			libutils.ExitIfErr(err, "failed to get args from pipe")
			utils.RequireArgs(args, cmd)

			input := strings.Join(args, " ")
			hash := crypto.Keccak256([]byte(input))
			fmt.Println(hex.EncodeToString(hash))
		},
	}

	return cmd
}
