package main

import (
	"fmt"
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/commands"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	fmt.Println("waiting for connection")

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConn(conn)
	}

}

func handleConn(conn net.Conn) {
	defer conn.Close()
	cl := commands.NewClient(conn)

	for {
		buffer := make([]byte, 128)

		_, err := conn.Read(buffer)
		if err != nil {
			fmt.Println(err)
			conn.Write([]byte(err.Error()))
			break
		}
		cl.HandleCommand(buffer)
	}
}
