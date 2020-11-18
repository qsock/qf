package qetcd

const (
	DefaultSchema = "default"
)

// etcd的配置
type Config struct {
	// etcd的地址 必须要
	EndPoints []string `toml:"endpoints" json:"endpoints"`
}

func (c *Config) check() bool {
	if len(c.EndPoints) == 0 {
		return false
	}
	return true
}

var (
	ActionDel = "del"
	ActionPut = "put"
)
