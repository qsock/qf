package qgrpc

import (
	"github.com/qsock/qf/net/ipaddr"
	"github.com/qsock/qf/service/qetcd"
	"sync"
)

var (
	O *Op
)

// 初始化
func Init(c *Config) error {
	// 初始化
	if !c.check() {
		return ErrConf
	}
	O = &Op{
		C:       c,
		LocalIp: ipaddr.GetLocalIp(),
		m:       make(map[string]*GrpcClient),
		lock:    new(sync.RWMutex),
	}
	var err error
	cfg := new(qetcd.Config)
	cfg.EndPoints = c.EndPoints
	if O.etcdClient, err = qetcd.New(cfg); err != nil {
		return err
	}
	err = O.register()
	if err != nil {
		return err
	}
	// 开始监听
	err = O.listen()
	if err != nil {
		return err
	}
	return nil
}
