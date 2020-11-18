package cache

import (
	"github.com/qsock/qvs/cache/keyhash"
	"github.com/qsock/qvs/cache/lru"
	"sync"
	"time"
)

const (
	//  total page cache
	Page = 1024
	// second for one day
	DaySecond = 86400
)

type item struct {
	d interface{}
	t int
}

type Cache struct {
	lock sync.RWMutex
	kvs  []*lru.LRU
}

// size is size*Page
func New(size int, callback ...func(interface{}, interface{})) *Cache {
	c := new(Cache)
	c.kvs = make([]*lru.LRU, 0, Page)
	var cb func(interface{}, interface{})
	if len(callback) > 0 {
		cb = callback[0]
	}
	for i := 0; i < Page; i++ {
		c.kvs = append(c.kvs, lru.New(size, cb))
	}
	go c.flush()
	return c
}

func (c *Cache) flush() {
	// every minute flush for once
	for range time.NewTicker(time.Minute).C {
		for i := 0; i < Page; i++ {
			page := c.kvs[i]
			keys := page.Keys()
			for _, k := range keys {
				_, _ = c.Get(k.(string))
			}
		}
	}
}

func (c *Cache) Clear() {
	c.lock.Lock()
	defer c.lock.Unlock()
	for i := 0; i < Page; i++ {
		page := c.kvs[i]
		page.Clear()
	}
}

func (c *Cache) Set(k string, v interface{}) bool {
	return c.SetEx(k, v, DaySecond)
}

func (c *Cache) SetEx(k string, v interface{}, exp int) bool {
	idx := hash(k)
	c.lock.Lock()
	defer c.lock.Unlock()

	kv := c.kvs[idx]
	return kv.Set(k, &item{d: v, t: int(time.Now().Unix()) + exp})
}

func (c *Cache) Del(k string) bool {
	idx := hash(k)
	kv := c.kvs[idx]

	c.lock.Lock()
	defer c.lock.Unlock()
	return kv.Del(k)
}

func (c *Cache) Get(k string) (interface{}, bool) {
	idx := hash(k)
	kv := c.kvs[idx]

	c.lock.Lock()
	defer c.lock.Unlock()
	val, ok := kv.Get(k)
	if !ok {
		return nil, false
	}
	v := val.(*item)
	if v.t >= int(time.Now().Unix()) {
		return v.d, true
	}
	_ = kv.Del(k)
	return nil, false
}

func (c *Cache) TTL(k string) int {
	idx := hash(k)
	kv := c.kvs[idx]

	c.lock.Lock()
	defer c.lock.Unlock()
	val, ok := kv.Get(k)
	if !ok {
		return 0
	}
	v := val.(*item)
	if v.t >= int(time.Now().Unix()) {
		return v.t
	}
	_ = kv.Del(k)
	return 0
}

func (c *Cache) Expires(k string, exp int) bool {
	idx := hash(k)
	kv := c.kvs[idx]

	c.lock.Lock()
	defer c.lock.Unlock()
	val, ok := kv.Get(k)
	if !ok {
		return false
	}
	v := val.(*item)
	v.t += exp

	if v.t < int(time.Now().Unix()) {
		_ = kv.Del(k)
		return false
	}
	return true
}

func (c *Cache) TryGet(k string) (interface{}, bool) {
	idx := hash(k)
	kv := c.kvs[idx]

	c.lock.RLock()
	defer c.lock.RUnlock()
	val, ok := kv.GetWithOutUpdate(k)
	if !ok {
		return nil, false
	}
	v := val.(*item)
	return v, true
}

func (c *Cache) Contains(k string) bool {
	idx := hash(k)
	kv := c.kvs[idx]
	c.lock.RLock()
	defer c.lock.RUnlock()
	return kv.Contains(k)
}

func (c *Cache) Len() int {
	var total int
	c.lock.RLock()
	defer c.lock.RUnlock()
	for i := 0; i < Page; i++ {
		page := c.kvs[i]
		total += page.Len()
	}
	return total
}

func (c *Cache) Keys() []string {
	keys := make([]string, 0)
	c.lock.RLock()
	defer c.lock.RUnlock()

	for _, kv := range c.kvs {
		for _, k := range kv.Keys() {
			keys = append(keys, k.(string))
		}
	}
	return keys
}

func hash(key string) int {
	return int(keyhash.Hash32([]byte(key)) % Page)
}
