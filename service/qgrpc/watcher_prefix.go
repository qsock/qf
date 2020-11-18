package qgrpc

import (
	"context"
	"encoding/json"
	"google.golang.org/grpc/resolver"
)

func (o *Op) prefixListenServer(serverPrefix string) error {
	listenerKey := GetListenerKey(serverPrefix)

	// 先得到所有server
	kvs, err := GetEtcdClient().Get(context.Background(), listenerKey)
	if err != nil {
		return err
	}
	// 先同步拉取一次
	for k, v := range kvs {
		m := new(RegisterModel)
		if err := json.Unmarshal([]byte(v), m); err != nil {
			continue
		}
		b := newBuilder(k)
		resolver.Register(b)
		if err := o.AddGrpcClient(k, m, b); err != nil {
			return err
		}

	}
	GetEtcdClient().SetStateChangeCallback(listenerKey, o.prefixStateChange)
	go GetEtcdClient().Watch(listenerKey)
	return nil
}
