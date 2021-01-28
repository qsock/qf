package qlog_test

import (
	"testing"

	"qifeiwu.com/qlog"
)

func Test_Info(t *testing.T) {
	qlog.Info("hello", qlog.Any("a", "b"))
}
