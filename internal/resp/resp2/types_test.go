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
	reply := NewArrayReply(
		false,
		NewSimpleStringReply("foo"),
		NewIntegerReply(10),
		NewErrorReply(errors.New("foo bar")),
		NewBulkStringReply(false, "bar"),
		NewBulkStringReply(true, ""),
	)
	testifyAssert.Equal(t, []byte("*5\r\n+foo\r\n:10\r\n-foo bar\r\n$3\r\nbar\r\n$-1\r\n"), reply.ToBytes())

	reply = NewArrayReply(false)
	testifyAssert.Equal(t, []byte("*0\r\n"), reply.ToBytes())

	reply = NewArrayReply(true)
	testifyAssert.Equal(t, []byte("*-1\r\n"), reply.ToBytes())
}
