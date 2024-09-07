package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

const (
	DEFAULT_ADDRESS = "localhost"
	DEFAULT_PORT    = "8080"
	RETRY_DURATION  = 1 * time.Minute // Try to reconnect for 1 minute before exiting
	RETRY_INTERVAL  = 5 * time.Second // Try to reconnect every 5 seconds
)

func main() {
	var conn net.Conn
	var err error

	conn, err = connectWithRetry()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to the server")

	// Listen to broadcasted messages
	go listen(conn)

	reader := bufio.NewReader(os.Stdin)
	// Enter messages
	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		_, err := conn.Write([]byte(input + "\n"))
		if err != nil {
			fmt.Println("Connection lost. Attempting to reconnect...")
			conn, err = connectWithRetry()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println("Connected to the server!")
			continue
		}
	}
}

func listen(conn net.Conn) {
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Connection lost. Attempting to reconnect...")
			conn, err = connectWithRetry()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println("Connected to the server!")
			continue
		}

		fmt.Print("\r" + message + "> ")
	}
}

func connectWithRetry() (net.Conn, error) {
	startTime := time.Now()

	address := os.Getenv("SERVER_ADDRESS")
	if address == "" {
		address = DEFAULT_ADDRESS
	}
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = DEFAULT_PORT
	}

	for {
		conn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", address, port))
		if err == nil {
			return conn, nil
		}

		if time.Since(startTime) > RETRY_DURATION {
			return nil, fmt.Errorf("Failed to reconnect: timeout after 1 minute")
		}

		fmt.Println("Failed to connect, retrying in 5 seconds")
		time.Sleep(RETRY_INTERVAL)
	}
}
