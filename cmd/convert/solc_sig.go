package convert

import (
	"fmt"
	"github.com/bcdevtools/devd/v2/constants"
	"regexp"
	"strings"

	"github.com/bcdevtools/devd/v2/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
)

// GetConvertSolcSignatureCmd creates a helper command that convert EVM method/event into keccak256 hash and 4 bytes signature
func GetConvertSolcSignatureCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "solc_sig [method or event]",
		Aliases: []string{"solc-sig"},
		Short:   "Convert Solidity method/event signature into hashed signature and 4 bytes signature.",
		Long: `Convert Solidity method/event signature into hashed signature and 4 bytes signature.
Output will be 4 lines:
1. Type: Method/Event
2. Method/Event interface used for apply keccak256 on
3. Hashed signature
4. Type based: Event will show hashed signature, Method will show 4 bytes signature
`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			utils.PrintfStdErr("WARN: from v3, this command will be renamed to `%s convert solc-sig` (`-` instead of '_')\n", constants.BINARY_NAME)

			_interface := strings.Join(args, " ")

			_4BytesSig, hash, finalInterface, err := getSignatureFromInterface(_interface)
			utils.ExitOnErr(err, "failed to get signature from interface")

			var _type, solcSig string
			if strings.ToUpper(string(finalInterface[0])) == string(finalInterface[0]) { // event
				_type = "Event"
				solcSig = hash.Hex()
			} else {
				_type = "Method"
				solcSig = _4BytesSig
			}

			fmt.Println(_type)
			fmt.Println(finalInterface)
			fmt.Println(hash.Hex())
			fmt.Println(solcSig)
		},
	}

	return cmd
}

func getSignatureFromInterface(_interface string) (_4BytesSig string, hash common.Hash, finalInterface string, err error) {
	_interface = normalizeEvmEventOrMethodInterface(_interface)

	fmt.Println(_interface)

	if !regexp.MustCompile(`^\w+\s*\(.*\)$`).MatchString(_interface) {
		err = fmt.Errorf("invalid EVM method/event interface, require format: `methodName(...)`")
		return
	}

	finalInterface, err = prepareInterfaceToHash(_interface)
	if err != nil {
		return
	}

	hash = common.BytesToHash(crypto.Keccak256([]byte(finalInterface)))

	_4BytesSig = fmt.Sprintf("0x%x", hash[:4])
	return
}

func normalizeEvmEventOrMethodInterface(_interface string) string {
	// drop part after ')' if any
	spl := strings.Split(_interface, ")")
	if len(spl) > 0 && spl[len(spl)-1] != "" {
		// drop the last part
		_interface = strings.Join(spl[:len(spl)-1], ")") + ")"
	}

	_interface = strings.TrimSpace(_interface)
	_interface = strings.TrimSuffix(_interface, ";")
	_interface = strings.TrimSpace(_interface)
	_interface = strings.TrimSuffix(_interface, "{")
	_interface = strings.TrimSpace(_interface)

	// remove indexed keyword
	//_interface = regexp.MustCompile(`[\s\t\n]+indexed[\s\t\n]+`).ReplaceAllString(_interface, " ")

	// remove extra spaces

	_interface = removeExtraSpaces(_interface)

	// remove space surrounding '(' & ')' & ','

	_interface = strings.ReplaceAll(_interface, " (", "(")
	_interface = strings.ReplaceAll(_interface, "( ", "(")
	_interface = strings.ReplaceAll(_interface, " )", ")")
	_interface = strings.ReplaceAll(_interface, ") ,", "),")
	_interface = strings.ReplaceAll(_interface, ") )", "))")
	_interface = strings.ReplaceAll(_interface, " ,", ",")
	_interface = strings.ReplaceAll(_interface, ", ", ",")

	// trim either prefix 'function ', 'event '
	if strings.HasPrefix(_interface, "function ") {
		_interface = strings.TrimPrefix(_interface, "function ")
	} else {
		_interface = strings.TrimPrefix(_interface, "event ")
	}
	_interface = strings.TrimSpace(_interface)

	// ...

	return strings.TrimSpace(_interface) // finalize
}

//func validateEvmEventOrMethodInterface(_interface string) (ok bool, desc string) {
//	// validate event/method interface
//}

func prepareInterfaceToHash(_interface string) (res string, err error) {
	defer func() {
		// remove all remaining spaces
		res = strings.ReplaceAll(res, " ", "")
	}()

	res = _interface

	// remove indexed keyword
	res = regexp.MustCompile(`[\s\t\n]+indexed[\s\t\n]+`).ReplaceAllString(res, " ")

	// remove any variable name
	if !strings.HasSuffix(res, ")") {
		err = fmt.Errorf("interface must ends with ')': %s", res)
		return
	}
	spl1 := strings.SplitN(res[:len(res)-1] /*remove suffix ')'*/, "(", 2)
	functionName := strings.TrimSpace(spl1[0])
	argsPart := strings.TrimSpace(spl1[1])

	var argsPartWithoutVariableName []string
	if len(argsPart) > 0 {
		var trimmedFragments []string
		for _, fragment := range strings.Split(argsPart, ",") {
			trimmedFragments = append(trimmedFragments, strings.TrimSpace(fragment))
		}

		argsPart = strings.Join(trimmedFragments, ",")
		argsPart += "," // add suffix ',' to simplify the logic
		var parenthesisLevel int
		var squareBracketLevel int
		var argName string
		var meetSpace bool
		for _, c := range argsPart {
			if c == ',' {
				if parenthesisLevel == 0 && squareBracketLevel == 0 {
					argsPartWithoutVariableName = append(argsPartWithoutVariableName, argName)
					argName = ""
					meetSpace = false
					continue
				}
			} else if c == '(' {
				parenthesisLevel++
			} else if c == ')' {
				parenthesisLevel--
			} else if c == '[' {
				squareBracketLevel++
			} else if c == ']' {
				squareBracketLevel--
			} else if c == ' ' {
				if parenthesisLevel == 0 && squareBracketLevel == 0 {
					meetSpace = true
				}
			}

			if meetSpace {
				continue
			}
			argName += string(c)
		}
	} else {
		argsPartWithoutVariableName = []string{}
	}

	res = fmt.Sprintf("%s(%s)", functionName, strings.Join(argsPartWithoutVariableName, ","))
	return
}

func removeExtraSpaces(str string) string {
	var passOne bool
	for {
		var replacedAny bool
		if strings.Contains(str, "  ") {
			passOne = false
			str = strings.ReplaceAll(str, "  ", " ")
			replacedAny = true
		}
		if strings.Contains(str, "\n") {
			passOne = false
			str = strings.ReplaceAll(str, "\n", " ")
			replacedAny = true
		}
		if strings.Contains(str, "\t") {
			passOne = false
			str = strings.ReplaceAll(str, "\t", " ")
			replacedAny = true
		}
		if replacedAny {
			continue
		}

		if !passOne {
			passOne = true
			continue // retry one more time
		}

		break
	}

	return str
}
