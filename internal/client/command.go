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
	"strings"

	"github.com/kasvith/kache/internal/protocol"
	"github.com/kasvith/kache/internal/resp/resp2"
)

// CommandFunc holds a function signature which can be used as a command.
type CommandFunc func(*Client, []string)

// Command holds a command structure which is used to execute a kache command
type Command struct {
	ModifyKeySpace bool
	Fn             CommandFunc
	MinArgs        int // 0
	MaxArgs        int // -1 ~ +inf, -1 mean infinite
	Args           []string
}

// GetCommand will fetch the command from command table
func GetCommand(cmd string) (*Command, error) {
	if v, ok := CommandTable[cmd]; ok {
		return &v, nil
	}

	return nil, &protocol.ErrUnknownCommand{Cmd: cmd}
}

// Execute a single command on the given database with args
func Execute(client *Client, cmd string, args []string) {
	command, err := GetCommand(cmd)
	if err != nil {
		if client.Multi {
			client.MultiError = true
			client.Commands = []*Command{}
		}
		client.WriteError(err)
		return
	}

	if argsLen := len(args); (command.MinArgs > 0 && argsLen < command.MinArgs) || (command.MaxArgs != -1 && argsLen > command.MaxArgs) {
		if client.Multi {
			client.MultiError = true
			client.Commands = []*Command{}
		}
		client.WriteError(&protocol.ErrWrongNumberOfArgs{Cmd: cmd})
		return
	}

	if client.Multi && strings.ToLower(cmd) != "exec" {
		// store args for later use
		command.Args = args
		client.Commands = append(client.Commands, command)
		client.WriteProtocolReply(resp2.NewSimpleStringReply("QUEUED"))
		return
	}

	// execute command directly
	command.Fn(client, args)
}
