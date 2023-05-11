package client

import (
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"io/ioutil"
	"net"
	"os"
	"strings"

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
		return formatErr(err)
	}
	c.dialer = initDialer(username, psw)
	c.session, err = c.initSession()
	if err != nil {
		return formatErr(err)
	}
	c.share, err = c.Mount(shareName)
	if err != nil {
		return formatErr(err)
	}
	return nil
}
func initConn(addressWithPort string) (net.Conn, error) {
	conn, err := net.Dial("tcp", addressWithPort)
	if err != nil {
		return nil, formatErr(err)
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
	if c == nil || c.dialer == nil {
		return nil, formatErr(errors.New("Client Dialer not initialized"))
	}
	s, err := c.dialer.Dial(c.conn)
	if err != nil {
		fmt.Println(err)
		return nil, formatErr(err)
	}
	return s, formatErr(err)
}

func (c *Client) IsConnected() bool {
	if c != nil && c.conn != nil && c.dialer != nil && c.share != nil {
		_, err := c.GetShares()
		return err == nil
	}
	return false
}
func (c *Client) CloseConn() {
	if c != nil {
		if c.share != nil {
			defer c.share.Umount()
		}
		if c.session != nil {
			defer c.session.Logoff()
		}
		if c.conn != nil {
			defer c.conn.Close()
		}
	}

}

func (c *Client) GetShares() ([]string, error) {

	if c == nil || c.session == nil {
		return nil, formatErr(errors.New("Client not initialized"))
	}
	names, err := c.session.ListSharenames()
	if err != nil {
		return nil, formatErr(err)
	}

	return names, nil

}

func (c *Client) Mount(shareName string) (*smb2.Share, error) {
	if c == nil || c.session == nil {
		return nil, formatErr(errors.New("Client not initialized"))
	}

	fs, err := c.session.Mount(shareName)
	if err != nil {
		return nil, formatErr(err)
	}
	return fs, nil

}

func (c *Client) WriteStringFromOffset(fileName string, strToWrite string, offset int64) error {
	return c.WriteFromOffset(fileName, []byte(strToWrite), offset)
}
func (c *Client) WriteFromOffset(fileName string, bytes []byte, offset int64) error {
	f, err := openOrCreate(c, fileName)
	if err != nil {
		fmt.Println(err)
		return formatErr(err)
	}
	defer f.Close()
	_, errStats := f.Stat()
	if errStats != nil {
		return formatErr(errStats)
	}
	_, err = f.WriteAt(bytes, offset)
	if err != nil {
		fmt.Println(err)
		return formatErr(err)
	}
	return nil
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
		return formatErr(err)
	}
	defer f.Close()
	stats, errStats := f.Stat()
	if errStats != nil {
		return formatErr(errStats)
	}

	if stats.Size() > 0 && newLine {
		bytes = []byte("\n" + string(bytes))
	}
	_, err = f.WriteAt(bytes, stats.Size())
	if err != nil {
		fmt.Println(err)
		return formatErr(err)
	}
	return nil
}

func (c *Client) ReadFile(fileName string) (string, error) {

	f, err := openFile(c, fileName)
	if err != nil {
		fmt.Println("errrr: s%", err)
		return "", formatErr(err)
	}

	defer f.Close()
	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		return "", formatErr(err)
	}
	bs, err := ioutil.ReadAll(f)
	if err != nil {
		return "", formatErr(err)
	}

	return string(bs), nil
}

func (c *Client) ReadFileWithOffsets(fileName string, offsetStart int64, dimesion int64) (string, error) {

	f, err := openFile(c, fileName)
	if err != nil {
		return "", formatErr(err)
	}

	defer f.Close()
	_, err = f.Seek(offsetStart, io.SeekStart)
	if err != nil {
		return "", formatErr(err)
	}
	bs, err := ioutil.ReadAll(f)
	if err != nil {
		return "", formatErr(err)
	}

	if dimesion > 0 {
		return string(bs[0 : dimesion+1]), nil
	}
	return string(bs), nil

}

func (c *Client) RemoveFile(fileName string) error {
	if c == nil || c.share == nil {
		return formatErr(errors.New("Client not initialized"))
	}
	err := c.share.Remove(fileName)
	return formatErr(err)

}

func (c *Client) RenameFile(pathOld string, pathNew string) error {
	if c == nil || c.share == nil {
		return formatErr(errors.New("Client not initialized"))
	}
	err := c.share.Rename(pathOld, pathNew)
	return formatErr(err)
}
func (c *Client) CheckIfFileExists(name string) (bool, error) {
	if c == nil || c.share == nil {
		return false, formatErr(errors.New("Client not initialized"))
	}
	_, err := c.share.Lstat(name)
	if err != nil {
		return false, formatErr(err)
	}
	return true, formatErr(err)
}

