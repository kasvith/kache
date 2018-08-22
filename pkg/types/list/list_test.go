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

package list

import (
	"strconv"
	"testing"

	testifyAssert "github.com/stretchr/testify/assert"
)

func TestPushSingleValue(t *testing.T) {
	assert := testifyAssert.New(t)
	l := New()

	val := "some"

	err := l.HPush([]string{val})
	assert.Nil(err)

	head := l.Head()
	assert.NotNil(head)

	s, ok := head.Value.(string)
	assert.True(ok)
	assert.Equal(val, s)
}

func TestLen(t *testing.T) {
	assert := testifyAssert.New(t)
	l := New()

	for i := 0; i < 10; i++ {
		l.HPush([]string{strconv.Itoa(i)})
	}

	assert.Equal(10, l.Len())
}

func TestFindAtIndexHead(t *testing.T) {
	assert := testifyAssert.New(t)
	l := New()

	for i := 0; i < 10; i++ {
		l.HPush([]string{strconv.Itoa(i)})
	}

	e1 := l.findAtIndex(0)
	assert.Equal("9", e1.Value)
}

func TestFindAtIndexTail(t *testing.T) {
	assert := testifyAssert.New(t)
	l := New()

	for i := 0; i < 10; i++ {
		l.HPush([]string{strconv.Itoa(i)})
	}

	e2 := l.findAtIndex(l.Len() - 1)
	assert.Equal("0", e2.Value)
}

func TestFindAtIndexNull(t *testing.T) {
	assert := testifyAssert.New(t)
	l := New()

	for i := 0; i < 10; i++ {
		l.HPush([]string{strconv.Itoa(i)})
	}

	e2 := l.findAtIndex(100)
	assert.Nil(e2)
}

func TestFindAtIndexMiddle(t *testing.T) {
	assert := testifyAssert.New(t)
	l := New()

	for i := 0; i < 10; i++ {
		l.HPush([]string{strconv.Itoa(i)})
	}

	e2 := l.findAtIndex(4)
	assert.Equal("5", e2.Value)
}

func TestTList_RangeAllItems(t *testing.T) {
	assert := testifyAssert.New(t)
	l := New()
	strs := make([]string, 10)

	for i := 0; i < 10; i++ {
		strs[i] = strconv.Itoa(9 - i)
		l.HPush([]string{strconv.Itoa(i)})
	}

	res := l.Range(0, -1)
	assert.Equal(strs, res)
}

func TestTList_RangeMinusDistance(t *testing.T) {
	assert := testifyAssert.New(t)
	l := New()
	strs := make([]string, 10)

	for i := 0; i < 10; i++ {
		strs[i] = strconv.Itoa(9 - i)
		l.HPush([]string{strconv.Itoa(i)})
	}

	res := l.Range(5, 2)
	assert.Equal([]string{}, res)
}

func TestTList_RangeStopOutOfBound(t *testing.T) {
	assert := testifyAssert.New(t)
	l := New()
	strs := make([]string, 10)

	for i := 0; i < 10; i++ {
		strs[i] = strconv.Itoa(9 - i)
		l.HPush([]string{strconv.Itoa(i)})
	}

	res := l.Range(0, 100)
	assert.Equal(strs, res)
}

func TestTList_RangeStartOutOfBound(t *testing.T) {
	assert := testifyAssert.New(t)
	l := New()

	for i := 0; i < 10; i++ {
		l.HPush([]string{strconv.Itoa(i)})
	}

	res := l.Range(100, -1)
	assert.Equal([]string{}, res)
}

func TestTList_RangeTwoFromTail(t *testing.T) {
	assert := testifyAssert.New(t)
	l := New()
	strs := make([]string, 10)

	for i := 0; i < 10; i++ {
		strs[i] = strconv.Itoa(9 - i)
		l.HPush([]string{strconv.Itoa(i)})
	}

	res := l.Range(-2, -1)
	assert.Equal([]string{"1", "0"}, res)
}

func TestTList_RangeTwoFromMiddle(t *testing.T) {
	assert := testifyAssert.New(t)
	l := New()
	strs := make([]string, 10)

	for i := 0; i < 10; i++ {
		strs[i] = strconv.Itoa(9 - i)
		l.HPush([]string{strconv.Itoa(i)})
	}

	res := l.Range(2, 5)
	assert.Equal([]string{"7", "6", "5", "4"}, res)
}

func TestTList_RangeOneFromHead(t *testing.T) {
	assert := testifyAssert.New(t)
	l := New()
	strs := make([]string, 10)

	for i := 0; i < 10; i++ {
		strs[i] = strconv.Itoa(9 - i)
		l.HPush([]string{strconv.Itoa(i)})
	}

	res := l.Range(0, 0)
	assert.Equal([]string{"9"}, res)
}

func TestTList_RangeOneFromTail(t *testing.T) {
	assert := testifyAssert.New(t)
	l := New()
	strs := make([]string, 10)

	for i := 0; i < 10; i++ {
		strs[i] = strconv.Itoa(9 - i)
		l.HPush([]string{strconv.Itoa(i)})
	}

	res := l.Range(-1, -1)
	assert.Equal([]string{"0"}, res)
}

func TestTList_HPop(t *testing.T) {
	assert := testifyAssert.New(t)
	l := New()
	strs := make([]string, 10)

	for i := 0; i < 10; i++ {
		strs[i] = strconv.Itoa(9 - i)
		l.HPush([]string{strconv.Itoa(i)})
	}

	elem := l.HPop()
	assert.Equal("9", elem)

	elem = l.HPop()
	assert.Equal("8", elem)
}

