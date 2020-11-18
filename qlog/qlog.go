package qlog

import (
	"encoding/json"
	"github.com/BurntSushi/toml"
	"github.com/qsock/qf/qlog/types"
	"sync"
)

var (
	lock    = new(sync.RWMutex)
	drivers = make(map[string]types.IDriver)
)

func Register(name string, driver types.IDriver) {
	lock.Lock()
	defer lock.Unlock()

	if len(name) == 0 {
		panic("ilog name error")
	}
	if driver == nil {
		panic("ilog nil error")
	}
	drivers[name] = driver
	return
}

func OpenToml(name string, rawConfig string) error {
	var kv map[string]interface{}
	if _, err := toml.Decode(rawConfig, &kv); err != nil {
		return err
	}
	return OpenKv(name, kv)
}

func OpenJson(name string, rawConfig string) error {
	var kv map[string]interface{}
	if err := json.Unmarshal([]byte(rawConfig), &kv); err != nil {
		return err
	}
	return OpenKv(name, kv)
}

func OpenKv(name string, kv map[string]interface{}) error {
	return drivers[name].Open(kv)
}

func Get(name ...string) types.IDriver {
	if len(name) > 0 {
		return drivers[name[0]]
	}

	for name, v := range drivers {
		// 默认的时候，优先不选择stdout
		if name == Name() {
			continue
		}
		return v
	}
	// 保证总有日志可以输出
	return drivers[Name()]
}

func Close() error {
	for _, v := range drivers {
		_ = v.Close()
	}
	return nil
}

func SetKey(key string) {
	for _, v := range drivers {
		v.CtxKey(key)
	}
}
