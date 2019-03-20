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
