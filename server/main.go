package main

import (
	"fmt"
	"net"
)

// TODO: Better handling of env variables
const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "8081"
	SERVER_TYPE = "tcp"
)

func main() {
	listener, err := net.Listen(SERVER_TYPE, fmt.Sprintf("%v:%v", SERVER_HOST, SERVER_PORT))
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}

	defer listener.Close()

	fmt.Printf("Listening on %v:%v", SERVER_HOST, SERVER_PORT)

	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}
		go processClient(connection)
	}
}

func processClient(conn net.Conn) {
	buffer := make([]byte, 1024)
	mLen, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading: ", err.Error())
		return
	}

	// TODO: Write to log file
	// TODO: Output message to all clients

	fmt.Println("Received: ", string(buffer[:mLen]))
	_, err = conn.Write([]byte("Message received: " + string(buffer[:mLen])))
	conn.Close() // TODO: Keep connection open
}
