package concurrent

import (
	"github.com/qsock/qf/encrypt"
	"sync"
)

const concurrentNum = 32

type IdMap struct {
	Ms    [concurrentNum]*sync.Map
	count *Int64
}

// 新建一个map
func NewIdMap() *IdMap {
	group := new(IdMap)
	for i := 0; i < concurrentNum; i++ {
		group.Ms[i] = new(sync.Map)
	}
	group.count = NewInt64(0)
	return group
}

func (g *IdMap) Set(id int64, item interface{}) {
	m := g.Ms[id%concurrentNum]
	m.Store(id, item)
	g.count.IncrementAndGet()
}

func (g *IdMap) Get(id int64) (interface{}, bool) {
	m := g.Ms[id%concurrentNum]
	return m.Load(id)
}

func (g *IdMap) Del(id int64) {
	m := g.Ms[id%concurrentNum]
	m.Delete(id)
	g.count.DecrementAndGet()
}

func (g *IdMap) SetS(id string, item interface{}) {
	hashId := encrypt.HashInt(id)
	m := g.Ms[hashId%concurrentNum]
	m.Store(id, item)
	g.count.IncrementAndGet()
}

func (g *IdMap) GetS(id string) (interface{}, bool) {
	hashId := encrypt.HashInt(id)
	m := g.Ms[hashId%concurrentNum]
	return m.Load(id)
}

func (g *IdMap) DelS(id string) {
	hashId := encrypt.HashInt(id)
	m := g.Ms[hashId%concurrentNum]
	m.Delete(id)
	g.count.DecrementAndGet()
}

func (g *IdMap) Count() int64 {
	return g.count.Get()
}

func (g *IdMap) Array() []interface{} {
	items := make([]interface{}, 0)
	for _, v := range g.Ms {
		v.Range(func(_, val interface{}) bool {
			items = append(items, val)
			return true
		})
	}
	return items
}
