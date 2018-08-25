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

// Set is a string set data structure implemented with a hashmap
type Set struct {
	m   map[string]int
	mux *sync.RWMutex
}

// New creates a new Set
func New() *Set {
	return &Set{m: make(map[string]int), mux: &sync.RWMutex{}}
}

// NewFromSlice creates a new set from a string slice
func NewFromSlice(data []string) *Set {
	m := make(map[string]int)
	for _, value := range data {
		m[value] = 1
	}

	return &Set{m: m, mux: &sync.RWMutex{}}
}

// getMap gets a copy of underlying map from Set
func (set *Set) getMap() map[string]int {
	set.mux.RLock()

	m := make(map[string]int)

	for key := range set.m {
		m[key] = 1
	}

	set.mux.RUnlock()
	return m
}

// Add keys to the set
func (set *Set) Add(keys []string) int {
	set.mux.Lock()

	added := 0

	for _, key := range keys {
		if _, found := set.m[key]; !found {
			set.m[key] = 1
			added++
		}
	}

	set.mux.Unlock()
	return added
}

// Card is the number of elements in the set
func (set *Set) Card() int {
	set.mux.RLock()
	length := len(set.m)
	set.mux.RUnlock()
	return length
}

// elems extracts keys from a map
func elems(m map[string]int) []string {
	res := make([]string, len(m))
	i := 0
	for key := range m {
		res[i] = key
		i++
	}

	return res
}

// Elems get all elements in the set
func (set *Set) Elems() []string {
	set.mux.RLock()
	elements := elems(set.m)
	set.mux.RUnlock()
	return elements
}

// duplicateMap is a utility function to duplicate a map
func duplicateMap(m map[string]int) map[string]int {
	dup := make(map[string]int)
	for key, value := range m {
		dup[key] = value
	}

	return dup
}

// Diff calculates set difference
func (set *Set) Diff(sets []Set) []string {
	set.mux.RLock()
	dup := duplicateMap(set.m)
	set.mux.RUnlock()

	for i := 0; i < len(sets); i++ {
		for _, key := range sets[i].Elems() {
			delete(dup, key)
		}
	}

	return elems(dup)
}

// DiffS calculate set diff and returns a Set
func (set *Set) DiffS(sets []Set) *Set {
	set.mux.RLock()
	dup := duplicateMap(set.m)
	set.mux.RUnlock()
	for i := 0; i < len(sets); i++ {
		for _, key := range sets[i].Elems() {
			delete(dup, key)
		}
	}

	return &Set{m: dup, mux: &sync.RWMutex{}}
}

// Exists find a key is in set
func (set *Set) Exists(key string) int {
	set.mux.RLock()

	if _, found := set.m[key]; found {
		set.mux.RUnlock()
		return 1
	}

	set.mux.RUnlock()
	return 0
}

// Intersection calculates the intersection of sets
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
		allIntersected := true
		for i := 0; i < len(sets); i++ {
			if sets[i].Exists(v) == 0 {
				allIntersected = false
				break
			}
		}
		if allIntersected {
			results = append(results, v)
		}
	}

	return results
}

// IntersectionS calculates intersection and returns a Set
func IntersectionS(sets []Set) *Set {
	return NewFromSlice(Intersection(sets))
}

// Move a key from one set to another
func Move(key string, src, dest *Set) int {
	src.mux.Lock()

	if _, found := src.m[key]; found {
		delete(src.m, key)
		dest.mux.Lock()
		dest.m[key] = 1
		dest.mux.Unlock()
		src.mux.Unlock()
		return 1
	}

	src.mux.Unlock()
	return 0
}

// Delete an element from a set
func (set *Set) Delete(keys []string) int {
	set.mux.Lock()

	deleted := 0
	for _, key := range keys {
		if _, ok := set.m[key]; ok {
			delete(set.m, key)
			deleted++
		}
	}

	set.mux.Unlock()
	return deleted
}

// Union of stes
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

// UnionS is Union which returns a Set
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

	return &Set{m: m, mux: &sync.RWMutex{}}
}

// TODO implement pop and randomelement
