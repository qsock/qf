package qgrpc

import (
	"context"
	"encoding/json"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"strings"
	"time"
)

const (
	MetaKey = "meta"
)

//根据server-name 获得对应的连接池
func GetConn(name string) *grpc.ClientConn {
	O.lock.RLock()
	defer O.lock.RUnlock()
	c := O.m[name]
	if c == nil {
		return nil
	}
	return c.conn
}

//100ms超时
func CallIn100ms(ctx context.Context, method string, args, reply interface{}) (err error) {
	return CallWithTimeout(ctx, method, args, reply, 100*time.Millisecond)
}

func CallIn200ms(ctx context.Context, method string, args, reply interface{}) (err error) {
	return CallWithTimeout(ctx, method, args, reply, 200*time.Millisecond)
}

func CallIn500ms(ctx context.Context, method string, args, reply interface{}) (err error) {
	return CallWithTimeout(ctx, method, args, reply, 500*time.Millisecond)
}

func CallWithTimeout(ctx context.Context, method string, args, reply interface{}, t time.Duration) (err error) {
	nctx, cancel := context.WithTimeout(ctx, t)
	defer cancel()
	return Call(nctx, method, args, reply)
}

func Call(ctx context.Context, method string, args, reply interface{}) (err error) {
	ss := strings.Split(method, "/")
	var (
		rpcServerName string
	)
	if len(ss) != 3 {
		return ErrNoSuchGrpcClient
	}
	rpcServerName = ss[1]
	conn := GetConn(rpcServerName)
	if conn == nil {
		return ErrNoSuchGrpcClient
	}

	return CallWithOption(ctx, conn, method, args, reply, grpc.WaitForReady(false))
}

func CallWithServerNameTimeout(ctx context.Context, serverName, method string, args, reply interface{}, t time.Duration) (err error) {
	nctx, cancel := context.WithTimeout(ctx, t)
	defer cancel()
	return CallWithServerName(nctx, serverName, method, args, reply)
}

func CallWithServerName(ctx context.Context, serverName, method string, args, reply interface{}) (err error) {
	conn := GetConn(serverName)
	if conn == nil {
		return ErrNoSuchGrpcClient
	}
	return CallWithOption(ctx, conn, method, args, reply, grpc.WaitForReady(false))
}

// method->/rpc.Rpc/Ping
func CallWithOption(ctx context.Context, conn *grpc.ClientConn, method string, args, reply interface{}, opts ...grpc.CallOption) (err error) {
	ctx = ctx2GrpcCtx(ctx)
	err = conn.Invoke(ctx, method, args, reply, opts...)
	return
}

func ctx2GrpcCtx(ctx context.Context) context.Context {
	meta := ctx.Value(MetaKey)
	if meta == nil {
		return ctx
	}
	b, _ := json.Marshal(meta)
	return metadata.NewOutgoingContext(
		ctx,
		metadata.New(
			map[string]string{MetaKey: string(b)}),
	)
}
