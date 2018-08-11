package db

import (
	"testing"
	"github.com/kasvith/kache/pkg/types/list"
	"reflect"
	"github.com/kasvith/kache/pkg/testsuite"
)

func TestReflection(t *testing.T)  {
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
