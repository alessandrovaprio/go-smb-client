package client

import (
	"net"

	"fmt"

	"github.com/hirochachacha/go-smb2"
)

type Client struct {
	conn    net.Conn
	dialer  *smb2.Dialer
	session *smb2.Session
}

func (c *Client) NewClient(addressWithPort string, username string, psw string) {
	c.conn = initConn(addressWithPort)
	c.dialer = initDialer(username, psw)
	c.session = c.initSession()
}
func initConn(addressWithPort string) net.Conn {
	conn, err := net.Dial("tcp", addressWithPort)
	if err != nil {
		panic(err)
	}
	return conn

}
func initDialer(user string, psw string) *smb2.Dialer {
	dialer := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     user,
			Password: psw,
		},
	}
	return dialer
}
func (c *Client) initSession() *smb2.Session {
	s, err := c.dialer.Dial(c.conn)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return s
}
func (c *Client) CloseConn() {
	defer c.session.Logoff()
	defer c.conn.Close()

}

func (c *Client) GetShares() ([]string, error) {
	// defer c.session.Logoff()

	names, err := c.session.ListSharenames()
	if err != nil {
		return nil, err
	}

	return names, nil
}
