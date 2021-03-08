package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
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
	var Addr string
	var Port string
	flag.StringVar(&Addr, "l", "0.0.0.0", "Specify IP address. Default is 0.0.0.0")
	flag.StringVar(&Port, "p", "12345", "Specify port. Default is 12345")

	flag.Parse()

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
	fmt.Println("Listening on ", Addr, Port)
	connectionHandler(ser)
}
