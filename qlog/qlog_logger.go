package qlog

import (
	"context"
	"github.com/qsock/qf/qlog/types"
)

func Ctx(ctx context.Context) types.ILog {
	return Get().Ctx(ctx)
}

func Debug(args ...interface{}) {
	Get().Logger().Debug(args...)
}

func Info(args ...interface{}) {
	Get().Logger().Info(args...)
}

func Warning(args ...interface{}) {
	Get().Logger().Warning(args...)
}

func Error(args ...interface{}) {
	Get().Logger().Error(args...)
}

func Fatal(args ...interface{}) {
	Get().Logger().Fatal(args...)
}

func Debugf(format string, args ...interface{}) {
	Get().Logger().Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	Get().Logger().Infof(format, args...)
}

func Warningf(format string, args ...interface{}) {
	Get().Logger().Warningf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	Get().Logger().Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	Get().Logger().Fatalf(format, args...)
}
