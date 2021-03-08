package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

const BufferSize = 1024

func sendFileToServer(fileName string, connection net.Conn) {
	var currentByte int64 = 0
	fmt.Println("send to client")
	fileBuffer := make([]byte, BufferSize)
	var err error
	file, _ := os.Open(strings.TrimSpace(fileName))
	_, _ = connection.Write([]byte(fileName))
	for err == nil || err != io.EOF {
		n := 0
		n, err = file.ReadAt(fileBuffer, currentByte)
		currentByte += BufferSize
		_, _ = connection.Write(fileBuffer[:n])
	}
	_ = file.Close()
	_ = connection.Close()
}

func main() {
	conn, _ := net.Dial("udp", "127.0.0.1:12345")
	sendFileToServer("send/awsm.txt", conn)
}
