package picovpnd

import (
	"encoding/json"
	"net"

	"github.com/anatolio-deb/picovpnd/common"
	"github.com/sirupsen/logrus"
)

func UserAdd(username, password string) common.Response {
	connection, err := net.Dial("unix", "picovpn.ru:5000")
	if err != nil {
		panic(err)
	}
	defer connection.Close()
	request := common.AddUserRequest{
		Username: username,
		Password: password,
	}
	b, err := json.Marshal(request)
	if err != nil {
		logrus.Error(err)
	}
	_, err = connection.Write(b)
	if err != nil {
		logrus.Error(err)
	}
	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer)
	if err != nil {
		logrus.Error(err)
	}
	var resp common.Response
	err = json.Unmarshal(buffer[:mLen], &resp)
	if err != nil {
		logrus.Error(err)
	}
	return resp
}
