package qlog

import (
	"context"
	"fmt"
	"github.com/qsock/qf/qlog/types"
	"runtime"
	"strings"
)

type Atom struct {
	line   int
	file   string
	format string
	level  types.LEVEL
	args   []interface{}
	ctx    context.Context
}

func (l *Atom) Ctx(ctx context.Context) types.ILog {
	l.ctx = ctx
	return l
}

func (l *Atom) Debug(args ...interface{}) {
	l.p(types.DEBUG, args...)
}

func (l *Atom) Info(args ...interface{}) {
	l.p(types.INFO, args...)
}

func (l *Atom) Warning(args ...interface{}) {
	l.p(types.WARNING, args...)
}

func (l *Atom) Error(args ...interface{}) {
	l.p(types.ERROR, args...)
}

func (l *Atom) Fatal(args ...interface{}) {
	l.p(types.FATAL, args...)
}

func (l *Atom) Debugf(format string, args ...interface{}) {
	l.pf(types.DEBUG, format, args...)
}

func (l *Atom) Infof(format string, args ...interface{}) {
	l.pf(types.INFO, format, args...)
}

func (l *Atom) Warningf(format string, args ...interface{}) {
	l.pf(types.WARNING, format, args...)
}

func (l *Atom) Errorf(format string, args ...interface{}) {
	l.pf(types.ERROR, format, args...)
}

func (l *Atom) Fatalf(format string, args ...interface{}) {
	l.pf(types.FATAL, format, args...)
}

func (l *Atom) p(level types.LEVEL, args ...interface{}) {
	if level < drv.level {
		return
	}
	l.file, l.line = l.getFileNameAndLine()
	l.level = level
	l.args = args
	b := drv.bytes(l)

	fmt.Printf(string(b))
	//fmt.Println()
}

func (l *Atom) pf(level types.LEVEL, format string, args ...interface{}) {
	l.file, l.line = l.getFileNameAndLine()
	l.level = level
	l.args = args
	l.format = format
	b := drv.bytes(l)
	fmt.Printf(string(b))
	fmt.Println()
}

func (l *Atom) getFileNameAndLine() (string, int) {
	var depth = 3 + l.line
	_, file, line, ok := runtime.Caller(depth)
	if !ok {
		return "???", 1
	}
	dirs := strings.Split(file, "/")
	if len(dirs) >= 2 {
		return dirs[len(dirs)-2] + "/" + dirs[len(dirs)-1], line
	}
	return file, line
}
