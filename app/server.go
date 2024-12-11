package main

import (
	"fmt"
	"net"
	"os"
)

func Response() string {
	return "HTTP/1.1 200 OK\r\n\r\n"
}

func handleConnection(c net.Conn) {
	c.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	c, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	handleConnection(c)
}
