package commands

import (
	"fmt"
	"net"
	"sync"

	"github.com/codecrafters-io/redis-starter-go/parser"
)

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
	lrange = "lrange"
)

type Client struct {
	cmd  parser.Command
	conn net.Conn
	data *sync.Map
}

func NewClient(conn net.Conn) Client {
	return Client{
		conn: conn,
		data: &sync.Map{},
	}
}

func (cl *Client) HandleCommand(buffer []byte) {
	cl.cmd = parser.Parse(buffer)
	fmt.Println(cl.cmd.Command, cl.cmd.Parameters)
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
		if len(cl.cmd.Parameters) == 1 {
			cl.handleLPop()
		} else {
			cl.handleLPopMulitpleEle()
		}
	case blpop:
		cl.handleblpop()
	default:
		_, err := cl.conn.Write([]byte("-ERROR command\r\n"))
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}
