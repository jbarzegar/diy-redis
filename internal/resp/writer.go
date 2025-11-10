package respparser

import (
	"fmt"
	"strings"
)

// Implements io.writer (kinda)
type RespWriter struct{}

func NewWriter() *RespWriter {
	return &RespWriter{}
}

func createLineStart(datatype DatatypePrimitive, length int) string {
	return fmt.Sprintf("%v%v", datatype, length)
}

func (w *RespWriter) WriteSimpleString(str string) string {
	length := len(str)
	blocks := []string{
		createLineStart(SIMPLE_STRING, length),
		str,
	}
	return strings.Join(blocks, LE) + LE
}

func (w *RespWriter) WriteArray(strs ...Command) string {
	arrLen := len(strs)
	value := createLineStart(ARRAY, arrLen) + LE

	for _, c := range strs {
		if c.Kind == DatatypeArray {
			fmt.Println("Sorry array writer doesn't support nested arrays yest")
			continue
		}
		value += createLineStart(c.Primitive, len(c.Value)) + LE
		value += c.Value + LE
	}

	if strings.HasSuffix(value, LE) {
		return value
	}

	return value + LE
}
