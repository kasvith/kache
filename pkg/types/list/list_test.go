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
	"errors"
	"strconv"
	"testing"

	"github.com/kasvith/kache/pkg/testsuite"
)

func TestPushSingleValue(t *testing.T) {
	l := New()

	val := "some"

	err := l.HPush([]string{val})

	if err != nil {
		t.Error(err)
	}

	head := l.Head()

	if head == nil {
		t.Error("Head is nil")
	}

	s, ok := head.Value.(string)
	if ok == false {
		t.Error("Not a string type")
	}

	if s != val {
		t.Error("Items are not same")
	}
}

func TestLen(t *testing.T) {
	l := New()

	for i := 0; i < 10; i++ {
		l.HPush([]string{strconv.Itoa(i)})
	}

	testsuite.AssertEqual(t, 10, l.Len())
}

func TestFindAtIndexHead(t *testing.T) {
	l := New()

	for i := 0; i < 10; i++ {
		l.HPush([]string{strconv.Itoa(i)})
	}

	e1 := l.findAtIndex(0)

	testsuite.AssertEqual(t, "9", e1.Value)
}

func TestFindAtIndexTail(t *testing.T) {
	l := New()

	for i := 0; i < 10; i++ {
		l.HPush([]string{strconv.Itoa(i)})
	}

	e2 := l.findAtIndex(l.Len() - 1)

	testsuite.AssertEqual(t, "0", e2.Value)
}

func TestFindAtIndexNull(t *testing.T) {
	l := New()

	for i := 0; i < 10; i++ {
		l.HPush([]string{strconv.Itoa(i)})
	}

	e2 := l.findAtIndex(100)
	testsuite.AssertNil(t, e2)
}

func TestFindAtIndexMiddle(t *testing.T) {
	l := New()

	for i := 0; i < 10; i++ {
		l.HPush([]string{strconv.Itoa(i)})
	}

	e2 := l.findAtIndex(4)

	testsuite.AssertEqual(t, "5", e2.Value)
}

func TestTList_RangeAllItems(t *testing.T) {
	l := New()
	strs := make([]string, 10)

	for i := 0; i < 10; i++ {
		strs[i] = strconv.Itoa(9 - i)
		l.HPush([]string{strconv.Itoa(i)})
	}

	res := l.Range(0, -1)
	testsuite.AssertStringSliceEqual(t, strs, res)
}

func TestTList_RangeMinusDistance(t *testing.T) {
	l := New()
	strs := make([]string, 10)

	for i := 0; i < 10; i++ {
		strs[i] = strconv.Itoa(9 - i)
		l.HPush([]string{strconv.Itoa(i)})
	}

	res := l.Range(5, 2)
	testsuite.AssertStringSliceEqual(t, []string{}, res)
}

func TestTList_RangeStopOutOfBound(t *testing.T) {
	l := New()
	strs := make([]string, 10)

	for i := 0; i < 10; i++ {
		strs[i] = strconv.Itoa(9 - i)
		l.HPush([]string{strconv.Itoa(i)})
	}

	res := l.Range(0, 100)
	testsuite.AssertStringSliceEqual(t, strs, res)
}

func TestTList_RangeStartOutOfBound(t *testing.T) {
	l := New()

	for i := 0; i < 10; i++ {
		l.HPush([]string{strconv.Itoa(i)})
	}

	res := l.Range(100, -1)
	testsuite.AssertStringSliceEqual(t, []string{}, res)
}

func TestTList_RangeTwoFromTail(t *testing.T) {
	l := New()
	strs := make([]string, 10)

	for i := 0; i < 10; i++ {
		strs[i] = strconv.Itoa(9 - i)
		l.HPush([]string{strconv.Itoa(i)})
	}

	res := l.Range(-2, -1)
	testsuite.AssertStringSliceEqual(t, []string{"1", "0"}, res)
}

func TestTList_RangeTwoFromMiddle(t *testing.T) {
	l := New()
	strs := make([]string, 10)

	for i := 0; i < 10; i++ {
		strs[i] = strconv.Itoa(9 - i)
		l.HPush([]string{strconv.Itoa(i)})
	}

	res := l.Range(2, 5)
	testsuite.AssertStringSliceEqual(t, []string{"7", "6", "5", "4"}, res)
}

func TestTList_RangeOneFromHead(t *testing.T) {
	l := New()
	strs := make([]string, 10)

	for i := 0; i < 10; i++ {
		strs[i] = strconv.Itoa(9 - i)
		l.HPush([]string{strconv.Itoa(i)})
	}

	res := l.Range(0, 0)
	testsuite.AssertStringSliceEqual(t, []string{"9"}, res)
}

func TestTList_RangeOneFromTail(t *testing.T) {
	l := New()
	strs := make([]string, 10)

	for i := 0; i < 10; i++ {
		strs[i] = strconv.Itoa(9 - i)
		l.HPush([]string{strconv.Itoa(i)})
	}

	res := l.Range(-1, -1)
	testsuite.AssertStringSliceEqual(t, []string{"0"}, res)
}

func TestTList_HPop(t *testing.T) {
	l := New()
	strs := make([]string, 10)

	for i := 0; i < 10; i++ {
		strs[i] = strconv.Itoa(9 - i)
		l.HPush([]string{strconv.Itoa(i)})
	}

	elem := l.HPop()
	testsuite.AssertEqual(t, "9", elem)

	elem = l.HPop()
	testsuite.AssertEqual(t, "8", elem)
}

func TestTList_TPop(t *testing.T) {
	l := New()

	for i := 0; i < 10; i++ {
		l.HPush([]string{strconv.Itoa(i)})
	}

	elem := l.TPop()
	testsuite.AssertEqual(t, "0", elem)

	elem = l.TPop()
	testsuite.AssertEqual(t, "1", elem)
}

