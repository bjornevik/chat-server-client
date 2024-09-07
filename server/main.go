package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

// TODO: Better handling of env variables
const (
	DEFAULT_ADDRESS = "localhost"
	DEFAULT_PORT    = "8080"
)

var (
	connections = make(map[net.Conn]bool)
	mu          sync.Mutex
	logFile     *os.File
)

func main() {
	address := os.Getenv("SERVER_ADDRESS")
	if address == "" {
		address = DEFAULT_ADDRESS
	}
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = DEFAULT_PORT
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", address, port))
	if err != nil {
		fmt.Println("Error starting server:", err.Error())
		return
	}
	defer listener.Close()
	fmt.Printf("Listening on %s:%s\n", address, port)

	logFileName := fmt.Sprintf("server_log.txt")
	logFile, err = os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Error creating log file: ", err)
		return
	}
	defer logFile.Close()

	logFile.WriteString("Server started at: " + time.Now().Format(time.RFC3339) + "\n")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}

		mu.Lock()
		connections[conn] = true
		mu.Unlock()

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	logMessage(fmt.Sprintf("[%s] A client connected", time.Now().Format(time.RFC3339)))
	// Remove from connections and close connection once this function stops running
	defer func() {
		mu.Lock()
		delete(connections, conn)
		mu.Unlock()
		conn.Close()
		logMessage(fmt.Sprintf("[%s] A client disconnected", time.Now().Format(time.RFC3339)))
	}()

	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		broadcastMessage(message, conn)
	}
}

func broadcastMessage(message string, sender net.Conn) {
	mu.Lock()
	defer mu.Unlock()

	msg := strings.TrimSpace(message)
	msg = fmt.Sprintf("[%s] %s", time.Now().Format(time.RFC3339), msg)

	// Only log messages when broadcasting
	logMessage(msg)

	for conn := range connections {
		// Don't send the message back to the sender
		if conn != sender {
			_, err := conn.Write([]byte(msg + "\n"))
			if err != nil {
				fmt.Println("Error broadcasting to connection:", err)
				conn.Close()
				delete(connections, conn)
			}
		}
	}
}

func logMessage(msg string) {
	fmt.Println(msg)
	if logFile != nil {
		_, err := logFile.WriteString(msg + "\n")
		if err != nil {
			fmt.Println("Error writing to log file:", err)
		}
	}
}
