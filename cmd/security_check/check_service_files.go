package security_check

//goland:noinspection SpellCheckingInspection
import (
	"fmt"
	"github.com/EscanBE/go-ienumerable/goe"
	"github.com/EscanBE/go-ienumerable/goe_helper"
	libutils "github.com/EscanBE/go-lib/utils"
	sectypes "github.com/bcdevtools/devd/cmd/security_check/types"
	"github.com/bcdevtools/devd/cmd/utils"
	"io/fs"
	"os"
	"path"
	"strings"
	"syscall"
)

//goland:noinspection GoSnakeCaseUsage
func checkServiceFiles_LinuxOnly() {
	const module = secureCheckModuleSystemdServiceFile

	systemdDir := "/etc/systemd/system"

	dirEntries, err := os.ReadDir(systemdDir)
	if err != nil {
		libutils.PrintlnStdErr("ERR: failed to listing entries in systemd directory:", err)
		securityReport.Add(
			sectypes.NewSecurityRecord(module, true, "Failed to listing entries in systemd directory"),
		)
		return
	}

	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			continue
		}

		fileName := dirEntry.Name()
		if !strings.EqualFold(".service", path.Ext(fileName)) {
			continue
		}

		fi, err := dirEntry.Info()
		if err != nil {
			libutils.PrintfStdErr("ERR: failed to get info of %s with error %v\n", fileName, err)
			securityReport.Add(
				sectypes.NewSecurityRecord(module, false, fmt.Sprintf("Failed to get info of service file %s", fileName)),
			)
			continue
		}

		// outer check

		if stat, ok := fi.Sys().(*syscall.Stat_t); ok {
			if stat.Uid == 0 && stat.Gid == 0 {
				// OK
			} else {
				sectypes.NewSecurityRecord(
					module,
					false,
					fmt.Sprintf("Current owner of service %s is %d:%d. Expected 0:0 (means root:root)", fileName, stat.Uid, stat.Gid),
				)
			}
		}

		perm := fi.Mode().Perm()

		owner, group, other := utils.ExtractPermissionParts(perm)
		if fi.Mode()&os.ModeSymlink == 0 { // not a symlink
			if owner < 4 {
				securityReport.Add(
					sectypes.NewSecurityRecord(
						module,
						false,
						fmt.Sprintf("Current permission of service file %s is %o, which is missing read permission for owner", fileName, perm),
					),
				)
			}
			if group > 5 {
				securityReport.Add(
					sectypes.NewSecurityRecord(
						module,
						false,
						fmt.Sprintf("Current permission of service file %s is %o, which grants write permission for group. The maximum recommended permission is Read", fileName, perm),
					),
				)
			}
			if other > 5 {
				securityReport.Add(
					sectypes.NewSecurityRecord(
						module,
						false,
						fmt.Sprintf("Current permission of service file %s is %o, which grants write permission for others. The maximum recommended permission is Read", fileName, perm),
					),
				)
			}
		}

		// content check

		filePath := path.Join(systemdDir, fileName)
		bz, err := os.ReadFile(filePath)
		if err != nil {
			libutils.PrintfStdErr("ERR: failed to read content of %s with error %v\n", fileName, err)
			securityReport.Add(
				sectypes.NewSecurityRecord(module, false, fmt.Sprintf("Failed to read %s for content detection", fileName)),
			)
		} else {
			content := string(bz)
			if strings.Contains(content, "Group=root") {
				securityReport.Add(
					sectypes.
						NewSecurityRecord(module, false, fmt.Sprintf("Service file %s configured Group=root, that seems danger", fileName)),
				)
			}

			ieLines := goe_helper.Select(goe.NewIEnumerable(strings.Split(content, "\n")...).Where(func(line string) bool {
				return !libutils.IsBlank(line)
			}), func(line string) string {
				return strings.TrimSpace(line)
			})

			ieUserLines := ieLines.Where(func(line string) bool {
				return strings.HasPrefix(line, "User=")
			})
			if ieUserLines.Any() {
				users := goe_helper.Select(ieUserLines, func(v string) string {
					spl := strings.SplitN(v, "=", 2)
					return strings.TrimSpace(spl[1])
				}).ToArray()
				for _, user := range users {
					if strings.EqualFold(user, "root") {
						securityReport.Add(
							sectypes.NewSecurityRecord(
								module,
								false,
								fmt.Sprintf("Service file %s configured User=root, that seems danger", fileName),
							),
						)
						continue
					}

					isSuper, err := utils.IsSuperUser(user)
					if err != nil {
						libutils.PrintfStdErr("ERR: failed to check if user %s is super user with error %v\n", user, err)
						continue
					}

					if isSuper {
						securityReport.Add(
							sectypes.NewSecurityRecord(
								module,
								false,
								fmt.Sprintf("Service file %s configured User=%s, which is a super user and that seems danger", fileName, user),
							),
						)
					}
				}
			}

			ieEnvironmentLines := goe_helper.Select(ieLines, func(line string) string {
				if //goland:noinspection SpellCheckingInspection
				!strings.Contains(line, "nvironment") {
					return ""
				}
				if !strings.Contains(line, "=") {
					return ""
				}
				spl := strings.SplitN(line, "=", 2)
				if !strings.EqualFold("Environment", spl[0]) {
					return ""
				}
				return strings.TrimSpace(spl[1])
			}).Where(func(v string) bool {
				return !libutils.IsBlank(v)
			})

			if ieEnvironmentLines.Any() {
				if ieEnvironmentLines.AnyBy(func(environment string) bool {
					return isDangerServiceFileContainsPassword(environment, perm)
				}) {
					securityReport.Add(
						sectypes.NewSecurityRecord(
							module,
							true,
							fmt.Sprintf("Service file %s include password/connection configuration, that seems incorrect", fileName),
						).WithGuide("sudo chmod 600 " + filePath),
					)
				}
			}
		}
	}
}

func isDangerServiceFileContainsPassword(content string, perm fs.FileMode) bool {
	if perm == 0o600 {
		return false
	}

	//goland:noinspection SpellCheckingInspection
	return strings.Contains(strings.ToLower(content), "password=")
}
