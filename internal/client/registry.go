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
	"sync"

	"github.com/kasvith/kache/internal/klogs"
)

// ConnectedClients represents connected clients
var ConnectedClients = Registry{clients: make(map[string]*Client)}

// Registry maintains the registry of connected clients
type Registry struct {
	clients map[string]*Client
	mux     sync.RWMutex
}

// Add a client to clients
func (cr *Registry) Add(client *Client) {
	cr.mux.Lock()
	cr.clients[client.RemoteAddr().String()] = client
	cr.mux.Unlock()

	// log about new client
	klogs.Logger.Debug("Connected:", client.RemoteAddr().String())
}

// Remove a client from clients
func (cr *Registry) Remove(remoteAddr string) {
	cr.mux.Lock()
	delete(cr.clients, remoteAddr)
	cr.mux.Unlock()

	// log about new client
	klogs.Logger.Debug("Disconnected:", remoteAddr)
}

// Count of the connected clients
func (cr *Registry) Count() (num int) {
	cr.mux.RLock()
	num = len(cr.clients)
	cr.mux.RUnlock()
	return
}

// LogClientCount to the logger
func (cr *Registry) LogClientCount() {
	count := cr.Count()
	if count > 0 {
		klogs.Logger.Debugf("%d connections are open", count)
	}
}

// Close all clients
func (cr *Registry) Close() error {
	cr.mux.Lock()
	for _, c := range cr.clients {
		c.Connection.Close()
	}
	cr.mux.Unlock()
	return nil
}
