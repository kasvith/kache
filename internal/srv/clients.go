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
	"net"
	"sync"

	"github.com/kasvith/kache/internal/klogs"
)

var ConnectedClients Clients

type Clients struct {
	numClients int
	mux        sync.Mutex
}

func (c *Clients) increase() {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.numClients++
}

func (c *Clients) decrease() {
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

func (c *Clients) logOnDisconnect(conn net.Conn) {
	klogs.Logger.Info("disconnected client from ", conn.RemoteAddr())
	c.decrease()
	logOpenedClients()
}

func (c *Clients) logOnConnect(conn net.Conn) {
	klogs.Logger.Info("connected client from ", conn.RemoteAddr())
	c.increase()
	logOpenedClients()
}
