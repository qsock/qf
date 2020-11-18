package concurrent

import "testing"

type T1 struct {
	K int
}

func TestConcurrentMap(t *testing.T) {

	ma := NewIdMap()
	obj0 := &T1{K: 50}
	ma.Set(50, obj0)

	obj1 := &T1{K: 10}
	ma.Set(10, obj1)
	ma.Del(10)

	obj2 := &T1{K: 30}
	ma.Set(30, obj2)

	obj3 := &T1{K: 20}
	ma.Set(20, obj3)

	obj4 := &T1{K: 40}
	ma.Set(40, obj4)

	if v, ok := ma.Get(40); ok {
		t.Log(true, 40, v, ma.Count())
	} else {
		t.Log(false, 40, v, ma.Count())
	}
}
