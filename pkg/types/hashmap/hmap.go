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

package hashmap

import (
	"errors"
	"strconv"
	"sync"
	"unicode/utf8"
)

// HashMap is a thread safe hashmap with RWMutex
type HashMap struct {
	m   map[string]string
	mux *sync.RWMutex
}

// New *HashMap is created
func New() *HashMap {
	return &HashMap{m: make(map[string]string), mux: &sync.RWMutex{}}
}

// Set key value tuple
// It will return 0 when key was already in the map, 1 when new key inserted
func (m *HashMap) Set(key, value string) int {
	m.mux.Lock()
	defer m.mux.Unlock()

	if value, found := m.m[key]; found {
		m.m[key] = value
		return 0
	}

	m.m[key] = value
	return 1
}

// Setx only sets when key is not in map
func (m *HashMap) Setx(key, value string) int {
	m.mux.Lock()
	defer m.mux.Unlock()

	if _, found := m.m[key]; found {
		return 0
	}

	m.m[key] = value
	return 1
}

// SetBulk sets key,value tuples
func (m *HashMap) SetBulk(fields []string) (string, error) {
	m.mux.Lock()
	defer m.mux.Unlock()

	if len(fields)%2 != 0 || len(fields) == 0 {
		return "", errors.New("invalid number of arguments")
	}

	for i := 0; i < len(fields); i += 2 {
		m.m[fields[i]] = fields[i+1]
	}

	return "OK", nil
}

// Get a value from a key
func (m *HashMap) Get(key string) string {
	m.mux.RLock()
	defer m.mux.RUnlock()

	return m.m[key]
}

// GetBulk returns an array of values for given keys, with nil values
func (m *HashMap) GetBulk(keys []string) []string {
	m.mux.RLock()
	defer m.mux.RUnlock()

	results := make([]string, len(keys))

	for i := 0; i < len(keys); i++ {
		results[i] = m.m[keys[i]]
	}

	return results
}

// Keys will return all keys in map
func (m *HashMap) Keys() []string {
	m.mux.RLock()
	defer m.mux.RUnlock()

	keys := make([]string, len(m.m))
	i := 0
	for k := range m.m {
		keys[i] = k
		i++
	}

	return keys
}

// Vals get all values of map
func (m *HashMap) Vals() []string {
	m.mux.RLock()
	defer m.mux.RUnlock()

	vals := make([]string, len(m.m))
	i := 0
	for _, v := range m.m {
		vals[i] = v
		i++
	}

	return vals
}

// Fields returns all key,value tuples as string array
func (m *HashMap) Fields() []string {
	m.mux.RLock()
	defer m.mux.RUnlock()

	paris := make([]string, len(m.m)*2)
	i := 0

	for key, val := range m.m {
		paris[i] = key
		i++
		paris[i] = val
		i++
	}
	return paris
}

// Delete set of keys
func (m *HashMap) Delete(keys []string) int {
	m.mux.Lock()
	defer m.mux.Unlock()

	deleted := 0
	for _, key := range keys {
		if _, found := m.m[key]; found {
			delete(m.m, key)
			deleted++
		}
	}

	return deleted
}

// Exists checks for key existancy
func (m *HashMap) Exists(key string) int {
	m.mux.RLock()
	defer m.mux.RUnlock()

	_, found := m.m[key]

	if found {
		return 1
	}

	return 0
}

// IncrementBy int amount for a key holding an int
func (m *HashMap) IncrementBy(key string, amount int) (int, error) {
	m.mux.Lock()
	defer m.mux.Unlock()

	target, found := m.m[key]

	if !found {
		m.m[key] = strconv.Itoa(amount)
		return amount, nil
	}

	// we have a hit
	targetVal, err := strconv.Atoi(target)

	if err != nil {
		return 0, errors.New("invalid type, excepted integer")
	}

	newVal := targetVal + amount
	m.m[key] = strconv.Itoa(newVal)

	return newVal, nil
}

// IncrementByFloat amount for a given key
func (m *HashMap) IncrementByFloat(key string, amount float64) (float64, error) {
	m.mux.Lock()
	defer m.mux.Unlock()

	target, found := m.m[key]

	if !found {
		m.m[key] = strconv.FormatFloat(amount, 'f', 6, 64)
		return amount, nil
	}

	// we have a hit
	targetVal, err := strconv.ParseFloat(target, 64)

	if err != nil {
		return 0, errors.New("invalid type, excepted float")
	}

	newVal := targetVal + amount
	m.m[key] = strconv.FormatFloat(newVal, 'f', 6, 64)

	return newVal, nil
}

// Len is length of the HashMap
func (m *HashMap) Len() int {
	m.mux.RLock()
	defer m.mux.RUnlock()

	return len(m.m)
}

// FLen is field length which returns the length of given key in bytes
func (m *HashMap) FLen(key string) int {
	m.mux.RLock()
	defer m.mux.RUnlock()

	return utf8.RuneCountInString(m.m[key])
}
