package ocserv

import (
	"os/exec"
	"time"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
	"github.com/sirupsen/logrus"
)

// ocpasswd -c /etc/ocserv/ocpasswd username
func UserAdd(username, password string) error {
	// return cryptInt("/etc/ocserv/ocpasswd", username, "*", password)
	cmd := exec.Command("ocpasswd", "-c", "/etc/ocserv/ocpasswd", username)

	err := cmd.Start()
	if err != nil {
		return err
	}

	time.Sleep(time.Second)
	keyboard.SimulateKeyPress(password)
	keyboard.SimulateKeyPress(keys.Enter)
	time.Sleep(time.Second)
	keyboard.SimulateKeyPress(password)
	keyboard.SimulateKeyPress(keys.Enter)

	return cmd.Wait()
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
