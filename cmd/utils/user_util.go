package utils

import (
	"fmt"
	"github.com/EscanBE/go-ienumerable/goe"
	libutils "github.com/EscanBE/go-lib/utils"
	"github.com/bcdevtools/devd/cmd/types"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
)

func GetOperationUserInfo() (operationUserInfo *types.OperationUserInfo, err error) {
	defer func() {
		if err != nil {
			operationUserInfo = nil
		}
	}()

	operationUserInfo = &types.OperationUserInfo{}

	operationUserInfo.EffectiveUserInfo, err = getEffectiveUserInfo()
	if err != nil {
		err = errors.Wrap(err, "failed to get effective user info")
		return
	}

	operationUserInfo.RealUserInfo, err = getRealUserInfo()
	if err != nil {
		err = errors.Wrap(err, "failed to get real user info")
		return
	}

	operationUserInfo.IsSameUser = operationUserInfo.EffectiveUserInfo.UserId == operationUserInfo.RealUserInfo.UserId
	operationUserInfo.IsSuperUser = operationUserInfo.EffectiveUserInfo.IsSuperUser
	operationUserInfo.OperatingAsSuperUser = operationUserInfo.EffectiveUserInfo.UserId == 0

	return
}

func IsSuperUser(username string) (isSuperUser bool, err error) {
	if username == "root" {
		isSuperUser = true
	} else {
		var groups []string
		groups, err = GetGroupsOfUser(username)
		if err == nil {
			isSuperUser = goe.NewIEnumerable(groups...).AnyBy(func(s string) bool {
				switch s {
				case "sudo":
					return true
				case "admin":
					return true
				case "google-sudoers":
					return true
				default:
					return false
				}
			})
		}
	}

	return
}

func GetGroupsOfUser(username string) ([]string, error) {
	var bz []byte
	bz, err := exec.Command("groups", username).CombinedOutput()
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to get groups of user %s", username))
	}

	return goe.NewIEnumerable(strings.Split(strings.TrimSpace(string(bz)), " ")...).ToArray(), nil
}

func getEffectiveUserInfo() (*types.UserInfo, error) {
	var err error
	var effectiveUserInfo types.UserInfo
	var usr *user.User

	usr, err = user.Current()
	if err != nil {
		err = errors.Wrap(err, "failed to get current user")
		return nil, err
	}

	effectiveUserInfo.UserId = os.Geteuid()
	effectiveUserInfo.Username = usr.Username
	effectiveUserInfo.HomeDir = usr.HomeDir
	effectiveUserInfo.User = usr
	effectiveUserInfo.IsSuperUser, err = IsSuperUser(usr.Username)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to check if user %s is super user", usr.Username))
	}

	return &effectiveUserInfo, nil
}

func getRealUsername() (realUsername string, err error) {
	sudoUser := os.Getenv("SUDO_USER")
	if !libutils.IsBlank(sudoUser) {
		realUsername = sudoUser
		return
	}

	if IsLinux() && HasBinaryName("who") {
		bz, runErr := exec.Command("who", "am", "i").Output()
		if runErr == nil {
			output := strings.TrimSpace(string(bz))
			if strings.Contains(output, " ") {
				realUsername = strings.Split(output, " ")[0]
				return
			}
		}
	} else if IsDarwin() {
		bz, runErr := exec.Command("ps", "-o", "user=", "-p", strconv.Itoa(os.Getpid())).Output()
		if runErr == nil {
			output := strings.TrimSpace(string(bz))
			if len(output) > 0 {
				realUsername = output
				return
			}
		}
	}

	var currentUser *user.User
	currentUser, err = user.Current()
	if err != nil {
		err = errors.Wrap(err, "failed to get current user")
	} else {
		realUsername = currentUser.Username
	}

	return
}

func getRealUserInfo() (*types.UserInfo, error) {
	var err error
	var realUserInfo types.UserInfo
	var usr *user.User

	realUserInfo.Username, err = getRealUsername()
	if err != nil {
		err = errors.Wrap(err, "failed to get real user name")
		return nil, err
	}

	usr, err = user.Lookup(realUserInfo.Username)
	if err != nil {
		err = errors.Wrap(err, "failed to get real user")
		return nil, err
	}

	var uid int64
	uid, err = strconv.ParseInt(usr.Uid, 10, 64)
	if err != nil {
		err = errors.Wrap(err, "failed to parse real user id "+usr.Uid)
		return nil, err
	}

	realUserInfo.UserId = int(uid)
	realUserInfo.Username = usr.Username
	realUserInfo.HomeDir = usr.HomeDir
	realUserInfo.User = usr
	realUserInfo.IsSuperUser, err = IsSuperUser(usr.Username)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to check if user %s is super user", usr.Username))
	}

	return &realUserInfo, nil
}
