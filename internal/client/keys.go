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
	"errors"

	"strconv"
	"time"

	"github.com/kasvith/kache/internal/resp/resp2"

	"github.com/kasvith/kache/internal/protocol"
	"github.com/kasvith/kache/internal/sys"
)

// Exists will check for key existence in given db
func Exists(cl *Client, args []string) {
	found := cl.Database.Exists(args[0])
	cl.WriteInteger(found)
}

// Del will delete set of keys and return number of deleted keys
func Del(cl *Client, args []string) {
	deleted := cl.Database.Del(args)
	cl.WriteInteger(deleted)
}

// Keys will return all keys of the db as a list
func Keys(cl *Client, args []string) {
	keys := cl.Database.Keys()

	switch cl.Protocol {
	case RESP2:
		// TODO do a proper RESP3
	case RESP3:
		arr := make([]protocol.Reply, len(keys))
		for i := 0; i < len(keys); i++ {
			arr[i] = *resp2.NewBulkStringReply(false, keys[i])
		}

		cl.WriteProtocolReply(resp2.NewArrayReply(false, arr))
		break
	}
}

// Expire a key
func Expire(cl *Client, args []string) {
	if v, ok := cl.Database.GetNode(args[0]); ok {
		val, err := strconv.Atoi(args[1])
		if err != nil {
			cl.WriteError(&protocol.ErrCastFailedToInt{Val: args[1]})
		}

		if val < 0 {
			cl.WriteError(errors.New("invalid seconds"))
		}

		ttl := sys.GetTTL(int64(val), time.Second)
		v.SetExpiration(ttl)
		cl.WriteInteger(1)
	}

	cl.WriteInteger(0)
}
