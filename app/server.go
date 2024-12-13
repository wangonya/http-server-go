package main

import (
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
)

func extractEchoString(content string) string {
	re := regexp.MustCompile(`GET /echo/(.*) `)
	match := re.FindStringSubmatch(content)
	return match[1]
}

func extractUserAgentString(content string) string {
	re := regexp.MustCompile(`GET /user-agent HTTP/1.1\r\nHost: localhost:4221\r\nUser-Agent: (.*)\r\n\r\n`)
	match := re.FindStringSubmatch(content)
	return match[1]
}

func handleConnection(conn net.Conn) {
	req := make([]byte, 1024)
	conn.Read(req)
	content := string(req)

	switch {
	case strings.HasPrefix(content, "GET / HTTP/1.1"):
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	case strings.HasPrefix(content, "GET /echo/"):
		echoString := extractEchoString(content)
		response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(echoString), echoString)
		conn.Write([]byte(response))
	case strings.HasPrefix(content, "GET /user-agent"):
		userAgentString := extractUserAgentString(content)
		response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(userAgentString), userAgentString)
		conn.Write([]byte(response))
	default:
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		conn.Close()
	}
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		fmt.Println(err)
		os.Exit(1)
	}
	defer l.Close()

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	handleConnection(conn)
}