func TestTList_TPop(t *testing.T) {
	assert := testifyAssert.New(t)
	l := New()

	for i := 0; i < 10; i++ {
		l.HPush([]string{strconv.Itoa(i)})
	}

	elem := l.TPop()
	assert.Equal("0", elem)

	elem = l.TPop()
	assert.Equal("1", elem)
}

func TestTList_TPush(t *testing.T) {
	assert := testifyAssert.New(t)
	l := New()
	strs := make([]string, 10)

	for i := 0; i < 10; i++ {
		strs[i] = strconv.Itoa(9 - i)
		l.TPush([]string{strconv.Itoa(i)})
	}

	elem := l.HPop()
	assert.Equal("0", elem)

	elem = l.HPop()
	assert.Equal("1", elem)
}

func TestTList_TPushListHead(t *testing.T) {
	assert := testifyAssert.New(t)
	l := New()
	strs := make([]string, 10)
	vals := make([]string, 10)

	for i := 0; i < 10; i++ {
		strs[i] = strconv.Itoa(9 - i)
		vals[i] = strconv.Itoa(i)
	}

	l.HPush(vals)

	res := l.Range(0, -1)
	assert.Equal(strs, res)

	l.HPush([]string{"0", "1"})
	res = l.Range(0, -1)
	assert.Equal([]string{"1", "0", "9", "8", "7", "6", "5", "4", "3", "2", "1", "0"}, res)
}

func TestTList_TPushListTail(t *testing.T) {
	assert := testifyAssert.New(t)
	l := New()
	strs := make([]string, 10)
	vals := make([]string, 10)

	for i := 0; i < 10; i++ {
		strs[i] = strconv.Itoa(9 - i)
		vals[i] = strconv.Itoa(i)
	}

	l.TPush(vals)

	res := l.Range(0, -1)
	assert.Equal(vals, res)

	l.TPush([]string{"0", "1"})
	res = l.Range(0, -1)
	assert.Equal([]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "1"}, res)
}

func TestTList_HPushNil(t *testing.T) {
	assert := testifyAssert.New(t)
	l := New()

	err := l.HPush([]string{})
	assert.EqualError(err, "no items to insert")
}

func TestTList_TPushNil(t *testing.T) {
	assert := testifyAssert.New(t)
	l := New()

	err := l.TPush([]string{})
	assert.EqualError(err, "no items to insert")
}

func TestTList_HPopNil(t *testing.T) {
	assert := testifyAssert.New(t)
	l := New()

	it := l.HPop()
	assert.Equal("", it)
}

func TestTList_TPopNil(t *testing.T) {
	assert := testifyAssert.New(t)
	l := New()

	it := l.TPop()
	assert.Equal("", it)
}

func TestTList_TrimNullList(t *testing.T) {
	assert := testifyAssert.New(t)
	l := New()

	l.Trim(0, 10)
	assert.Equal(0, l.Len())
}

func TestTList_TrimFullList(t *testing.T) {
	assert := testifyAssert.New(t)
	l := New()

	for i := 0; i < 10; i++ {
		l.HPush([]string{strconv.Itoa(i)})
	}

	// list now trimmed to 2 elements
	l.Trim(0, 1)
	assert.Equal(2, l.Len())
	assert.Equal([]string{"9", "8"}, l.Range(0, -1))
}

func TestTList_TrimAddMoreAndAssert(t *testing.T) {
	assert := testifyAssert.New(t)
	l := New()

	for i := 0; i < 10; i++ {
		l.HPush([]string{strconv.Itoa(i)})
	}

	// list now trimmed to 2 elements
	l.Trim(0, 1)
	assert.Equal(2, l.Len())
	assert.Equal([]string{"9", "8"}, l.Range(0, -1))

	// push new elem
	l.HPush([]string{"a"})

	// list now trimmed to 2 elements
	l.Trim(0, 1)

	// assert for 2 elements
	assert.Equal(2, l.Len())

	// check list updated
	assert.Equal([]string{"a", "9"}, l.Range(0, -1))
}

func TestTList_TrimMinusStart(t *testing.T) {
	assert := testifyAssert.New(t)
	l := New()

	for i := 0; i < 10; i++ {
		l.HPush([]string{strconv.Itoa(i)})
	}

	// list now trimmed to 2 elements
	l.Trim(-100, 1)
	assert.Equal(2, l.Len())
	assert.Equal([]string{"9", "8"}, l.Range(0, -1))
}

func TestTList_TrimExceededStopLimit(t *testing.T) {
	assert := testifyAssert.New(t)
	l := New()

	for i := 0; i < 10; i++ {
		l.HPush([]string{strconv.Itoa(i)})
	}

	// list now trimmed to 2 elements
	l.Trim(0, 10000)
	assert.Equal(10, l.Len())
	assert.Equal([]string{"9", "8", "7", "6", "5", "4", "3", "2", "1", "0"}, l.Range(0, -1))
}

func TestTList_TrimStopLesserThanStart(t *testing.T) {
	assert := testifyAssert.New(t)
	l := New()

	for i := 0; i < 10; i++ {
		l.HPush([]string{strconv.Itoa(i)})
	}

	// list now trimmed to 2 elements
	l.Trim(4, 2)
	assert.Equal(10, l.Len())
	assert.Equal([]string{"9", "8", "7", "6", "5", "4", "3", "2", "1", "0"}, l.Range(0, -1))
}

func TestTList_TrimHead(t *testing.T) {
	assert := testifyAssert.New(t)
	l := New()

	for i := 0; i < 10; i++ {
		l.HPush([]string{strconv.Itoa(i)})
	}

	// list now trimmed to 2 elements
	l.Trim(2, -1)
	assert.Equal(8, l.Len())
	assert.Equal([]string{"7", "6", "5", "4", "3", "2", "1", "0"}, l.Range(0, -1))
}
