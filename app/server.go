package main

import (
	"fmt"
	"net"
	"strings"
)

type Request struct {
	Method      string
	URL         string
	HTTPVersion string
	UserAgent   string
}

func getUserAgent(lines []string) string {
	var userAgent string
	for _, line := range lines {
		if strings.HasPrefix(line, "User-Agent:") {
			parts := strings.SplitN(line, " ", 2)
			if len(parts) == 2 {
				userAgent = parts[1]
			}
			break
		}
	}
	return userAgent
}

func parseRequest(rawRequest string) (*Request, error) {
	lines := strings.Split(rawRequest, "\r\n")
	if len(lines) < 1 {
		return nil, fmt.Errorf("bad request format")
	}

	requestLine := strings.Split(lines[0], " ")
	if len(requestLine) < 3 {
		return nil, fmt.Errorf("invalid request line")
	}

	method, url, httpVersion := requestLine[0], requestLine[1], requestLine[2]
	userAgent := getUserAgent(lines)

	return &Request{Method: method,
		URL:         url,
		HTTPVersion: httpVersion,
		UserAgent:   userAgent}, nil
}

func routeRequest(request *Request, connection net.Conn) {
	switch {
	case request.Method != "GET":
		handleNotFound(connection, request.HTTPVersion)

	case strings.HasPrefix(request.URL, "/echo/"):
		body := strings.TrimPrefix(request.URL, "/echo/")
		handleGET(connection, request.HTTPVersion, body)

	case request.URL == "/user-agent":
		handleGET(connection, request.HTTPVersion, request.UserAgent)

	case request.URL == "/":
		handleGETRoot(connection, request.HTTPVersion)

	default:
		handleNotFound(connection, request.HTTPVersion)
	}
}

func handleRequest(connection net.Conn) {

	defer connection.Close()

	buffer := make([]byte, 10240)
	n, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("error during read the http request!")
		return
	}

	rawRequest := string(buffer[:n])
	request, err := parseRequest(rawRequest)

	if err != nil {
		fmt.Println("Failed to parse request:", err)
		handleNotFound(connection, "HTTP/1.1")
		return
	}

	routeRequest(request, connection)

}

func handleGETRoot(connection net.Conn, httpVersion string) {
	response := fmt.Sprintf("%s 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 0\r\n\r\n", httpVersion)
	connection.Write([]byte(response))
}

func handleGET(connection net.Conn, httpVersion string, userAgent string) {
	response := fmt.Sprintf("%s 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", httpVersion, len(userAgent), userAgent)
	connection.Write([]byte(response))
}

func handleNotFound(connection net.Conn, httpVersion string) {
	response := fmt.Sprintf("%s 404 Not Found\r\n\r\n", httpVersion)
	connection.Write([]byte(response))
}

func main() {

	fmt.Println("Server is running on port 4221...")

	listener, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Can not bind the port for tcp server")
		return
	}

	defer listener.Close()

	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection:", err)
			continue
		}
		go handleRequest(connection)
	}
}
