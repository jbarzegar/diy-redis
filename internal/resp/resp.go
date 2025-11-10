package respparser

// LE is the standard line endings for RESP
const LE = "\r\n"

type Command struct {
	Kind      Datatype
	Primitive DatatypePrimitive
	Value     string
	// Only set when "Kind" is Array
	Args []Command
}

// DatatypePrimitive are literal symbols RESP uses to denote datatypes
type DatatypePrimitive string

const (
	SIMPLE_STRING DatatypePrimitive = "+"
	ERROR         DatatypePrimitive = "-"
	INTEGER       DatatypePrimitive = ":"
	BULK_STRING   DatatypePrimitive = "$"
	ARRAY         DatatypePrimitive = "*"
)

// Datatype is the human readable equivilent of each data type
type Datatype string

const (
	DatatypeSimpleString Datatype = "simple-string"
	DatatypeError        Datatype = "error"
	DatatypeInteger      Datatype = "integer"
	DatatypeBulkString   Datatype = "bulk_string"
	DatatypeArray        Datatype = "array"
)
