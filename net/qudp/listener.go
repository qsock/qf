package qudp

import (
	"github.com/qsock/qf/qlog"
	"github.com/qsock/qf/store/cache"
	"net"
)

type Listener struct {
	conn *net.UDPConn

	newConn chan *Session
	newErr  chan error

	cache *cache.Cache

	s *Server
}

func listen(conn *net.UDPConn, s *Server) *Listener {
	l := new(Listener)
	l.conn = conn
	l.newConn = make(chan *Session, maxClientSize)
	l.newErr = make(chan error, maxErrSize)
	l.cache = cache.New(bufferSize)
	l.s = s
	go l.listen()
	return l
}

func (l *Listener) listen() {
	data := make([]byte, maxBufferSize)
	for {
		n, addr, err := l.conn.ReadFromUDP(data)
		if err != nil {
			qlog.Error(err.Error())
			l.newErr <- err
			continue
		}
		client, ok := l.cache.Get(addr.String())
		if !ok {
			client = newClient(l.conn, addr, l.s)
			// 保持1天
			l.cache.SetEx(addr.String(), client, 86400)
			l.newConn <- client.(*Session)
		}
		b := make([]byte, n)
		copy(b, data[:n])
		client.(*Session).recvChan <- b
	}
}

func (l *Listener) Close() error {
	_ = l.conn.Close()
	close(l.newConn)
	close(l.newErr)
	for _, k := range l.cache.Keys() {
		client, ok := l.cache.Get(k)
		if ok {
			_ = client.(*Session).Close()
		}
		l.cache.Del(k)
	}
	l.cache.Clear()
	return nil
}
