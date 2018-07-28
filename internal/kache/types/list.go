package types

import (
	"container/list"
	"errors"
	"sync"
)

// List Linked list represenation in the memoery, it's thread safe
type List struct {
	list *list.List
	mux  sync.Mutex
}

// New Creates a new List
func New() *List {
	return &List{list: list.New()}
}

// Push Inserts an item to the head of the list
func (list *List) Push(val ...string) error {
	list.mux.Lock()
	defer list.mux.Unlock()

	if len(val) == 0 {
		return errors.New("No items to insert")
	} else if len(val) == 1 {
		list.list.PushFront(val)
		return nil
	} else {
		newList := buildValueList(true, val...)
		list.list.PushFrontList(newList)
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
