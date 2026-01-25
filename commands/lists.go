package commands

import (
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"
	"time"
)

// element is used to pass data about a newly pushed element to the `newPushedElement` channel.
// This is used to unblock clients waiting on a BLPOP command.
type element struct {
	key   string // The key of the list where the element was pushed.
	value string // The value of the pushed element.
}

// newPushedElement is a channel used to signal `BLPOP` handlers that a new element has been pushed to a list.
var newPushedElement = make(chan element)

// handleRPush appends one or more elements to the end of a list.
func (cl *Client) handleRPush() {
	if len(cl.cmd.Parameters) < 2 {
		_, err := cl.conn.Write([]byte("- error you should provide key and elements\r\n"))
		if err != nil {
			log.Println(err)
			return
		}
		return
	}
	key := cl.cmd.Parameters[0]
	list, ok := cl.ds.Data.Load(key)
	if !ok {
		// If the list does not exist, create a new one.
		list = make([]string, 0, 10)
	}

	switch l := list.(type) {
	case []string:
		for i := 1; i < len(cl.cmd.Parameters); i++ {
			l = append(l, cl.cmd.Parameters[i])
		}
		cl.ds.Data.Store(key, l)
		// Signal any waiting BLPOP clients that a new element has been pushed.
		go func() {
			newPushedElement <- element{
				key:   key,
				value: l[0],
			}
		}()
		res := fmt.Sprint(":", len(l), "\r\n")
		_, err := cl.conn.Write([]byte(res))
		if err != nil {
			log.Println(err)
			return
		}
	}
}

// handleLPush prepends one or more elements to the beginning of a list.
func (cl *Client) handleLPush() {
	if len(cl.cmd.Parameters) < 2 {
		_, err := cl.conn.Write([]byte("- error you should provide key and elements\r\n"))
		if err != nil {
			log.Println(err)
			return
		}
		return
	}
	key := cl.cmd.Parameters[0]
	list, ok := cl.ds.Data.Load(key)
	if !ok {
		// If the list does not exist, create a new one.
		list = make([]string, 0, 10)
	}
	switch l := list.(type) {
	case []string:
		newList := make([]string, 0, 10)

		for i := 1; i < len(cl.cmd.Parameters); i++ {
			newList = append(newList, cl.cmd.Parameters[i])
		}

		// Reverse the new elements to maintain the correct order when prepending.
		slices.Reverse(newList)
		newList = append(newList, l...)
		cl.ds.Data.Store(key, newList)

		// Signal any waiting BLPOP clients that a new element has been pushed.
		go func() {
			newPushedElement <- element{
				key:   key,
				value: newList[0],
			}
		}()

		res := fmt.Sprint(":", len(newList), "\r\n")
		_, err := cl.conn.Write([]byte(res))
		if err != nil {
			log.Println(err)
			return
		}
	}
}

// handleLRange returns a range of elements from a list.
func (cl *Client) handleLRange() {
	if len(cl.cmd.Parameters) < 3 {
		_, err := cl.conn.Write([]byte("- SERROR\r\n"))
		if err != nil {
			log.Println(err)
			return
		}
		return
	}
	key, start, end := cl.cmd.Parameters[0], cl.cmd.Parameters[1], cl.cmd.Parameters[2]
	startIdx, err := strconv.Atoi(start)
	endIdx, err := strconv.Atoi(end)
	if err != nil {
		_, err := cl.conn.Write([]byte("- NERROR\r\n"))
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(err)
		return
	}
	list, ok := cl.ds.Data.Load(key)
	if !ok {
		_, err := cl.conn.Write([]byte("*0\r\n"))
		if err != nil {
			log.Println(err)
			return
		}
		return
	}
	// Build the response array.
	arr := lRangeBuilder(list.([]string), startIdx, endIdx)
	_, err = cl.conn.Write(arr)
	if err != nil {
		log.Println(err)
		return
	}
}

// lRangeBuilder is a helper function to build the RESP array for the LRANGE command.
// It handles positive and negative indices.
func lRangeBuilder(list []string, start, end int) []byte {

	str := strings.Builder{}

	if len(list) == 0 {
		str.WriteString("*0\r\n")
		return []byte(str.String())
	}

	// Adjust negative and out-of-range indices.
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

	return []byte(str.String())
}

// handleLLen returns the length of a list.
func (cl *Client) handleLLen() {
	key := cl.cmd.Parameters[0]
	list, ok := cl.ds.Data.Load(key)
	if !ok {
		cl.conn.Write([]byte(":0\r\n"))
		return
	}
	res := fmt.Sprintf(":%d\r\n", len(list.([]string)))
	cl.conn.Write([]byte(res))
}

