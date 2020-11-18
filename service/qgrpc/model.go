package qgrpc

import (
	"github.com/qsock/qf/service/qetcd"
	"google.golang.org/grpc"
	"sync"
)

type RegisterModel struct {
	Name      string `json:"name,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	Addr      string `json:"addr,omitempty"`
	Hostname  string `json:"hostname,omitempty"`
	Ext       string `json:"ext,omitempty"`
}

type GrpcClient struct {
	name       string
	serverName string
	conn       *grpc.ClientConn
	watcher    *watcher
	register   *RegisterModel
}

type Op struct {
	etcdClient *qetcd.Op
	// 配置文件
	C *Config
	// 本地的地址
	LocalIp string
	m       map[string]*GrpcClient
	lock    *sync.RWMutex
}
