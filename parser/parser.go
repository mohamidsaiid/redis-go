// Package parser is responsible for parsing the Redis Serialization Protocol (RESP).
package parser

import (
	"bytes"
)

// Command represents a parsed Redis command, including the command itself and its parameters.
type Command struct {
	Command    string
	Parameters Parameters
}

// Parameters is a slice of strings representing the parameters of a Redis command.
type Parameters []string

// Parse takes a raw byte buffer from the client and parses it into a Command.
// This is the main entry point for the parser.
func Parse(buffer []byte) Command {
	// The RESP protocol uses "\r\n" as a separator.
	buf := bytes.Split(buffer, []byte("\r\n"))
	// The last element is always empty, so we remove it.
	buf = buf[:len(buf)-1]
	// The actual parsing is done in the `parse` function.
	newBuf := parse(buf)
	// The first element is the command, and the rest are parameters.
	return Command{
		Command:    string(newBuf[0]),
		Parameters: newBuf[1:],
	}
}

// parse is a helper function that takes a slice of byte slices (lines) and extracts the command and parameters.
// It filters out the RESP type prefixes (e.g., '$', ':', '*') to get the raw string values.
func parse(buf [][]byte) Parameters {
	parameters := make(Parameters, 0, 10)
	for _, val := range buf {
		// All commands and parameters are case-insensitive, so we convert them to lowercase.
		val = bytes.ToLower(val)
		// RESP uses prefixes to indicate the type of data. We are interested in the actual values,
		// so we skip the lines that contain these prefixes.
		if val[0] != '$' && val[0] != ':' && val[0] != '*' {
			parameters = append(parameters, string(val))
		}
	}
	return parameters
}
