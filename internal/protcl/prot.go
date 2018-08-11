package protcl

type Reply interface {
	Reply() string
}

type Message struct {
	Rep Reply
	Err error
}
