package core

import (
	"os"
	"os/exec"
	"time"

	"github.com/Netflix/go-expect"
)

func UserAdd(username, password string) error {
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
	_, err := exec.Command("ocpasswd", "--lock", username).CombinedOutput()
	return err
}

func UserUnlock(username string) error {
	_, err := exec.Command("ocpasswd", "--unlock", username).CombinedOutput()
	return err
}
