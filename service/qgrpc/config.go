package qgrpc

import "strings"

const (
	DefaultSchema = "default"
)

// etcd的配置
type Config struct {
	// etcd的地址 必须要
	Endpoints []string `toml:"endpoints" json:"endpoints"`
	// etcd的服务的前缀,需要通过prefix来区别正式测试等等环境 必须要
	Prefix string `toml:"prefix" json:"prefix"`

	// 注册部分
	// 服务名称，固定为 aa.Aa ,grpc的服务名称,如果服务多个，可以用 | 来分割
	ServerName string `toml:"serverName" json:"server_name"`
	// 你注册的服务地址,可以是外网地址，也可以是内网地址
	Addr string `toml:"addr" json:"addr"`
	// 额外存储的部分
	Ext string `toml:"ext" json:"ext"`

	Schema string `toml:"schema" json:"schema"`

	// 监听部分
	// 想要监听的服务地址
	// 用grpc做lb
	WatchServers []string `toml:"watchServers" json:"watch_servers"`

	// 这种就是默认全监听，收到之后，相当于每个服务，会单独的去监听
	// 适合长连接的类型，不需要lb的这种
	WatchPrefix []string `toml:"watchPrefix" json:"watch_prefix"`

	// 是否阻塞启动
	Block bool `toml:"block"`
}

func (c *Config) check() bool {
	if len(c.EndPoints) == 0 {
		return false
	}

	// 如果注册了服务名称，就必须要注册port
	if len(c.ServerName) > 0 && len(c.Addr) == 0 {
		return false
	}

	if len(c.Addr) > 0 && len(c.ServerName) == 0 {
		return false
	}
	if len(c.Schema) == 0 {
		c.Schema = DefaultSchema
	}

	for i := 0; i < len(c.WatchServers); i++ {
		watcherKey := c.WatchServers[i]
		if strings.Contains(watcherKey, "/") {
			watcherKey = strings.ReplaceAll(watcherKey, "/", "")
		}
		c.WatchServers[i] = watcherKey
	}

	for i := 0; i < len(c.WatchPrefix); i++ {
		watcherKey := c.WatchPrefix[i]
		if strings.Contains(watcherKey, "/") {
			watcherKey = strings.ReplaceAll(watcherKey, "/", "")
		}
		c.WatchPrefix[i] = watcherKey
	}

	return true
}
