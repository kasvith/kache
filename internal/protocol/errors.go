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

package protocol

import "fmt"

const (
	// PrefixWrongType WRONGTYP
	PrefixWrongType = "WRONGTYP"

	// PrefixErr ERR
	PrefixErr = "ERR"
)

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

// ErrConvertType is for convert value to another type fail
type ErrConvertType struct {
	Type  string
	Err   error
	Value interface{}
}

// Recoverable whether error is recoverable or not
func (e *ErrConvertType) Recoverable() bool {
	return true
}

func (e *ErrConvertType) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("convert %v to %s fail", e.Value, e.Type)
	}
	return fmt.Sprintf("convert %v to %s fail, because of %s", e.Value, e.Type, e.Err)
}

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

// ErrUnknownProtocol is used to indicate user that server cannot understand the protocl
type ErrUnknownProtocol struct {
}

func (ErrUnknownProtocol) Error() string {
	return "unknown protocol"
}

// Recoverable whether error is recoverable or not
func (ErrUnknownProtocol) Recoverable() bool {
	return true
}

// ErrExecAbortTransaction is used to indicate whether an error is occurred while preparing a multi transaction
type ErrExecAbortTransaction struct {
}

// Recoverable whether error is recoverable or not
func (ErrExecAbortTransaction) Recoverable() bool {
	return true
}

func (ErrExecAbortTransaction) Error() string {
	return "EXECABORT Transaction discarded because of previous errors."
}

// ErrExecWithoutMulti is used to indicate that multi is not enabled
type ErrExecWithoutMulti struct {
}

// Recoverable whether error is recoverable or not
func (ErrExecWithoutMulti) Recoverable() bool {
	return true
}

func (ErrExecWithoutMulti) Error() string {
	return "ERR EXEC without MULTI"
}
