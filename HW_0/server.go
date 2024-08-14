package main

import (
	"fmt"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":4545")

	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()
	fmt.Println("Server is listening...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			conn.Close()
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	_, err := conn.Write([]byte("ОК"))
	if err != nil {
		fmt.Println("Write error:", err)
	}

	fmt.Println("Connection handled, sent 'ОК'")

}
