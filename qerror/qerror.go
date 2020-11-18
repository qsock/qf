package qerror

import (
	"fmt"
	"runtime"
	"strings"
)

type Error struct {
	err      error
	content  string
	filename string
	line     int
}

func New(text string) error {
	if len(text) == 0 {
		return nil
	}
	filename, line := getFileNameAndLine()
	return &Error{
		content:  text,
		filename: filename,
		line:     line,
	}
}

func Newf(format string, args ...interface{}) error {
	if len(format) == 0 {
		return nil
	}
	return New(fmt.Sprintf(format, args...))
}

func Wrap(err error, text ...string) error {
	if err == nil {
		return nil
	}
	filename, line := getFileNameAndLine()
	var content string
	if len(text) > 0 {
		content = text[0]
	}

	return &Error{
		err:      err,
		content:  content,
		filename: filename,
		line:     line,
	}
}

func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return Wrap(err, fmt.Sprintf(format, args...))
}

func (e *Error) Error() string {
	ret := fmt.Sprintf("file:%s||line:%d", e.filename, e.line)
	if len(e.content) > 0 {
		ret += "||content:%s" + e.content
	}
	if e.err != nil {
		ret += "||err:%s" + e.err.Error()
	}
	return ret
}

// 堆栈信息
func getFileNameAndLine() (string, int) {
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		return "???", 1
	}
	dirs := strings.Split(file, "/")
	if len(dirs) >= 2 {
		return dirs[len(dirs)-2] + "/" + dirs[len(dirs)-1], line
	}
	return file, line
}
