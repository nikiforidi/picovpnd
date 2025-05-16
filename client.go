package picovpnd

import (
	"encoding/json"
	"net"

	"github.com/anatolio-deb/picovpnd/common"
)

type Client struct {
	Network string // TCP
	Address string // addr:port
	conn    net.Conn
	resp    common.Response
}

func New(network, address string) (Client, error) {
	c := Client{Network: network, Address: address, resp: common.Response{}}
	conn, err := net.Dial("tcp", "picovpn.ru:5000")
	if err != nil {
		return c, err
	}
	c.conn = conn
	return c, err
}

func (c Client) Send(req common.Request) common.Response {
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

func (c Client) UserAdd(username, password string) common.Response {
	defer c.conn.Close()
	// resp := common.Response{}
	payload := common.UserAddPayload{
		UserMixin: common.UserMixin{Username: username},
		Password:  password,
	}
	request := common.Request{
		Method:  common.UserAdd,
		Payload: payload,
	}
	return c.Send(request)
}

func (c Client) UserLock(username string) common.Response {
	defer c.conn.Close()
	payload := common.UserMixin{
		Username: username,
	}
	request := common.Request{
		Method:  common.UserLock,
		Payload: payload,
	}
	return c.Send(request)
}
