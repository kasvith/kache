/*
 * MIT License
 *
 * Copyright (c)  2018 Kasun Vithanage
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package srv

import (
	"bufio"
	"github.com/kasvith/kache/internal/arch"
	"github.com/kasvith/kache/internal/config"
	"github.com/kasvith/kache/internal/db"
	"github.com/kasvith/kache/internal/klogs"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Clients struct {
	numClients int
	mux        sync.Mutex
}

var ConnectedClients Clients

var DB = db.NewDB()
var dbCommand = &arch.DBCommand{}

func (c *Clients) Increase() {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.numClients++
}

func (c *Clients) Decrease() {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.numClients--
}

func logOpenedClients() {
	if ConnectedClients.numClients > 0 {
		klogs.Logger.Info(ConnectedClients.numClients, " connections are now open")
		return
	}

	klogs.Logger.Info("no connections are now open")
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
		read, err := rw.ReadString('\n')

		if err != nil && err == io.EOF {
			break
		}

		strs := strings.Split(strings.TrimSpace(read), " ")

		if len(strs) == 0 {
			rw.Flush()
			continue
		}

		rep := dbCommand.Execute(DB, strings.ToLower(strs[0]), strs[1:])

		if rep.Err == nil {
			rw.WriteString(rep.Rep.Reply())
		} else {
			err := rep.Err.Err()
			rw.WriteString(err.Error())
		}

		rw.Flush()
	}

	klogs.Logger.Info("Disconnected client from ", conn.RemoteAddr())
	ConnectedClients.Decrease()
	logOpenedClients()
}

func Start(config config.AppConfig) {
	addr := net.JoinHostPort(config.Host, strconv.Itoa(config.Port))
	listener, err := net.Listen("tcp", addr)

	if err != nil {
		klogs.Logger.Fatalf("error binding to port %d is already in use", config.Port)
		os.Exit(3)
	}

	klogs.Logger.Infof("application is ready to accept connections on port %d", config.Port)

	for {
		conn, err := listener.Accept()

		if err != nil {
			klogs.Logger.Error("Error on connection with", conn.RemoteAddr().String(), ":", err.Error())
			conn.Close()
			continue // we skip malformed user
		}

		// client connected
		klogs.Logger.Info("Connected client on ", conn.RemoteAddr())
		ConnectedClients.Increase()
		logOpenedClients()

		go handleConnection(conn)
	}
}
