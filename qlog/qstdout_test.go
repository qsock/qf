package qlog

import (
	"context"
	"github.com/qsock/qf/qlog/types"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	Get().SetLevel(types.WARNING)
	Get().CtxKey("hello")
	Get().Ctx(context.WithValue(context.TODO(), "hello", map[string]string{"a": "v", "v": "c"})).Error("hello")
	Get().Ctx(context.WithValue(context.TODO(), "hello", map[string]string{"a": "v", "v": "c"})).Info("hello")
	time.Sleep(time.Second)
	defer Close()
}

func BenchmarkLogger(b *testing.B) {
	Get().SetLevel(types.WARNING)
	Get().CtxKey("trace_id")

	for i := 0; i < 1000; i++ {
		Get().Ctx(context.WithValue(context.TODO(), "trace_id", map[string]string{"a": "v", "v": "c"})).Errorf("error||hello:%d", i)
		Get().Ctx(context.WithValue(context.TODO(), "trace_id", map[string]string{"a": "v", "v": "c"})).Infof("info||hello:%d", i)
	}

	// 有个连接释放的过程
	defer Close()
}
