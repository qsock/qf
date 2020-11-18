package qlog_test

import (
	"context"
	"github.com/qsock/qf/qlog"
	"testing"
)

func TestInfo(t *testing.T) {
	qlog.Info("1234", 4567)
	qlog.Get().Ctx(context.TODO()).Info("123")
}
