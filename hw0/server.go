package main

import (
	"fmt"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	defer listener.Close()

	fmt.Println("Сервер запущен")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Ошибка подключения:", err)
			continue
		}

		conn.Write([]byte("OK\n"))
		conn.Close()
	}
}
