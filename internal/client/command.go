package client

import (
	"fmt"

	"github.com/kasvith/kache/internal/protocol"
)

// Command holds a command structure which is used to execute a kache command
type Command struct {
	ModifyKeySpace bool
	Fn             CommandFunc
	MinArgs        int // 0
	MaxArgs        int // -1 ~ +inf, -1 mean infinite
}

// CommandFunc holds a function signature which can be used as a command.
type CommandFunc func(*Client, []string)

// DBCommand is a command that executes on a given db
type DBCommand struct {
}

// GetCommand will fetch the command from command table
func GetCommand(cmd string) (*Command, error) {
	if v, ok := CommandTable[cmd]; ok {
		return &v, nil
	}

	return nil, &protocol.ErrUnknownCommand{Cmd: cmd}
}

// Execute a single command on the given database with args
func Execute(client *Client, cmd string, args []string) {
	command, err := GetCommand(cmd)
	if err != nil {
		fmt.Println(err)
		client.WriteError(err)
		return
	}

	if argsLen := len(args); (command.MinArgs > 0 && argsLen < command.MinArgs) || (command.MaxArgs != -1 && argsLen > command.MaxArgs) {
		client.WriteError(&protocol.ErrWrongNumberOfArgs{Cmd: cmd})
		return
	}

	command.Fn(client, args)
}
