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

// CommandTable holds all commands that are supported by kache
var CommandTable = map[string]Command{
	// server
	"ping":  {ModifyKeySpace: false, Fn: Ping, MinArgs: 0, MaxArgs: 1},
	"multi": {ModifyKeySpace: true, Fn: Multi, MinArgs: 0, MaxArgs: 0},
	"exec":  {ModifyKeySpace: true, Fn: Exec, MinArgs: 0, MaxArgs: 0},

	// key space
	"exists": {ModifyKeySpace: false, Fn: Exists, MinArgs: 1, MaxArgs: 1},
	"del":    {ModifyKeySpace: true, Fn: Del, MinArgs: 1, MaxArgs: -1},
	"keys":   {ModifyKeySpace: false, Fn: Keys, MinArgs: 0, MaxArgs: 0},
	"expire": {ModifyKeySpace: false, Fn: Expire, MinArgs: 2, MaxArgs: 2},

	// strings
	"get":  {ModifyKeySpace: false, Fn: Get, MinArgs: 1, MaxArgs: 1},
	"set":  {ModifyKeySpace: true, Fn: Set, MinArgs: 2, MaxArgs: 2},
	"incr": {ModifyKeySpace: true, Fn: Incr, MinArgs: 1, MaxArgs: 1},
	"decr": {ModifyKeySpace: true, Fn: Decr, MinArgs: 1, MaxArgs: 1},
}
