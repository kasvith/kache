package cmds

import (
	"github.com/kasvith/kache/internal/db"
	"github.com/kasvith/kache/internal/protcl"
	"github.com/kasvith/kache/pkg/util"
)

func Get(d *db.DB, args []string) protcl.Message {
	if len(args) != 1 {
		return protcl.Message{Rep: nil, Err: &protcl.ErrInsufficientArgs{Cmd: "get"}}
	}

	val, err := d.Get(args[0])
	if err != nil {
		return protcl.Message{Rep: nil, Err: &protcl.ErrGeneric{Error: err}}
	}

	if val.Type != db.TypeString {
		return protcl.Message{Rep: nil, Err: &protcl.ErrWrongType{}}
	}

	return protcl.Message{Rep: protcl.NewBulkStringReply(false, util.ToString(val.Value)), Err: nil}
}

func Set(d *db.DB, args []string) protcl.Message {
	if len(args) != 2 {
		return protcl.Message{Rep: nil, Err: &protcl.ErrInsufficientArgs{Cmd: "set"}}
	}

	key := args[0]
	val := args[1]

	d.Set(key, db.NewDataNode(db.TypeString, -1, val))

	return protcl.Message{Rep: protcl.NewSimpleStringReply("OK"), Err: nil}
}

func Exists(d *db.DB, args []string) protcl.Message {
	if len(args) != 1 {
		return protcl.Message{Rep: nil, Err: &protcl.ErrInsufficientArgs{Cmd: "get"}}
	}
	found := d.Exists(args[0])

	return protcl.Message{Rep: protcl.NewIntegerReply(found), Err: nil}
}
