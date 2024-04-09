package utils

import (
	"fmt"
	libutils "github.com/EscanBE/go-lib/utils"
	"github.com/pkg/errors"
	"io/fs"
	"os"
	"strconv"
)

func ValidatePasswordFileMode(mode fs.FileMode) error {
	str := fmt.Sprintf("%o", int(mode))
	symbolNum, err := strconv.ParseInt(str, 10, 64)
	libutils.PanicIfErr(err, fmt.Sprintf("failed to parse %s", str))
	if symbolNum%100 != 0 {
		return fmt.Errorf("not allowed to have permission for group/other")
	}
	if symbolNum < 400 {
		return fmt.Errorf("require read permission")
	}
	return nil
}

func IsFileAndExists(path string) (exists bool, perm fs.FileMode, err error) {
	if libutils.IsBlank(path) {
		panic(fmt.Errorf("input is empty"))
	}

	var fi os.FileInfo
	fi, err = os.Stat(path)
	if err == nil {
		if !fi.IsDir() {
			exists = true
			perm = fi.Mode().Perm()
		}
	} else {
		if os.IsNotExist(err) {
			err = nil
		} else {
			err = errors.Wrap(err, fmt.Sprintf("problem while checking target file %s", path))
		}
	}

	return
}

func IsDirAndExists(path string) (exists bool, perm fs.FileMode, err error) {
	if libutils.IsBlank(path) {
		panic(fmt.Errorf("input is empty"))
	}

	var di os.FileInfo
	di, err = os.Stat(path)
	if err == nil {
		if di.IsDir() {
			exists = true
			perm = di.Mode().Perm()
		}
	} else {
		if os.IsNotExist(err) {
			err = nil
		} else {
			err = errors.Wrap(err, fmt.Sprintf("problem while checking target directory %s", path))
		}
	}

	return
}

func ExtractPermissionParts(perm fs.FileMode) (owner, group, other int) {
	str := fmt.Sprintf("%o", int(perm))
	symbolNum, err := strconv.ParseInt(str, 10, 64)
	libutils.PanicIfErr(err, fmt.Sprintf("failed to parse %s", str))
	sym := int(symbolNum)
	owner = sym / 100
	group = (sym % 100) / 10
	other = sym % 10
	return
}
