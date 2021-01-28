package main

import (
	"context"
	"github.com/qsock/qf/qlog"
	"go.uber.org/zap"
	"time"
)

func main() {
	config := qlog.Config{
		Name:  "default.log",
		Dir:   "/tmp",
		Level: "info",
		//Debug:     true,
		AddCaller: true,
	}
	logger := config.Build()

	logger.SetCtxParse(func(ctx context.Context) []zap.Field {
		val := ctx.Value("meta")
		return []zap.Field{zap.String("trace_id", val.(string))}
	})
	defer logger.Flush()
	ctx := context.WithValue(context.Background(), "meta", time.Now().String())
	logger.Ctx(ctx).Info("hello there")
	logger.SetLevel(qlog.DebugLevel)
	logger.Debug("debug", qlog.String("a", "b"))
	logger.Debugf("debug %s", "a")
	//logger.Panic("hello")
}
