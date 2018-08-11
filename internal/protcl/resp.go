package protcl

import "fmt"

// IntegerReply Represents an integer reply
type IntegerReply struct {
	 value int
}

// Reply method for integers
func (rep *IntegerReply) Reply() string {
	return fmt.Sprintf(":%d\r\n", rep.value)
}

// SimpleStringReply Binary unsafe strings
type SimpleStringReply struct {
	value string
}

// Reply method for integers
func (rep *SimpleStringReply) Reply() string {
	return fmt.Sprintf("+%s\r\n", rep.value)
}


