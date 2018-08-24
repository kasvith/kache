/*
 * MIT License
 *
 * Copyright (c)  2018 Kasun Vithanage
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package protcl

import (
	"fmt"
	"strings"
)

const (
	// RepSimpleString RESP simple string identifier
	RepSimpleString byte = '+'

	// RepInteger RESP integer
	RepInteger = ':'

	// RepBulkString RESP bulk string
	RepBulkString = '$'

	// RepError RESP error
	RepError = '-'

	// RepArray RESP array
	RepArray = '*'
)

const (
	// PrefixWrongType WRONGTYP
	PrefixWrongType = "WRONGTYP"

	// PrefixErr ERR
	PrefixErr = "ERR"
)

// Reply for RESP replies
type Reply interface {
	Reply() string
}

// Message is a RESP message
type Message struct {
	Reply
	Err error
}

// NewMessage creates a new message
func NewMessage(rep Reply, err error) *Message {
	return &Message{Reply: rep, Err: err}
}

// checks for resp prefix in errors
func hasRespPrefix(str string) bool {
	if strings.HasPrefix(str, PrefixWrongType) || strings.HasPrefix(str, PrefixErr) {
		return true
	}

	return false
}

// RespError creates a RESP error with prefix
func RespError(Err error) string {
	err := Err.Error()

	if !hasRespPrefix(err) {
		err = "ERR:" + err
	}

	return fmt.Sprintf("-%s\r\n", err)
}

// RespReply builds a RESP reply
func (msg *Message) RespReply() string {
	return msg.Reply.Reply()
}

// RespCommand represents a command that can be executed by the kache server
// It can also contain a multi command(pipelined)
type RespCommand struct {
	// Name of the command
	Name string

	// Args for command
	Args []string

	// Multi type command
	Multi bool

	// Commands array when its multi
	Commands []RespCommand
}

// NewIntegerReply creates a new integer reply
func NewIntegerReply(value int) *IntegerReply {
	return &IntegerReply{Value: value}
}

// IntegerReply Represents a RESP integer reply
type IntegerReply struct {
	Value int
}

// Reply method for integers
func (rep *IntegerReply) Reply() string {
	return fmt.Sprintf(":%d\r\n", rep.Value)
}

// NewSimpleStringReply creates a new SimpleStringReply
func NewSimpleStringReply(value string) *SimpleStringReply {
	return &SimpleStringReply{Value: value}
}

// SimpleStringReply Binary unsafe strings
type SimpleStringReply struct {
	Value string
}

// Reply method for integers
func (rep *SimpleStringReply) Reply() string {
	return fmt.Sprintf("+%s\r\n", rep.Value)
}

// NewBulkStringReply creates a new BulkStringReply
func NewBulkStringReply(isNil bool, value string) *BulkStringReply {
	return &BulkStringReply{Nil: isNil, Value: value}
}

// BulkStringReply is a binary safe string
type BulkStringReply struct {
	Value string
	Nil   bool
}

// Reply will return string reperesntation of bulk string reply in RESP
func (rep *BulkStringReply) Reply() string {
	if rep.Nil {
		return fmt.Sprintf("$-1\r\n")
	}

	return fmt.Sprintf("$%d\r\n%s\r\n", len(rep.Value), rep.Value)
}

// ArrayReply is used to reply with arrays
type ArrayReply struct {
	Elems []Reply
	Nil   bool
}

// NewArrayReply creates a new ArrayReply
func NewArrayReply(isNil bool, elems []Reply) *ArrayReply {
	return &ArrayReply{Elems: elems, Nil: isNil}
}

// Reply is RESP representation of ArrayReply
func (rep *ArrayReply) Reply() string {
	if rep.Nil {
		return "*-1\r\n"
	}

	length := len(rep.Elems)
	builder := strings.Builder{}

	builder.WriteString(fmt.Sprintf("*%d\r\n", length))

	for _, re := range rep.Elems {
		builder.WriteString(re.Reply())
	}

	return builder.String()
}