// handleLPop removes and returns the first element of a list.
func (cl *Client) handleLPop() {
	key := cl.cmd.Parameters[0]
	list, ok := cl.ds.Data.Load(key)
	if !ok {
		cl.conn.Write([]byte("$-1\r\n"))
		return
	}

	switch v := list.(type) {
	case []string:
		if len(v) == 0 {
			cl.conn.Write([]byte("$-1\r\n"))
			return
		}
		ele := v[0]
		// If the list becomes empty after popping, delete the key.
		if len(v) == 1 {
			cl.ds.Data.Delete(key)
		} else {
			v = v[1:]
			cl.ds.Data.Store(key, v)

		}
		res := fmt.Sprintf("$%d\r\n%s\r\n", len(ele), ele)
		_, err := cl.conn.Write([]byte(res))
		if err != nil {
			log.Println(err)
			return
		}
	}
}

// handleLPopMulitpleEle removes and returns multiple elements from the beginning of a list.
func (cl *Client) handleLPopMulitpleEle() {
	key := cl.cmd.Parameters[0]
	list, ok := cl.ds.Data.Load(key)
	if !ok {
		_, err := cl.conn.Write([]byte("$-1\r\n"))
		if err != nil {
			log.Println(err)
		}
		return
	}
	number, err := strconv.Atoi(cl.cmd.Parameters[1])
	if err != nil {
		_, err := cl.conn.Write([]byte("+ERROR\r\n"))
		if err != nil {
			log.Println(err)
		}
		return
	}

	switch v := list.(type) {
	case []string:
		if number > len(v) {
			number = len(v) - 1
		}
		res := strings.Builder{}
		res.WriteString(fmt.Sprintf("*%d\r\n", number))
		for i := 0; i < number; i++ {
			ele := v[0]
			v = v[1:]
			res.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(ele), ele))
		}
		// If the list becomes empty after popping, delete the key.
		if len(v) == 0 {
			cl.ds.Data.Delete(key)
		} else {
			cl.ds.Data.Store(key, v)
		}
		_, err := cl.conn.Write([]byte(res.String()))
		if err != nil {
			log.Println(err)
		}
		return
	}
}

// handleblpop is a blocking list pop operation.
// It blocks the connection if the list is empty and waits for a new element to be pushed.
// It also supports a timeout.
func (cl *Client) handleblpop() {
	if len(cl.cmd.Parameters) < 2 {
		_, err := cl.conn.Write([]byte("- error you should provide the key and time or 0\r\n"))
		if err != nil {
			log.Println(err)
			return
		}
		return
	}
	key := cl.cmd.Parameters[0]
	// res is a helper function to format the response.
	res := func(k, e string) string {
		return fmt.Sprintf("*2\r\n$%d\r\n%s\r\n$%d\r\n%s\r\n", len(k), k, len(e), e)
	}
	// First, check if the list already has elements.
	if list, ok := cl.ds.Data.Load(key); ok {
		switch l := list.(type) {
		case []string:
			if len(l) > 0 {
				ele := l[0]
				if len(l) == 1 {
					cl.ds.Data.Delete(key)
				} else {
					l = l[1:]
					cl.ds.Data.Store(key, l)

				}
				_, err := cl.conn.Write([]byte(res(key, ele)))
				if err != nil {
					log.Println(err)
					return
				}
			}
		}
		return
	}

	// If the list is empty, enter the blocking loop.
	for {
		// If the timeout is 0, block indefinitely until an element is pushed.
		if cl.cmd.Parameters[1] == "0" {
			// Wait for a new element to be pushed to any list.
			ele := <-newPushedElement
			// If the pushed element is not for the key we are waiting for, continue waiting.
			if ele.key != key {
				continue
			}
			list, _ := cl.ds.Data.Load(key)
			switch l := list.(type) {
			case []string:
				if len(l) == 1 {
					cl.ds.Data.Delete(key)
				} else {
					l = l[1:]
					cl.ds.Data.Store(key, l)

				}
				_, err := cl.conn.Write([]byte(res(key, ele.value)))
				if err != nil {
					log.Println(err)
					return
				}
			}
			break
		} else {
			// If a timeout is specified, use a select statement to wait for either a new element or the timeout.
			t := cl.cmd.Parameters[1]
			t += "s"
			timeout, err := time.ParseDuration(t)
			if err != nil {
				_, err := cl.conn.Write([]byte("- error: Invalid time support\r\n"))
				if err != nil {
					log.Println(err)
					return
				}
			}

			select {
			// A new element was pushed.
			case ele := <-newPushedElement:
				if ele.key != key {
					continue
				}
				list, _ := cl.ds.Data.Load(key)
				switch l := list.(type) {
				case []string:
					if len(l) == 1 {
						cl.ds.Data.Delete(key)
					} else {
						l = l[1:]
						cl.ds.Data.Store(key, l)

					}
					_, err := cl.conn.Write([]byte(res(key, ele.value)))
					if err != nil {
						log.Println(err)
						return
					}
				}
			// The timeout was reached.
			case <-time.After(timeout):
				_, err := cl.conn.Write([]byte("*-1\r\n"))
				if err != nil {
					log.Println(err)
					return
				}
			}
			return
		}
	}
}
