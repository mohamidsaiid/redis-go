package commands

import (
	"fmt"
	"time"
)

func (cl *Client) handleSet() {
	if len(cl.cmd.Parameters) == 2 {
		key, value := cl.cmd.Parameters[0], cl.cmd.Parameters[1]
		cl.ds.Data.Store(key, value)
		_, err := cl.conn.Write([]byte("+OK\r\n"))
		if err != nil {
			fmt.Println(err)
			return
		}
	} else if len(cl.cmd.Parameters) == 4 {
		key, value, expiry, t := cl.cmd.Parameters[0], cl.cmd.Parameters[1], cl.cmd.Parameters[2], cl.cmd.Parameters[3]

		cl.ds.Data.Store(key, value)

		switch expiry {
		case "ex":
			t += "s"
		case "px":
			t += "ms"
		}
		newTime, err := time.ParseDuration(t)
		if err != nil {
			_, err := cl.conn.Write([]byte("Error invalid time"))
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(err)
			return
		}
		go func() {
			time.Sleep(newTime)
			cl.ds.Data.Delete(key)
		}()
		_, err = cl.conn.Write([]byte("+OK\r\n"))
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

func (cl *Client) handleGet() {
	if len(cl.cmd.Parameters) == 1 {
		if val, ok := cl.ds.Data.Load(cl.cmd.Parameters[0]); ok {
			response := fmt.Sprintf("+%s\r\n", val)
			_, err := cl.conn.Write([]byte(response))
			if err != nil {
				fmt.Println(err)
				return
			}
		} else {
			response := "$-1\r\n"
			_, err := cl.conn.Write([]byte(response))
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	} else {
		_, err := cl.conn.Write([]byte("- Error get only recieve one key\r\n"))
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
