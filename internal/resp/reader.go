package respparser

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type RespParser struct {
	reader *bufio.Reader
}

func NewRespParser(rd io.Reader) *RespParser {
	return &RespParser{reader: bufio.NewReader(rd)}
}

// Implements RespParser's Read method
// Read will always read the initial data type and size
// then do byte shift twice over the first terminations
func (r *RespParser) Read() (Command, error) {
	// we do a first pass to get the type and size.
	// these values are needed to actually handle the incoming data
	t, size, err := r.initialParse()
	if err != nil {
		return Command{}, err
	}
	// Line endings
	r.reader.ReadByte()
	r.reader.ReadByte()

	// handle data types
	return r.handleDatatype(t, size)
}

// noramlizeString trims characters from incoming string.
// returning a "cleaner" instance of the same string
func noramlizeString(s string) string {
	return strings.Trim(s, "\x00")
}

// handleDatatype will switch over dataType and parse the present r.reader byte
// position as a "new line" This function should **NOT** call the reader itself and
// allow each datatype handler to call functions that shift byte positions on
// the reader
func (r *RespParser) handleDatatype(dataType DatatypePrimitive, size int64) (Command, error) {
	switch dataType {
	// When array is encountered a recursive loop is started.
	// This loop will run through all arrays and parse through each item in each array
	// in therory this should support nested array's however it hasn't been tested
	case ARRAY:
		value, err := r.parseArray(size)
		if err != nil {
			return Command{}, err
		}

		return Command{
			Kind:      DatatypeArray,
			Primitive: ARRAY,
			Value:     "",
			Args:      value,
		}, nil
	case SIMPLE_STRING:
		vba, err := r.readValue(size)
		if err != nil {
			return Command{}, err
		}
		return Command{
			Kind:      DatatypeSimpleString,
			Primitive: SIMPLE_STRING,
			Value:     noramlizeString(string(vba)),
		}, nil
	// TODO: Correctly handle bulk strings
	// right now this handles it the same way as simple strings
	case BULK_STRING:
		vba, err := r.readValue(size)
		if err != nil {
			return Command{}, err
		}
		return Command{
			Kind:      DatatypeBulkString,
			Primitive: BULK_STRING,
			Value:     noramlizeString(string(vba)),
		}, nil
	default:
		return Command{}, fmt.Errorf("Unsupported datatype: %v", dataType)
	}
}

// readValue reads the count of bytes denoted by `size`
// return a built byte array
// this function does not account for shifting terminations
func (r *RespParser) readValue(size int64) ([]byte, error) {
	valueB := make([]byte, size)
	for i := 0; i < int(size); i++ {
		b, err := r.reader.ReadByte()
		if err != nil {
			return []byte{}, err
		}
		valueB = append(valueB, b)
	}

	return valueB, nil
}

// parseSize converts a given string to a standard int64 used to parse commands
func parseSize(sizeB string) (int64, error) {
	size, err := strconv.ParseInt(sizeB, 16, 64)
	if err != nil {
		fmt.Println("failed parse int of size", string(sizeB))
		return -1, err
	}

	return size, nil
}

func (r *RespParser) initialParse() (DatatypePrimitive, int64, error) {
	//[0] Get datatype (shift first byte)
	t, err := r.reader.ReadByte()
	if err != nil {
		fmt.Println("failed to get dataType")
		return "", -1, err
	}
	dataType := DatatypePrimitive(t)

	// [1] Get data value size
	// -- if an array value size equates length of list
	// -- else len of value
	sizeB, err := r.reader.ReadByte()
	if err != nil {
		fmt.Println("failed to get size")
		return "", -1, err
	}
	// parse size byte array to string & convert to noramlized number
	size, err := parseSize(string(sizeB))
	if err != nil {
		fmt.Println("failed parse int of size")
		return "", -1, err
	}

	return dataType, size, nil
}

// parseArray starts a loop using the size of the incoming array parseArray will
// loop through recursively until the underlying "lines" are parsed.
// item parsing happens indirectly via r.handleDatatype
// this function assumes to shift between line endings
func (r *RespParser) parseArray(size int64) ([]Command, error) {
	// pre make the number of commands
	cmds := make([]Command, size)
	for i := 0; i < int(size); i++ {
		var dataType DatatypePrimitive
		var size int64
		var err error
		// Fetch the type and size of the current line
		dataType, size, err = r.initialParse()
		// After first arg the rest can be read conventionally
		if err != nil {
			fmt.Println("failed to read datatype2")
			return []Command{}, err
		}
		// shift past Line endings
		r.reader.ReadByte()
		r.reader.ReadByte()
		// Call handleData again to read & parse the inner argument
		cmd, err := r.handleDatatype(dataType, size)
		if err != nil {
			return []Command{}, err
		}
		// shift past Line endings
		r.reader.ReadByte()
		r.reader.ReadByte()
		// Set the cmd in a given index
		cmds[i] = cmd
	}
	return cmds, nil
}
