package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	data := make([]byte, 3)
	n, err := conn.Read(data)
	if err != nil {
		fmt.Println(err)
		return
	}
        if string(data[:n]) == "OK\n" {
	fmt.Print(string(data[:n]))
	} else {
		fmt.Println("wrong answer")
	}
}

