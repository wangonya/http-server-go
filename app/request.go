package main

import (
	"strconv"
	"strings"
)

// eg: GET /index.html HTTP/1.1\r\nHost: localhost:4221\r\nUser-Agent: curl/7.64.1\r\nAccept: */*\r\n\r\n
type HTTPRequest struct {
	Method   string
	URI      string
	Protocol string

	Host          string
	UserAgent     string
	Accept        string
	ContentType   string
	ContentLength int

	Body string
}

func ParseRequest(b []byte) (HTTPRequest, error) {
	bufferString := strings.Split(string(b), "\r\n")

	requestLine := bufferString[0]
	request := strings.Split(requestLine, " ")

	r := HTTPRequest{Method: request[0], URI: request[1], Protocol: request[2]}

	for line, v := range bufferString[1:] {
		if strings.Contains(v, "Host") {
			r.Host = strings.Split(v, " ")[1]
		}

		if strings.Contains(v, "User-Agent") {
			r.UserAgent = strings.Split(v, " ")[1]
		}

		if strings.Contains(v, "Accept") {
			r.Accept = strings.Split(v, " ")[1]
		}

		if strings.Contains(v, "Content-Type") {
			r.ContentType = strings.Split(v, " ")[1]
		}

		if strings.Contains(v, "Content-Length") {
			r.ContentLength, _ = strconv.Atoi(strings.Split(v, " ")[1])
		}

		if line == len(bufferString)-2 && !strings.HasSuffix(v, "\r\n") {
			r.Body = v
		}
	}

	return r, nil
}
