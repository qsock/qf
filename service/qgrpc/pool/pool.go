package pool

import (
	"errors"
	"github.com/qsock/qf/concurrent"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/keepalive"
	"sync"
	"time"
)

var (
	ErrNoConn = errors.New("no conn")
)

//GRPCPool pool info
type GRPCPool struct {
	timeout time.Duration
	sync.RWMutex
	conns   chan *grpcIdleConn
	factory func() (*grpc.ClientConn, error)
	close   func(*grpc.ClientConn) error
	targets []string
	n       *concurrent.AtomicUint64
}

type grpcIdleConn struct {
	conn *grpc.ClientConn
	t    time.Time
}

func Init(targets []string, timeout time.Duration) (*GRPCPool, error) {
	pool := new(GRPCPool)
	pool.n = concurrent.NewAtomicUint64(0)
	pool.timeout = timeout
	pool.targets = make([]string, 3*len(targets))
	for _, target := range targets {
		pool.targets = append(pool.targets, target, target, target)
	}
	pool.conns = make(chan *grpcIdleConn, len(targets)*3)
	pool.factory = pool.PoolFactory
	pool.close = func(c *grpc.ClientConn) error {
		return c.Close()
	}
	for i := 0; i < len(pool.targets); i++ {
		conn, err := pool.factory()
		if err != nil {
			return nil, err
		}
		pool.conns <- &grpcIdleConn{conn: conn, t: time.Now()}
	}
	return pool, nil
}

func (c *GRPCPool) PoolFactory() (*grpc.ClientConn, error) {
	target := c.nextTarget()
	var kacp = keepalive.ClientParameters{
		Time:                1 * time.Second,  // send pings every 1 seconds if there is no activity
		Timeout:             time.Second * 10, // wait 10 second for ping ack before considering the connection dead
		PermitWithoutStream: true,             // send pings even without active streams
	}

	conn, err := grpc.Dial(target,
		grpc.WithInsecure(),
		grpc.WithConnectParams(grpc.ConnectParams{MinConnectTimeout: 3 * time.Second}),
		grpc.WithKeepaliveParams(kacp),
	)
	return conn, err
}

// 下一个target
func (c *GRPCPool) nextTarget() string {
	n := c.n.IncrementAndGet()
	c.RLock()
	defer c.RUnlock()
	g := n % uint64(len(c.targets))
	return c.targets[int(g)]
}

//Get get from pool
func (c *GRPCPool) Get() (*grpc.ClientConn, error) {
	c.RLock()
	defer c.RUnlock()
	conns := c.conns

	if conns == nil {
		return nil, ErrNoConn
	}
	for {
		select {
		case wrapConn := <-conns:
			if wrapConn == nil {
				return nil, ErrNoConn
			}
			state := wrapConn.conn.GetState()
			if state == connectivity.Ready {
				defer func() {
					c.conns <- wrapConn
				}()
				return wrapConn.conn, nil
			}
		case <-time.After(c.timeout):
			conn, err := c.factory()
			if err != nil {
				return nil, err
			}
			return conn, nil
		}
	}
}
