package ws

import (
	"time"
)

type Config struct {
	// 写超时
	WriteTimeout time.Duration
	TTLTimeout   time.Duration
	// 一个消息最大的size
	MaxMessageSize int64
	// session消息最大的数量
	MessageBufferSize int
}

func DefaultConfig() *Config {
	c := new(Config)
	c.MaxMessageSize = 64 * 1024 // 64kb
	c.MessageBufferSize = 64
	c.WriteTimeout = 10 * time.Second
	c.TTLTimeout = 300 * time.Second
	return c
}