func TestTList_TPush(t *testing.T) {
	l := New()
	strs := make([]string, 10)

	for i := 0; i < 10; i++ {
		strs[i] = strconv.Itoa(9 - i)
		l.TPush([]string{strconv.Itoa(i)})
	}

	elem := l.HPop()
	testsuite.AssertEqual(t, "0", elem)

	elem = l.HPop()
	testsuite.AssertEqual(t, "1", elem)
}

func TestTList_TPushListHead(t *testing.T) {
	l := New()
	strs := make([]string, 10)
	vals := make([]string, 10)

	for i := 0; i < 10; i++ {
		strs[i] = strconv.Itoa(9 - i)
		vals[i] = strconv.Itoa(i)
	}

	l.HPush(vals)

	res := l.Range(0, -1)
	testsuite.AssertStringSliceEqual(t, strs, res)

	l.HPush([]string{"0", "1"})
	res = l.Range(0, -1)
	testsuite.AssertStringSliceEqual(t, []string{"1", "0", "9", "8", "7", "6", "5", "4", "3", "2", "1", "0"}, res)
}

func TestTList_TPushListTail(t *testing.T) {
	l := New()
	strs := make([]string, 10)
	vals := make([]string, 10)

	for i := 0; i < 10; i++ {
		strs[i] = strconv.Itoa(9 - i)
		vals[i] = strconv.Itoa(i)
	}

	l.TPush(vals)

	res := l.Range(0, -1)
	testsuite.AssertStringSliceEqual(t, vals, res)

	l.TPush([]string{"0", "1"})
	res = l.Range(0, -1)
	testsuite.AssertStringSliceEqual(t, []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "1"}, res)
}

func TestTList_HPushNil(t *testing.T) {
	l := New()
	err := l.HPush([]string{})

	testsuite.ExceptError(t, errors.New("no items to insert"), err)
}

func TestTList_TPushNil(t *testing.T) {
	l := New()
	err := l.TPush([]string{})

	testsuite.ExceptError(t, errors.New("no items to insert"), err)
}

func TestTList_HPopNil(t *testing.T) {
	l := New()
	it := l.HPop()

	testsuite.AssertEqual(t, "", it)
}

func TestTList_TPopNil(t *testing.T) {
	l := New()
	it := l.TPop()

	testsuite.AssertEqual(t, "", it)
}

func TestTList_TrimNullList(t *testing.T) {
	l := New()

	l.Trim(0, 10)

	testsuite.AssertEqual(t, 0, l.Len())
}

func TestTList_TrimFullList(t *testing.T) {
	l := New()

	for i := 0; i < 10; i++ {
		l.HPush([]string{strconv.Itoa(i)})
	}

	// list now trimmed to 2 elements
	l.Trim(0, 1)
	testsuite.AssertEqual(t, 2, l.Len())
	testsuite.AssertStringSliceEqual(t, []string{"9", "8"}, l.Range(0, -1))
}

func TestTList_TrimAddMoreAndAssert(t *testing.T) {
	l := New()

	for i := 0; i < 10; i++ {
		l.HPush([]string{strconv.Itoa(i)})
	}

	// list now trimmed to 2 elements
	l.Trim(0, 1)
	testsuite.AssertEqual(t, 2, l.Len())
	testsuite.AssertStringSliceEqual(t, []string{"9", "8"}, l.Range(0, -1))

	// push new elem
	l.HPush([]string{"a"})

	// list now trimmed to 2 elements
	l.Trim(0, 1)

	// assert for 2 elements
	testsuite.AssertEqual(t, 2, l.Len())

	// check list updated
	testsuite.AssertStringSliceEqual(t, []string{"a", "9"}, l.Range(0, -1))
}

func TestTList_TrimMinusStart(t *testing.T) {
	l := New()

	for i := 0; i < 10; i++ {
		l.HPush([]string{strconv.Itoa(i)})
	}

	// list now trimmed to 2 elements
	l.Trim(-100, 1)
	testsuite.AssertEqual(t, 2, l.Len())
	testsuite.AssertStringSliceEqual(t, []string{"9", "8"}, l.Range(0, -1))
}

func TestTList_TrimExceededStopLimit(t *testing.T) {
	l := New()

	for i := 0; i < 10; i++ {
		l.HPush([]string{strconv.Itoa(i)})
	}

	// list now trimmed to 2 elements
	l.Trim(0, 10000)
	testsuite.AssertEqual(t, 10, l.Len())
	testsuite.AssertStringSliceEqual(t, []string{"9", "8", "7", "6", "5", "4", "3", "2", "1", "0"}, l.Range(0, -1))
}

func TestTList_TrimStopLesserThanStart(t *testing.T) {
	l := New()

	for i := 0; i < 10; i++ {
		l.HPush([]string{strconv.Itoa(i)})
	}

	// list now trimmed to 2 elements
	l.Trim(4, 2)
	testsuite.AssertEqual(t, 10, l.Len())
	testsuite.AssertStringSliceEqual(t, []string{"9", "8", "7", "6", "5", "4", "3", "2", "1", "0"}, l.Range(0, -1))
}

func TestTList_TrimHead(t *testing.T) {
	l := New()

	for i := 0; i < 10; i++ {
		l.HPush([]string{strconv.Itoa(i)})
	}

	// list now trimmed to 2 elements
	l.Trim(2, -1)
	testsuite.AssertEqual(t, 8, l.Len())
	testsuite.AssertStringSliceEqual(t, []string{"7", "6", "5", "4", "3", "2", "1", "0"}, l.Range(0, -1))
}
