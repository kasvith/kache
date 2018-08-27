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
	"net"
	"time"

	"github.com/kasvith/kache/internal/protcl"
)

var c *cli

type cli struct {
	conn        net.Conn
	resp3Parser *protcl.Resp3Parser
	addr        string
}

// Write send string to server
func (r *cli) Write(s string) error {
	return r.write(s, true)
}

func (r *cli) write(s string, reconnect bool) error {
	n, err := c.conn.Write([]byte(s))
	if n == 0 && err != nil && reconnect {
		fmt.Println("reconnecting...")

		if err := Dial(r.addr); err != nil {
			return err
		}
		return r.write(s, false)
	}
	return err
}

// Dial conn kache server
func Dial(addr string) error {
	conn, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		return err
	}

	c = new(cli)
	c.conn = conn
	c.resp3Parser = protcl.NewResp3Parser(bufio.NewReader(conn))
	c.addr = addr

	return nil
}
