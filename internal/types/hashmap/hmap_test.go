package hashmap

import (
	"github.com/kasvith/kache/pkg/testsuite"
	"strconv"
	"testing"
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
		hm.Set(strconv.Itoa(i), strconv.Itoa(i))
		elements[i] = strconv.Itoa(i)
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

	deleted := hm.Delete("key1")

	testsuite.AssertEqual(t, 1, deleted)

	deleted = hm.Delete("key2", "key3")
	testsuite.AssertEqual(t, 2, deleted)

	deleted = hm.Delete("nonexistent")
	testsuite.AssertEqual(t, 0, deleted)
}

func TestHashMap_Exists(t *testing.T) {
	hm := New()

	hm.Set("key1", "val1")

	key1Exists := hm.Exists("key1")
	testsuite.AssertEqual(t, 1, key1Exists)

	unknown := hm.Exists("somethingelse")
	testsuite.AssertEqual(t, 0, unknown)
}
