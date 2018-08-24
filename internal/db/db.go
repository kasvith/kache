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

package db

import (
	"fmt"
	"sync"
)

// DB holds a thread safe struct for store data
type DB struct {
	file map[string]*DataNode
	mux  sync.Mutex
}

// KeyNotFoundError has the key which was not able to found in a DB
type KeyNotFoundError struct {
	key string
}

func (e *KeyNotFoundError) Error() string {
	return fmt.Sprintf("%s not found", e.key)
}

// NewDB returns a new *DB
func NewDB() *DB {
	return &DB{file: make(map[string]*DataNode)}
}

// Get the value of a key
func (db *DB) Get(key string) (*DataNode, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	if v, ok := db.file[key]; ok {
		return v, nil
	}

	return nil, &KeyNotFoundError{key: key}
}

// Set the value of a key
func (db *DB) Set(key string, val *DataNode) {
	db.mux.Lock()
	db.file[key] = val
	db.mux.Unlock()
}

// GetIfNotSet will try to get the key if not will set it to a given value
func (db *DB) GetIfNotSet(key string, val *DataNode) (value *DataNode, found bool) {
	db.mux.Lock()
	defer db.mux.Unlock()

	if v, found := db.file[key]; found {
		return v, true
	}
	db.file[key] = val

	return val, false
}

// Del will delete keys
func (db *DB) Del(keys []string) int {
	db.mux.Lock()
	del := 0
	for _, k := range keys {
		if _, ok := db.file[k]; ok {
			delete(db.file, k)
			del++
		}
	}
	db.mux.Unlock()

	return del
}

// Exists finds the existancy of a key
func (db *DB) Exists(key string) int {
	db.mux.Lock()
	defer db.mux.Unlock()
	if _, ok := db.file[key]; ok {
		return 1
	}

	return 0
}
