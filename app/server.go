package main

import (
	"fmt"
	"net"
	"os"
	"bufio"
	"strings"
)

var _ = net.Listen
var _ = os.Exit

func handleRequest(connection net.Conn){
	
	defer connection.Close()

	reader := bufio.NewReader(connection)
	requestLine , err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading request: " , err) 
		return 
	}

	requestLine = strings.TrimSpace(requestLine)

	parts := strings.Split(requestLine , " ") 
	if len(parts) < 3 {
		fmt.Println("Invalid request line:" , requestLine)
		return 
	}

	

	method := parts[0]
	url := parts[1]
	httpVersion := parts[2]

	if method == "GET" && url == "/"  {
		
		handleGET(connection , httpVersion)
	} else {
		handleNotFound(connection , httpVersion)
	}

}

func handleGET(connection net.Conn , httpVersion string) {
	response := fmt.Sprintf("%s 200 OK\r\n\r\n" , httpVersion) 
	connection.Write([]byte(response))
	defer connection.Close()
}

func handleNotFound(connection net.Conn , httpVersion string) {
	response := fmt.Sprintf("%s 404 Not Found\r\n\r\n" , httpVersion)
	connection.Write([]byte(response))
	defer connection.Close()
} 

func main() {
	
	fmt.Println("Logs from your program will appear here!")
	
	//create an tcp server !
	listener , err := net.Listen("tcp" , "0.0.0.0:4221")
	if err != nil {
		fmt.Println("A7a mfesh connection: " , err) 
		return 
	}

	//wait for a client connection 
	for {
		connection , err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: " , err.Error())
			continue
		}
		go handleRequest(connection)
	}



}
