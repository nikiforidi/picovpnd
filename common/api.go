package common

import (
	"os/exec"

	"github.com/sirupsen/logrus"
)

type AddUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

func UserAdd(username, password string) error {
	b, err := exec.Command("echo", password, "|", "ocpasswd", "-c", "/etc/ocserv/ocpasswd", username).CombinedOutput()
	if err != nil {
		return err
	}
	logrus.Debug(b)
	return nil
}
