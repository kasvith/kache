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
	"errors"
	"testing"

	"github.com/kasvith/kache/internal/protocol"
	testifyAssert "github.com/stretchr/testify/assert"
)

func TestSimpleStringReply(t *testing.T) {
	reply := NewSimpleStringReply("foo")
	testifyAssert.Implements(t, (*protocol.Reply)(nil), reply)
	testifyAssert.Equal(t, []byte("+foo\r\n"), reply.ToBytes())
}

func TestIntegerReply(t *testing.T) {
	reply := NewIntegerReply(100)
	testifyAssert.Equal(t, []byte(":100\r\n"), reply.ToBytes())

	reply = NewIntegerReply(-100)
	testifyAssert.Equal(t, []byte(":-100\r\n"), reply.ToBytes())
}

func TestErrorReply(t *testing.T) {
	reply := NewErrorReply(errors.New("something went wrong"))
	testifyAssert.Equal(t, []byte("-something went wrong\r\n"), reply.ToBytes())
}

func TestBulkStringReply(t *testing.T) {
	reply := NewBulkStringReply(false, "foobar")
	testifyAssert.Equal(t, []byte("$6\r\nfoobar\r\n"), reply.ToBytes())

	reply = NewBulkStringReply(false, "")
	testifyAssert.Equal(t, []byte("$0\r\n\r\n"), reply.ToBytes())

	reply = NewBulkStringReply(true, "")
	testifyAssert.Equal(t, []byte("$-1\r\n"), reply.ToBytes())
}

func TestArrayReply(t *testing.T) {
	reply := NewArrayReply(false, []protocol.Reply{
		NewSimpleStringReply("foo"),
		NewIntegerReply(10),
		NewErrorReply(errors.New("foo bar")),
		NewBulkStringReply(false, "bar"),
		NewBulkStringReply(true, ""),
	})

	testifyAssert.Equal(t, []byte("*5\r\n+foo\r\n:10\r\n-foo bar\r\n$3\r\nbar\r\n$-1\r\n"), reply.ToBytes())

	reply = NewArrayReply(false, []protocol.Reply{})
	testifyAssert.Equal(t, []byte("*0\r\n"), reply.ToBytes())

	reply = NewArrayReply(true, []protocol.Reply{})
	testifyAssert.Equal(t, []byte("*-1\r\n"), reply.ToBytes())
}
