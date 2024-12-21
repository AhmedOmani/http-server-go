package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	// Start listening on the required port
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	fmt.Println("Server is listening on port 4221")

	for {
		// Accept an incoming connection
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}

		// Handle the connection
		handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Read the incoming request
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading from connection:", err.Error())
		return
	}
	fmt.Printf("Received request:\n%s\n", string(buffer[:n]))

	// Send a basic HTTP 200 OK response
	response := "HTTP/1.1 200 OK\r\n" +
		"Content-Length: 0\r\n" +
		"\r\n"

	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing response:", err.Error())
		return
	}
	fmt.Println("Response sent")
}
