package wire

import (
	"bufio"
	"strings"

	"github.com/kasvith/kache/internal/protocol"
)

// Parser is used to parse wire protocol
type Parser struct {
	r *bufio.Reader
}

// NewParser creates a wire protocol parser
func NewParser(r *bufio.Reader) *Parser {
	return &Parser{r: r}
}

// Parse and return a Command and an error
func (p Parser) Parse() (*protocol.Command, error) {
	str, err := p.r.ReadString('\n')
	if err != nil {
		return nil, err
	}

	strLen := len(str)
	remLen := strLen
	if strLen > 0 {
		if str[strLen-1] == '\n' {
			remLen--
		}

		if strLen > 1 && str[strLen-2] == '\r' {
			remLen--
		}
	}

	str = str[:remLen]
	tokens := strings.Split(str, " ")

	if len(tokens) > 0 {
		return &protocol.Command{Name: strings.ToLower(tokens[0]), Args: tokens[1:]}, nil
	}

	return &protocol.Command{}, nil
}
