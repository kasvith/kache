package protcl

import "fmt"

type ErrorInsufficientArgs struct {
	Cmd string
}

func (e *ErrorInsufficientArgs) Error() string {
	return fmt.Sprintf("%s has insufficent args", e.Cmd)
}

type ErrorWrongType struct {}

func (ErrorWrongType) Error() string {
	return "wrong type"
}