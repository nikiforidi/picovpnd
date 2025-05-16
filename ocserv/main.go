package ocserv

import (
	"os/exec"

	"github.com/sirupsen/logrus"
)

func UserAdd(username, password string) error {
	return cryptInt("/etc/ocserv/ocpasswd", username, "*", password)
}

func UserLock(username string) error {
	b, err := exec.Command("ocpasswd", "--lock", username).CombinedOutput()
	if err != nil {
		return err
	}
	logrus.Debug(string(b))
	return nil
}

func UserUnlock(username string) error {
	b, err := exec.Command("ocpasswd", "--unlock", username).CombinedOutput()
	if err != nil {
		return err
	}
	logrus.Debug(string(b))
	return nil
}
