// Package commands implements the handlers for the supported Redis commands.
package commands

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/Datastore"
	"github.com/codecrafters-io/redis-starter-go/parser"
)

// Supported Redis commands
const (
	echo   = "echo"
	ping   = "ping"
	get    = "get"
	set    = "set"
	rpush  = "rpush"
	lpush  = "lpush"
	rpop   = "rpop"
	llen   = "llen"
	lpop   = "lpop"
	blpop  = "blpop"
	Type   = "type"
	lrange = "lrange"
)

// Client represents a connected client and holds the state for handling commands.
type Client struct {
	cmd  parser.Command // The parsed command from the client.
	conn net.Conn       // The network connection to the client.
	ds   *datastore.Datastore // The shared datastore.
}

// NewClient creates a new Client for a given connection and datastore.
func NewClient(conn net.Conn, ds *datastore.Datastore) Client {
	return Client{
		conn: conn,
		ds:   ds,
	}
}

// HandleCommand parses the command from the buffer and calls the appropriate handler.
func (cl *Client) HandleCommand(buffer []byte) {
	// Parse the raw command from the buffer.
	cl.cmd = parser.Parse(buffer)
	// A switch statement to dispatch the command to the correct handler.
	switch cl.cmd.Command {
	case ping:
		cl.handlePing()
	case echo:
		cl.handleEcho()
	case set:
		cl.handleSet()
	case get:
		cl.handleGet()
	case rpush:
		cl.handleRPush()
	case lrange:
		cl.handleLRange()
	case lpush:
		cl.handleLPush()
	case llen:
		cl.handleLLen()
	case lpop:
		// LPOP can take an optional count argument.
		if len(cl.cmd.Parameters) == 1 {
			cl.handleLPop()
		} else {
			cl.handleLPopMulitpleEle()
		}
	case blpop:
		cl.handleblpop()
	case Type:
		cl.handleType()
	default:
		// If the command is not recognized, return an error.
		_, err := cl.conn.Write([]byte("-ERROR command\r\n"))
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}
