package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"path"
	"strconv"
	"strings"
)

const BufferSize = 1024

func getFileFromClient(connection net.Conn, fileName string) {
	var currentByte int64 = 0
	fileBuffer := make([]byte, BufferSize)
	var err error
	file, _ := os.Create(path.Join("receive", fileName))
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

	cleanedBuffer := bytes.Trim(buffer, "\x00")
	fileName := strings.TrimSpace(string(cleanedBuffer))
	getFileFromClient(connection, fileName)
}

func main() {
	var Addr = ""
	var Port = ""
	flag.StringVar(&Addr, "l", "", "Specify IP address. eg. '0.0.0.0'")
	flag.StringVar(&Port, "p", "", "Specify port. eg. '12345'")
	flag.Parse()
	if (Addr == "") || Port == "" {
		fmt.Println("Use --help")
		os.Exit(1)
	}
	PortNumber, err := strconv.Atoi(Port)
	if err != nil {
		fmt.Println("Use a valid port: ", err)
		os.Exit(1)
	}
	addr := net.UDPAddr{
		Port: PortNumber,
		IP:   net.ParseIP(Addr),
	}
	ser, _ := net.ListenUDP("udp", &addr)
	fmt.Println("Listening on ", Addr, ":", Port)
	connectionHandler(ser)
}
