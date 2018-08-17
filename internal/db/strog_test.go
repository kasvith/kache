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

package db

import (
	"github.com/kasvith/kache/pkg/testsuite"
	"github.com/kasvith/kache/pkg/types/list"
	"reflect"
	"testing"
)

func TestReflection(t *testing.T) {
	dn := NewDataNode(TypeList, -1, list.New())

	ty := reflect.TypeOf(dn.Value).String()
	testsuite.AssertEqual(t, "*list.TList", ty)

	reflect.ValueOf(dn.Value).MethodByName("HPush").Call([]reflect.Value{reflect.ValueOf("a")})
	v := reflect.ValueOf(dn.Value).MethodByName("Len").Call([]reflect.Value{})
	testsuite.AssertEqual(t, 1, v[0].Interface())
	reflect.ValueOf(dn.Value).MethodByName("HPush").Call([]reflect.Value{reflect.ValueOf("b")})
	v = reflect.ValueOf(dn.Value).MethodByName("Len").Call([]reflect.Value{})
	testsuite.AssertEqual(t, 2, v[0].Interface())
}
