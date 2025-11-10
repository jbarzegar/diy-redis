package respparser_test

import (
	respparser "diy-redis/internal/resp"
	"fmt"
	"testing"
)

func TestWriterSimpleString(t *testing.T) {
	str := "TEST"
	expected := fmt.Sprintf("+%v\r\n%v\r\n", len(str), str)
	writer := respparser.NewWriter()

	result := writer.WriteSimpleString(str)

	if result != expected {
		fmt.Println("Simple string doesn't match")
		fmt.Printf("Got %q\n", result)
		t.Fail()
	}
}

func TestWriterArrays(t *testing.T) {
	arrayItems := []respparser.Command{}
	strs := []string{
		"COMMAND",
		"DOCS",
	}
	expected := "*2\r\n"

	for _, str := range strs {
		// Append item to expected
		expected += fmt.Sprintf("$%v\r\n%v\r\n", len(str), str)
		// Build the array item to test in turn
		arrayItems = append(arrayItems, respparser.Command{
			Kind:      respparser.DatatypeBulkString,
			Primitive: respparser.BULK_STRING,
			Value:     str,
		})
	}
	writer := respparser.NewWriter()
	// call WriteArray
	result := writer.WriteArray(arrayItems...)

	if result != expected {
		fmt.Println("array doesn't match")
		fmt.Printf("Got %q\n", result)
		fmt.Printf("Expected %q \n", expected)
		t.Fail()
	}
}