func (c *Client) GetFileStats(name string) (fs.FileInfo, error) {
	if c == nil || c.share == nil {
		return nil, formatErr(errors.New("Client not initialized"))
	}
	stats, err := c.share.Lstat(name)
	if err != nil {
		return nil, formatErr(err)
	}
	return stats, formatErr(err)
}
func (c *Client) GetStats(path string) (fs.FileInfo, error) {
	isDir, err := c.IsDir(path)
	if err != nil {
		return nil, formatErr(err)
	}
	if isDir {
		return c.GetFolderStats(path)
	}
	return c.GetFileStats(path)
}
func (c *Client) GetStatsAsString(path string) (string, error) {
	stats, errStats := c.GetStats(path)
	if errStats != nil {
		fmt.Println(errStats)
		return "", errStats
	}
	out, errJ := json.Marshal(stats)
	if errJ != nil {
		return "", errJ
	}
	return string(out), nil
}
func (c *Client) CreateFolder(name string) error {
	if c == nil || c.share == nil {
		return formatErr(errors.New("Client not initialized"))
	}
	err := c.share.Mkdir(name, os.ModeDir)
	return formatErr(err)

}
func (c *Client) RenameFolder(oldPath string, newPath string) error {
	if c == nil || c.share == nil {
		return formatErr(errors.New("Client not initialized"))
	}
	err := c.share.Rename(oldPath, newPath)
	return formatErr(err)
}
func (c *Client) CheckIfFolderExists(name string) (bool, error) {
	if c == nil || c.share == nil {
		return false, formatErr(errors.New("Client not initialized"))
	}
	_, err := c.share.ReadDir(name)
	if err != nil {
		return false, formatErr(err)
	}
	return true, formatErr(err)
}
func (c *Client) GetFolderStats(name string) (fs.FileInfo, error) {
	if c == nil || c.share == nil {
		return nil, formatErr(errors.New("Client not initialized"))
	}
	stats, err := c.share.Stat(name)
	if err != nil {
		return nil, formatErr(err)
	}
	return stats, formatErr(err)
}
func (c *Client) DeleteFolder(name string) error {
	if c == nil || c.share == nil {
		return formatErr(errors.New("Client not initialized"))
	}
	fInfo, err := c.share.ReadDir(name)
	if err != nil {
		return formatErr(err)
	}
	if len(fInfo) > 0 {
		return formatErr(errors.New("directory is not empty"))
	}
	err = c.share.Remove(name)
	if err != nil {
		return formatErr(err)
	}
	return formatErr(err)
}

func (c *Client) IsDir(name string) (bool, error) {
	if c == nil || c.share == nil {
		return false, formatErr(errors.New("Client not initialized"))
	}
	fInfo, err := c.share.Stat(name)
	if err != nil {
		return false, formatErr(err)
	}
	return fInfo.IsDir(), nil
}
func (c *Client) ListFilesInFolder(name string) ([]string, error) {
	if c == nil || c.share == nil {
		return nil, formatErr(errors.New("Client not initialized"))
	}
	fInfo, err := c.share.ReadDir(name)
	if err != nil {
		return nil, formatErr(err)
	}
	if len(fInfo) == 0 {
		return nil, formatErr(err)
	}
	var retArray []string
	for _, info := range fInfo {
		retArray = append(retArray, info.Name())
	}
	return retArray, err

}
func (c *Client) DeleteFile(name string) error {
	if c == nil || c.share == nil {
		return formatErr(errors.New("Client not initialized"))
	}
	err := c.share.Remove(name)
	return formatErr(err)
}
func openOrCreate(c *Client, fileName string) (*smb2.File, error) {
	if c == nil || c.share == nil {
		return nil, formatErr(errors.New("Client not initialized"))
	}
	f, err := c.share.OpenFile(fileName, os.O_APPEND, os.ModeAppend)
	if err != nil {
		f, err = c.share.Create(fileName)
	}
	return f, formatErr(err)
}
func openFile(c *Client, fileName string) (*smb2.File, error) {
	if c == nil || c.share == nil {
		return nil, formatErr(errors.New("Client not initialized"))
	}
	f, err := c.share.OpenFile(fileName, os.O_APPEND, os.ModeAppend)
	if err != nil {
		return nil, formatErr(err)
	}
	return f, err
}

// formatErr add #ERROR# prefix to identify error in retur string
func formatErr(err error) error {
	if err != nil && (!strings.Contains(err.Error(), "#ERROR#")) {
		return errors.New("#ERROR# " + err.Error())
	}
	return err
}

// FormatErrStr add #ERROR# prefix to identify error in retur string
func (c *Client) FormatErrStr(err error) string {
	if err != nil {
		return (errors.New("#ERROR# " + err.Error())).Error()
	}
	return ""
}
