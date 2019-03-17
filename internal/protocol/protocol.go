package protocol

type Reply interface {
	ToBytes() []byte
}
