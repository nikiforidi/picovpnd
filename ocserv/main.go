package ocserv

import (
	"os"
	"os/exec"
	"time"

	expect "github.com/Netflix/go-expect"
	"github.com/sirupsen/logrus"
)

func UserAdd(username, password string) error {
	// return cryptInt("/etc/ocserv/ocpasswd", username, "*", password)
	c, err := expect.NewConsole(expect.WithStdout(os.Stdout))
	if err != nil {
		return err
	}
	defer c.Close()

	cmd := exec.Command("ocpasswd", "-c", "/etc/ocserv/passwd", username)
	cmd.Stdin = c.Tty()
	cmd.Stdout = c.Tty()
	cmd.Stderr = c.Tty()

	err = cmd.Start()
	if err != nil {
		return err
	}

	go func() {
		c.ExpectString("Enter password:")
	}()

	time.Sleep(time.Second)
	c.SendLine(password)

	go func() {
		c.ExpectString("Re-enter password:")
	}()

	time.Sleep(time.Second)
	c.SendLine(password)

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
