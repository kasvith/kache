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

package hashmap

import (
	"strconv"
	"testing"

	"github.com/kasvith/kache/pkg/testsuite"
)

func TestHashMap_Set(t *testing.T) {
	hm := New()

	rep := hm.Set("mykey", "myval")
	testsuite.AssertEqual(t, 1, len(hm.m))
	testsuite.AssertEqual(t, 1, rep)

	rep = hm.Set("mykey", "updated")
	testsuite.AssertEqual(t, 1, len(hm.m))
	testsuite.AssertEqual(t, 0, rep)

	rep = hm.Set("mykey1", "myval")
	testsuite.AssertEqual(t, 2, len(hm.m))
	testsuite.AssertEqual(t, 1, rep)
}

func TestHashMap_Get(t *testing.T) {
	hm := New()

	hm.Set("mykey", "myval")

	mykey := hm.Get("mykey")

	testsuite.AssertEqual(t, "myval", mykey)
}

func TestHashMap_GetNullKey(t *testing.T) {
	hm := New()

	hm.Set("mykey", "myval")

	mykey := hm.Get("mykey1")

	testsuite.AssertEqual(t, "", mykey)
}

func TestHashMap_KeysNullMap(t *testing.T) {
	hm := New()

	keys := hm.Keys()

	testsuite.AssertStringSliceEqual(t, []string{}, keys)
}

func TestHashMap_Keys(t *testing.T) {
	hm := New()

	elements := make([]string, 10)
	for i := 0; i < 10; i++ {
		hm.Set("key"+strconv.Itoa(i), strconv.Itoa(i))
		elements[i] = "key" + strconv.Itoa(i)
	}

	keys := hm.Keys()

	testsuite.AssertEqual(t, 10, len(keys))
	testsuite.ContainsElements(t, elements, keys)
}

func TestHashMap_Fields(t *testing.T) {
	hm := New()
	arr := make([]string, 20)

	for i := 0; i < 10; i++ {
		hm.Set("key"+strconv.Itoa(i), strconv.Itoa(i))
		arr[i] = "key" + strconv.Itoa(i)
		arr[i*2] = strconv.Itoa(i)
	}

	res := hm.Fields()

	testsuite.ContainsElements(t, arr, res)
}

func TestHashMap_Delete(t *testing.T) {
	hm := New()

	hm.Set("key1", "val1")
	hm.Set("key2", "val2")
	hm.Set("key3", "val3")

	deleted := hm.Delete([]string{"key1"})

	testsuite.AssertEqual(t, 1, deleted)

	deleted = hm.Delete([]string{"key2", "key3"})
	testsuite.AssertEqual(t, 2, deleted)

	deleted = hm.Delete([]string{"nonexistent"})
	testsuite.AssertEqual(t, 0, deleted)
}

func TestHashMap_Exists(t *testing.T) {
	hm := New()

	hm.Set("key1", "val1")

	key1Exists := hm.Exists("key1")
	testsuite.AssertEqual(t, 1, key1Exists)

	unknown := hm.Exists("nonexistent")
	testsuite.AssertEqual(t, 0, unknown)
}

func TestHashMap_IncrementBy(t *testing.T) {
	hm := New()

	// Test for non existent key
	res, err := hm.IncrementBy("counter", 1)
	testsuite.AssertEqual(t, nil, err)
	testsuite.AssertEqual(t, 1, res)

	res, err = hm.IncrementBy("counter", 4)
	testsuite.AssertEqual(t, nil, err)
	testsuite.AssertEqual(t, 5, res)

	hm.Set("key", "val")
	res, err = hm.IncrementBy("key", 1)
	testsuite.AssertEqual(t, "invalid type, excepted integer", err.Error())
	testsuite.AssertEqual(t, 0, res)
}

func TestHashMap_IncrementByFloat(t *testing.T) {
	hm := New()

	// Test for non existent key
	res, err := hm.IncrementByFloat("counter", 10)
	testsuite.AssertEqual(t, nil, err)
	testsuite.AssertEqual(t, float64(10), res)

	fl, err := hm.IncrementByFloat("counter", 0.5)
	testsuite.AssertEqual(t, nil, err)
	testsuite.AssertEqual(t, 10.5, fl)

	fl, err = hm.IncrementByFloat("counter", -5)
	testsuite.AssertEqual(t, nil, err)
	testsuite.AssertEqual(t, 5.5, fl)

	hm.Set("key", "val")
	fl, err = hm.IncrementByFloat("key", 1.5)
	testsuite.AssertEqual(t, "invalid type, excepted float", err.Error())
	testsuite.AssertEqual(t, float64(0), fl)
}

func TestHashMap_Len(t *testing.T) {
	hm := New()

	testsuite.AssertEqual(t, 0, hm.Len())
	hm.Set("key", "val")
	testsuite.AssertEqual(t, 1, hm.Len())
	hm.Set("key2", "val")
	testsuite.AssertEqual(t, 2, hm.Len())
	hm.Delete([]string{"key2"})
	testsuite.AssertEqual(t, 1, hm.Len())
}

func TestHashMap_SetBulk(t *testing.T) {
	hm := New()

	res, err := hm.SetBulk([]string{"field1", "val1", "field2", "val2"})
	testsuite.AssertEqual(t, nil, err)
	testsuite.AssertEqual(t, "OK", res)
	testsuite.AssertEqual(t, 2, hm.Len())

	res, err = hm.SetBulk([]string{"field1", "val1", "field2"})
	testsuite.AssertEqual(t, "invalid number of arguments", err.Error())
	testsuite.AssertEqual(t, "", res)
}

func TestHashMap_GetBulk(t *testing.T) {
	hm := New()

	hm.Set("key1", "val1")
	hm.Set("key2", "val2")

	res := hm.GetBulk([]string{"key1", "nonexistent", "key2"})
	testsuite.AssertStringSliceEqual(t, []string{"val1", "", "val2"}, res)
}

func TestHashMap_Setx(t *testing.T) {
	hm := New()

	rep := hm.Setx("key1", "val1")
	val := hm.Get("key1")
	testsuite.AssertEqual(t, 1, rep)
	testsuite.AssertEqual(t, "val1", val)

	rep = hm.Setx("key1", "val2")
	val = hm.Get("key1")
	testsuite.AssertEqual(t, 0, rep)
	testsuite.AssertEqual(t, "val1", val)
}

func TestHashMap_FLen(t *testing.T) {
	hm := New()

	hm.Set("key1", "val1")
	n := hm.FLen("key1")
	testsuite.AssertEqual(t, 4, n)

	n = hm.FLen("nonexistent")
	testsuite.AssertEqual(t, 0, n)
}

func TestHashMap_Vals(t *testing.T) {
	hm := New()

	elements := make([]string, 10)
	for i := 0; i < 10; i++ {
		hm.Set("key"+strconv.Itoa(i), strconv.Itoa(i))
		elements[i] = strconv.Itoa(i)
	}

	vals := hm.Vals()

	testsuite.AssertEqual(t, 10, len(vals))
	testsuite.ContainsElements(t, elements, vals)
}
