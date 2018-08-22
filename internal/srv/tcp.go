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
	"io"
	"net"
	"os"
	"strconv"

	"github.com/kasvith/kache/internal/arch"
	"github.com/kasvith/kache/internal/config"
	"github.com/kasvith/kache/internal/db"
	"github.com/kasvith/kache/internal/klogs"
	"github.com/kasvith/kache/internal/protcl"
)

var DB = db.NewDB()
var dbCommand = &arch.DBCommand{}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		reader := protcl.NewReader(conn)

		// TODO determine client type by first issued command to kache, this can improve performance
		w := bufio.NewWriter(conn)
		command, err := reader.ParseMessage()

		if err != nil {
			// if eof stop now
			if err == io.EOF {
				break
			}

			// anything else should be sent to client with prefix ERR
			klogs.Logger.Debug(conn.RemoteAddr(), ": ", err.Error())
			w.WriteString(protcl.RespError(err))
			w.Flush()
			continue
		}

		message := dbCommand.Execute(DB, command.Name, command.Args)

		if message.Err == nil {
			w.WriteString(message.RespReply())
		} else {
			w.WriteString(protcl.RespError(message.Err))
		}

		w.Flush()
	}

	ConnectedClients.logOnDisconnect(conn)
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
		ConnectedClients.logOnConnect(conn)

		go handleConnection(conn)
	}
}
