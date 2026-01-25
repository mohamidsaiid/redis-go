package commands

import (
	"fmt"
)

// handleEcho replies to the client with the first argument of the command.
// This is used for testing the connection and for debugging.
func (cl *Client) handleEcho() {
	if len(cl.cmd.Parameters) == 1 {
		arg := cl.cmd.Parameters[0]
		response := fmt.Sprintf("$%d\r\n%s\r\n", len(arg), arg)
		_, err := cl.conn.Write([]byte(response))
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		_, err := cl.conn.Write([]byte("+ERROR\r\n"))
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

// handlePing replies to the client with "PONG".
// This is used to test if the server is alive.
func (cl *Client) handlePing() {
	if len(cl.cmd.Parameters) == 0 {
		_, err := cl.conn.Write([]byte("+PONG\r\n"))
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		_, err := cl.conn.Write([]byte("+ERROR\r\n"))
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
