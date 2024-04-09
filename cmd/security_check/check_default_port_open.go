package security_check

//goland:noinspection SpellCheckingInspection
import (
	"fmt"
	sectypes "github.com/bcdevtools/devd/cmd/security_check/types"
	"net"
	"time"
)

//goland:noinspection SpellCheckingInspection

func checkDefaultPortOpen() {
	const module = secureCheckModuleDefaultPortOpen

	type portCheck struct {
		name string
		port uint16
	}

	portChecks := []portCheck{
		{
			name: "SSH",
			port: 22,
		},
		{
			name: "PostgreSQL",
			port: 5432,
		},
		{
			name: "Redis",
			port: 6379,
		},
	}

	for _, pc := range portChecks {
		isDefaultPortOpen, _ := isPortOpen(pc.port)
		if isDefaultPortOpen {
			securityReport.Add(
				sectypes.NewSecurityRecord(module, false, fmt.Sprintf("Default %s port %d is open, should be changed to another port", pc.name, pc.port)),
			)
		}
	}
}

func isPortOpen(port uint16) (open bool, err error) {
	var conn net.Conn

	conn, err = net.DialTimeout("tcp", net.JoinHostPort("localhost", fmt.Sprintf("%d", port)), 1*time.Second)

	if conn != nil {
		defer func() {
			_ = conn.Close()
		}()
		open = true
	}

	return
}
