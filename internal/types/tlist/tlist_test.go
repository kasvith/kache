package tlist

import (
	"testing"
	"strconv"
	"github.com/kasvith/kache/pkg/testsuite"
)

func TestPushSingleValue(t *testing.T) {
	l := New()

	val := "some"

	err := l.HPush(val)

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
		l.HPush(strconv.Itoa(i))
	}

	if l.Len() != 10 {
		t.Error("Assert Length is unequal expected 10, got", l.Len())
	}
}

func TestRange(t *testing.T) {
	l := New()

	for i := 0; i < 10; i++ {
		l.HPush(strconv.Itoa(i))
	}
}

func TestToString(t *testing.T) {
	var i interface{}
	i = "str"

	testsuite.AssertEqual(t, i, "str")
}

func TestFindAtIndexHead(t *testing.T) {
	l := New()

	for i := 0; i < 10; i++ {
		l.HPush(strconv.Itoa(i))
	}

	e1 := l.findAtIndex(0)

	testsuite.AssertEqual(t, e1.Value, "9")
}


func TestFindAtIndexTail(t *testing.T) {
	l := New()

	for i := 0; i < 10; i++ {
		l.HPush(strconv.Itoa(i))
	}

	e2 := l.findAtIndex(l.Len() - 1)

	testsuite.AssertEqual(t, e2.Value, "0")
}

func TestFindAtIndexNull(t *testing.T)  {
	l := New()

	for i := 0; i < 10; i++ {
		l.HPush(strconv.Itoa(i))
	}

	e2 := l.findAtIndex(100)
	testsuite.AssertNil(t, e2)
}

func TestFindAtIndexMiddle(t *testing.T) {
	l := New()

	for i := 0; i < 10; i++ {
		l.HPush(strconv.Itoa(i))
	}

	e2 := l.findAtIndex(4)

	testsuite.AssertEqual(t, e2.Value, "5")
}

func TestTList_RangeAllItems(t *testing.T) {
	l := New()
	strs := make([]string, 10)

	for i := 0; i < 10; i++ {
		strs[i] = strconv.Itoa(9 - i)
		l.HPush(strconv.Itoa(i))
	}

	res := l.Range(0, -1)
	testsuite.AssertStringSliceEqual(t, strs, res)
}

func TestTList_RangeTwoFromTail(t *testing.T) {
	l := New()
	strs := make([]string, 10)

	for i := 0; i < 10; i++ {
		strs[i] = strconv.Itoa(9 - i)
		l.HPush(strconv.Itoa(i))
	}

	res := l.Range(-2, -1)
	testsuite.AssertStringSliceEqual(t, []string{"1", "0"}, res)
}

func TestTList_RangeTwoFromMiddle(t *testing.T) {
	l := New()
	strs := make([]string, 10)

	for i := 0; i < 10; i++ {
		strs[i] = strconv.Itoa(9 - i)
		l.HPush(strconv.Itoa(i))
	}

	res := l.Range(2, 5)
	testsuite.AssertStringSliceEqual(t, []string{"7", "6", "5", "4"}, res)
}


func TestTList_RangeOneFromHead(t *testing.T) {
	l := New()
	strs := make([]string, 10)

	for i := 0; i < 10; i++ {
		strs[i] = strconv.Itoa(9 - i)
		l.HPush(strconv.Itoa(i))
	}

	res := l.Range(0, 0)
	testsuite.AssertStringSliceEqual(t, []string{"9"}, res)
}

func TestTList_RangeOneFromTail(t *testing.T) {
	l := New()
	strs := make([]string, 10)

	for i := 0; i < 10; i++ {
		strs[i] = strconv.Itoa(9 - i)
		l.HPush(strconv.Itoa(i))
	}

	res := l.Range(-1, -1)
	testsuite.AssertStringSliceEqual(t, []string{"0"}, res)
}