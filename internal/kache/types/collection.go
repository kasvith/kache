package types

// Collection data strucure
type Collection interface {
	Add(interface{}) error
	Delete(interface{}) error
	Update(interface{}) error
	Find(interface{}) (interface{}, error)
	Count() uint8
	Clear() error
	Sort(func(a interface{}, b interface{}))
}
