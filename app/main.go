package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	echo   = "echo"
	ping   = "ping"
	get    = "get"
	set    = "set"
	rpush  = "rpush"
	rpop   = "rpop"
	lrange = "lrange"
	sep    = "\r\n"
)

type parameters []string

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
				case rpush:
					handleRPush(cmds, data, conn)
				case lrange:
					handleLRange(cmds, data, conn)
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

func parse(buf [][]byte) parameters {
	cmds := make(parameters, 0, 10)
	for _, val := range buf {
		val = bytes.ToLower(val)
		if val[0] != '$' && val[0] != ':' && val[0] != '*' {
			cmds = append(cmds, string(val))
		}
	}
	return cmds
}

func handlePing(cmds parameters, conn net.Conn) {
	if len(cmds) == 1 {
		conn.Write([]byte("+PONG\r\n"))
	} else {
		conn.Write([]byte("+ERROR\r\n"))
	}
}

func handleEcho(cmds parameters, conn net.Conn) {
	if len(cmds) > 1 && len(cmds) < 3 {
		arg := cmds[1]
		response := fmt.Sprintf("$%d\r\n%s\r\n", len(arg), arg)
		conn.Write([]byte(response))
	} else {
		conn.Write([]byte("+ERROR\r\n"))
	}
}

func handleSet(cmds parameters, data *sync.Map, conn net.Conn) {
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

func handleGet(cmds parameters, data *sync.Map, conn net.Conn) {
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

func handleRPush(cmds parameters, data *sync.Map, conn net.Conn) {
	key := cmds[1]
	list, ok := data.Load(key)
	if !ok {
		list = make([]string, 0, 10)
	}
	for i := 2; i < len(cmds); i++ {
		list = append(list.([]string), cmds[i])
	}
	data.Store(key, list)
	res := fmt.Sprint(":", len(list.([]string)), "\r\n")
	conn.Write([]byte(res))
}

func handleLRange(cmds parameters, data *sync.Map, conn net.Conn) {
	if len(cmds) < 4 {
		conn.Write([]byte("+SERROR\r\n"))
		return
	}
	key, start, end := cmds[1], cmds[2], cmds[3]
	startIdx, err := strconv.Atoi(start)
	endIdx, err := strconv.Atoi(end)
	if err != nil {
		conn.Write([]byte("+NERROR\r\n"))
		return
	}
	list, ok := data.Load(key)
	if !ok {
		conn.Write([]byte("*0\r\n"))
		return
	}
	arr := LRangeBuilder(list.([]string), startIdx, endIdx)
	conn.Write(arr)
}

func LRangeBuilder(list []string, start, end int) []byte {

	str := strings.Builder{}

	if len(list) == 0 {
		str.WriteString("*0\r\n")
		return []byte(str.String())
	}

	if end >= len(list) {
		end = len(list) - 1
	}
	if end*-1 > len(list) {
		end = 0
	} else if end < 0 {
		end += len(list)
	}
	if start*-1 > len(list) {
		start = 0
	} else if start < 0 {
		start += len(list)
	}
	str.WriteString(fmt.Sprintf("*%d\r\n", end-start+1))
	for i := start; i <= end; i++ {
		s := fmt.Sprintf("$%d\r\n%s\r\n", len(list[i]), list[i])
		str.WriteString(s)
	}

	fmt.Println(start, end)

	fmt.Println(str.String())
	return []byte(str.String())
}
