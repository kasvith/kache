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

package client

import (
	"github.com/kasvith/kache/internal/protocol"
	"github.com/kasvith/kache/internal/resp/resp2"
)

// Ping will return PONG when no argument found or will echo the given argument
func Ping(client *Client, args []string) {
	if len(args) == 0 {
		client.WriteProtocolReply(resp2.NewSimpleStringReply("PONG"))
		return
	}

	client.WriteProtocolReply(resp2.NewSimpleStringReply(args[0]))
}

// Multi command will put client in multi mode where can execute multiple commands at once
func Multi(client *Client, args []string) {
	client.Multi = true
	client.WriteProtocolReply(resp2.NewSimpleStringReply("OK"))
}

// Exec command will execute a multi transaction
func Exec(client *Client, args []string) {
	if !client.Multi {
		client.WriteError(protocol.ErrExecWithoutMulti{})
		return
	}

	client.Multi = false
	if client.MultiError {
		client.MultiError = false
		client.Commands = []*Command{}
		client.WriteError(protocol.ErrExecAbortTransaction{})
		return
	}

	for _, cmd := range client.Commands {
		cmd.Fn(client, cmd.Args)
	}

	// clear all commands
	client.Commands = []*Command{}
}
