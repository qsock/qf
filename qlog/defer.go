package qlog

import (
	"github.com/qsock/qf/qlog/internal"
)

var (
	globalDefers = internal.NewStack()
)

// deferRegister 注册一个defer函数
func deferRegister(fns ...func() error) {
	globalDefers.Push(fns...)
}

// deferClean 清除
func deferClean() {
	globalDefers.Clean()
}
