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

package testsuite

import (
	"reflect"
	"testing"
)

//AssertEqual Will do a typical assertation
func AssertEqual(t *testing.T, expected interface{}, given interface{}) {
	if expected != given {
		t.Errorf("Assertion failed, excepted [%T] %v : given [%T] %v", expected, expected, given, given)
	}
}

func AssertStringSliceEqual(t *testing.T, expected []string, given []string) {
	if !reflect.DeepEqual(expected, given) {
		t.Errorf("Slice failed, excepted [%T] %v : given [%T] %v", expected, expected, given, given)
	}
}

func AssertNil(t *testing.T, i interface{}) {
	if !reflect.ValueOf(i).IsNil() {
		t.Errorf("Given [%T] %v is not nil", i, i)
	}
}

func ExceptError(t *testing.T, excepted error, given error) {
	if excepted.Error() != given.Error() {
		t.Errorf("Excepted [%T] %v, found [%T] %v", excepted, excepted, given, given)
	}
}

func ContainsElements(t *testing.T, expected []string, given []string) {
	for elem := range expected {
		found := false
		for test := range given {
			if elem == test {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("Element %v(%T) was not found in given array %v(%T)", elem, elem, given, given)
		}
	}
}
