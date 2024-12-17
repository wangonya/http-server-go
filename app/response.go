package main

import "fmt"

const (
	OK       = "HTTP/1.1 200 OK\r\n"
	NotFound = "HTTP/1.1 404 Not Found\r\n"
	Created  = "HTTP/1.1 201 Created\r\n"
)

// e.g HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 5\r\n\r\nhello
type HTTPResponse struct {
	StatusLine string

	ContentType   string
	ContentLength int

	Body string
}

func (r HTTPResponse) String() string {
	response := r.StatusLine

	if len(r.ContentType) != 0 {
		response += fmt.Sprintf("Content-Type: %s\r\n", r.ContentType)
	}

	if r.ContentLength != 0 {
		response += fmt.Sprintf("Content-Length: %d\r\n", r.ContentLength)
	}

	response += "\r\n"

	if len(r.Body) != 0 {
		response += r.Body
	}

	return response
}
