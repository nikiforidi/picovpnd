package ocserv

import (
	"os"
	"os/exec"
	"time"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
	expect "github.com/Netflix/go-expect"
	"github.com/sirupsen/logrus"
)

// ocpasswd -c /etc/ocserv/ocpasswd username
func UserAdd(username, password string) error {
	// return cryptInt("/etc/ocserv/ocpasswd", username, "*", password)
	c, err := expect.NewConsole(expect.WithStdout(os.Stdout))
	if err != nil {
		return err
	}
	defer c.Close()

	cmd := exec.Command("ocpasswd", "-c", "/etc/ocserv/ocpasswd", username)
	cmd.Stdin = c.Tty()
	cmd.Stdout = c.Tty()
	cmd.Stderr = c.Tty()

	go func() {
		c.ExpectEOF()
	}()

	err = cmd.Start()
	if err != nil {
		return err
	}

	time.Sleep(time.Second)
	c.Send(password)
	keyboard.SimulateKeyPress(keys.Enter)
	time.Sleep(time.Second)
	c.Send(password)
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
