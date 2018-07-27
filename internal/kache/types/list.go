package types

import (
	"errors"
	"sync"
)

// List Data structure
type List struct {
	Capacity uint8
	data     []string
	sync.Mutex
}

// Add an item to the list
func (list *List) Add(value interface{}) error {
	s, err := value.(string)

	if err == true {
		return errors.New("Invalid type")
	}

	list.Lock()
	defer list.Unlock()
	list.data = append(list.data, s)
	return nil
}
