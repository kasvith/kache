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
	"sync"

	"github.com/kasvith/kache/internal/klogs"
)

// ConnectedClients represents connected clients
var ConnectedClients = Clients{clients: make(map[string]*Client)}

// Clients is the struct for keep track of clients
type Clients struct {
	clients map[string]*Client
	mux     sync.RWMutex
}

// Add a client to clients
func (clients *Clients) Add(client *Client) {
	clients.mux.Lock()
	clients.clients[client.RemoteAddr().String()] = client
	clients.mux.Unlock()

	// log about new client
	klogs.Logger.Debug("Connected:", client.RemoteAddr().String())
}

// Remove a client from clients
func (clients *Clients) Remove(remoteAddr string) {
	clients.mux.Lock()
	delete(clients.clients, remoteAddr)
	clients.mux.Unlock()

	// log about new client
	klogs.Logger.Debug("Disconnected:", remoteAddr)
}

// Count of the connected clients
func (clients *Clients) Count() (num int) {
	clients.mux.RLock()
	num = len(clients.clients)
	clients.mux.RUnlock()
	return
}

// LogClientCount to the logger
func (clients *Clients) LogClientCount() {
	count := clients.Count()
	if count > 0 {
		klogs.Logger.Infof("%d connections are open", count)
	}
}

// Close all clients
func (clients *Clients) Close() error {
	clients.mux.Lock()
	for _, c := range clients.clients {
		c.Connection.Close()
	}
	clients.mux.Unlock()
	return nil
}
