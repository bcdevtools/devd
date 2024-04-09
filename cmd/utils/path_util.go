package utils

import (
	"os/exec"
)

func HasBinaryName(binaryName string) bool {
	_, err := exec.LookPath(binaryName)
	return err == nil
}
