package wire

import (
	"bufio"
	"strings"
	"testing"

	testifyAssert "github.com/stretchr/testify/assert"
)

func getParser(str string) *Parser {
	return NewParser(bufio.NewReader(strings.NewReader(str)))
}

func TestParser_ParseCRLF(t *testing.T) {
	p := getParser("welcome to kache\r\n")
	cmd, err := p.Parse()
	testifyAssert.Nil(t, err)

	testifyAssert.Equal(t, "welcome", cmd.Name)
	testifyAssert.Equal(t, "to", cmd.Args[0])
	testifyAssert.Equal(t, "kache", cmd.Args[1])
}

func TestParser_ParseLF(t *testing.T) {
	p := getParser("welcome to kache\n")
	cmd, err := p.Parse()
	testifyAssert.Nil(t, err)

	testifyAssert.Equal(t, "welcome", cmd.Name)
	testifyAssert.Equal(t, "to", cmd.Args[0])
	testifyAssert.Equal(t, "kache", cmd.Args[1])
}

func TestParser_ParseEmpty(t *testing.T) {
	p := getParser("\n")
	cmd, err := p.Parse()
	testifyAssert.Nil(t, err)
	testifyAssert.Equal(t, "", cmd.Name)
	testifyAssert.Equal(t, 0, len(cmd.Args))
}
