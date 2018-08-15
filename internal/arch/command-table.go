package arch

import (
	"github.com/kasvith/kache/internal/cmds"
	"github.com/kasvith/kache/internal/db"
	"github.com/kasvith/kache/internal/protcl"
)

type CommandFunc func(*db.DB, []string) protcl.Message

type Command struct {
	ModifyKeySpace bool
	Fn             CommandFunc
}

var CommandTable = map[string]Command{
	"get":    {ModifyKeySpace: false, Fn: cmds.Get},
	"set":    {ModifyKeySpace: true, Fn: cmds.Set},
	"exists": {ModifyKeySpace: false, Fn: cmds.Exists},
}

type DBCommand struct {
}

func getCommand(cmd string) (*Command, protcl.Err) {
	if v, ok := CommandTable[cmd]; ok {
		return &v, nil
	}

	return nil, &protcl.ErrUnknownCommand{}
}

func (DBCommand) Execute(db *db.DB, cmd string, args []string) protcl.Message {
	command, err := getCommand(cmd)
	if err != nil {
		return protcl.Message{Rep: nil, Err: err}
	}

	return command.Fn(db, args)
}
