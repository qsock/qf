package qlog

import (
	"context"
	"go.uber.org/zap"
)

// defaultLogger default logger
// Biz Log
// debug=true as default, will be
var defaultLogger *Logger

func SetCfg(cfg *Config) *Logger {
	defaultLogger = cfg.Build()
	return defaultLogger
}

// Auto ...
func Auto(err error) Func {
	if err != nil {
		return defaultLogger.With(zap.Any("err", err.Error())).Error
	}

	return defaultLogger.Info
}

// Info ...
func Info(msg string, fields ...Field) {
	defaultLogger.Info(msg, fields...)
}

// Debug ...
func Debug(msg string, fields ...Field) {
	defaultLogger.Debug(msg, fields...)
}

// Warn ...
func Warn(msg string, fields ...Field) {
	defaultLogger.Warn(msg, fields...)
}

// Error ...
func Error(msg string, fields ...Field) {
	defaultLogger.Error(msg, fields...)
}

// Panic ...
func Panic(msg string, fields ...Field) {
	defaultLogger.Panic(msg, fields...)
}

// DPanic ...
func DPanic(msg string, fields ...Field) {
	defaultLogger.DPanic(msg, fields...)
}

// Fatal ...
func Fatal(msg string, fields ...Field) {
	defaultLogger.Fatal(msg, fields...)
}

// Debugw ...
func Debugw(msg string, keysAndValues ...interface{}) {
	defaultLogger.Debugw(msg, keysAndValues...)
}

// Infow ...
func Infow(msg string, keysAndValues ...interface{}) {
	defaultLogger.Infow(msg, keysAndValues...)
}

// Warnw ...
func Warnw(msg string, keysAndValues ...interface{}) {
	defaultLogger.Warnw(msg, keysAndValues...)
}

// Errorw ...
func Errorw(msg string, keysAndValues ...interface{}) {
	defaultLogger.Errorw(msg, keysAndValues...)
}

// Panicw ...
func Panicw(msg string, keysAndValues ...interface{}) {
	defaultLogger.Panicw(msg, keysAndValues...)
}

// DPanicw ...
func DPanicw(msg string, keysAndValues ...interface{}) {
	defaultLogger.DPanicw(msg, keysAndValues...)
}

// Fatalw ...
func Fatalw(msg string, keysAndValues ...interface{}) {
	defaultLogger.Fatalw(msg, keysAndValues...)
}

// Debugf ...
func Debugf(msg string, args ...interface{}) {
	defaultLogger.Debugf(msg, args...)
}

// Infof ...
func Infof(msg string, args ...interface{}) {
	defaultLogger.Infof(msg, args...)
}

// Warnf ...
func Warnf(msg string, args ...interface{}) {
	defaultLogger.Warnf(msg, args...)
}

// Errorf ...
func Errorf(msg string, args ...interface{}) {
	defaultLogger.Errorf(msg, args...)
}

// Panicf ...
func Panicf(msg string, args ...interface{}) {
	defaultLogger.Panicf(msg, args...)
}

// DPanicf ...
func DPanicf(msg string, args ...interface{}) {
	defaultLogger.DPanicf(msg, args...)
}

// Fatalf ...
func Fatalf(msg string, args ...interface{}) {
	defaultLogger.Fatalf(msg, args...)
}

// Log ...
func (fn Func) Log(msg string, fields ...Field) {
	fn(msg, fields...)
}

// With ...
func With(fields ...Field) *Logger {
	return defaultLogger.With(fields...)
}

// SetCtxParse ...
func SetCtxParse(parser CtParserFunc) {
	defaultLogger.SetCtxParse(parser)
}

// Context ...
func Ctx(ctx context.Context) *Logger {
	return defaultLogger.Ctx(ctx)
}

// Flush ...
func Flush() error {
	return defaultLogger.Flush()
}
