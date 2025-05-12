package main

import (
	"encoding/json"
	"net"

	"github.com/anatolio-deb/picovpnd/common"
	"github.com/sirupsen/logrus"
)

var HOSTNAME string

func main() {
	ip := GetPublicIP()
	if ip == ""{
		ip = GetLocalHostname()
	}
	addr := ip+":"+"5000"
	server, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	logrus.Debugf("Listening on %s", addr)
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