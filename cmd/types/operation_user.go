package types

import (
	"github.com/bcdevtools/devd/cmd/utils"
	"os"
)

type OperationUserInfo struct {
	EffectiveUserInfo *UserInfo
	RealUserInfo      *UserInfo
	// IsSameUser indicate that operation user is the same as real user
	IsSameUser bool
	// IsSuperUser indicate that operation user is "root" or a user that belong to a super-user group
	IsSuperUser bool
	// OperatingAsSuperUser indicates that operation user is "root" or a user runs command with "sudo"
	OperatingAsSuperUser bool
}

// GetDefaultWorkingUser returns:
//
// - Effective user if is not super-user and not operating as super-user
//
// - Real user otherwise
func (o *OperationUserInfo) GetDefaultWorkingUser() *UserInfo {
	if o.IsSameUser {
		return o.RealUserInfo
	} else if !o.IsSuperUser && !o.OperatingAsSuperUser {
		return o.EffectiveUserInfo
	} else {
		return o.RealUserInfo
	}
}

func (o *OperationUserInfo) RequireSuperUser() {
	if !o.IsSuperUser {
		utils.PrintlnStdErr("ERR: operation user must be a super user")
		os.Exit(1)
	}
}

func (o *OperationUserInfo) RequireOperatingAsSuperUser() {
	if !o.OperatingAsSuperUser {
		utils.PrintlnStdErr("ERR: must be ran as super user")
		utils.PrintlnStdErr("*** Hint: Try again with 'sudo' ***")
		os.Exit(1)
	}
}

func (o *OperationUserInfo) RequireOperatingAsNonSuperUser() {
	if o.OperatingAsSuperUser {
		utils.PrintlnStdErr("ERR: can not be run as super user or with 'sudo'")
		os.Exit(1)
	}
}
