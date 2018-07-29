package tlist

import (
	"testing"
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
		l.HPush(string(i))
	}

	if l.Len() != 10 {
		t.Error("Assert Length is unequal expected 10, got", l.Len())
	}
}
