package protcl

type ErrWrongType struct {
}

func (ErrWrongType) Err() ErrorReply {
	return ErrorReply{Err: "Invalid operation against key holding invalid type of value", Prefix: WRONGTYP}
}

type ErrGeneric struct {
	Error error
}

func (e *ErrGeneric) Err() ErrorReply {
	return ErrorReply{Err: e.Error.Error(), Prefix: ERR}
}

type ErrInsufficientArgs struct {
	Cmd string
}

func (e *ErrInsufficientArgs) Err() ErrorReply {
	return ErrorReply{Err: e.Cmd + " has insufficient arguments", Prefix: ERR}
}

type ErrUnknownCommand struct {
	cmd string
}

func (e *ErrUnknownCommand) Err() ErrorReply {
	return ErrorReply{Err: "unknown command " + e.cmd, Prefix: ERR}
}
