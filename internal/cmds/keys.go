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

package cmds

import (
	"errors"
	"strconv"
	"time"

	"github.com/kasvith/kache/internal/db"
	"github.com/kasvith/kache/internal/protocol"
	"github.com/kasvith/kache/internal/sys"
)

// Exists will check for key existence in given db
func Exists(d *db.DB, args []string) *protocol.Resp3 {
	found := d.Exists(args[0])
	return &protocol.Resp3{Type: protocol.Resp3Number, Integer: found}
}

// Del will delete set of keys and return number of deleted keys
func Del(d *db.DB, args []string) *protocol.Resp3 {
	deleted := d.Del(args)
	return &protocol.Resp3{Type: protocol.Resp3Number, Integer: deleted}
}

// Keys will return all keys of the db as a list
func Keys(d *db.DB, args []string) *protocol.Resp3 {
	keys := d.Keys()
	return &protocol.Resp3{Type: protocol.Resp3Array, Elems: keys}
}

// Expire a key
func Expire(d *db.DB, args []string) *protocol.Resp3 {
	if v, ok := d.GetNode(args[0]); ok {
		val, err := strconv.Atoi(args[1])
		if err != nil {
			return &protocol.Resp3{Type: protocol.Resp3SimpleError, Err: &protocol.ErrCastFailedToInt{Val: args[1]}}
		}

		if val < 0 {
			return &protocol.Resp3{Type: protocol.Resp3SimpleError, Err: errors.New("invalid seconds")}
		}

		ttl := sys.GetTTL(int64(val), time.Second)
		v.SetExpiration(ttl)
		return &protocol.Resp3{Type: protocol.Resp3Number, Integer: 1}
	}

	return &protocol.Resp3{Type: protocol.Resp3Number, Integer: 0}
}
