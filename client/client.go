package client

import (
	"io"
	"io/ioutil"
	"net"
	"os"

	"fmt"

	"github.com/hirochachacha/go-smb2"
)

type Client struct {
	conn    net.Conn
	dialer  *smb2.Dialer
	session *smb2.Session
	share   *smb2.Share
}

func (c *Client) NewClient(addressWithPort string, username string, psw string, shareName string) error {
	var err error
	c.conn, err = initConn(addressWithPort)
	if err != nil {
		return err
	}
	c.dialer = initDialer(username, psw)
	c.session, err = c.initSession()
	if err != nil {
		return err
	}
	c.share, err = c.Mount(shareName)
	if err != nil {
		return err
	}
	return nil
}
func initConn(addressWithPort string) (net.Conn, error) {
	conn, err := net.Dial("tcp", addressWithPort)
	if err != nil {
		return nil, err
	}
	return conn, nil

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
func (c *Client) initSession() (*smb2.Session, error) {
	s, err := c.dialer.Dial(c.conn)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return s, nil
}
func (c *Client) CloseConn() {
	defer c.share.Umount()
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

func (c *Client) Mount(shareName string) (*smb2.Share, error) {
	fs, err := c.session.Mount(shareName)
	if err != nil {
		return nil, err
	}

	return fs, nil
}

func (c *Client) AppendLine(fileName string, strToWrite string) error {
	return c.AppendBytes(fileName, []byte(strToWrite), true)
}
func (c *Client) AppendString(fileName string, strToWrite string) error {
	return c.AppendBytes(fileName, []byte(strToWrite), false)
}

func (c *Client) AppendBytes(fileName string, bytes []byte, newLine bool) error {

	f, err := openOrCreate(c, fileName)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer f.Close()
	stats, errStats := f.Stat()
	if errStats != nil {
		return errStats
	}

	if stats.Size() > 0 && newLine {
		bytes = []byte("\n" + string(bytes))
	}
	_, err = f.WriteAt(bytes, stats.Size())
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (c *Client) ReadFile(fileName string) (string, error) {

	f, err := openFile(c, fileName)
	if err != nil {
		fmt.Println("errrr: s%", err)
		return "", err
	}

	defer f.Close()
	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		panic(err)
	}
	bs, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	return string(bs), nil
}

func (c *Client) ReadFileWithOffsets(fileName string, offsetStart int64, dimesion int64) (string, error) {

	f, err := openFile(c, fileName)
	if err != nil {
		fmt.Println("errrr: s%", err)
		return "", err
	}

	defer f.Close()
	_, err = f.Seek(offsetStart, io.SeekStart)
	if err != nil {
		panic(err)
	}
	bs, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	if dimesion > 0 {
		return string(bs[0 : dimesion+1]), nil
	}
	return string(bs), nil

}

func (c *Client) RemoveFile(fileName string) error {
	err := c.share.Remove(fileName)
	return err
}
func (c *Client) CreateFolder(name string) error {
	err := c.share.Mkdir(name, os.ModeDir)
	return err
}
func (c *Client) CheckIfFolderExists(name string) (bool, error) {
	_, err := c.share.ReadDir(name)
	if err != nil {
		return false, err
	}
	return true, err
}
func openOrCreate(c *Client, fileName string) (*smb2.File, error) {
	f, err := c.share.OpenFile(fileName, os.O_APPEND, os.ModeAppend)
	if err != nil {
		f, err = c.share.Create(fileName)
	}
	return f, err
}
func openFile(c *Client, fileName string) (*smb2.File, error) {
	f, err := c.share.OpenFile(fileName, os.O_APPEND, os.ModeAppend)
	if err != nil {
		return nil, err
	}
	return f, err
}
