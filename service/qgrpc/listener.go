package qgrpc

import (
	"context"
	"encoding/json"
	"github.com/qsock/qf/qlog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/resolver"
	"path"
	"time"
)

func (o *Op) ping() {
	// 每秒ping
	o.lock.RLock()
	servers := o.m
	o.lock.RUnlock()
	for range time.NewTicker(time.Second).C {
		for _, m := range servers {
			req, resp := new(NoArgs), new(NoArgs)
			method := path.Join("/", m.serverName, "/Ping")
			if err := CallWithServerNameTimeout(context.Background(), m.name, method, req, resp, 100*time.Millisecond); err != nil {
				qlog.Get().Logger().Errorf("grpc||name:%s||err:%v", m.name, method)
			}
		}
	}
}

func (o *Op) listen() error {
	//  监听
	// 监听服务
	for _, serverName := range o.C.WatchServers {
		o.listenServer(serverName)
	}
	for _, serverPrefix := range o.C.WatchPrefix {
		o.prefixListenServer(serverPrefix)
	}
	if len(o.C.WatchServers) > 0 || len(o.C.WatchPrefix) > 0 {
		go o.ping()
	}
	return nil
}

// 因为用户是唯一的,而他对于同一个name的轮询策略是决定好了的
// 这里只能自己做轮询管理，如果委托给grpc管理，则没办法
func (o *Op) listenServer(serverName string) error {
	listenerKey := GetListenerKey(serverName)
	b := newBuilder(listenerKey)
	resolver.Register(b)
	return o.AddGrpcClient(serverName, nil, b)
}

func (o *Op) DeleteGrpcClient(serverName string) {
	o.lock.Lock()
	defer o.lock.Unlock()
	delete(O.m, serverName)
}

func (o *Op) AddGrpcClient(serverName string, target *RegisterModel, w *watcher) error {
	var listedAddr = o.C.Schema + "://grpc/" + serverName

	mam := map[string]interface{}{
		"LoadBalancingPolicy": roundrobin.Name,
	}
	jsonB, _ := json.Marshal(mam)

	dialOpt := []grpc.DialOption{
		grpc.WithConnectParams(grpc.ConnectParams{MinConnectTimeout: 3 * time.Second}),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(string(jsonB)),
	}
	if o.C.Block {
		dialOpt = append(dialOpt, grpc.WithBlock())
	}
	cli, err := grpc.Dial(listedAddr, dialOpt...)
	if err != nil {
		return err
	}
	o.lock.Lock()
	defer o.lock.Unlock()
	cliM := new(GrpcClient)
	cliM.watcher = w
	cliM.serverName = serverName
	if target != nil {
		cliM.serverName = target.Name
	}
	cliM.conn = cli
	cliM.register = target
	cliM.name = serverName
	O.m[serverName] = cliM
	return nil
}
