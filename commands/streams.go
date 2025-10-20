package commands

import (
	"log"
)

func (cl *Client) handleType() {
	key := cl.cmd.Parameters[0]
	data, ok := cl.ds.Data.Load(key)
	if !ok {
		_, err := cl.conn.Write([]byte("+none\r\n"))
		if err != nil {
			log.Println(err)
			return
		}
		return
	}

	switch data.(type) {
	case []string:
		_, err := cl.conn.Write([]byte("+list\r\n"))
		if err != nil {
			log.Println(err)
			return
		}
	case string:
		_, err := cl.conn.Write([]byte("+string\r\n"))
		if err != nil {
			log.Println(err)
			return
		}
	}
}
