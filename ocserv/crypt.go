package ocserv

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	cryptlib "github.com/tredoe/crypt"
	md5 "github.com/tredoe/crypt/md5_crypt"
	sha256 "github.com/tredoe/crypt/sha256_crypt"
)

const SALT_SIZE = 16
const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func cryptInt(fpasswd, username, groupname, passwd string) error {
	var salt [SALT_SIZE]byte
	var saltStr strings.Builder
	var crPasswd string

	// Generate random salt
	if _, err := exec.Command("sh", "-c", "head -c 16 /dev/urandom").Output(); err != nil {
		return err
	}

	saltStr.WriteString("$1$") // Change to "$5$" for SHA2
	for i := 0; i < SALT_SIZE; i++ {
		saltStr.WriteByte(alphabet[salt[i]%byte(len(alphabet)-1)])
	}
	saltStr.WriteByte('$')

	crPasswd, err := crypt(passwd, saltStr.String(), cryptlib.SHA256)
	if err != nil {
		saltStr_ := strings.Replace(saltStr.String(), "1", "5", 1) // Try MD5
		crPasswd, err = crypt(passwd, saltStr_, cryptlib.MD5)
		if err != nil {
			return err
		}
	}
	if crPasswd == "" {
		return errors.New("error in crypt()")
	}

	tmpPasswd := fmt.Sprintf("%s.tmp", fpasswd)
	if _, err := os.Stat(tmpPasswd); err == nil {
		return err
	}

	fd2, err := os.Create(tmpPasswd)
	if err != nil {
		return err
	}
	defer fd2.Close()

	fd, err := os.Open(fpasswd)
	if err != nil {
		return err
	} else {
		found := false
		lines, _ := io.ReadAll(fd)
		for _, line := range strings.Split(string(lines), "\n") {
			if line == "" {
				continue
			}
			p := strings.Index(line, ":")
			if p == -1 {
				continue
			}
			if len(line[:p]) == len(username) && line[:p] == username {
				fmt.Fprintf(fd2, "%s:%s:%s\n", username, groupname, crPasswd)
				found = true
			} else {
				fmt.Fprintln(fd2, line)
			}
		}
		fd.Close()

		if !found {
			fmt.Fprintf(fd2, "%s:%s:%s\n", username, groupname, crPasswd)
		}
	}
	return os.Rename(tmpPasswd, fpasswd)
}

func crypt(passwd, salt string, algo cryptlib.Crypt) (string, error) {
	var crypter cryptlib.Crypter
	var magic string
	switch algo {
	case cryptlib.SHA256:
		crypter = sha256.New()
		magic = sha256.MagicPrefix
	case cryptlib.MD5:
		crypter = md5.New()
		magic = md5.MagicPrefix
	}

	hash, err := crypter.Generate(
		[]byte(passwd),
		[]byte(magic+salt),
	)
	if err != nil {
		return hash, err
	}
	return hash, crypter.Verify(hash, []byte(passwd))
}
