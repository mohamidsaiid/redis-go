package parser

import (
	"bytes"
)

type Command struct {
	Command    string
	Parameters Parameters
}

type Parameters []string

func Parse(buffer []byte) Command {
	buf := bytes.Split(buffer, []byte("\r\n"))
	buf = buf[:len(buf)-1]
	newBuf := parse(buf)
	return Command{
		Command: string(newBuf[0]),
		Parameters: newBuf[1:],
	}
}

func parse(buf [][]byte) Parameters {
	parameters := make(Parameters, 0, 10)
	for _, val := range buf {
		val = bytes.ToLower(val)
		if val[0] != '$' && val[0] != ':' && val[0] != '*' {
			parameters = append(parameters, string(val))
		}
	}
	return parameters
}
