package security_check

//goland:noinspection SpellCheckingInspection
import (
	libutils "github.com/EscanBE/go-lib/utils"
	sectypes "github.com/bcdevtools/devd/cmd/security_check/types"
	"os/exec"
	"strings"
)

//goland:noinspection GoSnakeCaseUsage
func checkUfwFirewall_LinuxOnly() {
	const module = secureCheckModuleUfw

	bz, err := exec.Command("ufw", "status", "verbose").Output()
	if err != nil {
		libutils.PrintlnStdErr("ERR: failed to check ufw status:", err)
		securityReport.Add(
			sectypes.NewSecurityRecord(module, true, "Failed to check firewall status"),
		)
		return
	}

	content := string(bz)
	if !strings.Contains(content, "Status: active") {
		securityReport.Add(
			sectypes.NewSecurityRecord(module, false, "Firewall is not enabled. Anyway, make sure setup correct rules before enable it, at least allow SSH `sudo ufw allow ssh`"),
		)
		return
	}

	if !strings.Contains(content, "Logging: on") {
		securityReport.Add(
			sectypes.
				NewSecurityRecord(module, false, "Firewall logging is probably not enabled, suggest enable it at low level").
				WithGuide("sudo ufw logging on"),
		)
	}
}
