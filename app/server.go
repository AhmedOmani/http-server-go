package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

var _ = net.Listen
var _ = os.Exit

func handleRequest(connection net.Conn) {

	defer connection.Close()

	buffer := make([]byte, 4096)
	n, err := connection.Read(buffer)

	if err != nil {
		fmt.Println("Error reading request:", err)
		return
	}

	httpRequest := string(buffer[:n])

	lines := strings.Split(httpRequest, "\r\n")

	if len(lines) < 1 {
		fmt.Println("Invlid request structure", err)
		return
	}

	requestLine := lines[0]
	requestParts := strings.Split(requestLine, " ")

	if len(requestParts) < 3 {
		fmt.Println("Invalid request format", err)
		return
	}

	method, url, httpVersion := requestParts[0], requestParts[1], requestParts[2]

	if method == "GET" && url == "/" {
		handleGET(connection, httpVersion)
	} else {
		handleNotFound(connection, httpVersion)
	}

}

func handleGET(connection net.Conn, httpVersion string) {
	response := fmt.Sprintf("%s 200 OK\r\n\r\n", httpVersion)
	connection.Write([]byte(response))
}

func handleNotFound(connection net.Conn, httpVersion string) {
	response := fmt.Sprintf("%s 404 Not Found\r\n\r\n", httpVersion)
	connection.Write([]byte(response))
}

func main() {

	fmt.Println("Logs from your program will appear here!")

	//create an tcp server !
	listener, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("A7a mfesh connection: ", err)
		return
	}

	//wait for a client connection
	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}
		go handleRequest(connection)
	}

}
