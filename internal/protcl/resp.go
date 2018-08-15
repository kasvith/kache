package protcl

import (
	"fmt"
	"strings"
)

// IntegerReply Represents an integer reply
type IntegerReply struct {
	Value int
}

// Reply method for integers
func (rep *IntegerReply) Reply() string {
	return fmt.Sprintf(":%d\r\n", rep.Value)
}

func NewIntegerReply(value int) *IntegerReply {
	return &IntegerReply{Value: value}
}

// SimpleStringReply Binary unsafe strings
type SimpleStringReply struct {
	Value string
}

// Reply method for integers
func (rep *SimpleStringReply) Reply() string {
	return fmt.Sprintf("+%s\r\n", rep.Value)
}

func NewSimpleStringReply(value string) *SimpleStringReply {
	return &SimpleStringReply{Value: value}
}

type BulkStringReply struct {
	Value string
	Nil   bool
}

func (rep *BulkStringReply) Reply() string {
	if rep.Nil == true {
		return fmt.Sprintf("$-1\r\n")
	}

	return fmt.Sprintf("$%d\r\n%s\r\n", len(rep.Value), rep.Value)
}

func NewBulkStringReply(isNil bool, value string) *BulkStringReply {
	return &BulkStringReply{Nil: isNil, Value: value}
}

type ArrayReply struct {
	Elems []Reply
}

func NewArrayReply(elems []Reply) *ArrayReply {
	return &ArrayReply{Elems: elems}
}

func (rep *ArrayReply) Reply() string {
	length := len(rep.Elems)
	builder := strings.Builder{}

	builder.WriteString(fmt.Sprintf("*%d\r\n", length))

	for _, re := range rep.Elems {
		builder.WriteString(re.Reply() + "\r\n")
	}

	return builder.String()
}

type ErrorReply struct {
	Prefix string
	Err    string
}

func (rep *ErrorReply) Error() string {
	return fmt.Sprintf("-%s %s\r\n", rep.Prefix, rep.Err)
}
