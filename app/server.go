package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
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

	return &Request{
		Method:      method,
		URL:         url,
		HTTPVersion: httpVersion,
		UserAgent:   userAgent}, nil
}

func routeRequest(request *Request, connection net.Conn, directory string) {
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

	case strings.HasPrefix(request.URL, "/files/"):
		path := strings.TrimPrefix(request.URL, "/files/")
		fullPath := filepath.Join(directory, path)
		handleFileRequest(connection, request, fullPath)

	default:
		handleNotFound(connection, request.HTTPVersion)
	}
}

func handleFileRequest(connection net.Conn, request *Request, path string) {

	info, err := os.Stat(path)
	if err != nil {
		response := fmt.Sprintf("%s 404 Not Found\r\n\r\n", request.HTTPVersion)
		connection.Write([]byte(response))
		return
	}

	if info.IsDir() {
		response := fmt.Sprintf("%s 404 Not Found\r\n\r\n", request.HTTPVersion)
		connection.Write([]byte(response))
		return
	}

	file, err := os.Open(path)
	if err != nil {
		response := fmt.Sprintf("%s 404 Not Found\r\n\r\n", request.HTTPVersion)
		connection.Write([]byte(response))
		return
	}

	defer file.Close()

	size, _ := file.Seek(0, io.SeekEnd)
	file.Seek(0, io.SeekStart)
	data := make([]byte, size+10)
	_, _ = file.Read(data)
	content := string(data)
	fmt.Println(content)
	response := fmt.Sprintf("%s 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", request.HTTPVersion, size, content)
	connection.Write([]byte(response))
}

func handleRequest(connection net.Conn, directory string) {

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

	routeRequest(request, connection, directory)

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

	var directory string
	flag.StringVar(&directory, "directory", ".", "directory from which to serve files")
	flag.Parse()

	info, err := os.Stat(directory)
	if err != nil {
		fmt.Printf("Failed to check directory path: %v\n", err)
		os.Exit(1)
	}

	if !info.IsDir() {
		fmt.Printf("Invalid directory path %s\n", directory)
		os.Exit(1)
	}

	fmt.Println("Server is running on port 4221...")
	listener, err := net.Listen("tcp", "0.0.0.0:4221")

	if err != nil {
		fmt.Println("Can not bind the port for tcp server")
		return
	}

	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection:", err)
			continue
		}
		go handleRequest(connection, directory)
	}
}
