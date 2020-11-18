package types

import (
	"context"
	"io"
)

type ILog interface {
	// set context
	Ctx(context.Context) ILog

	Debug(args ...interface{})
	Debugf(format string, args ...interface{})

	Info(args ...interface{})
	Infof(format string, args ...interface{})

	Warning(args ...interface{})
	Warningf(format string, args ...interface{})

	Error(args ...interface{})
	Errorf(format string, args ...interface{})

	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
}

type IDriver interface {
	io.Closer
	// open logger driver
	Open(kv map[string]interface{}) error
	// set log level
	SetLevel(LEVEL) IDriver
	CtxKey(string) IDriver
	// set context
	Ctx(context.Context) ILog
	// get logger
	Logger() ILog
}
