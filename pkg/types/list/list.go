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
	"container/list"
	"errors"
	"sync"

	"github.com/kasvith/kache/pkg/util"
)

// TList Linked list representation in the memory, it's thread safe
type TList struct {
	list *list.List
	mux  *sync.RWMutex
}

// New Creates a new List
func New() *TList {
	return &TList{list: list.New(), mux: &sync.RWMutex{}}
}

func buildValueList(front bool, val []string) *list.List {
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

// HPush Inserts an item to the head of the list
func (list *TList) HPush(val []string) error {
	if len(val) == 0 {
		return errors.New("no items to insert")
	}

	list.mux.Lock()
	defer list.mux.Unlock()

	if len(val) == 1 {
		list.list.PushFront(val[0])
		return nil
	} else {
		newList := buildValueList(true, val)
		list.list.PushFrontList(newList)
		return nil
	}
}

// TPush Inserts an item to the tail of the list
func (list *TList) TPush(val []string) error {
	if len(val) == 0 {
		return errors.New("no items to insert")
	}

	list.mux.Lock()
	defer list.mux.Unlock()

	if len(val) == 1 {
		list.list.PushBack(val[0])
		return nil
	} else {
		newList := buildValueList(false, val)
		list.list.PushBackList(newList)
		return nil
	}
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

	return util.ToString(list.list.Remove(list.Head()))
}

// TPop pops out the element from tail
func (list *TList) TPop() string {
	list.mux.Lock()
	defer list.mux.Unlock()

	if list.Tail() == nil {
		return ""
	}

	return util.ToString(list.list.Remove(list.Tail()))
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

// Range will output set of keys based on index query
func (list *TList) Range(start, stop int) []string {
	list.mux.RLock()
	defer list.mux.RUnlock()

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

	res := make([]string, dist+1)

	for i, j, e := start, 0, list.findAtIndex(start); e != nil && i <= stop; i, j, e = i+1, j+1, e.Next() {
		res[j] = util.ToString(e.Value)
	}

	return res[:]
}

func (list *TList) Trim(start, stop int) {
	list.mux.Lock()

	start = list.convertPos(start)
	stop = list.convertPos(stop)

	if start > list.Len()-1 {
		list.mux.Unlock()
		return
	}

	if start < 0 {
		start = 0
	}

	if stop > list.Len() {
		stop = list.Len() - 1
	}

	if stop < start {
		list.mux.Unlock()
		return
	}

	itemsForRemoveFromHead := start
	itemsForRemoveFromTail := list.Len() - 1 - stop

	list.mux.Unlock()

	// TODO We need a more optimized method
	for i := itemsForRemoveFromHead; i > 0; i-- {
		list.HPop()
	}

	for i := itemsForRemoveFromTail; i > 0; i-- {
		list.TPop()
	}
}
