package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

const BufferSize = 1024

func sendFileToServer(fileName string, connection net.Conn) {
	var currentByte int64 = 0
	fmt.Println("sending")
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
	fmt.Println("sent")
}

func main() {
	var Addr string
	flag.StringVar(&Addr, "d", "127.0.0.1:12345", "Specify destination. Default is 127.0.0.1")

	flag.Parse()

	conn, _ := net.Dial("udp", Addr)
	sendFileToServer("send/awsm.txt", conn)
}
