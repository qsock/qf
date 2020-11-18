package lru

import "testing"

func TestLRU_Set(t *testing.T) {
	callback := func(k interface{}, v interface{}) { t.Log(k, v) }
	m := New(10, callback)
	for i := 0; i < 100; i++ {
		t.Log(m.Set(i, i))
	}
}
