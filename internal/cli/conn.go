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

package cli

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/kasvith/kache/internal/protcl"
)

var c *cli

type cli struct {
	conn   net.Conn
	reader *bufio.Reader
}

func (r *cli) write(s string) error {
	send := make([]byte, len(s)+2)
	copy(send[:len(s)], s)
	copy(send[len(s):], []byte{'\r', '\n'})

	_, err := c.conn.Write(send)
	return err
}

func (r *cli) parseResp() (string, error) {
	buf, err := r.reader.ReadBytes('\n')
	if err != nil {
		if err == io.EOF {
			return "", nil
		}
		return "", err
	}

	if err := protcl.EndWithCRLF(buf); err != nil {
		return "", err
	}

	switch buf[0] {
	case '*':
		strs, err := protcl.ParseMultiBulkReply(r.reader, buf)
		if err != nil {
			return "", err
		}
		builder := strings.Builder{}
		for key, value := range strs {
			builder.WriteString(fmt.Sprintf("%d) %s\n", key+1, value))
		}
		return builder.String(), nil
	case '$':
		return protcl.ParseBulkString(r.reader, buf)
	case '+':
		return fmt.Sprintf("%q", buf[1:len(buf)-2]), nil
	case ':':
		return fmt.Sprintf("(integer) %s", buf[1:len(buf)-2]), nil
	case '-':
		return fmt.Sprintf("(error) %s", buf[1:len(buf)-2]), nil
	}

	return "", fmt.Errorf("unsupport resp type: %s", []byte{buf[0]})
}

// Dial conn kache server
func Dial(addr string) error {
	conn, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		return err
	}

	c = new(cli)
	c.conn = conn
	c.reader = bufio.NewReader(conn)

	return nil
}
