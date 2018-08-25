package cli

import (
	"strconv"
	"testing"
	"time"

	testifyAssert "github.com/stretchr/testify/assert"

	"github.com/kasvith/kache/internal/config"
	"github.com/kasvith/kache/internal/klogs"
	"github.com/kasvith/kache/internal/srv"
)

func initTestServerClient(t *testing.T) {
	testPort := 9999

	conf := config.AppConfig{
		Host:               "127.0.0.1",
		Port:               testPort,
		LogType:            "default",
		MaxMultiBulkLength: config.DefaultMaxMultiBulkLength,
	}
	klogs.InitLoggers(conf)
	go srv.Start(conf)

	for i := 0; i < 3; i++ {
		if err := Dial("127.0.0.1:" + strconv.Itoa(testPort)); err == nil {
			return
		}
		time.Sleep(time.Second)
	}

	t.Fatalf("connect to server failed")
}

func runTestSendRecv(t *testing.T, send, recv string) {
	assert := testifyAssert.New(t)

	assert.Nil(c.write(send))
	resp, err := c.parseResp()
	assert.Nil(err)
	assert.Equal(recv, resp)
}

func TestCli(t *testing.T) {
	assert := testifyAssert.New(t)
	initTestServerClient(t)

	// ping
	runTestSendRecv(t, "ping", `"PONG"`)

	// strings
	{
		// get not found
		runTestSendRecv(t, "get a", "(error) ERR: a not found")

		// set
		runTestSendRecv(t, "set a 1", `"OK"`)
		runTestSendRecv(t, "set b 2", `"OK"`)

		// get exist
		runTestSendRecv(t, "get a", "1")

		// incr decr
		runTestSendRecv(t, "incr b", "(integer) 3")
		runTestSendRecv(t, "incr b", "(integer) 4")
		runTestSendRecv(t, "decr b", "(integer) 3")
		runTestSendRecv(t, "decr b", "(integer) 2")
	}

	// key space
	{
		// keys
		assert.Nil(c.write("keys"))
		resp, err := c.parseResp()
		assert.Nil(err)
		assert.Contains([]string{"1) b\n2) a\n", "1) a\n2) b\n"}, resp)

		// exists
		runTestSendRecv(t, "exists a", "(integer) 1")
		runTestSendRecv(t, "exists not-found", "(integer) 0")

		// del
		runTestSendRecv(t, "del a", "(integer) 1")
		runTestSendRecv(t, "del a", "(integer) 0")
		runTestSendRecv(t, "exists a", "(integer) 0")
	}
}
