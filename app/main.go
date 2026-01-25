package main

import (
	"fmt"
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/Datastore"
	"github.com/codecrafters-io/redis-starter-go/commands"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Listen for incoming TCP connections on all network interfaces on port 6379.
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	fmt.Println("waiting for connection")

	// Create a new in-memory datastore that will be shared across all connections.
	ds := datastore.NewDataStore()
	// The main server loop. It continuously accepts new connections.
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		// Handle each connection in a new goroutine to allow for concurrent clients.
		go handleConn(conn, ds)
	}

}

// handleConn handles a single client connection.
func handleConn(conn net.Conn, ds *datastore.Datastore) {
	// Ensure the connection is closed when the function returns.
	defer conn.Close()
	// Create a new client to handle commands for this connection.
	cl := commands.NewClient(conn, ds)

	// The command reading loop for a single connection.
	for {
		buffer := make([]byte, 128)

		// Read data from the connection.
		_, err := conn.Read(buffer)
		if err != nil {
			fmt.Println(err)
			conn.Write([]byte(err.Error()))
			break
		}
		// Handle each command in a new goroutine to allow for pipelining of commands.
		go cl.HandleCommand(buffer)
	}
}
