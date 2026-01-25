package commands

import (
	"log"
)

// handleType returns the type of the data structure stored at the given key.
func (cl *Client) handleType() {
	key := cl.cmd.Parameters[0]
	data, ok := cl.ds.Data.Load(key)
	if !ok {
		// If the key does not exist, the type is "none".
		_, err := cl.conn.Write([]byte("+none\r\n"))
		if err != nil {
			log.Println(err)
			return
		}
		return
	}

	// Check the type of the data and respond accordingly.
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
