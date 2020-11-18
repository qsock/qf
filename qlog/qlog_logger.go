package qlog

import (
	"context"
	"github.com/qsock/qf/qlog/types"
)

func Ctx(ctx context.Context) types.ILog {
	return Get().Ctx(ctx)
}

func Debug(args ...interface{}) {
	Get().Logger(1).Debug(args...)
}

func Info(args ...interface{}) {
	Get().Logger(1).Info(args...)
}

func Warning(args ...interface{}) {
	Get().Logger(1).Warning(args...)
}

func Error(args ...interface{}) {
	Get().Logger(1).Error(args...)
}

func Fatal(args ...interface{}) {
	Get().Logger(1).Fatal(args...)
}

func Debugf(format string, args ...interface{}) {
	Get().Logger(1).Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	Get().Logger(1).Infof(format, args...)
}

func Warningf(format string, args ...interface{}) {
	Get().Logger(1).Warningf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	Get().Logger(1).Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	Get().Logger(1).Fatalf(format, args...)
}
