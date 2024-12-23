package main

import (
	"fmt"
	"net"
	"strings"
)

func getUserAgent(lines []string) string {
	var userAgent string 
	for _ , line := range lines {
		if strings.HasPrefix(line , "User-Agent:") {
			parts := strings.SplitN(line , " " , 2)
			if len(parts) == 2 {
				userAgent = parts[1]
			}
			break 
		}
	}
	return userAgent 
}

func handleRequest(connection net.Conn) {

	defer connection.Close()

	buffer := make([]byte, 4096)
	//take the request as stream of bytes
	n, err := connection.Read(buffer)

	if err != nil {
		fmt.Println("error during read the http request!")
		return
	}

	//convert the stream of bytes to string
	httpRequest := string(buffer[:n])

	lines := strings.Split(httpRequest, "\r\n")

	if len(lines) < 1 {
		fmt.Println("so bad format of request")
		return
	}
	
	//Get userAgent
	userAgent := getUserAgent(lines)

	requestLine := lines[0]
	requestParts := strings.Split(requestLine, " ")

	if len(requestParts) < 3 {
		fmt.Println("Invalid request structure", err)
		return
	}

	method, url, httpVersion := requestParts[0], requestParts[1], requestParts[2]
	
	if method != "GET" {
		handleNotFound(connection , httpVersion)
		return 
	}

	if strings.HasPrefix(url , "/echo/") {
		body := strings.TrimPrefix(url , "/echo/")
		handleGETEcho(connection , httpVersion , body)
		return 
	}

	if strings.HasPrefix(url , "/user-agent") {
		handleGETUserAgent(connection , httpVersion , userAgent)
		return
	}

	if strings.HasPrefix(url , "/") {
		body := strings.TrimPrefix(url , "/")
		if body != "" {
			handleNotFound(connection , httpVersion)
		}
		handleGETRoot(connection , httpVersion)
	}

	handleNotFound(connection , httpVersion)
	
}

func handleGETRoot(connection net.Conn , httpVersion string) {
	response := fmt.Sprintf("%s 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 0\r\n\r\n", httpVersion)
	connection.Write([]byte(response))
}

func handleGETUserAgent(connection net.Conn , httpVersion string , userAgent string) {
	response := fmt.Sprintf("%s 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", httpVersion , len(userAgent) , userAgent)
	connection.Write([]byte(response))
}

func handleNotFound(connection net.Conn, httpVersion string) {
	response := fmt.Sprintf("%s 404 Not Found\r\n\r\n", httpVersion)
	connection.Write([]byte(response))
}

func handleGETEcho(connection net.Conn ,httpVersion string , body string) {
	response := fmt.Sprintf("%s 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", httpVersion , len(body) , body)
	connection.Write([]byte(response))
}

func main() {

	fmt.Println("Logs from your program will appear here!")

	listener, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Can not bind the port for tcp server")
		return
	}

	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println("cant accept this connection")
			return
		}
		handleRequest(connection)
	}
}