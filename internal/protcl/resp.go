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
	REP_SIMPLE_STRING = "+"
	REP_INTEGER       = ":"
	REP_BULKSTRING    = "$"
	REP_ERROR         = "-"
)

const (
	WRONGTYP = "WRONGTYP"
	ERR      = "ERR"
)

type Reply interface {
	Reply() string
}

type Err interface {
	Err() ErrorReply
}

type Message struct {
	Rep Reply
	Err Err
}

type ErrorReply struct {
	Prefix string
	Err    string
}

func (rep *ErrorReply) Error() string {
	return fmt.Sprintf("-%s %s\r\n", rep.Prefix, rep.Err)
}

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
	IsNil bool
}

func NewArrayReply(isNil bool, elems []Reply) *ArrayReply {
	return &ArrayReply{Elems: elems, IsNil: isNil}
}

func (rep *ArrayReply) Reply() string {
	if rep.IsNil {
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
