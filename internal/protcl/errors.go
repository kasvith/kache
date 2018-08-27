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

import "fmt"

// RecoverableError indicates an error is recoverable, if it's not it leads for critical actions like disconnecting a
// client
type RecoverableError interface {
	// Recoverable whether error is recoverable or not
	Recoverable() bool
}

// ErrCastFailedToInt for cast fails to ints
type ErrCastFailedToInt struct {
	Val interface{}
}

func (e *ErrCastFailedToInt) Error() string {
	return fmt.Sprintf("%s: error casting %v to int", PrefixErr, e.Val)
}

// Recoverable whether error is recoverable or not
func (ErrCastFailedToInt) Recoverable() bool {
	return true
}

// ErrWrongType is for wrong type errors
type ErrWrongType struct {
}

func (ErrWrongType) Error() string {
	return fmt.Sprintf("%s: invalid operation against key holding invalid type of value", PrefixWrongType)
}

// Recoverable whether error is recoverable or not
func (ErrWrongType) Recoverable() bool {
	return true
}

// ErrGeneric for generic errors
type ErrGeneric struct {
	Err error
}

func (e *ErrGeneric) Error() string {
	return fmt.Sprintf("%s: %s", PrefixErr, e.Err)
}

// Recoverable whether error is recoverable or not
func (ErrGeneric) Recoverable() bool {
	return true
}

// ErrWrongNumberOfArgs for wrong argument count errors
type ErrWrongNumberOfArgs struct {
	Cmd string
}

func (e *ErrWrongNumberOfArgs) Error() string {
	return fmt.Sprintf("%s: %s has wrong number of arguments", PrefixWrongType, e.Cmd)
}

// Recoverable whether error is recoverable or not
func (ErrWrongNumberOfArgs) Recoverable() bool {
	return true
}

// ErrUnknownCommand for unknown commands
type ErrUnknownCommand struct {
	Cmd string
}

func (e *ErrUnknownCommand) Error() string {
	return fmt.Sprintf("%s: unknown command %s", PrefixErr, e.Cmd)
}

// Recoverable whether error is recoverable or not
func (ErrUnknownCommand) Recoverable() bool {
	return true
}

// ErrProtocolType is for protocol type error
type ErrProtocolType struct {
	Type byte
}

// Recoverable whether error is recoverable or not
func (e *ErrProtocolType) Recoverable() bool {
	return true
}

func (e *ErrProtocolType) Error() string {
	return fmt.Sprintf("unknown protocol type: %s", string(e.Type))
}

// ErrUnexpectString is for unexpect string error
type ErrUnexpectString struct {
	Str string
}

// Recoverable whether error is recoverable or not
func (e *ErrUnexpectString) Recoverable() bool {
	return true
}

func (e *ErrUnexpectString) Error() string {
	return fmt.Sprintf("unexpect string: %s", e.Str)
}
