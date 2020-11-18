package qetcd

import (
	"context"
	"go.etcd.io/etcd/client/v3"
	"time"
)

func (o *Op) checkRegister(ctx context.Context, kvs map[string]string) bool {
	for k, _ := range kvs {
		if _, ok := o.kvs[k]; ok {
			return false
		}
	}

	return true
}

// 去etcd注册
func (o *Op) Register(ctx context.Context, kvs map[string]string, ttl int64) error {
	if !o.checkRegister(ctx, kvs) {
		return ErrHasRegistered
	}
	for k, v := range kvs {
		m := opKvs{k: k, v: v, ttl: ttl}
		o.kvs[k] = m
	}

	// 创建一个ttl秒的租约
	resp, err := o.Cli.Grant(context.Background(), ttl)
	if err != nil {
		return err
	}
	o.Id = resp.ID

	// 用同一个租约，创建多个注册
	for k, v := range kvs {
		if _, err := o.Cli.Put(ctx, k, v, clientv3.WithLease(resp.ID)); err != nil {
			return err
		}
	}

	// 续约
	ch, err := o.Cli.KeepAlive(context.Background(), resp.ID)
	if err != nil {
		return err
	}
	go o.keepAlive(ch)
	return nil
}

func (o *Op) keepAlive(ch <-chan *clientv3.LeaseKeepAliveResponse) {
	for {
		select {
		case <-o.Cli.Ctx().Done():
			{
				return
			}
		case _, ok := <-ch:
			{
				if !ok {
					// 释放
					_ = o.Revoke()
					// 不直接重连
					time.Sleep(time.Second)
					// 重新注册一次
					ma := make(map[string]opKvs)
					for k, v := range o.kvs {
						ma[k] = v
					}

					// 清空
					o.kvs = make(map[string]opKvs)
					for k, v := range ma {
						// 重新注册
						_ = o.Register(context.Background(), map[string]string{k: v.v}, v.ttl)
					}
					return
				}
			}
		}
	}
}

func (o *Op) Revoke() error {
	_, err := o.Cli.Revoke(context.Background(), o.Id)
	return err
}
