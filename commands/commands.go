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
	// case rpush:
	// 	cl.handleRPush(cmds, data, conn)
	// case lrange:
	// 	cl.handleLRange(cmds, data, conn)
	// // case lpush:
	// 	cl.handleLPush(cmds, data, conn)
	// case llen:
	// 	cl.handleLLen(cmds, data, conn)
	// case lpop:
	// 	if len(cmds) == 2 {
	// 		cl.handleLPop(cmds, data, conn)
	// 	} else {
	// 		cl.handleLPopMulitpleEle(cmds, data, conn)
	// 	}
	// case blpop:
	// 	cl.handleblpop(cmds, data, conn)
	default:
		_, err := cl.conn.Write([]byte("-ERROR command\r\n"))
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}
