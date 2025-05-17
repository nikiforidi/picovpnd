package common

import (
	"encoding/json"
	"fmt"

	"github.com/anatolio-deb/picovpnd/ocserv"
	"github.com/sirupsen/logrus"
)

type Method string

const (
	CertificateFile = "/etc/letsencrypt/live/picovpn.ru/fullchain.pem"
	ListenAddress   = "picovpn.ru:5000"
	UserAdd         = Method("user_add")
	UserLock        = Method("user_lock")
	UserUnlock      = Method("user_unlock")
)

type Request struct {
	Method  Method `json:"method"`
	Payload any    `json:"payload"`
}

func (r *Request) SetPayload() error {
	b, err := json.Marshal(r.Payload)
	if err != nil {
		return err
	}
	r.Payload = b
	return nil
}

type UserMixin struct {
	Username string `json:"username"`
}

type UserAddPayload struct {
	UserMixin
	Password string `json:"password"`
}

// type UserLockPayload struct {
// 	UserMixin
// }

// type UserUnlockPayload struct {
// 	UserMixin
// }

func PayloadDispatcher(req Request) error {
	logrus.Infof("Dispatching %s request", req.Method)
	logrus.Infof("Request payload: %v", req.Payload)
	switch req.Method {
	case UserAdd:
		payload, ok := req.Payload.(map[string]string)
		if ok {
			username, ok := payload["username"]
			if !ok {
				return fmt.Errorf("bad request: %s", req.Method)
			}
			logrus.Infof("Request create user %s", username)
			password, ok := payload["password"]
			if !ok {
				return fmt.Errorf("bad request: %s", req.Method)
			}

			return ocserv.UserAdd(username, password)
		} else {
			return fmt.Errorf("bad request: %s", req.Method)
		}
	case UserLock:
		p, ok := req.Payload.(UserMixin)
		if ok {
			logrus.Infof("Request lock user %s", p.Username)
			return ocserv.UserLock(p.Username)
		} else {
			return fmt.Errorf("bad request: %s", req.Method)
		}
	case UserUnlock:
		p, ok := req.Payload.(UserMixin)
		if ok {
			logrus.Infof("Request unlock user %s", p.Username)
			return ocserv.UserUnlock(p.Username)
		} else {
			return fmt.Errorf("bad request: %s", req.Method)
		}
	default:
		return fmt.Errorf("bad request: %s", req.Method)
	}
	return nil
}

type Response struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}
