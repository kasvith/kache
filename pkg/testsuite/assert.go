package testsuite

import (
	"testing"
	"reflect"
)

//AssertEqual Will do a typical assertation
func AssertEqual(t *testing.T, expected interface{}, given interface{})  {
	if expected != given {
		t.Errorf("Assertion failed, excepted %v : given %v", expected, given)
	}
}

func AssertStringSliceEqual(t *testing.T, expected []string, given []string)  {
	if !reflect.DeepEqual(expected, given) {
		t.Errorf("Slice failed, excepted %v : given %v", expected, given)
	}
}

func AssertNil(t *testing.T, i interface{})  {
	if !reflect.ValueOf(i).IsNil() {
		t.Errorf("Given %v is not nil", i)
	}
}

func ExceptError(t *testing.T, excepted error, given error)  {
	if excepted.Error() != given.Error() {
		t.Errorf("Excepted %v, found %v", excepted, given)
	}
}
