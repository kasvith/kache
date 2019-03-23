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
	"fmt"
	"strconv"

	"github.com/kasvith/kache/internal/resp/resp2"

	"github.com/kasvith/kache/internal/db"
	"github.com/kasvith/kache/internal/protocol"
	"github.com/kasvith/kache/pkg/util"
)

// Get will find the value of a given string key and return it
func Get(cl *Client, args []string) {
	fmt.Println(cl.Protocol)
	val, err := cl.Database.Get(args[0])
	if err != nil {
		switch cl.Protocol {
		case RESP2, RESP3:
			cl.WriteProtocolReply(resp2.NewBulkStringReply(true, ""))
		}
		return
	}

	if val.Type != db.TypeString {
		cl.WriteError(&protocol.ErrWrongType{})
	}

	// TODO handle RESP3 value also
	cl.WriteProtocolReply(resp2.NewBulkStringReply(false, util.ToString(val.Value)))
}

// Set will create a new string key value pair
func Set(cl *Client, args []string) {
	key := args[0]
	val := args[1]

	cl.Database.Set(key, db.NewDataNode(db.TypeString, -1, val))

	cl.WriteProtocolReply(resp2.NewSimpleStringReply("OK"))
}

// Incr will increment a given string key by 1
// If key not found it will be set to 0 and will do operation
// If key type is invalid it will return an error
func Incr(cl *Client, args []string) {
	accumulateBy(cl, args[0], 1, true)
}

// Decr will decrement a given string key by 1
// If key not found it will be set to 0 and will do operation
// If key type is invalid it will return an error
func Decr(cl *Client, args []string) {
	accumulateBy(cl, args[0], -1, true)
}

// accumulateBy will accumulate the value of key by given amount
func accumulateBy(cl *Client, key string, v int, incr bool) {
	val, found := cl.Database.GetIfNotSet(key, db.NewDataNode(db.TypeString, -1, strconv.Itoa(v)))

	if !found {
		cl.WriteInteger(v)
	}

	if val.Type != db.TypeString {
		cl.WriteProtocolReply(resp2.NewErrorReply(&protocol.ErrWrongType{}))
	}

	i, err := strconv.Atoi(util.ToString(val.Value))

	if err != nil {
		cl.WriteError(&protocol.ErrCastFailedToInt{Val: val.Value})
	}

	var n int
	if incr {
		n = i + v
	} else {
		n = i - v
	}

	cl.Database.Set(key, db.NewDataNode(db.TypeString, -1, strconv.Itoa(n)))

	cl.WriteInteger(n)
}
