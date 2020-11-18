package qsyslog

import (
	"context"
	"github.com/qsock/qf/qlog"
	"github.com/qsock/qf/qlog/types"
	"testing"
	"time"
)

var (
	cfg = `tag       = "test"`
)

func TestLogger(t *testing.T) {
	_ = qlog.OpenToml(Name(), cfg)
	qlog.Get().SetLevel(types.WARNING)
	qlog.Get().CtxKey("hello")
	qlog.Get().Ctx(context.WithValue(context.TODO(), "hello", map[string]string{"a": "v", "v": "c"})).Error("hello")
	qlog.Get().Ctx(context.WithValue(context.TODO(), "hello", map[string]string{"a": "v", "v": "c"})).Info("hello")
	time.Sleep(time.Second)
	defer qlog.Close()
}

func BenchmarkLogger(b *testing.B) {
	_ = qlog.OpenToml(Name(), cfg)
	qlog.Get().SetLevel(types.WARNING)
	qlog.Get().CtxKey("trace_id")
	time.Sleep(time.Second)
	for i := 0; i < 1000; i++ {
		qlog.Get().Ctx(context.WithValue(context.TODO(), "trace_id", map[string]string{"a": "v", "v": "c"})).Errorf("info||hello:%d", i)
		qlog.Get().Ctx(context.WithValue(context.TODO(), "trace_id", map[string]string{"a": "v", "v": "c"})).Warningf("error||hello:%d", i)
	}
	// 有个连接释放的过程
	defer qlog.Close()
}
