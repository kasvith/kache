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
	"errors"
	"fmt"
	"github.com/kasvith/kache/internal/config"
	"github.com/kasvith/kache/internal/klogs"
	"io"
	"strconv"
	"strings"
)

type ParseError struct {
}

func (ParseError) Error() string {
	return "-ERR:parse error"
}

type Reader struct {
	br *bufio.Reader
}

func NewReader(r io.Reader) *Reader {
	return &Reader{br: bufio.NewReader(r)}
}

func (r *Reader) ParseMessage() (*RespCommand, error) {
	return parse(r.br)
}

func parse(r *bufio.Reader) (*RespCommand, error) {
	// we have two kind of messages to parse
	// a redis array is an acceptable command
	// a simple string with space separated is also acceptable

	// read until first delimiter occurrence
	buf, err := r.ReadBytes('\n')

	if err != nil {
		klogs.Logger.Debugln(err)
		return nil, err
	}

	// Clients require to send commands with CRLF
	switch buf[0] {
	case REP_ARR:
		// this is an array of redis strings
		// arr len is in the buffer

		fmt.Println(buf)

		if err := hasCRLF(buf); err != nil {
			return nil, err
		}

		mblkLen, err := strconv.Atoi(string(buf[1 : len(buf)-2]))

		if err != nil {
			return nil, errors.New("value out of range")
		}

		// we now have multibulk length, now need to loop that amount
		// TODO check for maximum number of array elements to process to handle memory issues

		fmt.Println("arr", mblkLen)

		strs := make([]string, mblkLen)
		for i := 0; i < mblkLen; i++ {
			str, err := parseBlkString(r)
			if err != nil {
				return nil, err
			}

			strs[i] = str
			fmt.Println(str)
		}

		fmt.Printf("%v\n", mblkLen)

		return nil, errors.New("ERR:not yet implemented")
	default:
		// probably the read bytes contains the string
		strCmd := string(buf)

		// trim the trailing \r from cmd
		trimmed := strings.TrimSpace(strCmd)

		// split args by space
		args := strings.Split(trimmed, " ")

		if len(args) == 0 {
			return nil, errors.New("ERR:no command")
		}

		return &RespCommand{Name: strings.ToLower(args[0]), Args: args[1:]}, nil
	}

	return nil, errors.New("ERR:not yet implemented")
}

func hasCRLF(buf []byte) error {
	if len(buf) >= 2 && buf[len(buf)-1] == '\n' && buf[len(buf)-2] == '\r' {
		return nil
	}

	return errors.New("invalid line end")
}

func parseBlkString(r *bufio.Reader) (string, error) {
	// read a byte
	buf, err := r.ReadBytes('\n')
	if err != nil {
		return "", err
	}

	fmt.Println(buf)

	if len(buf) > 0 && buf[0] != REP_BULKSTRING {
		fmt.Println(buf[0])
		return "", errors.New("invalid bulk string")
	}

	if err = hasCRLF(buf); err != nil {
		return "", err
	}

	llen, err := strconv.Atoi(string(buf[1 : len(buf)-2]))
	if err != nil {
		return "", err
	}

	if llen > config.AppConf.MaxMultiBlkLength {
		return "", errors.New("buffer exceeded")
	}

	fmt.Println("llen", llen)

	// we need to read exactly llen bytes from the stream
	strBuf := make([]byte, llen)
	n, err := r.Read(strBuf)

	fmt.Println("buf", strBuf)

	if err != nil {
		return "", err
	}

	if n != llen || (llen >= 2 && (strBuf[llen-2] == '\r' || strBuf[llen-1] == '\n')) {
		return "", errors.New("error parsing bulk string")
	}

	// read trailing CRLF
	rest, err := r.Peek(2)
	if err != nil {
		if err == io.EOF {
			return "", err
		}

		return "", errors.New("error reading line end")
	}

	if rest[0] != '\r' && rest[1] != '\n' {
		return "", errors.New("error reading line end")
	}

	// discard crlf
	_, err = r.Discard(2)

	if err != nil {
		return "", err
	}

	return string(strBuf), nil
}
