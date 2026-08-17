package main

import (
	"fmt"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	defer listener.Close()

	fmt.Println("Run server")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection failed:", err)
			continue
		}

		go func(conn net.Conn) {
			defer conn.Close()

			_, err := conn.Write([]byte("OK\n"))
			if err != nil {
				fmt.Println("Write failed:", err)
				return
			}
		}(conn)
	}
}
