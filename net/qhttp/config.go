package qhttp

import (
	"github.com/qsock/qf/net/ipaddr"
	"strconv"
	"time"
)

type Config struct {
	Schema   string        `json:"schema" toml:"schema"`
	Name     string        `json:"name" toml:"name"`
	Host     string        `json:"host" toml:"host"`
	Port     int           `json:"port" toml:"port"`
	Addrs    []string      `json:"addrs" toml:"addrs"`
	Timeout  time.Duration `json:"timeout" toml:"timeout"`
	PoolSize int           `json:"pool_size" toml:"poolSize"`
}

func (c *Config) check() bool {
	if c.Port <= 0 && len(c.Schema) == 0 {
		return false
	}
	if len(c.Name) == 0 {
		return false
	}
	if len(c.Host) == 0 {
		return false
	}
	if c.Timeout == 0 {
		c.Timeout = time.Second * 3
	}

	if len(c.Schema) == 0 {
		c.Schema = "http"
	}
	if c.Port == 0 {
		if c.Schema == "http" {
			c.Port = 80
		} else if c.Schema == "https" {
			c.Port = 443
		}
	}
	if len(c.Addrs) == 0 {
		c.Addrs, _ = ipaddr.GetRemoteIp(c.Host)
		for i := 0; i < len(c.Addrs); i++ {
			addr := c.Addrs[i]
			addr += ":" + strconv.FormatInt(int64(c.Port), 10)
			c.Addrs[i] = addr
		}
	}
	if c.PoolSize == 0 {
		c.PoolSize = len(c.Addrs)
	}
	return true
}
