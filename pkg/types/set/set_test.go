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

package set

import (
	"testing"

	testifyAssert "github.com/stretchr/testify/assert"
)

func TestSet_Add(t *testing.T) {
	assert := testifyAssert.New(t)
	set := New()

	res := set.Add([]string{"hello", "world", "hello"})
	assert.Equal(2, res)
}

func TestSet_Card(t *testing.T) {
	assert := testifyAssert.New(t)
	set := New()

	set.Add([]string{"hello", "world", "hello", "bye"})
	assert.Equal(3, set.Card())
}

func TestSet_Elems(t *testing.T) {
	assert := testifyAssert.New(t)
	set := New()

	set.Add([]string{"hello", "world", "hello", "bye"})
	res := set.Elems()

	assert.ElementsMatch([]string{"hello", "world", "bye"}, res)
}

func TestSet_Diff(t *testing.T) {
	assert := testifyAssert.New(t)
	set1 := New()
	set2 := New()
	set3 := New()

	set1.Add([]string{"a", "b", "c", "d"})
	set2.Add([]string{"c"})
	set3.Add([]string{"a", "b"})

	el := set1.Diff([]Set{*set2, *set3})
	assert.ElementsMatch([]string{"d"}, el)
}

func TestSet_DiffS(t *testing.T) {
	assert := testifyAssert.New(t)
	set1 := New()
	set2 := New()
	set3 := New()

	set1.Add([]string{"a", "b", "c", "d"})
	set2.Add([]string{"c"})
	set3.Add([]string{"a", "b"})

	el := set1.DiffS([]Set{*set2, *set3})
	assert.ElementsMatch([]string{"d"}, el.Elems())
}

func TestIntersection(t *testing.T) {
	assert := testifyAssert.New(t)
	set1 := New()
	set2 := New()
	set3 := New()

	set1.Add([]string{"a", "b", "c", "d"})
	set2.Add([]string{"b"})
	set3.Add([]string{"a", "b"})

	inter := Intersection([]Set{*set1, *set2, *set3})
	assert.Subset(inter, []string{"b"})

	set4 := New()
	inter = Intersection([]Set{*set1, *set2, *set3, *set4})
	assert.ElementsMatch([]string{}, inter)

	set4.Add([]string{"x"})
	inter = Intersection([]Set{*set1, *set2, *set3, *set4})
	assert.ElementsMatch([]string{}, inter)
}

func TestIntersectionS(t *testing.T) {
	assert := testifyAssert.New(t)
	set1 := New()
	set2 := New()
	set3 := New()

	set1.Add([]string{"a", "b", "c", "d"})
	set2.Add([]string{"b"})
	set3.Add([]string{"a", "b"})

	inter := IntersectionS([]Set{*set1, *set2, *set3})
	assert.Subset(inter.Elems(), []string{"b"})

	set4 := New()
	inter = IntersectionS([]Set{*set1, *set2, *set3, *set4})
	assert.ElementsMatch([]string{}, inter.Elems())

	set4.Add([]string{"x"})
	inter = IntersectionS([]Set{*set1, *set2, *set3, *set4})
	assert.ElementsMatch([]string{}, inter.Elems())
}

func TestSet_Exists(t *testing.T) {
	assert := testifyAssert.New(t)
	set := New()

	set.Add([]string{"hello", "world", "hello", "bye"})
	rep := set.Exists("hello")
	assert.Equal(1, rep)
	rep = set.Exists("nonexistent")
	assert.Equal(0, rep)
}

func TestMove(t *testing.T) {
	assert := testifyAssert.New(t)
	set1 := New()
	set2 := New()

	set1.Add([]string{"a", "b", "c", "d"})
	set2.Add([]string{"b"})

	res := Move("a", set1, set2)
	assert.Equal(1, res)
	assert.ElementsMatch([]string{"b", "c", "d"}, set1.Elems())
	assert.ElementsMatch([]string{"a", "b"}, set2.Elems())

	res = Move("unknown", set1, set2)
	assert.Equal(0, res)
}

func TestSet_Delete(t *testing.T) {
	assert := testifyAssert.New(t)
	set := New()

	set.Add([]string{"hello", "world", "hello", "bye"})
	rep := set.Delete([]string{"hello", "world", "nonexistent"})
	assert.Equal(2, rep)
}

func TestUnion(t *testing.T) {
	assert := testifyAssert.New(t)
	set1 := New()
	set2 := New()
	set3 := New()

	set1.Add([]string{"b", "g"})
	set2.Add([]string{"a", "b", "c", "d"})
	set3.Add([]string{"a", "b", "f"})

	union := Union([]Set{*set1, *set2, *set3})
	assert.ElementsMatch([]string{"a", "b", "c", "d", "g", "f"}, union)
}

func TestUnionS(t *testing.T) {
	assert := testifyAssert.New(t)
	set1 := New()
	set2 := New()
	set3 := New()

	set1.Add([]string{"b", "g"})
	set2.Add([]string{"a", "b", "c", "d"})
	set3.Add([]string{"a", "b", "f"})

	union := UnionS([]Set{*set1, *set2, *set3})
	assert.ElementsMatch([]string{"a", "b", "c", "d", "g", "f"}, union.Elems())
}
