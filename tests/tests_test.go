package tests

import (
	"fmt"
	"strings"
	"testing"

	Client "github.com/alessandrovaprio/go-smb-client/client"
)

func TestConnect(t *testing.T) {
	client := new(Client.Client)
	err := client.NewClient("127.0.0.1:445", "rio", "letsdance", "Data")
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	client.CloseConn()

}
func TestFolder(t *testing.T) {
	client := new(Client.Client)
	folderName := "myfolder"
	err := client.NewClient("127.0.0.1:445", "rio", "letsdance", "Data")
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	defer client.CloseConn()

	errF := client.CreateFolder(folderName)
	if errF != nil {
		fmt.Println(errF)
		t.Fail()
	}
	stats, errStats := client.GetStatsAsString(folderName)
	if errStats != nil {
		fmt.Println(errStats)
		t.Fail()
	}
	fmt.Println(stats)

}

func TestWriteAndDelete(t *testing.T) {
	fileName := "myfile"
	client := new(Client.Client)
	err := client.NewClient("127.0.0.1:445", "rio", "letsdance", "Data")
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	defer client.CloseConn()
	errW := client.AppendLine(fileName, "myrow")
	if errW != nil {
		fmt.Println(errW)
		t.Fail()
	}
	retOK, errCheck := client.CheckIfFileExists(fileName)
	if errCheck != nil {
		fmt.Println(errCheck)
		t.Fail()
	}
	fmt.Println(retOK)
	stats, errStats := client.GetStatsAsString(fileName)
	if errStats != nil {
		fmt.Println(errStats)
		t.Fail()
	}
	fmt.Println(stats)
	data, errC := client.ReadFile(fileName)
	if errC != nil {
		fmt.Println(errC)
		t.Fail()
	}
	if len(data) == 0 || (!strings.Contains(data, "myrow")) {
		fmt.Println("Not contains 'myrow' string")
		t.Fail()
	}
	fmt.Println("Write Success, file contains:" + data)
	errD := client.RemoveFile(fileName)
	if errD != nil {
		fmt.Println(errD)
		t.Fail()
	}

}
