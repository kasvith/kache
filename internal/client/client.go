/*
 * MIT License
 *
 * Copyright (c) 2019 Kasun Vithanage
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
 *
 */

package client

import (
	"bufio"
	"net"

	"github.com/kasvith/kache/internal/resp/resp2"
	"github.com/kasvith/kache/internal/wire"

	"io"

	"github.com/kasvith/kache/internal/db"
	"github.com/kasvith/kache/internal/klogs"
	"github.com/kasvith/kache/internal/protocol"
)

// TODO: Need refactoring this to allow multiple DBs for use

// DB is the database used
var dbase = db.NewDB()

const (
	// RESP2 represents protocol version resp2
	RESP2 = "resp2"
	// RESP3 represents protocol version resp3
	RESP3 = "resp3"
)

// Client represents a structure to manage connected client
type Client struct {
	// Connection for client
	Connection net.Conn

	// DB selected database for client, default to 0
	Database *db.DB

	// Parser is used for parsing a request
	Parser protocol.CommandParser

	// Protocol
	Protocol string

	// Multi will indicate that client is in multi mode or not
	Multi bool

	// MultiError indicates whether a command failed during a transaction
	MultiError bool

	// Commands store a list of queued commands in a multi transaction
	Commands []*Command

	// Writer is used to write out data to client connection
	*bufio.Writer
}

// NewClient creates a new client object
// Note all clients will be initialized to use RESP2 as the default reply protocol
// This can be changed in future
func NewClient(conn net.Conn) *Client {
	return &Client{Connection: conn, Protocol: RESP2, Writer: bufio.NewWriter(conn), Database: dbase}
}

// RemoteAddr returns remote address of client
func (client *Client) RemoteAddr() net.Addr {
	return client.Connection.RemoteAddr()
}

// Handle the client
func (client *Client) Handle() {
	err := client.detectParser()

	if err != nil {
		klogs.Logger.Error(client.RemoteAddr(), ": ", err.Error())
		client.logAndRemove()
		return
	}

	for {
		command, err := client.Parser.Parse()

		// handle any parse errors
		if err != nil {
			recoverable, isRecoverableErr := err.(protocol.RecoverableError)

			if isRecoverableErr {
				if recoverable.Recoverable() {
					// log the error, inform client continue loop
					// anything else should be sent to client with prefix PrefixErr
					klogs.Logger.Debug(client.RemoteAddr(), ": ", err.Error())

					client.WriteError(err)
					continue
				}
			}

			// If not recoverable or does not implement the interface, then its a critical error
			// break from the loop to close connection, well we ignore EOF in normal mode
			if err != io.EOF {
				klogs.Logger.Debug(client.RemoteAddr(), ": ", err.Error())
			}
			break
		}

		// executes the command
		Execute(client, command.Name, command.Args)
	}

	client.logAndRemove()
}

func (client *Client) logAndRemove() {
	ConnectedClients.Remove(client.RemoteAddr().String())
	_ = client.Connection.Close()
	ConnectedClients.LogClientCount()
}

func (client *Client) detectParser() error {
	reader := bufio.NewReader(client.Connection)
	b, err := reader.ReadByte()
	if err != nil {
		return err
	}
	err = reader.UnreadByte()
	if err != nil {
		return err
	}

	switch b {
	case resp2.TypeArray:
		// we have resp2
		client.Parser = resp2.NewParser(reader)
	default:
		// use wire protocol
		client.Parser = wire.NewParser(reader)
	}

	return nil
}

// WriteError will write an error message to the connection
func (client *Client) WriteError(err error) {
	switch client.Protocol {
	case RESP2, RESP3:
		client.WriteProtocolReply(resp2.NewErrorReply(err))
		break
	}
}

// WriteInteger will write an int value to the client
func (client *Client) WriteInteger(n int) {
	switch client.Protocol {
	case RESP2, RESP3:
		client.WriteProtocolReply(resp2.NewIntegerReply(n))
		break
	}
}

// WriteProtocolReply will write a protocol reply
func (client *Client) WriteProtocolReply(reply protocol.Reply) {
	// ok we are clear to send
	_, err := client.Write(reply.ToBytes())
	if err != nil {
		client.logAndRemove()
		return
	}
	err = client.Flush()
	if err != nil {
		client.logAndRemove()
		return
	}
}
