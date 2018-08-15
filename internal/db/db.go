package db

import (
	"fmt"
	"sync"
)

type DB struct {
	file map[string]*DataNode
	mux  sync.Mutex
}

type KeyNotFoundError struct {
	key string
}

func (e *KeyNotFoundError) Error() string {
	return fmt.Sprintf("%s not found", e.key)
}

func NewDB() *DB {
	return &DB{file: make(map[string]*DataNode)}
}

func (db *DB) Get(key string) (*DataNode, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	if v, ok := db.file[key]; ok {
		return v, nil
	}

	return nil, &KeyNotFoundError{key: key}
}

func (db *DB) Set(key string, val *DataNode) {
	db.mux.Lock()
	defer db.mux.Unlock()

	db.file[key] = val
}

func (db *DB) Del(keys []string) int {
	db.mux.Lock()
	defer db.mux.Unlock()

	del := 0
	for _, k := range keys {
		if _, ok := db.file[k]; ok {
			delete(db.file, k)
			del++
		}
	}

	return del
}

func (db *DB) Exists(key string) int {
	db.mux.Lock()
	defer db.mux.Unlock()

	if _, ok := db.file[key]; ok {
		return 1
	}

	return 0
}
