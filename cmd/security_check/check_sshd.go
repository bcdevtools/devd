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
func checkSshdConfig_LinuxOnly() {
	const module = secureCheckModuleSshd
	const sshdConfigFile = "/etc/ssh/sshd_config"
	fmt.Println("Checking from SSHD config")
	toCheckOptions := []string{
		"PasswordAuthentication",
		"KbdInteractiveAuthentication",
		"ChallengeResponseAuthentication",
		"UsePAM",
		"PermitRootLogin",
		"PermitEmptyPasswords",
		"AllowUsers",
	}

	bz, err := exec.Command(
		"bash",
		"-c",
		fmt.Sprintf("sudo sshd -T | egrep -i '%s'", strings.Join(toCheckOptions, "|")),
	).CombinedOutput()

	if err != nil {
		libutils.PrintlnStdErr("ERR: failed to check SSHD config with error:", err)
		securityReport.Add(sectypes.NewSecurityRecord(module, true, "Failed to check SSHD config"))
		return
	}

	output := string(bz)
	spl := strings.Split(output, "\n")
	ieConfigs := goe_helper.Select(goe.NewIEnumerable(spl...), func(v string) string {
		return strings.ToLower(strings.TrimSpace(v))
	}).Where(func(v string) bool {
		return len(v) > 0
	})
	type config struct {
		name  string
		value string
	}
	//goland:noinspection GoRedundantConversion
	configs := map[string]string(goe_helper.ToDictionary(goe_helper.Select(ieConfigs, func(v string) config {
		spl2 := strings.SplitN(v, " ", 2)
		cfg := config{
			name: spl2[0],
		}
		if len(spl2) > 1 {
			cfg.value = spl2[1]
		}
		return cfg
	}), func(source config) string {
		return source.name
	}, func(source config) string {
		return source.value
	}))

	if //goland:noinspection SpellCheckingInspection
	passwordAuthentication, found := configs["passwordauthentication"]; found {
		if passwordAuthentication != "no" {
			securityReport.Add(
				sectypes.
					NewSecurityRecord(module, true, "PasswordAuthentication is enabled, must disable it").
					WithGuide("Update in " + sshdConfigFile + ", make sure 'PasswordAuthentication no', without comment"),
			)
		}
	} else {
		securityReport.Add(
			sectypes.
				NewSecurityRecord(module, true, "Cannot lookup PasswordAuthentication in SSHD config").
				WithGuide("Check in " + sshdConfigFile + ", make sure 'PasswordAuthentication no', without comment"),
		)
	}

	//goland:noinspection SpellCheckingInspection
	kbdInteractiveAuthentication, foundKbdInteractiveAuthentication := configs["kbdinteractiveauthentication"]
	if foundKbdInteractiveAuthentication {
		if kbdInteractiveAuthentication != "no" {
			securityReport.Add(
				sectypes.
					NewSecurityRecord(module, true, "KbdInteractiveAuthentication is enabled, must disable it").
					WithGuide("Update in " + sshdConfigFile + ", make sure 'KbdInteractiveAuthentication no', without comment"),
			)
		}
	} else {
		securityReport.Add(
			sectypes.
				NewSecurityRecord(module, true, "Cannot lookup KbdInteractiveAuthentication in SSHD config").
				WithGuide("Check in " + sshdConfigFile + ", make sure 'KbdInteractiveAuthentication no', without comment"),
		)
	}

	if //goland:noinspection SpellCheckingInspection
	challengeResponseAuthentication, found := configs["challengeresponseauthentication"]; found {
		if challengeResponseAuthentication != "no" {
			securityReport.Add(
				sectypes.
					NewSecurityRecord(module, true, "ChallengeResponseAuthentication is enabled, must disable it").
					WithGuide("Update in " + sshdConfigFile + ", make sure 'ChallengeResponseAuthentication no', without comment"),
			)
		}
	} else if !foundKbdInteractiveAuthentication {
		securityReport.Add(
			sectypes.
				NewSecurityRecord(module, true, "Cannot lookup ChallengeResponseAuthentication in SSHD config").
				WithGuide("Check in " + sshdConfigFile + ", make sure 'ChallengeResponseAuthentication no', without comment"),
		)
	}

	if //goland:noinspection SpellCheckingInspection
	usePAM, found := configs["usepam"]; found {
		if usePAM != "no" {
			securityReport.Add(
				sectypes.
					NewSecurityRecord(module, true, "UsePAM is enabled, must disable it").
					WithGuide("Update in " + sshdConfigFile + ", make sure 'UsePAM no', without comment"),
			)
		}
	} else {
		securityReport.Add(
			sectypes.
				NewSecurityRecord(module, true, "Cannot lookup UsePAM in SSHD config").
				WithGuide("Check in " + sshdConfigFile + ", make sure 'UsePAM no', without comment"),
		)
	}

	if //goland:noinspection SpellCheckingInspection
	permitRootLogin, found := configs["permitrootlogin"]; found {
		if permitRootLogin != "without-password" {
			securityReport.Add(
				sectypes.
					NewSecurityRecord(module, true, "PermitRootLogin is set to "+permitRootLogin+", must prohibit password").
					WithGuide("Update in " + sshdConfigFile + ", make sure 'PermitRootLogin prohibit-password', without comment"),
			)
		}
	} else {
		securityReport.Add(
			sectypes.
				NewSecurityRecord(module, true, "Cannot lookup PermitRootLogin in SSHD config").
				WithGuide("Check in " + sshdConfigFile + ", make sure 'PermitRootLogin prohibit-password', without comment"),
		)
	}

	if //goland:noinspection SpellCheckingInspection
	permitEmptyPasswords, found := configs["permitemptypasswords"]; found {
		if permitEmptyPasswords != "no" {
			securityReport.Add(
				sectypes.
					NewSecurityRecord(module, true, "PermitEmptyPasswords is enabled, must disable it").
					WithGuide("Update in " + sshdConfigFile + ", make sure 'PermitEmptyPasswords no', without comment"),
			)
		}
	} else {
		securityReport.Add(
			sectypes.
				NewSecurityRecord(module, true, "Cannot lookup PermitEmptyPasswords in SSHD config").
				WithGuide("Check in " + sshdConfigFile + ", make sure 'PermitEmptyPasswords no', without comment"),
		)
	}

	if //goland:noinspection SpellCheckingInspection
	allowUsers, found := configs["allowusers"]; found {
		if allowUsers != "" {
			securityReport.Add(
				sectypes.
					NewSecurityRecord(module, false, "AllowUsers is set to "+allowUsers+", should disable it to allow any user login").
					WithGuide("Update in " + sshdConfigFile + ", make sure 'AllowUsers' is commented to disable it"),
			)
		}
	}
}
