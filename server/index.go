package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
)

const BufferSize = 1024

func getFileFromClient(connection net.Conn) {
	var currentByte int64 = 0
	fileBuffer := make([]byte, BufferSize)
	var err error
	file, _ := os.Create("receive/" + "awsm.txt")
	for err == nil || err != io.EOF {
		_, _ = connection.Read(fileBuffer)
		cleanedFileBuffer := bytes.Trim(fileBuffer, "\x00")
		_, err = file.WriteAt(cleanedFileBuffer, currentByte)
		if len(string(fileBuffer)) != len(string(cleanedFileBuffer)) {
			break
		}
		currentByte += BufferSize
	}
	_ = connection.Close()
	_ = file.Close()
}

func connectionHandler(connection net.Conn) {
	buffer := make([]byte, BufferSize)
	_, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("There is an err reading from connection", err.Error())
		return
	}
	fmt.Println("command received: " + string(buffer))
	fmt.Println("getting a file")
	getFileFromClient(connection)
}

func main() {
	addr := net.UDPAddr{
		Port: 12345,
		IP:   net.ParseIP("127.0.0.1"),
	}
	ser, _ := net.ListenUDP("udp", &addr)
	connectionHandler(ser)
}
