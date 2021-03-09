package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"path"
)

const BufferSize = 1024

func sendFileToServer(filePath string, connection net.Conn) {
	var currentByte int64 = 0
	fmt.Println("sending")
	fileBuffer := make([]byte, BufferSize)
	var err error
	dir, fileName := path.Split(path.Clean(filePath))
	file, _ := os.Open(path.Join(dir, fileName))
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
	var Addr = ""
	var filePath = ""
	flag.StringVar(&Addr, "d", "", "Specify destination. eg. '127.0.0.1:12345'")
	flag.StringVar(&filePath, "f", "", "Specify file path. eg. '"+path.Join("send", "awsm.txt")+"'")
	flag.Parse()
	if (Addr == "") || filePath == "" {
		fmt.Println("Use --help")
		os.Exit(1)
	}
	conn, _ := net.Dial("udp", Addr)
	sendFileToServer(filePath, conn)
}
