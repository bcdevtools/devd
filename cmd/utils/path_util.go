package utils

import (
	"fmt"
	libutils "github.com/EscanBE/go-lib/utils"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
)

func HasToolSshPass() bool {
	cmdApp := exec.Command("sshpass", "-V")
	if err := cmdApp.Run(); err != nil {
		return false
	}
	return true
}

func HasBinaryName(binaryName string) bool {
	_, err := exec.LookPath(binaryName)
	return err == nil
}

func EnsureBinaryNameExists(binaryName string) {
	if !HasBinaryName(binaryName) {
		libutils.PrintlnStdErr("ERR: missing required binary:", binaryName)
		os.Exit(1)
	}
}

func UpdateOwner(path, username, group string, recursive bool) error {
	usr, err := user.Lookup(username)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to lookup user %s", username))
	}

	grp, err := user.LookupGroup(group)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to lookup group %s", group))
	}

	uid, err := strconv.ParseInt(usr.Uid, 10, 64)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to parse uid of user %s with value %s", usr.Username, usr.Uid))
	}
	gid, err := strconv.ParseInt(grp.Gid, 10, 64)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to parse gid of group %s with value %s", grp.Name, grp.Gid))
	}

	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return errors.Wrap(err, fmt.Sprintf("failed to stat %s", path))
	}

	if fi.IsDir() && recursive {
		err = filepath.Walk(path, func(p string, fi os.FileInfo, err error) error {
			if err == nil {
				err = os.Chown(p, int(uid), int(gid))
			}
			return err
		})

		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to recursively update owner of %s to %s:%s", path, usr.Username, grp.Name))
		}
	} else {
		err = os.Chown(path, int(uid), int(gid))

		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to update owner of %s to %s:%s", path, usr.Username, grp.Name))
		}
	}

	return nil
}
