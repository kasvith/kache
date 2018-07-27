package types

// Integer Holds a typical integer value
type Integer int

// UInt8 holds a uint8
type UInt8 uint8

type String string

type HashMap struct {
	Capacity uint8
	data     map[string]string
}
