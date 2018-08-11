package cmds

import (
	"github.com/kasvith/kache/internal/db"
	"github.com/kasvith/kache/internal/protcl"
)

func Get(d *db.DB, key string) protcl.Message {
	val, err := d.Get(key)
	if err != nil {
		return protcl.Message{Rep: nil, Err: err}
	}

	if val.Type != db.TypeString {
		return protcl.Message{Rep: nil, Err: &protcl.ErrorWrongType{}}
	}

	return protcl.Message{Rep: &protcl., Err: nil}
}

func Set(d *db.DB, args []string) protcl.Message {
	if len(args) < 2 {
		return protcl.Message{Rep: nil, Err: &protcl.ErrorInsufficientArgs{Cmd: "set"}}
	}

	key := args[0]
	val := args[1]

	d.Set(key, db.NewDataNode(db.TypeString, -1, val))

	return protcl.Message{Rep: "OK", Err: nil}
}

func Exists(d *db.DB, key string) protcl.Message {
	found := d.Exists(key)

	return protcl.Message{Rep: found, Err: nil}
}
