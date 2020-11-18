package qetcd

import (
	"context"
	"github.com/qsock/qf/qlog"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.etcd.io/etcd/client/v3"
)

func (o *Op) Get(ctx context.Context, prefix string) (map[string]string, error) {
	resp, err := o.Cli.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	m := make(map[string]string)
	for _, ev := range resp.Kvs {
		k := string(ev.Key)
		v := string(ev.Value)
		m[k] = v
		o.PutItem(prefix, k, v)
	}
	return m, nil
}

func (o *Op) Watch(prefix string) {
	rch := o.Cli.Watch(context.Background(), prefix, clientv3.WithPrefix())
	for wResp := range rch {
		for _, ev := range wResp.Events {
			switch ev.Type {
			case mvccpb.PUT:
				{
					o.PutItem(prefix, string(ev.Kv.Key), string(ev.Kv.Value))
				}
			case mvccpb.DELETE:
				{
					o.DelItem(prefix, string(ev.Kv.Key))
				}
			}
		}
	}
}

func (o *Op) GetWatch(prefix string) map[string]string {
	o.lock.RLock()
	defer o.lock.RUnlock()
	return o.items[prefix]
}

func (o *Op) PutItem(prefix, k, v string) {
	qlog.Infof("PutItem||pre:%s||k:%s||v:%s", prefix, k, v)
	o.lock.Lock()
	defer o.lock.Unlock()

	val, ok := o.items[prefix]
	if !ok {
		val = make(map[string]string)
	}
	val[k] = v
	o.items[prefix] = val
	if changeFunc, ok := o.stateChange[prefix]; ok {
		changeFunc(prefix, ActionPut, val)
	}
}

func (o *Op) DelItem(prefix, k string) {
	qlog.Infof("DelItem||pre:%s||k:%s", prefix, k)
	o.lock.Lock()
	defer o.lock.Unlock()
	val := o.items[prefix]
	delete(val, k)
	o.items[prefix] = val

	if changeFunc, ok := o.stateChange[prefix]; ok {
		changeFunc(prefix, ActionDel, val)
	}
}
