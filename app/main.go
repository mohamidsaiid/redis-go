package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
)

const (
	echo = "echo"
	ping = "ping"
	get  = "get"
	set  = "set"
	sep  = "\r\n"
)

type commands []string

var data = map[string]string{}

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

	for {
		buffer := make([]byte, 128)

		_, err := conn.Read(buffer)
		if err != nil {
			fmt.Println(err)
			conn.Write([]byte(err.Error()))
			break
		}
		fmt.Println(string(buffer))
		buf := bytes.Split(buffer, []byte(sep))
		buf = buf[:len(buf)-1]
		if buf[0][0] == '*' {
			// this should be an array
			cmds := parse(buf)
			switch cmds[0] {
			case ping:
				if len(cmds) == 1 {
					conn.Write([]byte("+PONG\r\n"))
				} else {
					conn.Write([]byte("+ERROR\r\n"))
				}
			case echo:
				if len(cmds) > 1 && len(cmds) < 3 {
					arg := cmds[1]
					response := fmt.Sprintf("$%d\r\n%s\r\n", len(arg), arg)
					conn.Write([]byte(response))
				} else {
					conn.Write([]byte("+ERROR\r\n"))
				}
			case set:
				if len(cmds) == 3 {
					data[cmds[1]] = cmds[2]
					conn.Write([]byte("+OK\r\n"))
				} else {
					conn.Write([]byte("+ERROR\r\n"))
				}
			case get:
				if len(cmds) == 2 {
					if val, ok := data[cmds[1]]; ok {
						response := fmt.Sprintf("+%s\r\n", val)
						conn.Write([]byte(response))
					} else {
						response := "+(nil)\r\n"
						conn.Write([]byte(response))
					}
				}
			default:
				conn.Write([]byte("+ERROR\r\n"))
				continue
			}
		} else {
			conn.Write([]byte("+ERROR\r\n"))
		}

	}
}

func parse(buf [][]byte) commands {
	cmds := make(commands, 0, 10)
	for _, val := range buf {
		val = bytes.ToLower(val)
		if val[0] >= 'a' && val[0] <= 'z' {
			cmds = append(cmds, string(val))
		}
	}
	return cmds
}
