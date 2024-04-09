package security_check

//goland:noinspection SpellCheckingInspection
import (
	"fmt"
	"github.com/EscanBE/go-ienumerable/goe"
	"github.com/EscanBE/go-ienumerable/goe_helper"
	libutils "github.com/EscanBE/go-lib/utils"
	sectypes "github.com/bcdevtools/devd/cmd/security_check/types"
	"os/exec"
	"strings"
)

//goland:noinspection GoSnakeCaseUsage
func checkLinuxUsersPasswdDisabled_LinuxOnly() {
	const module = secureCheckModuleUserPassword

	bz, err := exec.Command("passwd", "-a", "-S").Output()
	if err != nil {
		libutils.PrintlnStdErr("ERR: failed to check user password status:", err)
		securityReport.Add(
			sectypes.NewSecurityRecord(module, true, "Failed to user password status"),
		)
		return
	}

	//goland:noinspection GoRedundantConversion
	usersInfo := [][]string(goe_helper.Select(goe.NewIEnumerable(strings.Split(string(bz), "\n")...).Where(func(s string) bool {
		return !libutils.IsBlank(s)
	}), func(v string) []string {
		return strings.Split(strings.TrimSpace(v), " ")
	}).ToArray())

	for _, userInfo := range usersInfo {
		if len(userInfo) != 7 {
			fmt.Println("WARN: abnormal user info:", userInfo)
		}

		if len(userInfo) < 2 {
			securityReport.Add(
				sectypes.NewSecurityRecord(module, true, "Abnormal user info: "+strings.Join(userInfo, " ")),
			)
			continue
		}

		userName := userInfo[0]
		status := userInfo[1]
		switch status {
		case "L":
			if userName == "root" {
				securityReport.Add(
					sectypes.NewSecurityRecord(module, false, "'root' user should have password enabled to prevent login incident"),
				)
			}
			break
		case "P":
			if userName != "root" {
				securityReport.Add(
					sectypes.
						NewSecurityRecord(module, true, fmt.Sprintf("User '%s' has password enabled, should only allow connect via SSH", userName)).
						WithGuide(fmt.Sprintf("sudo usermod -p '*' %s", userName)),
				)
			}
			break
		case "NP":
			securityReport.Add(
				sectypes.
					NewSecurityRecord(module, true, fmt.Sprintf("User '%s' has no password", userName)).
					WithGuide(fmt.Sprintf("sudo usermod -p '*' %s", userName)),
			)
			break
		default:
			securityReport.Add(
				sectypes.
					NewSecurityRecord(module, true, fmt.Sprintf("Un-expected value '%s' of password status of user '%s'", status, userName)),
			)
			break
		}
	}
}
