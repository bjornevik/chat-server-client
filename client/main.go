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
	conn, err := net.Dial(SERVER_TYPE, fmt.Sprintf("%v:%v", SERVER_HOST, SERVER_PORT))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	_, err = conn.Write(([]byte("Hello, world")))
	if err != nil {
		fmt.Println(err)
		return
	}
}
