package qlog

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config ...
type Config struct {
	// Dir 日志输出目录
	Dir string `json:"dir" toml:"dir"`
	// Name 日志文件名称
	Name string `json:"name" toml:"name"`
	// Level 日志初始等级
	Level string `json:"level" toml:"level"`
	// 日志初始化字段
	Fields []zap.Field `json:"fields" toml:"fields"`
	// 是否添加调用者信息
	AddCaller bool `json:"add_caller" toml:"add_caller"`
	// 日志前缀
	Prefix string `json:"prefix" toml:"prefix"`
	// 日志输出文件最大长度，超过改值则截断
	MaxSize   int `json:"max_size" toml:"max_size"`
	MaxAge    int `json:"max_age" toml:"max_age"`
	MaxBackup int `json:"max_backup" toml:"max_backup"`
	// 日志磁盘刷盘间隔
	Interval      time.Duration          `json:"internal" toml:"internal"`
	CallerSkip    int                    `json:"caller_skip" toml:"caller_skip"`
	Async         bool                   `json:"async" toml:"async"`
	Queue         bool                   `json:"queue" toml:"queue"`
	QueueSleep    time.Duration          `json:"queue_sleep" toml:"queue_sleep"`
	Core          zapcore.Core           `json:"core" toml:"core"`
	Debug         bool                   `json:"debug" toml:"debug"`
	EncoderConfig *zapcore.EncoderConfig `json:"encoder_config" toml:"encoder_config"`
}

// Filename ...
func (config *Config) Filename() string {
	return fmt.Sprintf("%s/%s", config.Dir, config.Name)
}

// DefaultConfig ...
func DefaultConfig() *Config {
	return &Config{
		Name:          "default.log",
		Dir:           ".",
		Level:         "info",
		MaxSize:       500, // 500M
		MaxAge:        1,   // 1 day
		MaxBackup:     10,  // 10 backup
		Interval:      24 * time.Hour,
		CallerSkip:    1,
		AddCaller:     false,
		Async:         true,
		Queue:         false,
		QueueSleep:    100 * time.Millisecond,
		EncoderConfig: DefaultZapConfig(),
	}
}

// Build ...
func (config Config) Build() *Logger {
	if config.EncoderConfig == nil {
		config.EncoderConfig = DefaultZapConfig()
	}
	if config.Debug {
		config.EncoderConfig.EncodeLevel = DebugEncodeLevel
	}
	logger := newLogger(&config)
	return logger
}
