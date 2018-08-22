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
	"testing"

	"github.com/kasvith/kache/pkg/testsuite"
)

func TestIntegerReply_Reply(t *testing.T) {
	rep := NewIntegerReply(1995)
	testsuite.AssertEqual(t, ":1995\r\n", rep.Reply())
}

func TestNewSimpleStringReply(t *testing.T) {
	rep := NewSimpleStringReply("foo")
	testsuite.AssertEqual(t, "+foo\r\n", rep.Reply())
}

func TestBulkStringReply_Reply(t *testing.T) {
	// test for nil strings
	nilRep := NewBulkStringReply(true, "")
	testsuite.AssertEqual(t, "$-1\r\n", nilRep.Reply())

	// test for normal strings
	rep := NewBulkStringReply(false, "bar")
	testsuite.AssertEqual(t, "$3\r\nbar\r\n", rep.Reply())
}

func TestArrayReply_Reply(t *testing.T) {
	// nil array
	nilRep := NewArrayReply(true, []Reply{})
	testsuite.AssertEqual(t, "*-1\r\n", nilRep.Reply())

	// normal array
	replies := []Reply{NewBulkStringReply(false, "foo"), NewBulkStringReply(false, "foobar")}
	arrRep := NewArrayReply(false, replies)
	targetStr := "*2\r\n$3\r\nfoo\r\n$6\r\nfoobar\r\n"
	testsuite.AssertEqual(t, targetStr, arrRep.Reply())

	// array of arrays
	arr1 := NewArrayReply(false, []Reply{NewBulkStringReply(false, "foo"), NewIntegerReply(1)})
	arr2 := NewArrayReply(false, []Reply{NewSimpleStringReply("bar")})
	arrOfArrReps := []Reply{arr1, arr2}
	arrOfArrs := NewArrayReply(false, arrOfArrReps)
	targetRep := "*2\r\n*2\r\n$3\r\nfoo\r\n:1\r\n*1\r\n+bar\r\n"
	testsuite.AssertEqual(t, targetRep, arrOfArrs.Reply())
}
