package qlog_test

import (
	"github.com/qsock/qf/qlog"
	"testing"
)

func Test_Info(t *testing.T) {
	qlog.Info("hello", qlog.Any("a", "b"))
}
