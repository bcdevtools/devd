package security_check

//goland:noinspection SpellCheckingInspection
import (
	"fmt"
	libutils "github.com/EscanBE/go-lib/utils"
	sectypes "github.com/bcdevtools/devd/cmd/security_check/types"
	"github.com/bcdevtools/devd/cmd/utils"
	"github.com/bcdevtools/devd/constants"
	"github.com/pkg/errors"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func checkSshDir(homeDir string) {
	const module = secureCheckModuleDotSsh

	username, isRoot, isSuperUser, err := tryExtractUserInfoFromHomeDir(homeDir)
	var reportNonSuperUserNotConfigureDotSshAndAuthorizedKeys bool
	if err != nil {
		libutils.PrintfStdErr("WARN: failed to extract user info from home dir %s with error: %v\n", homeDir, err)
		reportNonSuperUserNotConfigureDotSshAndAuthorizedKeys = true
	}

	fmt.Println("Checking .ssh directory in", homeDir)

	sshDir := path.Join(homeDir, ".ssh")
	exists, perm, err := utils.IsDirAndExists(sshDir)
	if err != nil {
		libutils.PrintlnStdErr("ERR: failed to check if .ssh directory exists:", err)
		securityReport.Add(
			sectypes.NewSecurityRecord(module, true, fmt.Sprintf("Failed to check permission of %s", sshDir)),
		)
		return
	}

	if !exists {
		if isRoot {
			securityReport.Add(
				sectypes.NewSecurityRecord(module, true, fmt.Sprintf("Root's '.ssh' directory %s does not exist", sshDir)),
			)
			return
		}

		if isSuperUser {
			securityReport.Add(
				sectypes.NewSecurityRecord(module, true, fmt.Sprintf("Super user %s's '.ssh' directory %s does not exist", username, sshDir)),
			)
			return
		}

		if reportNonSuperUserNotConfigureDotSshAndAuthorizedKeys {
			securityReport.Add(
				sectypes.NewSecurityRecord(module, false, fmt.Sprintf("User %s's '.ssh' directory %s does not exist", username, sshDir)),
			)
		}
		return
	}

	if perm != constants.REQUIRE_SSH_DIR_PERMISSION {
		securityReport.Add(
			sectypes.NewSecurityRecord(
				module, true, fmt.Sprintf(
					"Directory %s has BAD permission %o, must change to %o",
					sshDir,
					perm,
					constants.REQUIRE_SSH_DIR_PERMISSION,
				),
			).WithGuide(fmt.Sprintf("sudo chmod %o %s", constants.REQUIRE_SSH_DIR_PERMISSION, sshDir)),
		)
	}

	// authorized_keys
	authorizedKeysFile := path.Join(sshDir, "authorized_keys")
	exists, perm, err = utils.IsFileAndExists(authorizedKeysFile)
	if err != nil {
		libutils.PrintlnStdErr("ERR: failed to check if authorized_keys file exists:", err)
		securityReport.Add(
			sectypes.NewSecurityRecord(module, true, fmt.Sprintf("Failed to check existence of %s", authorizedKeysFile)),
		)
	} else if exists {
		if perm != constants.REQUIRE_SSH_AUTHORIZED_KEYS_FILE_PERMISSION {
			securityReport.Add(
				sectypes.NewSecurityRecord(
					module, true, fmt.Sprintf(
						"File %s has BAD permission %o, must change to %o",
						authorizedKeysFile,
						perm,
						constants.REQUIRE_SSH_AUTHORIZED_KEYS_FILE_PERMISSION,
					),
				).WithGuide(fmt.Sprintf("sudo chmod %o %s", constants.REQUIRE_SSH_AUTHORIZED_KEYS_FILE_PERMISSION, authorizedKeysFile)),
			)
		}
	} else if !exists {
		if isRoot {
			securityReport.Add(
				sectypes.NewSecurityRecord(module, true, fmt.Sprintf("Root's SSH authorized keys file %s does not exist", authorizedKeysFile)),
			)
		} else if isSuperUser {
			securityReport.Add(
				sectypes.NewSecurityRecord(module, true, fmt.Sprintf("Super user %s's SSH authorized keys file %s does not exist", username, authorizedKeysFile)),
			)
		} else if reportNonSuperUserNotConfigureDotSshAndAuthorizedKeys {
			securityReport.Add(
				sectypes.NewSecurityRecord(module, false, fmt.Sprintf("User %s's SSH authorized keys file %s does not exist", username, authorizedKeysFile)),
			)
		}
	}

	// other files
	_ = filepath.Walk(sshDir, func(path string, _ os.FileInfo, _ error) error {
		fi, err := os.Stat(path)
		if err != nil {
			libutils.PrintlnStdErr("ERR: failed to stat file", path, "with error:", err)
			sectypes.NewSecurityRecord(
				module, false, fmt.Sprintf("Failed to stat %s", path),
			)
			return err
		}

		if fi.IsDir() {
			return nil
		}

		if fi.Name() == "authorized_keys" { // checked above
			return nil
		}

		perm = fi.Mode().Perm()

		isPubKey := strings.HasSuffix(fi.Name(), ".pub")
		if isPubKey {
			if perm != constants.RECOMMEND_PUBLIC_KEY_FILE_PERMISSION {
				sectypes.NewSecurityRecord(
					module, false, fmt.Sprintf(
						"Public key file %s has non-generic permission %o, suggest change to %o",
						path,
						perm,
						constants.RECOMMEND_PUBLIC_KEY_FILE_PERMISSION,
					),
				).WithGuide(fmt.Sprintf("sudo chmod %o %s", constants.RECOMMEND_PUBLIC_KEY_FILE_PERMISSION, path))
			}
			return nil
		}

		isKnownHosts := fi.Name() == "known_hosts"
		if isKnownHosts {
			if perm != constants.RECOMMEND_KNOWN_HOSTS_FILE_PERMISSION {
				sectypes.NewSecurityRecord(
					module, false, fmt.Sprintf(
						"known_hosts file at %s has non-generic permission %o, suggest change to %o",
						path,
						perm,
						constants.RECOMMEND_KNOWN_HOSTS_FILE_PERMISSION,
					),
				).WithGuide(fmt.Sprintf("sudo chmod %o %s", constants.RECOMMEND_KNOWN_HOSTS_FILE_PERMISSION, path))
			}
			return nil
		}

		isConfigFile := fi.Name() == "config" || strings.HasSuffix(fi.Name(), "_config")
		if isConfigFile {
			if perm != constants.RECOMMEND_SSH_CONFIG_FILE_PERMISSION {
				sectypes.NewSecurityRecord(
					module, false, fmt.Sprintf(
						"SSH config file at %s has non-generic permission %o, suggest change to %o",
						path,
						perm,
						constants.RECOMMEND_SSH_CONFIG_FILE_PERMISSION,
					),
				).WithGuide(fmt.Sprintf("sudo chmod %o %s", constants.RECOMMEND_SSH_CONFIG_FILE_PERMISSION, path))
			}
			if fi.Name() == "config" { // general file name, no need to check content
				return nil
			}
		}

		bz, err := os.ReadFile(path)
		if err != nil {
			libutils.PrintlnStdErr("ERR: failed to read file", path, "with error:", err)
			sectypes.NewSecurityRecord(
				module, true, fmt.Sprintf("Failed to read %s (checking if private key)", path),
			)
			return err
		}

		isPrivateKey := strings.Contains(string(bz), "PRIVATE KEY")
		if isPrivateKey {
			if perm == constants.REQUIRE_SSH_PRIVATE_KEY_FILE_PERMISSION {
				// OK
			} else if perm == constants.ACCEPTABLE_SSH_PRIVATE_KEY_FILE_PERMISSION {
				sectypes.NewSecurityRecord(
					module, false, fmt.Sprintf(
						"Private key file %s has permission %o which is not recommended, suggest change to %o",
						path,
						perm,
						constants.REQUIRE_SSH_PRIVATE_KEY_FILE_PERMISSION,
					),
				).WithGuide(fmt.Sprintf("sudo chmod %o %s", constants.REQUIRE_SSH_PRIVATE_KEY_FILE_PERMISSION, path))
			} else {
				sectypes.NewSecurityRecord(
					module, true, fmt.Sprintf(
						"Private key file %s has BAD permission %o, must change to %o",
						path,
						perm,
						constants.REQUIRE_SSH_PRIVATE_KEY_FILE_PERMISSION,
					),
				).WithGuide(fmt.Sprintf("sudo chmod %o %s", constants.REQUIRE_SSH_PRIVATE_KEY_FILE_PERMISSION, path))
			}
			return nil
		}

		sectypes.NewSecurityRecord(
			module, false, fmt.Sprintf("Un-expected existence of file %s, unable to detect type", path),
		)

		return nil
	})
}

func tryExtractUserInfoFromHomeDir(homeDir string) (username string, isRoot, isSuperUser bool, err error) {
	var success bool
	username, success = utils.TryExtractUserNameFromHomeDir(homeDir)

	if success {
		if username == "root" {
			isRoot = true
			isSuperUser = true
			err = nil
		} else {
			isSuperUser, err = utils.IsSuperUser(username)
			if err == nil {
				isRoot = false
				// isSuperUser = ...
				// err = nil
			} else {
				isRoot = false
				isSuperUser = false
				err = errors.Wrap(err, fmt.Sprintf("failed to check if user %s is super user", username))
			}
		}
	} else {
		isRoot = false
		isSuperUser = false
		err = fmt.Errorf("failed to extract username from home directory %s", homeDir)
	}

	return
}
