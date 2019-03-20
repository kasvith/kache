/*
 * MIT License
 *
 * Copyright (c) 2019 Kasun Vithanage
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
 *
 */

package db

import (
	"fmt"
	"sync"
)

// DB holds a thread safe struct for store data
type DB struct {
	file map[string]*DataNode
	mux  sync.RWMutex
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

// GetNode will clear the key if its expired
func (db *DB) GetNode(key string) (*DataNode, bool) {
	db.mux.RLock()

	if v, ok := db.file[key]; ok {
		db.mux.RUnlock()
		if v.IsExpired() {
			db.mux.Lock()
			delete(db.file, key)
			db.mux.Unlock()

			return nil, false
		}

		return v, true
	}

	db.mux.RUnlock()
	return nil, false
}

// Get the value of a key
func (db *DB) Get(key string) (*DataNode, error) {
	if v, ok := db.GetNode(key); ok {
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

// GetIfNotSet will try to GetNode the key if not will set it to a given value
func (db *DB) GetIfNotSet(key string, val *DataNode) (value *DataNode, found bool) {
	if v, found := db.GetNode(key); found {
		return v, true
	}

	db.mux.Lock()
	db.file[key] = val
	db.mux.Unlock()

	return val, false
}

// Del will delete keys
func (db *DB) Del(keys []string) int {
	db.mux.Lock()
	del := 0
	for _, k := range keys {
		if v, ok := db.file[k]; ok {
			// dont count already deleted keys aka expired
			if !v.IsExpired() {
				del++
			}
			delete(db.file, k)
		}
	}
	db.mux.Unlock()

	return del
}

// Exists finds the existence of a key
func (db *DB) Exists(key string) int {
	if _, ok := db.GetNode(key); ok {
		return 1
	}

	return 0
}

// Keys returns all keys of the db
func (db *DB) Keys() []string {
	db.mux.RLock()

	keys := make([]string, 0)
	for key, val := range db.file {
		if !val.IsExpired() {
			keys = append(keys, key)
		}

	}

	db.mux.RUnlock()
	return keys
}

// SetExpire time for a key
func (db *DB) SetExpire(key string, ttl int64) bool {
	if ttl < 0 && ttl != -1 {
		return false
	}

	if v, ok := db.GetNode(key); ok {
		v.SetExpiration(ttl)
	}

	return false
}
