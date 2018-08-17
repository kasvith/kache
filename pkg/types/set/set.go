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

package set

import (
	"sync"
)

type Set struct {
	m   map[string]int
	mux *sync.Mutex
}

func New() *Set {
	return &Set{m: make(map[string]int), mux: &sync.Mutex{}}
}

func NewFromSlice(data []string) *Set {
	m := make(map[string]int)
	for _, value := range data {
		m[value] = 1
	}

	return &Set{m: m, mux: &sync.Mutex{}}
}

func (set *Set) getMap() map[string]int {
	set.mux.Lock()
	defer set.mux.Unlock()

	m := make(map[string]int)

	for key := range set.m {
		m[key] = 1
	}

	return m
}

func (set *Set) Add(keys []string) int {
	set.mux.Lock()
	defer set.mux.Unlock()

	added := 0

	for _, key := range keys {
		if _, found := set.m[key]; !found {
			set.m[key] = 1
			added++
		}
	}

	return added
}

// Card this will return number of elements in the set
func (set *Set) Card() int {
	set.mux.Lock()
	defer set.mux.Unlock()

	return len(set.m)
}

func elems(m map[string]int) []string {
	res := make([]string, len(m))
	i := 0
	for key := range m {
		res[i] = key
		i++
	}

	return res
}

func (set *Set) Elems() []string {
	set.mux.Lock()
	defer set.mux.Unlock()

	return elems(set.m)
}

func duplicateMap(m map[string]int) map[string]int {
	dup := make(map[string]int)
	for key, value := range m {
		dup[key] = value
	}

	return dup
}

func (set *Set) Diff(sets []Set) []string {
	set.mux.Lock()
	defer set.mux.Unlock()

	dup := duplicateMap(set.m)

	for i := 0; i < len(sets); i++ {
		for _, key := range sets[i].Elems() {
			delete(dup, key)
		}
	}

	return elems(dup)
}

func (set *Set) DiffS(sets []Set) *Set {
	set.mux.Lock()
	defer set.mux.Unlock()
	dup := duplicateMap(set.m)

	for i := 0; i < len(sets); i++ {
		for _, key := range sets[i].Elems() {
			delete(dup, key)
		}
	}

	return &Set{m: dup, mux: &sync.Mutex{}}
}

func (set *Set) Exists(key string) int {
	set.mux.Lock()
	defer set.mux.Unlock()

	if _, found := set.m[key]; found {
		return 1
	}

	return 0
}

func Intersection(sets []Set) []string {
	minSetIdx := 0
	for i := 0; i < len(sets); i++ {
		if sets[minSetIdx].Card() > sets[i].Card() {
			minSetIdx = i
		}
	}

	minSet := &sets[minSetIdx]
	sets = append(sets[:minSetIdx], sets[minSetIdx+1:]...)

	// an empty set means empty intersection
	if minSet.Card() == 0 {
		return []string{}
	}

	// iterate through minimum set to find the intersection
	results := make([]string, 0)
	for _, v := range minSet.Elems() {
		for i := 0; i < len(sets); i++ {
			if sets[i].Exists(v) == 1 {
				results = append(results, v)
			}
		}
	}

	return results
}

func IntersectionS(sets []Set) *Set {
	return NewFromSlice(Intersection(sets))
}

func Move(key string, src, dest *Set) int {
	src.mux.Lock()
	defer src.mux.Unlock()
	if _, found := src.m[key]; found {
		delete(src.m, key)
		dest.mux.Lock()
		dest.m[key] = 1
		dest.mux.Unlock()
		return 1
	}

	return 0
}

func (set *Set) Delete(keys []string) int {
	set.mux.Lock()
	defer set.mux.Unlock()

	deleted := 0

	for _, key := range keys {
		if _, ok := set.m[key]; ok {
			delete(set.m, key)
			deleted++
		}
	}

	return deleted
}

func Union(sets []Set) []string {
	maxSetIdx := 0
	for i := 0; i < len(sets); i++ {
		if sets[maxSetIdx].Card() < sets[i].Card() {
			maxSetIdx = i
		}
	}

	m := sets[maxSetIdx].getMap()
	sets = append(sets[:maxSetIdx], sets[maxSetIdx+1:]...)

	for i := 0; i < len(sets); i++ {
		for _, v := range sets[i].Elems() {
			m[v] = 1
		}
	}

	return elems(m)
}

func UnionS(sets []Set) *Set {
	maxSetIdx := 0
	for i := 0; i < len(sets); i++ {
		if sets[maxSetIdx].Card() < sets[i].Card() {
			maxSetIdx = i
		}
	}

	m := sets[maxSetIdx].getMap()
	sets = append(sets[:maxSetIdx], sets[maxSetIdx+1:]...)

	for i := 0; i < len(sets); i++ {
		for _, v := range sets[i].Elems() {
			m[v] = 1
		}
	}

	return &Set{m: m, mux: &sync.Mutex{}}
}

// TODO implement pop and randomelement
