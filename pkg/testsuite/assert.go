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
