package picovpnd

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net"

	"github.com/anatolio-deb/picovpnd/common"
)

type client struct {
	Address string // addr:port
	conn    net.Conn
	resp    common.Response
}

func New(address string, port int, cert []byte) (*client, error) {
	// cert, err := os.ReadFile(common.CertificateFile)
	// if err != nil {
	// 	return nil, err
	// }
	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(cert); !ok {
		return nil, fmt.Errorf("unable to parse cert from %s", cert)
	}
	config := &tls.Config{RootCAs: certPool}

	conn, err := tls.Dial("tcp", fmt.Sprintf(address, ":", port), config)
	if err != nil {
		return nil, err
	}
	c := client{Address: address, resp: common.Response{}, conn: conn}
	return &c, err
}

func (c client) Send(req common.Request) common.Response {
	b, err := json.Marshal(req)
	if err != nil {
		c.resp.Error = err.Error()
		c.resp.Code = 1
		return c.resp
	}
	_, err = c.conn.Write(b)
	if err != nil {
		c.resp.Error = err.Error()
		c.resp.Code = 1
		return c.resp
	}
	buffer := make([]byte, 1024)
	mLen, err := c.conn.Read(buffer)
	if err != nil {
		c.resp.Error = err.Error()
		c.resp.Code = 1
		return c.resp
	}
	err = json.Unmarshal(buffer[:mLen], &c.resp)
	if err != nil {
		c.resp.Error = err.Error()
		c.resp.Code = 1
		return c.resp
	}
	return c.resp
}

func (c client) UserAdd(username, password string) common.Response {
	defer c.conn.Close()
	// resp := common.Response{}
	// payload := common.UserAddPayload{
	// 	UserMixin: common.UserMixin{Username: username},
	// 	Password:  password,
	// }
	request := common.Request{
		Method: common.UserAdd,
		Payload: map[string]string{
			"username": username,
			"password": password,
		},
	}
	return c.Send(request)
}

func (c client) UserLock(username string) common.Response {
	defer c.conn.Close()
	// payload := common.UserMixin{
	// 	Username: username,
	// }
	request := common.Request{
		Method: common.UserLock,
		Payload: map[string]string{
			"username": username,
		},
	}
	return c.Send(request)
}
