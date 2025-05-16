package picovpnd

import (
	"encoding/json"
	"net"

	"github.com/anatolio-deb/picovpnd/common"
)

func UserAdd(username, password string) common.Response {
	resp := common.Response{}
	connection, err := net.Dial("tcp", "picovpn.ru:5000")
	if err != nil {
		resp.Error = err.Error()
		resp.Code++
		return resp
	}
	defer connection.Close()
	request := common.AddUserRequest{
		Username: username,
		Password: password,
	}
	b, err := json.Marshal(request)
	if err != nil {
		resp.Error = err.Error()
		resp.Code++
		return resp
	}
	_, err = connection.Write(b)
	if err != nil {
		resp.Error = err.Error()
		resp.Code++
		return resp
	}
	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer)
	if err != nil {
		resp.Error = err.Error()
		resp.Code++
		return resp
	}
	err = json.Unmarshal(buffer[:mLen], &resp)
	if err != nil {
		resp.Error = err.Error()
		resp.Code++
		return resp
	}
	return resp
}
