package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

type Request struct {
	request_api_key     int16
	request_api_version int16
	client_id           *string
	tag_buffer          []int32
	message_size        int32
	correlation_id      int32
}

type responseMessage struct {
	messageSize   int32
	correlationId int32
	error_code    int16
}

func (r *responseMessage) Byte() []byte {
	buff := bytes.Buffer{}
	binary.Write(&buff, binary.BigEndian, r.messageSize)
	binary.Write(&buff, binary.BigEndian, r.correlationId)
	binary.Write(&buff, binary.BigEndian, r.error_code)

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

	message_size := buff[0:4]
	corelation_id := buff[8:12]

	response := responseMessage{
		messageSize:   int32(binary.BigEndian.Uint32(message_size)),
		correlationId: int32(binary.BigEndian.Uint32(corelation_id)),
		error_code:    35,
	}

	_, err = conn.Write(response.Byte())
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
