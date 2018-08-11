package db

type DataType int

const (
	TypeString  DataType = 1
	TypeList    DataType = 2
	TypeHashMap DataType = 3
	TypeSet     DataType = 4
)

type DataNode struct {
	Type  DataType
	ExpiresAt int
	Value interface{}
}

func NewDataNode(t DataType,exp int, val interface{}) *DataNode {
	return &DataNode{t,exp, val}
}
