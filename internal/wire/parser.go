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
