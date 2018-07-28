package tlist

import (
	"container/list"
	"errors"
	"sync"
)

// TList Linked list represenation in the memoery, it's thread safe
type TList struct {
	list *list.List
	mux  sync.Mutex
}

// New Creates a new List
func New() *TList {
	return &TList{list: list.New()}
}

// HPush Inserts an item to the head of the list
func (list *TList) HPush(val ...string) error {
	list.mux.Lock()
	defer list.mux.Unlock()

	if len(val) == 0 {
		return errors.New("No items to insert")
	} else if len(val) == 1 {
		list.list.PushFront(val[0])
		return nil
	} else {
		newList := buildValueList(true, val...)
		list.list.PushFrontList(newList)
		return nil
	}
}

// TPush Inserts an item to the tail of the list
func (list *TList) TPush(val ...string) error {
	list.mux.Lock()
	defer list.mux.Unlock()

	if len(val) == 0 {
		return errors.New("No items to insert")
	} else if len(val) == 1 {
		list.list.PushBack(val[0])
		return nil
	} else {
		newList := buildValueList(false, val...)
		list.list.PushBackList(newList)
		return nil
	}
}

func buildValueList(front bool, val ...string) *list.List {
	l := list.New()

	if front == true {
		for _, v := range val {
			l.PushFront(v)
		}

		return l
	}

	for _, v := range val {
		l.PushBack(v)
	}

	return l
}

// Head Gets head of the list
func (list *TList) Head() *list.Element {
	list.mux.Lock()
	defer list.mux.Unlock()

	return list.list.Front()
}

// Tail Gets tail of the list
func (list *TList) Tail() *list.Element {
	list.mux.Lock()
	defer list.mux.Unlock()

	return list.list.Back()
}
