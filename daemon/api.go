package main

import (
	"os/exec"

	"github.com/sirupsen/logrus"
)


func userAdd(username, password string) error {
	b, err := exec.Command("echo", password, "|", "ocpasswd", "-c", "/etc/ocserv/ocpasswd", username).CombinedOutput()
	if err != nil {
		return err
	}
	logrus.Debug(b)
	return nil
}