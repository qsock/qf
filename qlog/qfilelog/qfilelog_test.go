package qfilelog

import (
	"context"
	"github.com/qsock/qf/qlog"
	"github.com/qsock/qf/qlog/types"
	"testing"
	"time"
)

var (
	cfg = `
 	# 文件名称
    filename       = "debug"
    # 文件数量
    filenum        = 50
    # 文件大小Mb
    filesize       = 256
    # 日志输出等级，详细请看
    level          = 2
    # 日志的输出文件夹
    dir            = "./logs"
    #  是否开启gzip压缩
    use_gzip       = true`
)

func TestLogger(t *testing.T) {
	_ = qlog.OpenToml(Name(), cfg)
	qlog.Get().SetLevel(types.WARNING)
	qlog.Get().CtxKey("hello")
	qlog.Get().Ctx(context.WithValue(context.TODO(), "hello", map[string]string{"a": "v", "v": "c"})).Error("hello:" + time.Now().String())
	qlog.Get().Ctx(context.WithValue(context.TODO(), "hello", map[string]string{"a": "v", "v": "c"})).Info("hello:" + time.Now().String())
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
	}
	time.Sleep(time.Second)
	defer qlog.Close()
}
