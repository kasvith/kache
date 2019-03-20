package protocol

// Reply is used to transform a data type into a protocol compatible representation
type Reply interface {
	// ToBytes returns a byte representation of given type
	ToBytes() []byte
}
