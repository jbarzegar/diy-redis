package respparser

import (
	"fmt"
	"strings"
	"testing"
)

func makeCMD(dataType DatatypePrimitive, value string) string {
	return fmt.Sprintf("%v%v\r\n%v\r\n", dataType, len(value), value)
}

func makeArray(values ...string) string {
	return fmt.Sprintf("*%v\r\n%v\r\n", len(values), strings.Join(values, "\r\n"))
}

func TestParserSimpleString(t *testing.T) {
	str := makeCMD(SIMPLE_STRING, "SOMETHING")

	reader := NewRespParser(strings.NewReader(str))
	parsed, err := reader.Read()
	if err != nil {
		fmt.Print("reader failed")
		t.Fatal(err)
	}

	if parsed.Kind != "simple-string" {
		fmt.Println("Should be simple string")
		t.Fail()
	}

	if parsed.Value != "SOMETHING" {
		fmt.Printf("value not right %q\n", parsed.Value)
		t.Fail()
	}
}

func TestParserArray(t *testing.T) {
	str := "*2\r\n$7\r\nCOMMAND\r\n$4\r\nDOCS\r\n"

	reader := NewRespParser(strings.NewReader(str))
	parsed, err := reader.Read()
	if err != nil {
		t.Fatal(err)
	}

	if parsed.Kind != "array" {
		fmt.Println("Should be simple string")
		t.Fail()
	}

	// fmt.Println(parsed.Args)
	fmt.Println("TODO: Ensure data format can be handled correctly")
	// t.Fail()
	if len(parsed.Args) != 2 {
		fmt.Println("Unexpected number of args")
		fmt.Println("wanted 2")
		t.Fail()
	}
	for _, cmd := range parsed.Args {
		fmt.Printf("%q%q\n", cmd.Kind, cmd.Value)
	}

}
