package kv

import (
	"errors"
	"github.com/go-redis/redis"
	"net"
	"time"
)

var x map[string]*redis.Client = make(map[string]*redis.Client)

type Config struct {
	Addrs       []string `toml:"addrs" json:"addrs"`
	Pwd         string   `toml:"pwd" json:"pwd"`
	PoolSize    int      `toml:"pool_size" json:"pool_size"`
	ReadTimeout int      `toml:"read_timeout" json:"read_timeout"` //单位ms
}

type UnixConfig struct {
	Addr     string `toml:"addr" json:"addr"`
	PoolSize int    `toml:"pool_size" json:"pool_size"`
}

func AddByUnix(name string, c *UnixConfig) error {
	p := redis.NewClient(&redis.Options{
		Network:  "unix",
		Addr:     c.Addr,
		PoolSize: c.PoolSize,
	})

	x[name] = p
	return nil
}

func Add(name string, c *Config) error {
	p, err := NewRedisClientWithConfig(c)
	if err != nil {
		return err
	}
	x[name] = p
	return nil
}

func GetRedisConn(name string) (*redis.Client, error) {
	conn, ok := x[name]
	if ok {
		return conn, nil
	}
	return nil, errors.New("no redis conn")
}

func (c *Config) Parse() (*redis.Options, error) {
	redisNum := len(c.Addrs)
	if redisNum < 1 {
		return nil, errors.New("redis addrs is empty")
	}
	ch := make(chan []string, redisNum)
	for i := 0; i < redisNum; i++ {
		list := make([]string, redisNum)
		for j := 0; j < redisNum; j++ {
			list[j] = c.Addrs[(i+j)%redisNum]
		}
		ch <- list
	}

	if c.ReadTimeout < 100 {
		c.ReadTimeout = 100
	}

	options := &redis.Options{
		Password:     c.Pwd,
		PoolSize:     c.PoolSize,
		ReadTimeout:  time.Duration(c.ReadTimeout) * time.Millisecond,
		WriteTimeout: time.Second,
		PoolTimeout:  100 * time.Millisecond,
		IdleTimeout:  600 * time.Second,
		Dialer: func() (net.Conn, error) {
			list := <-ch
			ch <- list
			for _, addr := range list {
				c, err := net.DialTimeout("tcp", addr, 100*time.Millisecond)
				if err == nil {
					return c, nil
				}
			}
			return nil, errors.New("all redis down")
		},
	}
	return options, nil
}

func NewRedisClientWithConfig(c *Config) (*redis.Client, error) {
	options, err := c.Parse()
	if err != nil {
		return nil, err
	}
	return redis.NewClient(options), nil
}
