package utils

import "runtime"

func IsLinux() bool {
	//goland:noinspection GoBoolExpressions
	return runtime.GOOS == "linux"
}

func IsDarwin() bool {
	//goland:noinspection GoBoolExpressions
	return runtime.GOOS == "darwin"
}

// GetUserGeneralRootDirLocation returns "/Users" on macOS, "/home" on Linux, panic on Windows and others
//
//goland:noinspection GoBoolExpressions
func GetUserGeneralRootDirLocation() string {
	if IsDarwin() {
		return "/Users"
	} else if runtime.GOOS == "windows" {
		panic("not support windows")
	} else if IsLinux() {
		return "/home"
	} else {
		panic("not support " + runtime.GOOS)
	}
}
