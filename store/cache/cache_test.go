package cache

import (
	"strconv"
	"testing"
	"time"
)

func TestCache_Get(t *testing.T) {
	start := time.Now()
	defer func(s time.Time) {
		t.Log(time.Since(s))
	}(start)
	cache := New(1, func(i interface{}, i2 interface{}) {
		t.Log(i, i2)
	})

	for i := int64(0); i < 1000; i++ {
		k := strconv.FormatInt(i, 10)
		if cache.SetEx(k, i, 1) {
			t.Log("exict", k)
		}
	}

	time.Sleep(time.Second)
	for i := int64(0); i < 1000; i++ {
		t.Log(cache.Get(strconv.FormatInt(i, 10)))
	}
}
