package resp2

import (
	"bytes"
	"fmt"
	"github.com/kasvith/kache/internal/protocol"
	"strconv"
	"strings"
)

const (
	// PrefixWrongType WRONGTYP
	PrefixWrongType = "WRONGTYP"

	// PrefixErr ERR
	PrefixErr = "ERR"

	TypeSimpleString = '+'

	TypeError = '-'

	TypeInteger = ':'

	TypeBulkString = '$'

	TypeArray = '*'

	CRLF = "\r\n"

	CR = '\r'

	LF = '\n'
)

// Command represents a command that can be executed
type Command struct {
	// Name of the command
	Name string
	// Args for command
	Args []string
}

type SimpleStringReply struct {
	Str string
}

func NewSimpleStringReply(str string) SimpleStringReply {
	return SimpleStringReply{Str: str}
}

func (s SimpleStringReply) ToBytes() []byte {
	return []byte(fmt.Sprintf("%c%s%s", TypeSimpleString, s.Str, CRLF))
}

type ErrorReply struct {
	Err error
}

func NewErrorReply(err error) ErrorReply {
	return ErrorReply{Err: err}
}

func (e ErrorReply) ToBytes() []byte {
	return []byte(fmt.Sprintf("%c%s%s", TypeError, e.Err.Error(), CRLF))
}

type IntegerReply struct {
	Value int
}

func NewIntegerReply(val int) IntegerReply {
	return IntegerReply{Value: val}
}

func (v IntegerReply) ToBytes() []byte {
	return []byte(fmt.Sprintf("%c%d%s", TypeInteger, v.Value, CRLF))
}

type BulkStringReply struct {
	Str    string
	IsNull bool
}

func NewBulkStringReply(isNull bool, str string) BulkStringReply {
	return BulkStringReply{Str: str, IsNull: isNull}
}

func (s BulkStringReply) ToBytes() []byte {
	if s.IsNull {
		return []byte("$-1\r\n")
	}

	builder := strings.Builder{}
	builder.WriteByte(TypeBulkString)
	builder.WriteString(strconv.Itoa(len(s.Str)))
	builder.WriteString(CRLF)
	builder.WriteString(s.Str)
	builder.WriteString(CRLF)

	return []byte(builder.String())
}

type ArrayReply struct {
	Reps   []protocol.Reply
	IsNull bool
}

func NewArrayReply(isNull bool, reps ...protocol.Reply) ArrayReply {
	return ArrayReply{IsNull: isNull, Reps: reps}
}

func (a ArrayReply) ToBytes() []byte {
	if a.IsNull {
		return []byte("*-1\r\n")
	}

	buf := bytes.Buffer{}
	buf.WriteByte(TypeArray)
	buf.WriteString(strconv.Itoa(len(a.Reps)))
	buf.WriteString(CRLF)
	for _, value := range a.Reps {
		buf.Write(value.ToBytes())
	}

	return buf.Bytes()
}
