package main

import (
	"fmt"
	"net"
	"strings"
)

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
	requestLine := lines[0]
	requestParts := strings.Split(requestLine, " ")

	if len(requestParts) < 3 {
		fmt.Println("Invalid request structure", err)
		return
	}
	method, url, httpVersion := requestParts[0], requestParts[1], requestParts[2]

	if method == "GET" {
		s := strings.Split(url, "/")
		if len(s) == 2 && s[1] == "" { 
			handleRoot(connection , httpVersion)
		} else if len(s) >= 2 && s[1] == "echo" {
			var body string
			if len(s) > 2 {
				body = s[2]
			} 
			handleGET(connection , httpVersion , body) 
		} else {
			handleNotFound(connection, httpVersion)
		}
		
	} else {
		handleNotFound(connection, httpVersion)
	}

}
func handleRoot(connection net.Conn , httpVersion string) {
	response := fmt.Sprintf("%s 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 0\r\n\r\n", httpVersion)
	connection.Write([]byte(response))
}
func handleGET(connection net.Conn, httpVersion string , body string) {
	response := fmt.Sprintf("%s 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", httpVersion , len(body) , body)
	connection.Write([]byte(response))
}

func handleNotFound(connection net.Conn, httpVersion string) {
	response := fmt.Sprintf("%s 404 Not Found\r\n\r\n", httpVersion)
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
