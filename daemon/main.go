package main

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"net"

	"github.com/anatolio-deb/picovpnd/common"
	"github.com/sirupsen/logrus"
)

func main() {
	cert, err := tls.LoadX509KeyPair(common.CertificateFile, "/etc/letsencrypt/live/picovpn.ru/privkey.pem")
	if err != nil {
		log.Fatal(err)
	}
	config := &tls.Config{Certificates: []tls.Certificate{cert}}
	l, err := tls.Listen("tcp", common.ListenAddress, config)
	if err != nil {
		logrus.Fatal(err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("accepted connection from %s\n", conn.RemoteAddr())

		go handler(conn)
	}

}

func handler(connection net.Conn) {
	defer connection.Close()
	var err error
	resp := common.Response{}

	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer)
	if err == nil {
		req := common.Request{}
		err = json.Unmarshal(buffer[:mLen], &req)
		if err == nil {
			err := common.PayloadDispatcher(req)
			if err != nil {
				logrus.Error(err)
			}
		}
	}
	if err != nil {
		resp.Code++
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
	logrus.Debug("leaving handler")
}
