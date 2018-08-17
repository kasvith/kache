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

type ErrWrongType struct {
}

func (ErrWrongType) Err() ErrorReply {
	return ErrorReply{Err: "Invalid operation against key holding invalid type of value", Prefix: WRONGTYP}
}

type ErrGeneric struct {
	Error error
}

func (e *ErrGeneric) Err() ErrorReply {
	return ErrorReply{Err: e.Error.Error(), Prefix: ERR}
}

type ErrWrongNumberOfArgs struct {
	Cmd string
}

func (e *ErrWrongNumberOfArgs) Err() ErrorReply {
	return ErrorReply{Err: e.Cmd + " has wrong number of arguments", Prefix: ERR}
}

type ErrUnknownCommand struct {
	cmd string
}

func (e *ErrUnknownCommand) Err() ErrorReply {
	return ErrorReply{Err: "unknown command " + e.cmd, Prefix: ERR}
}
