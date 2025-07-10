package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

type responseHeader struct {
	correlationId int32
}

type responseMessage struct {
	messageSize int32
	header      responseHeader
}

func Byte(message_size, corelation_id []byte) []byte {
	buff := bytes.Buffer{}
	binary.Write(&buff, binary.BigEndian, message_size)
	binary.Write(&buff, binary.BigEndian, corelation_id)

	return buff.Bytes()
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	buff := make([]byte, 1024)
	_, err := conn.Read(buff)
	if err != nil {
		fmt.Println("Error reading from connection: ", err.Error())
		return
	}

	message_size := make([]byte, 4)

	corelation_id := buff[8:12]

	_, err = conn.Write(Byte(message_size, corelation_id))
	if err != nil {
		fmt.Println("Error writing into conn: ", err.Error())
	}
}

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:9092")
	if err != nil {
		fmt.Println("Failed to bind to port 9092")
		os.Exit(1)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConn(conn)
	}
}
