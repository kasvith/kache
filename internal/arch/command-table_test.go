package arch

import (
	"testing"

	testifyAssert "github.com/stretchr/testify/assert"

	"github.com/kasvith/kache/internal/db"
	"github.com/kasvith/kache/internal/protcl"
)

func testRespError(t *testing.T, err error, resp3 *protcl.Resp3) {
	assert := testifyAssert.New(t)
	if err == nil {
		assert.NotEqual(protcl.Resp3SimpleError, resp3.Type)
		assert.NotEqual(protcl.Resp3BolbError, resp3.Type)
		return
	}

	assert.Equal(err.Error(), resp3.Str)
}

// TestCommandArgsCountValidator will validate the command args count field
func TestCommandArgsCountValidator(t *testing.T) {
	cmd := &DBCommand{}
	db := db.NewDB()

	// ping at most 1
	{
		testRespError(t, nil, cmd.Execute(db, "ping", nil))
		testRespError(t, nil, cmd.Execute(db, "ping", []string{"1"}))
		testRespError(t, &protcl.ErrWrongNumberOfArgs{Cmd: "ping"}, cmd.Execute(db, "ping", []string{"1", "2"}))
	}

	// del at least 1
	{
		testRespError(t, &protcl.ErrWrongNumberOfArgs{Cmd: "del"}, cmd.Execute(db, "del", nil))
		testRespError(t, nil, cmd.Execute(db, "del", []string{"1"}))
		testRespError(t, nil, cmd.Execute(db, "del", []string{"1", "2"}))
	}

	// set equal 2
	{
		testRespError(t, &protcl.ErrWrongNumberOfArgs{Cmd: "set"}, cmd.Execute(db, "set", nil))
		testRespError(t, &protcl.ErrWrongNumberOfArgs{Cmd: "set"}, cmd.Execute(db, "set", []string{"1"}))
		testRespError(t, nil, cmd.Execute(db, "set", []string{"1", "2"}))
		testRespError(t, &protcl.ErrWrongNumberOfArgs{Cmd: "set"}, cmd.Execute(db, "set", []string{"1", "2", "3"}))
	}

	// equal 1: get exists incr decr
	for _, command := range []string{"get", "exists", "incr", "decr"} {
		testRespError(t, &protcl.ErrWrongNumberOfArgs{Cmd: command}, cmd.Execute(db, command, nil))
		testRespError(t, nil, cmd.Execute(db, command, []string{"1"}))
		testRespError(t, &protcl.ErrWrongNumberOfArgs{Cmd: command}, cmd.Execute(db, command, []string{"1", "2"}))
	}
}
