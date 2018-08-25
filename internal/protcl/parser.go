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
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/kasvith/kache/internal/config"
	"github.com/kasvith/kache/pkg/util"
)

// ErrValueOutOfRange for out of range errors
type ErrValueOutOfRange struct {
}

// Recoverable whether error is recoverable or not
func (ErrValueOutOfRange) Error() string {
	return "value out of range"
}

// Recoverable whether error is recoverable or not
func (ErrValueOutOfRange) Recoverable() bool {
	return true
}

//ErrInvalidCommand for invalid commands
type ErrInvalidCommand struct {
}

func (ErrInvalidCommand) Error() string {
	return "invalid command"
}

// Recoverable whether error is recoverable or not
func (ErrInvalidCommand) Recoverable() bool {
	return true
}

// ErrBufferExceeded for buffer exceeds
type ErrBufferExceeded struct {
}

func (ErrBufferExceeded) Error() string {
	return "buffer exceeded"
}

// Recoverable whether error is recoverable or not
func (ErrBufferExceeded) Recoverable() bool {
	return true
}

// ErrUnexpectedLineEnd for unexpected line ends(no CRLF)
type ErrUnexpectedLineEnd struct {
}

func (ErrUnexpectedLineEnd) Error() string {
	return "unexpected line end"
}

// Recoverable whether error is recoverable or not
func (ErrUnexpectedLineEnd) Recoverable() bool {
	return true
}

// ErrInvalidToken for invalid tokens
type ErrInvalidToken struct {
	Token byte
}

func (e *ErrInvalidToken) Error() string {
	return fmt.Sprintf("excepted $, found %c", e.Token)
}

// Recoverable whether error is recoverable or not
func (ErrInvalidToken) Recoverable() bool {
	return true
}

// ErrInvalidBlkStringLength raised when bulk string length mismatch
type ErrInvalidBlkStringLength struct {
	Excepted, Given int
}

func (e *ErrInvalidBlkStringLength) Error() string {
	return fmt.Sprintf("invalid bulk string length, excepted %d processed %d", e.Excepted, e.Given)
}

// Recoverable whether error is recoverable or not
func (ErrInvalidBlkStringLength) Recoverable() bool {
	return true
}

// Reader for parser
type Reader struct {
	br *bufio.Reader
}

// NewReader creates a reader for parser
func NewReader(r io.Reader) *Reader {
	return &Reader{br: bufio.NewReader(r)}
}

// ParseMessage will parse a message from reader
func (r *Reader) ParseMessage() (*RespCommand, error) {
	return parse(r.br)
}

func parse(r *bufio.Reader) (*RespCommand, error) {
	// TODO: these reads can easily overflow the system buffer and crash the program, they need to be max buffer protected
	// for now we use default bufio package, we need a custom reader

	// we have two kind of messages to parse
	// a redis array is an acceptable command
	// a simple string with space separated is also acceptable

	// read until first delimiter occurrence
	buf, err := r.ReadBytes('\n')

	if err != nil {
		return nil, err
	}

	// Clients require to send commands with CRLF
	switch buf[0] {
	case RepArray:
		// this is an array of redis strings
		// arr len is in the buffer

		strs, err := ParseMultiBulkReply(r, buf)
		if err != nil {
			return nil, err
		} else if len(strs) == 0 {
			return nil, &ErrInvalidCommand{}
		}

		return &RespCommand{Name: strings.ToLower(strs[0]), Args: strs[1:]}, nil
	default:
		// probably the read bytes contains the string
		strCmd := string(buf)

		// trim the trailing \r from cmd
		trimmed := strings.TrimSpace(strCmd)

		// split args by space
		args, err := util.SplitSpacesWithQuotes(trimmed)

		// error is unbalanced quote
		if err != nil {
			return nil, err
		}

		if len(args) == 0 {
			return nil, &ErrInvalidCommand{}
		}

		return &RespCommand{Name: strings.ToLower(args[0]), Args: args[1:]}, nil
	}
}

// ParseMultiBulkReply parse multi bulk reply message
func ParseMultiBulkReply(r *bufio.Reader, buf []byte) ([]string, error) {
	if err := EndWithCRLF(buf); err != nil {
		// not a EOF, safe to return
		return nil, err
	}

	mblkLen, err := strconv.Atoi(string(buf[1 : len(buf)-2]))
	if err != nil {
		return nil, &ErrValueOutOfRange{}
	} else if mblkLen == 0 {
		return nil, nil
	}

	// we now have multibulk length, now need to loop that amount
	// TODO check for maximum number of array elements to process to handle memory issues

	strs := make([]string, mblkLen)
	for i := 0; i < mblkLen; i++ {
		str, err := parseBulkString(r)
		if err != nil {
			return nil, err
		}

		strs[i] = str
	}

	return strs, nil
}

// ParseBulkString parse bulk string message
func ParseBulkString(r *bufio.Reader, buf []byte) (string, error) {
	if len(buf) > 0 && buf[0] != RepBulkString {
		return "", &ErrInvalidToken{Token: buf[0]}
	}

	if err := EndWithCRLF(buf); err != nil {
		return "", err
	}

	llen, err := strconv.Atoi(string(buf[1 : len(buf)-2]))
	if err != nil {
		return "", err
	}

	if llen > config.AppConf.MaxMultiBulkLength {
		return "", &ErrBufferExceeded{}
	}

	// we need to read exactly llen bytes from the stream
	buf, err = r.ReadBytes('\n')

	if err != nil {
		return "", err
	}

	// error is not EOF
	err = EndWithCRLF(buf)
	if err != nil {
		return "", err
	}

	str := buf[:len(buf)-2]

	if len(str) != llen {
		return "", &ErrInvalidBlkStringLength{Excepted: llen, Given: len(str)}
	}

	return string(str), nil
}

func parseBulkString(r *bufio.Reader) (string, error) {
	// read a byte
	buf, err := r.ReadBytes('\n')
	if err != nil {
		return "", err
	}

	return ParseBulkString(r, buf)
}

// EndWithCRLF returns an error when buffer end is not CRLF
func EndWithCRLF(buf []byte) error {
	if len(buf) >= 2 && buf[len(buf)-1] == '\n' && buf[len(buf)-2] == '\r' {
		return nil
	}

	return &ErrUnexpectedLineEnd{}
}
