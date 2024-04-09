package types

import "os/user"

type UserInfo struct {
	UserId   int
	Username string
	HomeDir  string
	User     *user.User
	// IsSuperUser indicate that the user is "root" or a user that belong to a super-user group
	IsSuperUser bool
}
