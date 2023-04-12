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
)

// global client
var gClient *Client.Client

func main() {
	shares, err := getAllShares("myadd", "myuser", "mypsw")
	if err != nil {
		panic(err)
	}
	for _, name := range shares {
		fmt.Println(name)
	}

}

//export InitClient
func InitClient(addressWithPort *C.char, user *C.char, psw *C.char, shareName *C.char) *C.char {
	goAddressWithPort := C.GoString(addressWithPort)
	goUser := C.GoString(user)
	goPsw := C.GoString(psw)
	goShareName := C.GoString(shareName)
	gClient = new(Client.Client)
	gClient.NewClient(goAddressWithPort, goUser, goPsw, goShareName)
	return nil
}

//export CloseConn
func CloseConn() {
	gClient.CloseConn()
}

//export ListShares
func ListShares() *C.char {

	names, err := gClient.GetShares()
	retErr := ""
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

//export AppendLine
func AppendLine(fileName *C.char, strToWrite *C.char) *C.char {
	goFileName := C.GoString(fileName)
	goStrToWrite := C.GoString(strToWrite)
	err := gClient.AppendLine(goFileName, goStrToWrite)
	retErr := ""
	if err != nil {
		retErr = err.Error()
		return C.CString(fmt.Sprintf("%s", retErr))
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

// //export ListEnrolledFingers
// func ListEnrolledFingers(busType *C.char, destination *C.char, destInterface *C.char, operation *C.char, username operation *C.char) *C.char {
// 	goBusType := C.GoString(busType)
// 	goDestination := C.GoString(destination)
// 	goOperation := C.GoString(operation)
// 	goDestInterface := C.GoString(destInterface)
// 	goUsername := C.GoString(username)
// 	retVal := listEnrolledFingers(goBusType, goDestination, goDestInterface, goOperation, goUsername)
// 	tmp := C.CString(fmt.Sprintf("%s", retVal))
// 	return tmp
// }

// func listEnrolledFingers(busType string, destination string, destInterface string, operation string, username string) string {
// 	conn, err := dbus.SessionBus()
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(conn.Names())
// 	// func (conn *Conn) Object(dest string, path ObjectPath) *Object
// 	obj := conn.Object(destination, dbus.ObjectPath(destInterface))
// 	defer conn.Close()
// 	var busReturn []dbus.ObjectPath
// 	retCall := obj.Call(helpers.METHOD_ENROLL_START, 0, username)
// 	if retCall != nil {
// 		fmt.Fprintln(os.Stderr, "Failed to call "+operation+" function (is the server example running?):", retCall)
// 		os.Exit(1)
// 	}
// 	var values []string
// 	for _, v := range busReturn {
// 		tmp := fmt.Sprintf("%v", v)
// 		values = append(values, tmp)
// 	}
// 	retVal := strings.Join(values, ",")
// 	return retVal
// }

// // func claimDevice(busType string, destination string, destInterface string, operation string, username string) *dbus.Call {
// // 	obj, conn := getDBusObj(busType, destination, destInterface)
// // 	defer conn.Close()
// // 	err := obj.Call(operation, 0, "")
// // 	if err != nil {
// // 		fmt.Fprintln(os.Stderr, "Failed to call "+operation+" function (is the server example running?):", err)
// // 		return err
// // 	}
// // 	return nil

// // }
// // func releaseDevice(busType string, destination string, destInterface string, operation string) error {
// // 	obj, conn := getDBusObj(busType, destination, destInterface)
// // 	defer conn.Close()
// // 	var busReturn []dbus.ObjectPath
// // 	err := obj.Call(operation, 0).Store(&busReturn)
// // 	if err != nil {
// // 		fmt.Fprintln(os.Stderr, "Failed to call "+operation+" function (is the server example running?):", err)
// // 		return err
// // 	}
// // 	return nil

// // }

// func enroll(busType string, destination string, destInterface string, operation string, fingerName string) string {
// 	conn, err := dbus.SessionBus()
// 	if err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}

// 	obj := conn.Object(destination, dbus.ObjectPath(destInterface))
// 	obj.Call(helpers.METHOD_CLAIM, 0, "")

// 	defer conn.Close()
// 	// var busReturn []dbus.ObjectPath
// 	obj.Call(helpers.METHOD_ENROLL_START, 0, fingerName)
// 	obj.Call(helpers.METHOD_ENROLL_START, 0, fingerName)
// 	// if retCall != nil {
// 	// 	fmt.Fprintln(os.Stderr, "Failed to call "+operation+" function (is the server example running?):", retCall)
// 	// 	os.Exit(1)
// 	// }

// 	retVal := ""
// 	return retVal
// }
// func getDBusObj(conn *dbus.Conn, busType string, destination string, destInterface string) dbus.BusObject {
// 	path, dbusPathErr := helpers.GetDbusPath(destInterface)
// 	if dbusPathErr != nil {
// 		fmt.Fprintln(os.Stderr, "Failed to connect "+dbusPathErr.Error())
// 		os.Exit(1)
// 	}
// 	obj := conn.Object(destination, path)
// 	return obj
// }

// func workingExample() {

// 	conn, err := dbus.SessionBus()
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(conn.Names())
// 	// func (conn *Conn) Object(dest string, path ObjectPath) *Object
// 	obj := conn.Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")
// 	fmt.Println(conn.Names())
// 	// Interface from the specification:
// 	// UINT32 org.freedesktop.Notifications.Notify (STRING app_name, UINT32 replaces_id, STRING app_icon, STRING summary, STRING body, ARRAY actions, DICT hints, INT32 expire_timeout);

// 	// func (o *Object) Call(method string, flags Flags, args ...interface{}) *Call
// 	call := obj.Call("org.freedesktop.Notifications.Notify", 0, "c¼h", uint32(0), "", "Hallo Chaostreff!", "Ich begrüße euch herzlich zu meiner c¼h!", []string{}, map[string]dbus.Variant{}, int32(1000))
// 	if call.Err != nil {
// 		panic(call.Err)
// 	}
// }
