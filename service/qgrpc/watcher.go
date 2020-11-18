package qgrpc

import (
	"context"
	"fmt"
	"github.com/qsock/qf/net/ipaddr"
	"google.golang.org/grpc/resolver"
	"strings"
)

type watcher struct {
	prefix string
	cc     resolver.ClientConn
}

func (w *watcher) GetPrefix() string {
	if w == nil {
		return ""
	}
	return w.prefix
}

func newBuilder(prefix string) *watcher {
	b := &watcher{prefix: prefix}
	GetEtcdClient().SetStateChangeCallback(b.prefix, b.stateChange)
	return b
}

func (w *watcher) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	w.cc = cc
	_, err := O.etcdClient.Get(context.Background(), w.prefix)
	if err != nil {
		return nil, err
	}
	go O.etcdClient.Watch(w.prefix)
	return w, nil
}

func (w *watcher) Scheme() string {
	return O.C.Schema
}

func (w *watcher) ResolveNow(rn resolver.ResolveNowOptions) {

}

func (w *watcher) Close() {

}

func (w *watcher) resolveAddresses(changedAddrs []string) []resolver.Address {
	addrs := make([]resolver.Address, 0)
	addrsMap := make(map[string]bool)
	for _, addr := range changedAddrs {
		addrsMap[addr] = true
	}
	for _, addr := range changedAddrs {
		// 原始ip
		addrsMap[addr] = true
		// 如果地址包含了本地ip，多建立几个连接，跟本地
		addr = strings.Replace(addr, fmt.Sprintf("%s:", ipaddr.GetLocalIp()), "127.0.0.1:", -1)
		addrsMap[addr] = true
		addr = strings.Replace(addr, fmt.Sprintf("%s:", ipaddr.GetLocalIp()), "localhost:", -1)
		addrsMap[addr] = true
		addr = strings.Replace(addr, fmt.Sprintf("%s:", ipaddr.GetLocalIp()), ipaddr.GetHostname()+":", -1)
		addrsMap[addr] = true
	}
	for k, _ := range addrsMap {
		addrs = append(addrs, resolver.Address{Addr: k})
	}
	return addrs
}

func (w *watcher) stateChange(prefix, action string, kvs map[string]string) {
	addresses := make([]string, 0)
	for _, v := range kvs {
		m := ExtractRegisterModel(v)
		addresses = append(addresses, m.Addr)
	}
	w.cc.UpdateState(resolver.State{Addresses: w.resolveAddresses(addresses)})
}

// 如果是prefix的部分
// 因为有可能发生的是
func (o *Op) prefixStateChange(prefix, action string, kvs map[string]string) {

	for k, v := range kvs {
		if _, ok := o.m[k]; ok {
			continue
		}
		target := ExtractRegisterModel(v)
		// 如果有了就添加上
		b := newBuilder(k)
		resolver.Register(b)
		_ = o.AddGrpcClient(k, target, b)
	}
}
