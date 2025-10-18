package commands

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

type element struct {
	key   string
	value string
}

var newPushedElement = make(chan element)

func (cl *Client) handleRPush() {
	key := cl.cmd.Parameters[0]
	list, ok := cl.data.Load(key)
	if !ok {
		list = make([]string, 0, 10)
	}

	switch l := list.(type) {
	case []string:
		for i := 1; i < len(cl.cmd.Parameters); i++ {
			l = append(l, cl.cmd.Parameters[i])
		}
		cl.data.Store(key, l)
		go func() {
			newPushedElement <- element{
				key:   key,
				value: l[0],
			}
		}()
		res := fmt.Sprint(":", len(l), "\r\n")
		_, err := cl.conn.Write([]byte(res))
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func (cl *Client) handleLPush() {
	key := cl.cmd.Parameters[0]
	list, ok := cl.data.Load(key)
	if !ok {
		list = make([]string, 0, 10)
	}
	switch l := list.(type) {
	case []string:
		newList := make([]string, 0, 10)

		for i := 1; i < len(cl.cmd.Parameters); i++ {
			newList = append(newList, cl.cmd.Parameters[i])
		}

		slices.Reverse(newList)
		newList = append(newList, l...)
		cl.data.Store(key, newList)

		go func() {
			newPushedElement <- element{
				key:   key,
				value: newList[0],
			}
		}()

		res := fmt.Sprint(":", len(newList), "\r\n")
		_, err := cl.conn.Write([]byte(res))
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
func (cl *Client) handleLRange() {
	if len(cl.cmd.Parameters) < 3 {
		_, err := cl.conn.Write([]byte("+SERROR\r\n"))
		if err != nil {
			fmt.Println(err)
			return
		}
		return
	}
	key, start, end := cl.cmd.Parameters[0], cl.cmd.Parameters[1], cl.cmd.Parameters[2]
	startIdx, err := strconv.Atoi(start)
	endIdx, err := strconv.Atoi(end)
	if err != nil {
		_, err := cl.conn.Write([]byte("+NERROR\r\n"))
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(err)
		return
	}
	list, ok := cl.data.Load(key)
	if !ok {
		_, err := cl.conn.Write([]byte("*0\r\n"))
		if err != nil {
			fmt.Println(err)
			return
		}
		return
	}
	arr := lRangeBuilder(list.([]string), startIdx, endIdx)
	_, err = cl.conn.Write(arr)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func lRangeBuilder(list []string, start, end int) []byte {

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

func (cl *Client) handleLLen() {
	key := cl.cmd.Parameters[0]
	list, ok := cl.data.Load(key)
	if !ok {
		cl.conn.Write([]byte(":0\r\n"))
		return
	}
	res := fmt.Sprintf(":%d\r\n", len(list.([]string)))
	cl.conn.Write([]byte(res))
}

func (cl *Client) handleLPop() {
	key := cl.cmd.Parameters[0]
	list, ok := cl.data.Load(key)
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
		cl.data.Store(key, v[1:])
		res := fmt.Sprintf("$%d\r\n%s\r\n", len(ele), ele)
		cl.conn.Write([]byte(res))
	}

}

func (cl *Client) handleLPopMulitpleEle() {
	key := cl.cmd.Parameters[0]
	list, ok := cl.data.Load(key)
	if !ok {
		cl.conn.Write([]byte("$-1\r\n"))
		return
	}
	number, err := strconv.Atoi(cl.cmd.Parameters[1])
	if err != nil {
		cl.conn.Write([]byte("+ERROR\r\n"))
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
		cl.data.Store(key, v)
		cl.conn.Write([]byte(res.String()))
		return
	}
}

func (cl *Client) handleblpop() {

}
