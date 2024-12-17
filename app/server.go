package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	MaxBufferSize = 1024

	Host = "0.0.0.0"
	Port = 4221
)

func main() {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", Host, Port))
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		fmt.Println(err)
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnection(conn)
	}
}

func handleConnection(c net.Conn) {
	defer c.Close()

	buffer := make([]byte, MaxBufferSize)
	bytes, err := c.Read(buffer)
	if err != nil {
		fmt.Println("Error reading request: ", err.Error())
		return
	}

	request, err := ParseRequest(buffer[:bytes])
	if err != nil {
		fmt.Println("Error parsing request: ", err.Error())
		return
	}

	response := handleRequest(request)
	c.Write([]byte(response.String()))
}

func handleRequest(r HTTPRequest) HTTPResponse {
	response := HTTPResponse{StatusLine: NotFound}

	switch {
	case r.URI == "/":
		response.StatusLine = OK
	case strings.Contains(r.URI, "/echo/"):
		echoString := strings.Split(r.URI, "/")[2]
		response.StatusLine = OK
		response.ContentType = "text/plain"
		response.ContentLength = len(echoString)
		response.Body = echoString
	case r.UserAgent != "":
		response.StatusLine = OK
		response.ContentType = "text/plain"
		response.ContentLength = len(r.UserAgent)
		response.Body = r.UserAgent
	case strings.Contains(r.URI, "/files/"):
		filePath := strings.Split(r.URI, "/")[2]

		if r.Method == "GET" {
			data, err := os.ReadFile(os.Args[2] + filePath)
			if err != nil {
				break
			}

			response.StatusLine = OK
			response.ContentType = "application/octet-stream"
			response.ContentLength = len(data)
			response.Body = string(data)
		} else {
			err := os.WriteFile(os.Args[2]+filePath, []byte(r.Body), 0644)
			if err != nil {
				break
			}

			response.StatusLine = Created
		}
	}

	return response
}
