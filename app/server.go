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

func extractFilePath(content string) string {
	re := regexp.MustCompile(`GET /files/(.*) `)
	match := re.FindStringSubmatch(content)
	return match[1]
}

func readFile(path string) ([]byte, int, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, 0, err
	}
	defer f.Close()

	data := make([]byte, 2048)
	count, err := f.Read(data)

	return data, count, err
}

func handleConnection(conn net.Conn) {
	req := make([]byte, 1024)
	conn.Read(req)
	content := string(req)

	defer conn.Close()

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
	case strings.HasPrefix(content, "GET /files"):
		filePath := extractFilePath(content)
		fmt.Println(os.Args[2] + filePath)
		f, i, err := readFile(os.Args[2] + filePath)

		if err != nil {
			conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
			return
		}
		response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", i, f)
		conn.Write([]byte(response))
	default:
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
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

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnection(conn)
	}
}
