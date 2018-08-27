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
	"net"

	"io"

	"github.com/kasvith/kache/internal/arch"
	"github.com/kasvith/kache/internal/db"
	"github.com/kasvith/kache/internal/klogs"
	"github.com/kasvith/kache/internal/protcl"
)

// TODO: Need refactoring this to allow multiple DBs for use

// DB is the database used
var DB = db.NewDB()

var commander = &arch.DBCommand{}

// Client represents a structure to manage connected client
type Client struct {
	// Connection for client
	Connection net.Conn

	// DB selected database for client, default to 0
	DB int
}

// RemoteAddr returns remote address of client
func (client *Client) RemoteAddr() net.Addr {
	return client.Connection.RemoteAddr()
}

// Handle the client
func (client *Client) Handle() {
	// TODO determine client type by first issued command to kache, this can improve performance

	resp3Parser := protcl.NewResp3Parser(bufio.NewReader(client.Connection))
	writer := bufio.NewWriter(client.Connection)

	for {
		command, err := resp3Parser.Commands()

		// handle any parse errors
		if err != nil {
			recoverable, isRecoverableErr := err.(protcl.RecoverableError)

			if isRecoverableErr {
				if recoverable.Recoverable() {
					// log the error, inform client continue loop
					// anything else should be sent to client with prefix PrefixErr
					klogs.Logger.Debug(client.RemoteAddr(), ": ", err.Error())
					writer.WriteString(protcl.RespError(err))
					writer.Flush()
					continue
				}
			}

			// If not recoverable or does not implement the interface, then its a critical error
			// break from the loop to close connection, well we ignore EOF in normal mode
			if err == io.EOF {
				klogs.Logger.Debug(client.RemoteAddr(), ": ", err.Error())
			} else {
				klogs.Logger.Error(client.RemoteAddr(), ": ", err.Error())
			}
			break
		}

		// executes the command
		message := commander.Execute(DB, command.Name, command.Args)

		writer.WriteString(message.ProtocolString())

		writer.Flush()
	}

	ConnectedClients.Remove(client.RemoteAddr().String())
	client.Connection.Close()
	ConnectedClients.LogClientCount()
}
