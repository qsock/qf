package qgrpc

import "errors"

var (
	ErrNoSuchGrpcClient = errors.New("no such grpc client")
	ErrConf             = errors.New("conf error")
)
