package qlog_test

import (
	"github.com/qsock/qf/qlog"
	"testing"
)

func Test_Info(t *testing.T) {
	cfg := qlog.DefaultConfig()
	qlog.SetCfg(cfg)
	qlog.Info("hello", qlog.Any("a", "b"))
}
