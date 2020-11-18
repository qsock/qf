package qgrpc

import (
	"context"
)

// 去etcd注册
func (o *Op) register() error {
	if len(o.C.ServerName) == 0 {
		return nil
	}
	kvs := make(map[string]string)
	for _, serverName := range GetRegisterServerNames() {
		k := GetRegisterKey(serverName)
		v := GetRegisterValue()
		kvs[k] = v
	}
	if err := o.etcdClient.Register(context.Background(), kvs, 5); err != nil {
		return err
	}
	return nil
}
