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

package util

import (
	"errors"
	"strings"
)

var (
	// ErrUnbalancedQuotes raised when quotes are not balanced in a string
	ErrUnbalancedQuotes = errors.New("unbalanced quotes")
)

// ToString Convert an interface to string
func ToString(i interface{}) string {
	if s, ok := i.(string); ok {
		return s
	}

	return ""
}

// SplitSpacesWithQuotes will split the string by spaces and preserve texts inside " " marks
// error is returned when an unbalanced quote was found in the string
func SplitSpacesWithQuotes(s string) ([]string, error) {
	var ret []string
	var buf = new(strings.Builder) // not in quote string buffer
	var scanned string
	var err error

	for pos := 0; pos < len(s); pos++ {
		char := s[pos]

		switch char {
		case ' ':
			if buf.Len() > 0 {
				ret = append(ret, buf.String())
				buf.Reset()
			}
		case '"':
			pos, scanned, err = scanForByte(s, pos, '"')
			if err != nil {
				return nil, err
			}
			ret = append(ret, scanned)
		default:
			buf.WriteByte(char)
		}
	}

	if buf.Len() > 0 {
		ret = append(ret, buf.String())
	}

	return ret, nil
}

func scanForByte(s string, pos int, r byte) (int, string, error) {
	var ret = new(strings.Builder)
	for pos++; pos < len(s); pos++ {
		char := s[pos]

		switch char {
		case '\\':
			if pos >= len(s)-1 {
				return 0, "", ErrUnbalancedQuotes
			}
			pos++
			ret.WriteByte(s[pos])
		case r:
			return pos, ret.String(), nil
		default:
			ret.WriteByte(char)
		}
	}

	return 0, "", ErrUnbalancedQuotes
}
