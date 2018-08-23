package arch

import (
	"testing"

	testifyAssert "github.com/stretchr/testify/assert"

	"github.com/kasvith/kache/internal/db"
	"github.com/kasvith/kache/internal/protcl"
)

func TestCommandArgsCountValidator(t *testing.T) {
	assert := testifyAssert.New(t)
	cmd := &DBCommand{}
	db := db.NewDB()

	// ping at most 1
	{
		assert.Nil(cmd.Execute(db, "ping", nil).Err)
		assert.Nil(cmd.Execute(db, "ping", []string{"1"}).Err)
		assert.Equal(&protcl.ErrWrongNumberOfArgs{Cmd: "ping"}, cmd.Execute(db, "ping", []string{"1", "2"}).Err)
	}

	// del at least 1
	{
		assert.Equal(&protcl.ErrWrongNumberOfArgs{Cmd: "del"}, cmd.Execute(db, "del", nil).Err)
		assert.Nil(cmd.Execute(db, "del", []string{"1"}).Err)
		assert.Nil(cmd.Execute(db, "del", []string{"1", "2"}).Err)
	}

	// set equal 2
	{
		assert.Equal(&protcl.ErrWrongNumberOfArgs{Cmd: "set"}, cmd.Execute(db, "set", nil).Err)
		assert.Equal(&protcl.ErrWrongNumberOfArgs{Cmd: "set"}, cmd.Execute(db, "set", []string{"1"}).Err)
		assert.Nil(cmd.Execute(db, "set", []string{"1", "2"}).Err)
		assert.Equal(&protcl.ErrWrongNumberOfArgs{Cmd: "set"}, cmd.Execute(db, "set", []string{"1", "2", "3"}).Err)
	}

	// equal 1: get exists incr decr
	for _, command := range []string{"get", "exists", "incr", "decr"} {
		assert.Equal(&protcl.ErrWrongNumberOfArgs{Cmd: command}, cmd.Execute(db, command, nil).Err)
		assert.Nil(cmd.Execute(db, command, []string{"1"}).Err)
		assert.Equal(&protcl.ErrWrongNumberOfArgs{Cmd: command}, cmd.Execute(db, command, []string{"1", "2"}).Err)
	}
}
