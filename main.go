package main

/*
#cgo CFLAGS: -DPNG_DEBUG=1
#cgo amd64 386 CFLAGS: -DX86=1
#cgo LDFLAGS: -lstdc++ -lm
#include <stdio.h>
#include <stdlib.h>
#include <errno.h>
*/
import (
	"C"
	"net"

	Client "github.com/alessandrovaprio/go-smb-client/client"
	"github.com/hirochachacha/go-smb2"
)
import (
	"fmt"
	"strings"
)

// global client
var gClient *Client.Client

func main() {

}

//export InitClient
func InitClient(addressWithPort *C.char, user *C.char, psw *C.char, shareName *C.char) *C.char {
	goAddressWithPort := C.GoString(addressWithPort)
	goUser := C.GoString(user)
	goPsw := C.GoString(psw)
	goShareName := C.GoString(shareName)
	gClient = new(Client.Client)
	err := gClient.NewClient(goAddressWithPort, goUser, goPsw, goShareName)
	if err != nil {
		return C.CString(err.Error())
	}
	return nil
}

//export CloseConn
func CloseConn() {
	gClient.CloseConn()
}

//export IsConnected
func IsConnected() bool {
	return gClient.IsConnected()
}

//export ListShares
func ListShares() *C.char {

	names, err := gClient.GetShares()
	retErr := ""
	if err != nil {
		retErr = err.Error()
		return C.CString(retErr)
	}
	var sb strings.Builder
	for idx, name := range names {
		sb.WriteString(name)
		if idx != len(names) {
			sb.WriteString(",")
		}
		fmt.Println(name)
	}
	return C.CString(sb.String())
	// return C.CString(fmt.Sprintf("%s", ""))
}

//export WriteStringFromOffset
func WriteStringFromOffset(fileName *C.char, strToWrite *C.char, offset int64) *C.char {
	goFileName := C.GoString(fileName)
	goStrToWrite := C.GoString(strToWrite)

	err := gClient.WriteStringFromOffset(goFileName, goStrToWrite, offset)
	if err != nil {
		return C.CString(err.Error())
	}
	return nil
}

//export AppendLine
func AppendLine(fileName *C.char, strToWrite *C.char) *C.char {
	goFileName := C.GoString(fileName)
	goStrToWrite := C.GoString(strToWrite)
	err := gClient.AppendLine(goFileName, goStrToWrite)
	retErr := ""
	if err != nil {
		retErr = err.Error()
		return C.CString(retErr)
	}
	return nil
}

//export ReadFile
func ReadFile(fileName *C.char) *C.char {
	goFileName := C.GoString(fileName)
	mystr, err := gClient.ReadFile(goFileName)
	retErr := ""
	if err != nil {
		retErr = err.Error()
		return C.CString(retErr)
	}
	return C.CString(mystr)
}

//export ReadFileWithOffsets
func ReadFileWithOffsets(fileName *C.char, offset int64, chunckSize int64) *C.char {
	goFileName := C.GoString(fileName)
	mystr, err := gClient.ReadFileWithOffsets(goFileName, offset, chunckSize)
	retErr := ""
	if err != nil {
		retErr = err.Error()
		return C.CString(fmt.Sprintf("%s", retErr))
	}
	fmt.Println("mystr: %s", mystr)
	return C.CString(fmt.Sprintf("%s", mystr))
}

//export RemoveFile
func RemoveFile(fileName *C.char) *C.char {
	goFileName := C.GoString(fileName)
	err := gClient.RemoveFile(goFileName)
	retErr := ""
	if err != nil {
		retErr = err.Error()
		return C.CString(fmt.Sprintf("%s", retErr))
	}
	return nil
}

//export RenameFile
func RenameFile(oldPath *C.char, newPath *C.char) *C.char {
	oldP := C.GoString(oldPath)
	newP := C.GoString(newPath)
	err := gClient.RenameFile(oldP, newP)
	if err != nil {
		return C.CString(err.Error())
	}
	return nil
}

//export CreateFolder
func CreateFolder(name *C.char) *C.char {
	goFileName := C.GoString(name)
	err := gClient.CreateFolder(goFileName)

	if err != nil {
		return C.CString(err.Error())
	}
	return nil
}

//export DeleteFolder
func DeleteFolder(name *C.char) *C.char {
	goFileName := C.GoString(name)
	err := gClient.DeleteFolder(goFileName)

	if err != nil {
		return C.CString(err.Error())
	}
	return nil
}

//export CheckIfFolderExists
func CheckIfFolderExists(name *C.char) *C.char {
	goFileName := C.GoString(name)
	exists, err := gClient.CheckIfFolderExists(goFileName)

	if err != nil {
		return C.CString(err.Error())
	}
	return C.CString(fmt.Sprintf("%t", exists))
}

//export ListAllShares
func ListAllShares(addressWithPort *C.char, user *C.char, psw *C.char) *C.char {
	goAddressWithPort := C.GoString(addressWithPort)
	fmt.Println(goAddressWithPort)
	goUser := C.GoString(user)
	fmt.Println(goUser)
	goPsw := C.GoString(psw)
	fmt.Println(goPsw)
	retErr := ""
	names, err := getAllShares(goAddressWithPort, goUser, goPsw)
	if err != nil {
		retErr = err.Error()
		return C.CString(fmt.Sprintf("%s", retErr))
	}
	for _, name := range names {
		fmt.Println(name)
	}
	return nil
	// return C.CString(fmt.Sprintf("%s", ""))
}

func getAllShares(addressWithPort string, user string, psw string) ([]string, error) {
	conn, err := net.Dial("tcp", addressWithPort)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	defer conn.Close()

	d := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     user,
			Password: psw,
		},
	}
	s, err := d.Dial(conn)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer s.Logoff()

	names, err := s.ListSharenames()
	if err != nil {
		return nil, err
	}

	return names, nil

}
