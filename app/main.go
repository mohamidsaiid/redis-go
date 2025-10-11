package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

const (
	echo = "echo"
	ping = "ping"
	get  = "get"
	set  = "set"
	sep  = "\r\n"
)

type commands []string

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

	var data = &sync.Map{}
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
			if len(cmds) > 0 {
				switch cmds[0] {
				case ping:
					handlePing(cmds, conn)
				case echo:
					handleEcho(cmds, conn)
				case set:
					handleSet(cmds, data, conn)
				case get:
					handleGet(cmds, data, conn)
				default:
					conn.Write([]byte("+ERROR\r\n"))
					continue
				}
			} else {
				conn.Write([]byte("+ERROR\r\n"))
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
		if (val[0] >= 'a' && val[0] <= 'z') || (val[0] >= '0' && val[0] <= '9') {
			cmds = append(cmds, string(val))
		}
	}
	return cmds
}

func handlePing(cmds commands, conn net.Conn) {
	if len(cmds) == 1 {
		conn.Write([]byte("+PONG\r\n"))
	} else {
		conn.Write([]byte("+ERROR\r\n"))
	}
}

func handleEcho(cmds commands, conn net.Conn) {
	if len(cmds) > 1 && len(cmds) < 3 {
		arg := cmds[1]
		response := fmt.Sprintf("$%d\r\n%s\r\n", len(arg), arg)
		conn.Write([]byte(response))
	} else {
		conn.Write([]byte("+ERROR\r\n"))
	}
}

func handleSet(cmds commands, data *sync.Map, conn net.Conn) {
	fmt.Println(len(cmds))
	if len(cmds) == 3 {
		key, value := cmds[1], cmds[2]
		data.Store(key, value)
		conn.Write([]byte("+OK\r\n"))
	} else if len(cmds) == 5 {
		key, value, expiry, t := cmds[1], cmds[2], cmds[3], cmds[4]

		data.Store(key, value)

		switch expiry {
		case "ex":
			t += "s"
		case "px":
			t += "ms"
		}
		newTime, err := time.ParseDuration(t)
		if err != nil {
			conn.Write([]byte("Error invalid time"))
			return
		}
		go func() {
			time.Sleep(newTime)
			data.Delete(key)
		}()
		conn.Write([]byte("+OK\r\n"))
	} else {
		conn.Write([]byte("+ERROR\r\n"))
	}
}

func handleGet(cmds commands, data *sync.Map, conn net.Conn) {
	if len(cmds) == 2 {
		if val, ok := data.Load(cmds[1]); ok {
			response := fmt.Sprintf("+%s\r\n", val)
			conn.Write([]byte(response))
		} else {
			response := "$-1\r\n"
			conn.Write([]byte(response))
		}
	}
}
