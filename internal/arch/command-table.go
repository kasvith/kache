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

package arch

import (
	"github.com/kasvith/kache/internal/cmds"
	"github.com/kasvith/kache/internal/db"
	"github.com/kasvith/kache/internal/protcl"
)

type CommandFunc func(*db.DB, []string) *protcl.Message

type Command struct {
	ModifyKeySpace bool
	Fn             CommandFunc
	MinArgs        int // 0
	MaxArgs        int // -1 ~ +inf, -1 mean infinite
}

var CommandTable = map[string]Command{
	// server
	"ping": {ModifyKeySpace: false, Fn: cmds.Ping, MinArgs: 0, MaxArgs: 1},

	// key space
	"exists": {ModifyKeySpace: false, Fn: cmds.Exists, MinArgs: 1, MaxArgs: 1},
	"del":    {ModifyKeySpace: true, Fn: cmds.Del, MinArgs: 1, MaxArgs: -1},

	// strings
	"get":  {ModifyKeySpace: false, Fn: cmds.Get, MinArgs: 1, MaxArgs: 1},
	"set":  {ModifyKeySpace: true, Fn: cmds.Set, MinArgs: 2, MaxArgs: 2},
	"incr": {ModifyKeySpace: true, Fn: cmds.Incr, MinArgs: 1, MaxArgs: 1},
	"decr": {ModifyKeySpace: true, Fn: cmds.Decr, MinArgs: 1, MaxArgs: 1},
}

type DBCommand struct {
}

func getCommand(cmd string) (*Command, error) {
	if v, ok := CommandTable[cmd]; ok {
		return &v, nil
	}

	return nil, &protcl.ErrUnknownCommand{Cmd: cmd}
}

// Execute executes a single command on the given database
func (DBCommand) Execute(db *db.DB, cmd string, args []string) *protcl.Message {
	command, err := getCommand(cmd)
	if err != nil {
		return protcl.NewMessage(nil, err)
	}

	if argsLen := len(args); (command.MinArgs > 0 && argsLen < command.MinArgs) || (command.MaxArgs != -1 && argsLen > command.MaxArgs) {
		return protcl.NewMessage(nil, &protcl.ErrWrongNumberOfArgs{Cmd: cmd})
	}

	return command.Fn(db, args)
}
