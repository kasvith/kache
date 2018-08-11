package srv

import (
	"bufio"
	"github.com/kasvith/kache/internal/config"
	"github.com/kasvith/kache/internal/errh"
	"log"
	"net"
	"strconv"
	"sync"
)

type Clients struct {
	ConnectedClients int
	mux              sync.Mutex
}

var ConnectedClients Clients

func (c *Clients) Increase() {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.ConnectedClients++
}

func (c *Clients) Decrease() {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.ConnectedClients--
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
		read, _ := rw.ReadString('\n')
		rw.WriteString(read)
		rw.Flush()
	}
}

func Start(config config.AppConfig) {
	addr := net.JoinHostPort(config.Host, strconv.Itoa(config.Port))
	listener, err := net.Listen("tcp", addr)

	if err != nil {
		errh.LogErrorAndExit(err, 3)
	}

	for {
		conn, err := listener.Accept()

		if err != nil {
			errh.LogError("Error on connection with", conn.RemoteAddr().String(), ":", err.Error())
			conn.Close()
			continue // we skip malformed user
		}

		// client connected
		log.Println("Connected client on", conn.RemoteAddr())
		ConnectedClients.Increase()
		log.Println(ConnectedClients.ConnectedClients, "connections are now open")

		go handleConnection(conn)
	}
}
