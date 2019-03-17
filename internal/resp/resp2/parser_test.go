package resp2

import (
	"bufio"
	"strings"
	"testing"

	"github.com/kasvith/kache/internal/protocol"

	testifyAssert "github.com/stretchr/testify/assert"
)

func getParser(input string) *Parser {
	return NewParser(bufio.NewReader(strings.NewReader(input)))
}

func TestFailOnWireProtocol(t *testing.T) {
	parser := getParser("this is not a resp")
	cmd, err := parser.Parse()

	testifyAssert.Nil(t, cmd)
	testifyAssert.Equal(t, &protocol.ErrUnknownProtocol{}, err)
}

func TestParsesRESPEmptyArrays(t *testing.T) {
	parser := getParser("*0\r\n")
	cmd, err := parser.Parse()

	testifyAssert.Nil(t, err)
	testifyAssert.Empty(t, cmd.Name)
	testifyAssert.Empty(t, cmd.Args)
}

func TestParsesRESPNullArrays(t *testing.T) {
	parser := getParser("*-1\r\n")
	cmd, err := parser.Parse()

	testifyAssert.Nil(t, err)
	testifyAssert.Nil(t, cmd)
}

func TestParsesRESPBulkStringArraysWithOneElement(t *testing.T) {
	parser := getParser("*1\r\n$3\r\nfoo\r\n")
	cmd, err := parser.Parse()

	testifyAssert.Nil(t, err)
	testifyAssert.Equal(t, "foo", cmd.Name)
}

func TestParsesRESPBulkStringArraysWithMultipleElements(t *testing.T) {
	parser := getParser("*3\r\n$3\r\nfoo\r\n$3\r\nbar\r\n$6\r\nfoobar\r\n")
	cmd, err := parser.Parse()

	testifyAssert.Nil(t, err)
	testifyAssert.Equal(t, "foo", cmd.Name)
	testifyAssert.Equal(t, "bar", cmd.Args[0])
	testifyAssert.Equal(t, "foobar", cmd.Args[1])
}

func TestDoesNotParsesRESPArraysWithNonBulkStrings(t *testing.T) {
	parser := getParser("*3\r\n$3\r\nfoo\r\n$3\r\nbar\r\n:100\r\n")
	cmd, err := parser.Parse()

	testifyAssert.Nil(t, cmd)
	testifyAssert.Equal(t, &protocol.ErrWrongType{}, err)
}

func TestDoesNotParsesRESPArraysWithMalformedBulkStrings(t *testing.T) {
	parser := getParser("*2\r\n$3\r\nfo\r\n$3\r\nbar\r\n")
	cmd, err := parser.Parse()

	testifyAssert.Nil(t, cmd)
	testifyAssert.Equal(t, &protocol.ErrUnexpectedLineEnd{}, err)
}

func TestDoesNotParsesRESPArraysWithMalformedBulkStringLengths(t *testing.T) {
	parser := getParser("*2\r\n$x\r\nfoo\r\n$3\r\nbar\r\n")
	cmd, err := parser.Parse()

	testifyAssert.Nil(t, cmd)
	testifyAssert.Equal(t, &protocol.ErrCastFailedToInt{Val: "x"}, err)
}

func TestDoesNotParsesRESPArraysWithoutCRLF(t *testing.T) {
	parser := getParser("*2\r\n$3\r\nfoo$3\r\nbar\r\n")
	cmd, err := parser.Parse()

	testifyAssert.Nil(t, cmd)
	testifyAssert.Equal(t, &protocol.ErrUnexpectedLineEnd{}, err)
}

func TestDoesNotParsesRESPArraysWithoutArrayLength(t *testing.T) {
	parser := getParser("*c\r\n$x\r\nfoo\r\n$3\r\nbar\r\n")
	cmd, err := parser.Parse()

	testifyAssert.Nil(t, cmd)
	testifyAssert.Equal(t, &protocol.ErrCastFailedToInt{Val: "c"}, err)
}
