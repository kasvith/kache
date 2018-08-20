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
	"bytes"
	"errors"
)

var (
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

// TODO: need to handle special cases like \" \'
func SplitSpacesWithQuotes(s string) ([]string, error) {
	var ret []string
	const (
		space    = ' '
		dblQuote = '"'
	)

	var (
		buf          bytes.Buffer
		pos          = 0
		insideQuotes = false
	)

	for pos < len(s) {
		char := s[pos]

		if char == dblQuote {
			insideQuotes = !insideQuotes
			ret = appendIfBufferNotEmpty(&buf, ret)
			pos++
			continue
		}

		if char == space && !insideQuotes {
			// well skip it
			ret = appendIfBufferNotEmpty(&buf, ret)
			pos++
			continue
		}

		buf.WriteByte(char)
		pos++
	}

	// we have unbalanced quotes
	if insideQuotes {
		return []string{}, ErrUnbalancedQuotes
	}

	ret = appendIfBufferNotEmpty(&buf, ret)

	return ret, nil
}

func appendIfBufferNotEmpty(buf *bytes.Buffer, list []string) []string {
	if len(buf.String()) > 0 {
		list = append(list, buf.String())
		buf.Reset()
	}

	return list
}
