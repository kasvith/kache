package protcl

const (
	REP_SIMPLE_STRING = "+"
	REP_INTEGER = ":"
	REP_BULKSTRING = "$"
	REP_ERROR = "-"
)

type Reply interface {
	Reply() string
}

type Err interface {
	Err() ErrorReply
}

type Message struct {
	Rep Reply
	Err Err
}

const (
	WRONGTYP = "WRONGTYP"
	ERR = "ERR"
)
