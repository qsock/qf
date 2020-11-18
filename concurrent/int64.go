package concurrent

import (
	"fmt"
	"sync/atomic"
)

type Int64 int64

func NewInt64(val int64) *Int64 {
	a := Int64(val)
	return &a
}

// 得到该值
func (a *Int64) Get() int64 {
	return int64(*a)
}

// 将值设置进去
func (a *Int64) Set(val int64) {
	atomic.StoreInt64((*int64)(a), val)
}

// 对比并且设置，原子操作
func (a *Int64) CompareAndSet(expect, update int64) bool {
	return atomic.CompareAndSwapInt64((*int64)(a), expect, update)
}

// 设置新值，并返回旧值
func (a *Int64) GetAndSet(val int64) int64 {
	for {
		current := a.Get()
		if a.CompareAndSet(current, val) {
			return current
		}
	}
}

func (a *Int64) GetAndIncrement() int64 {

	for {
		current := a.Get()
		next := current + 1
		if a.CompareAndSet(current, next) {
			return current
		}
	}
}

func (a *Int64) GetAndDecrement() int64 {
	for {
		current := a.Get()
		next := current - 1
		if a.CompareAndSet(current, next) {
			return current
		}
	}
}

func (a *Int64) GetAndAdd(val int64) int64 {
	for {
		current := a.Get()
		next := current + val
		if a.CompareAndSet(current, next) {
			return current
		}
	}
}

func (a *Int64) IncrementAndGet() int64 {
	for {
		current := a.Get()
		next := current + 1
		if a.CompareAndSet(current, next) {
			return next
		}
	}
}

func (a *Int64) DecrementAndGet() int64 {
	for {
		current := a.Get()
		next := current - 1
		if a.CompareAndSet(current, next) {
			return next
		}
	}
}

func (a *Int64) AddAndGet(val int64) int64 {
	for {
		current := a.Get()
		next := current + val
		if a.CompareAndSet(current, next) {
			return next
		}
	}
}

func (a *Int64) String() string {
	return fmt.Sprintf("%d", a.Get())
}
