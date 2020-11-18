package lru

import (
	"container/list"
)

// when the item is remove, it will call to tell the user
type PurgeCallback func(interface{}, interface{})

type LRU struct {
	size      int
	purgeList *list.List
	items     map[interface{}]*list.Element
	onPurge   PurgeCallback
}

// entry is used to hold a value in the purgeList
type kv struct {
	k interface{}
	v interface{}
}

func New(size int, onPurge PurgeCallback) *LRU {
	if size <= 0 {
		return nil
	}
	c := new(LRU)
	c.size, c.purgeList, c.items, c.onPurge =
		size, list.New(), make(map[interface{}]*list.Element), onPurge
	return c
}

// ret if some value is replaced
func (c *LRU) Set(k, v interface{}) bool {
	if item, ok := c.items[k]; ok {
		// update item to front
		c.purgeList.MoveToFront(item)
		item.Value.(*kv).v = v
		return false
	}

	item := &kv{k, v}
	entry := c.purgeList.PushFront(item)
	c.items[k] = entry

	// to verify if the list should purge
	flag := c.purgeList.Len() > c.size
	if flag {
		c.delOldest()
	}
	return flag
}

func (c *LRU) delOldest() {
	// get the last item
	item := c.purgeList.Back()
	if item != nil {
		// delete
		c.delElement(item)
	}
}

func (c *LRU) delElement(e *list.Element) {
	c.purgeList.Remove(e)
	kv := e.Value.(*kv)
	delete(c.items, kv.k)
	if c.onPurge != nil {
		c.onPurge(kv.k, kv.v)
	}
}

func (c *LRU) Del(key interface{}) bool {
	if item, ok := c.items[key]; ok {
		c.delElement(item)
		return true
	}
	return false
}

func (c *LRU) Clear() {
	// clear all items
	for k, v := range c.items {
		if c.onPurge != nil {
			c.onPurge(k, v.Value.(*kv).v)
		}
		delete(c.items, k)
	}
	c.purgeList.Init()
}

func (c *LRU) Get(key interface{}) (interface{}, bool) {
	if item, ok := c.items[key]; ok {
		// set it a fresh k
		c.purgeList.MoveToFront(item)
		if item.Value.(*kv) == nil {
			return nil, false
		}
		return item.Value.(*kv).v, true
	}
	return nil, false
}

func (c *LRU) GetWithOutUpdate(key interface{}) (interface{}, bool) {
	if item, ok := c.items[key]; ok {
		if item.Value.(*kv) == nil {
			return nil, false
		}
		return item.Value.(*kv).v, true
	}
	return nil, false
}

func (c *LRU) Contains(key interface{}) (ok bool) {
	_, ok = c.items[key]
	return ok
}

func (c *LRU) DelOldest() (key, value interface{}, ok bool) {
	ent := c.purgeList.Back()
	if ent != nil {
		c.delElement(ent)
		kv := ent.Value.(*kv)
		return kv.k, kv.v, true
	}
	return nil, nil, false
}

// GetOldest returns the oldest entry
func (c *LRU) GetOldest() (key, value interface{}, ok bool) {
	ent := c.purgeList.Back()
	if ent != nil {
		kv := ent.Value.(*kv)
		return kv.k, kv.v, true
	}
	return nil, nil, false
}

func (c *LRU) Keys() []interface{} {
	keys := make([]interface{}, len(c.items))
	i := 0
	for item := c.purgeList.Back(); item != nil; item = item.Prev() {
		keys[i] = item.Value.(*kv).k
		i++
	}
	return keys
}

func (c *LRU) Len() int {
	return c.purgeList.Len()
}

func (c *LRU) Resize(size int) (evicted int) {
	diff := c.Len() - size
	if diff < 0 {
		diff = 0
	}
	for i := 0; i < diff; i++ {
		c.delOldest()
	}
	c.size = size
	return diff
}
