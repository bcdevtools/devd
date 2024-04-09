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
