package commands

import (
	"fmt"
	"time"
)

// handleSet sets a key-value pair in the datastore.
// It can also handle optional expiry times using the EX (seconds) and PX (milliseconds) options.
func (cl *Client) handleSet() {
	// Simple SET key value
	if len(cl.cmd.Parameters) == 2 {
		key, value := cl.cmd.Parameters[0], cl.cmd.Parameters[1]
		cl.ds.Data.Store(key, value)
		_, err := cl.conn.Write([]byte("+OK\r\n"))
		if err != nil {
			fmt.Println(err)
			return
		}
	} else if len(cl.cmd.Parameters) == 4 { // SET key value EX seconds or PX milliseconds
		key, value, expiry, t := cl.cmd.Parameters[0], cl.cmd.Parameters[1], cl.cmd.Parameters[2], cl.cmd.Parameters[3]

		cl.ds.Data.Store(key, value)

		// Determine the expiry duration based on the provided option.
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
		// Delete the key after the specified duration in a new goroutine.
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

// handleGet retrieves the value for a given key from the datastore.
func (cl *Client) handleGet() {
	if len(cl.cmd.Parameters) == 1 {
		// Load the value from the datastore.
		if val, ok := cl.ds.Data.Load(cl.cmd.Parameters[0]); ok {
			response := fmt.Sprintf("+%s\r\n", val)
			_, err := cl.conn.Write([]byte(response))
			if err != nil {
				fmt.Println(err)
				return
			}
		} else {
			// If the key does not exist, return a null bulk string.
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
