package tlist

import (
	"container/list"
	"errors"
	"sync"
)

// TList Linked list representation in the memory, it's thread safe
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
		return errors.New("no items to insert")
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
		return errors.New("no items to insert")
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
	return list.list.Front()
}

// Tail Gets tail of the list
func (list *TList) Tail() *list.Element {
	return list.list.Back()
}

// Len returns the length of the list
func (list *TList) Len() int {
	return list.list.Len()
}

// HPop pops out the element from head
func (list *TList) HPop() string {
	list.mux.Lock()
	defer list.mux.Unlock()

	if list.Head() == nil {
		return ""
	}

	return ToString(list.list.Remove(list.Head()))
}

// TPop pops out the element from tail
func (list *TList) TPop() string {
	list.mux.Lock()
	defer list.mux.Unlock()

	if list.Tail() == nil {
		return ""
	}

	return ToString(list.list.Remove(list.Tail()))
}

func (list *TList) convertPos(pos int) int {
	// get the real index from the negative one
	if pos < 0 {
		pos = pos + list.Len()
	}

	return pos
}

func (list *TList) findAtIndex(pos int) *list.Element {
	if pos < 0 || pos > list.Len() {
		return nil
	}
	
	if pos == 0 {
		return list.Head()
	}

	if pos == list.Len()-1 {
		return list.Tail()
	}

	for i, e := 1, list.Head().Next(); e != nil && i <= pos; i, e = i+1, e.Next() {
		if i == pos {
			return e
		}
	}

	return nil
}

// ToString Convert an interface to string
func ToString(i interface{}) string {
	if s, ok := i.(string); ok {
		return s
	}

	return ""
}

// Range will output set of keys based on index query
func (list *TList) Range(start, stop int) []string {
	list.mux.Lock()
	defer list.mux.Unlock()

	start = list.convertPos(start)
	stop = list.convertPos(stop)

	if start > list.Len()-1 || start < 0 {
		return []string{}
	}

	if stop > list.Len()-1 {
		stop = list.Len() - 1
	}

	dist := stop - start

	if dist < 0 {
		return []string{}
	}

	res := make([]string, dist + 1)


	for i, j , e := start, 0, list.findAtIndex(start); e != nil && i <= stop; i,j, e = i+1, j + 1, e.Next() {
		res[j] = ToString(e.Value)
	}

	return res[:]
}
