package qetcd

import (
	"go.etcd.io/etcd/client/v3"
	"sync"
	"time"
)

type StateChangeCallBackFunc func(string, string, map[string]string)

type Op struct {
	// 配置文件
	C   *Config
	Id  clientv3.LeaseID
	Cli *clientv3.Client

	lock *sync.RWMutex

	items map[string]map[string]string
	// 注册的内容
	kvs         map[string]opKvs
	stateChange map[string]StateChangeCallBackFunc
}

type opKvs struct {
	k   string
	v   string
	ttl int64
}

// 初始化
func New(c *Config) (*Op, error) {
	// 初始化
	if !c.check() {
		return nil, ErrConf
	}
	op := &Op{
		C:           c,
		lock:        new(sync.RWMutex),
		items:       make(map[string]map[string]string),
		kvs:         make(map[string]opKvs),
		stateChange: make(map[string]StateChangeCallBackFunc),
	}
	var err error
	op.Cli, err = clientv3.New(clientv3.Config{
		Endpoints:        c.EndPoints,
		DialTimeout:      3 * time.Second,
		AutoSyncInterval: 60 * time.Second, // 每60秒同步一次集群
	})

	if err != nil {
		return nil, err
	}
	return op, nil
}

func (o *Op) SetStateChangeCallback(prefix string, stateChange StateChangeCallBackFunc) {
	o.lock.Lock()
	o.lock.Unlock()
	o.stateChange[prefix] = stateChange
}

func (o *Op) Close() error {
	return o.Cli.Close()
}
