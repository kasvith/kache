/*
 * MIT License
 *
 * Copyright (c) 2019 Kasun Vithanage
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
 *
 */

package resp2

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/kasvith/kache/internal/protocol"
)

const (
	// PrefixWrongType WRONGTYP
	PrefixWrongType = "WRONGTYP"

	// PrefixErr ERR
	PrefixErr = "ERR"

	// TypeSimpleString represents a RESP2 simple string
	TypeSimpleString = '+'

	// TypeError represents a RESP2 error
	TypeError = '-'

	// TypeInteger represents a RESP2 integer
	TypeInteger = ':'

	// TypeBulkString represents a RESP2 bulk string
	TypeBulkString = '$'

	// TypeArray reperesnts a RESP2 array
	TypeArray = '*'

	// CRLF represents line ending \r\n
	CRLF = "\r\n"

	// CR represents Carriage Return
	CR = '\r'

	// LF represents Line Follow
	LF = '\n'
)

// SimpleStringReply represents a RESP2 simple string reply
type SimpleStringReply struct {
	Str string
}

// NewSimpleStringReply creates a SimpleStringReply
func NewSimpleStringReply(str string) *SimpleStringReply {
	return &SimpleStringReply{Str: str}
}

// ToBytes returns the byte representation of the SimpleStirngReply
func (s SimpleStringReply) ToBytes() []byte {
	return []byte(fmt.Sprintf("%c%s%s", TypeSimpleString, s.Str, CRLF))
}

// ErrorReply is used to indicate an error processing the request
type ErrorReply struct {
	Err error
}

// NewErrorReply creates a new ErrorReply struct
func NewErrorReply(err error) *ErrorReply {
	return &ErrorReply{Err: err}
}

// ToBytes returns byte representation of the ErrorReply
func (e ErrorReply) ToBytes() []byte {
	return []byte(fmt.Sprintf("%c%s%s", TypeError, e.Err.Error(), CRLF))
}

// IntegerReply used to return an integer from server
type IntegerReply struct {
	Value int
}

// NewIntegerReply creates a new IntegerReply
func NewIntegerReply(val int) *IntegerReply {
	return &IntegerReply{Value: val}
}

// ToBytes returns byte representation of IntegerReply
func (v IntegerReply) ToBytes() []byte {
	return []byte(fmt.Sprintf("%c%d%s", TypeInteger, v.Value, CRLF))
}

// BulkStringReply is used to store and return a bulk string
type BulkStringReply struct {
	Str    string
	IsNull bool
}

// NewBulkStringReply is used to create a new BulkStringReply
func NewBulkStringReply(isNull bool, str string) *BulkStringReply {
	return &BulkStringReply{Str: str, IsNull: isNull}
}

// ToBytes returns byte representation of BulkStringReply
func (s BulkStringReply) ToBytes() []byte {
	// handle null strings
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

// ArrayReply is used to store a RESP array
type ArrayReply struct {
	Reps   []protocol.Reply
	IsNull bool
}

// NewArrayReply creates a new ArrayReply
func NewArrayReply(isNull bool, reps []protocol.Reply) *ArrayReply {
	return &ArrayReply{IsNull: isNull, Reps: reps}
}

// ToBytes returns byte representation of ArrayReply
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
