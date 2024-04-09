package gen

import (
	"bufio"
	"fmt"
	libutils "github.com/EscanBE/go-lib/utils"
	"github.com/bcdevtools/devd/constants"
	"github.com/spf13/cobra"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	flagComment      = "comment"
	flagTemp         = "temp"
	temporaryComment = "Temp allow from laptop"
)

// GenerateUfwAllowCommand registers a sub-tree of commands
func GenerateUfwAllowCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ufw-allow [?IP] [?port]",
		Aliases: []string{"ufw"},
		Long: fmt.Sprintf(`Generate UFW allow command in 3 cases:
Case 1: No arguments: basic SSH setup commands (allow :22, deny incoming, allow outgoing, enable, reload, status)
Case 2. Specified a port: allow specified port from anywhere
> %s gen ufw-allow 22 # allow SSH from anywhere
Case 3. Specified an IP: allow connection to any port from a specified IP
> %s gen ufw-allow 8.8.8.8 # Allow all connections from 8.8.8.8 to any port
Case 4. Specified an IP, then a port: allow connects to a specified port from a specified IP
> %s gen ufw-allow 8.8.8.8 5432 # Allow connection to PostgreSQL at :5432 from 8.8.8.8

Pre-defined ports: http=80, https=443, ssh=22, db=5432, grpc=9090, rpc=26657, evm=8545, p2p=26656, rest=1317
`, constants.BINARY_NAME, constants.BINARY_NAME, constants.BINARY_NAME),
		Args: cobra.RangeArgs(0, 2),
		Run:  generateUfwAllow,
	}

	cmd.Flags().StringP(flagComment, "c", "", "comment for the rule")
	cmd.Flags().Bool(flagTemp, false, fmt.Sprintf("mark temporary rule in comment. If comment is absent, comment will be '%s'", temporaryComment))

	return cmd
}

func generateUfwAllow(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("sudo ufw allow ssh comment 'Allow SSH connections'")
		fmt.Println("sudo ufw default deny incoming")
		fmt.Println("sudo ufw default allow outgoing")
		fmt.Println("sudo ufw enable")
		fmt.Println("sudo ufw status numbered")
		return
	}

	userComment := cmd.Flag(flagComment).Value.String()
	var isTemp bool
	if cmd.Flag(flagTemp).Changed {
		if len(userComment) < 1 {
			userComment = temporaryComment
		}
		isTemp = true
	}

	var port uint16
	var ip string

	if len(args) == 1 {
		if isValidFirewallIpAddress(args[0]) {
			ip = args[0]
		} else if port16, ok := extractPort(args[0]); ok {
			port = port16
		} else {
			libutils.PrintlnStdErr("ERR: specified argument is neither a valid IP address nor port", args[0])
			os.Exit(1)
		}
	} else {
		if isValidFirewallIpAddress(args[0]) {
			ip = args[0]
		} else {
			libutils.PrintlnStdErr("ERR: in-case multiple arguments supplied, first argument must be a valid IP address", args[0])
			os.Exit(1)
		}
		if port16, ok := extractPort(args[1]); ok {
			port = port16
		} else {
			libutils.PrintlnStdErr("ERR: in-case multiple arguments supplied, second argument must be a valid port", args[1])
			os.Exit(1)
		}
	}

	if len(userComment) < 1 {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Comment:")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		if len(text) < 1 {
			libutils.PrintlnStdErr("ERR: comment is required")
			os.Exit(1)
		}
		userComment = text
	}

	var comment string
	if isTemp {
		comment = "*TEMP* "
	}
	comment += strings.ReplaceAll(userComment, "'", "\"")

	fmt.Println()
	if port > 0 && len(ip) < 1 {
		fmt.Printf("sudo ufw allow proto any to 0.0.0.0/0 port %d comment '%s'", port, comment+fmt.Sprintf(" (allow anywhere to :%d)", port))
	} else if port < 1 && len(ip) > 0 {
		fmt.Printf("sudo ufw allow from %s comment '%s'", ip, comment+fmt.Sprintf(" (allow from %s)", ip))
	} else if port > 0 && len(ip) > 0 {
		fmt.Printf("sudo ufw allow from %s to any port %d comment '%s'", ip, port, comment+fmt.Sprintf(" (allow to :%d from %s)", port, ip))
	} else {
		panic("unreachable")
	}
	fmt.Println(" && sudo ufw status numbered")
}

func isValidFirewallIpAddress(ip string) bool {
	return regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}(/(8|16|24|32))?$`).MatchString(ip)
}

func extractPort(port string) (uint16, bool) {
	switch strings.ToLower(port) {
	case "":
		return 0, false
	case "http":
		return 80, true
	case "https":
		return 443, true
	case "ssh":
		return 22, true
	case "db", "postgres", "psql", "pgsql", "database", "sql", "postgresql", "pg":
		return 5432, true
	case "grpc":
		return 9090, true
	case "rpc":
		return 26657, true
	case "evm", "ethereum", "eth", "json-rpc", "jsonrpc":
		return 8545, true
	case "p2p", "tendermint", "tm", "tendermint-p2p", "tendermint-p2p-rpc":
		return 26656, true
	case "rest", "rest-api", "restapi":
		return 1317, true
	}

	p, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		return 0, false
	}
	if p < 1 || p > 65535 {
		return 0, false
	}
	return uint16(p), true
}
