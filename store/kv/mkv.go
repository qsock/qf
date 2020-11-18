package kv

import (
	"errors"

	"github.com/bradfitz/gomemcache/memcache"
)

var ErrNoClient error = errors.New("no memcache client")
var mClient *memcache.Client

//memcache 的配置
type MConfig struct {
	Addrs    []string `toml:"addrs" json:"addrs"`
	PoolSize int      `toml:"pool_size" json:"pool_size"`
}

func (c *MConfig) GetAddrs() []string {
	if c != nil {
		return c.Addrs
	}
	return []string{}
}

func (c *MConfig) GetPoolSize() int {
	if c != nil {
		return c.PoolSize
	}
	return 5
}

func InitMcClient(c *MConfig) {
	mClient = memcache.New(c.GetAddrs()...)
	mClient.MaxIdleConns = c.GetPoolSize()
	return
}

func GetMcClient() *memcache.Client {
	return mClient
}

func Get(key string) ([]byte, error) {
	client := GetMcClient()
	if client == nil {
		return nil, ErrNoClient
	}

	item, err := client.Get(key)
	if err != nil {
		return nil, err
	}
	return item.Value, nil
}

func Del(key string) error {
	client := GetMcClient()
	if client == nil {
		return ErrNoClient
	}
	return client.Delete(key)
}
