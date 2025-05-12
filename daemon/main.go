package main

import (
	"encoding/json"
	"net"

	"github.com/anatolio-deb/picovpnd/common"
	"github.com/sirupsen/logrus"
)

var HOSTNAME string

func main() {
	server, err := net.Listen("tcp", "picovpn.ru:5000")
	if err != nil {
		panic(err)
	}
	defer server.Close()
	for {
		connection, err := server.Accept()
		if err != nil {
			panic(err)
		}
		go handler(connection)
	}
}

func handler(connection net.Conn) {
	defer connection.Close()
	var err error
	resp := common.Response{}

	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer)
	if err == nil {
		req := common.AddUserRequest{}
		err = json.Unmarshal(buffer[:mLen], &req)
		if err == nil {
			err = common.UserAdd(req.Username, req.Password)
		}
		logrus.Infof("Request create user %s", req.Username)
	}
	if err != nil {
		resp.Code = 1
		resp.Error = err.Error()
	}

	b, err := json.Marshal(resp)
	if err != nil {
		logrus.Error(err)
	}
	_, err = connection.Write(b)
	if err != nil {
		logrus.Error(err)
	}
}