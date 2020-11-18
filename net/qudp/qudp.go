package qudp

import (
	"github.com/qsock/qf/qlog"
	"net"
)

type Server struct {
	addr *net.UDPAddr
	l    *Listener

	//回调方法
	messageHandler     HandleMessageFunc
	messageSentHandler HandleMessageFunc

	errorHandler   HandleErrorFunc
	connectHandler HandleSessionFunc
}

func New(addr *net.UDPAddr) *Server {
	s := new(Server)
	s.addr = addr

	s.messageHandler = func(*Session, []byte) {}
	s.messageSentHandler = func(*Session, []byte) {}

	s.errorHandler = func(*Session, error) {}
	s.connectHandler = func(*Session) {}

	return s
}

func (s *Server) Listen() error {
	conn, err := net.ListenUDP("udp", s.addr)
	if err != nil {
		return err
	}
	s.l = listen(conn, s)
	s.run()
	return nil
}

func (s *Server) run() {
	for {
		select {
		case conn, ok := <-s.l.newConn:
			{
				if ok {
					s.connectHandler(conn)
				}
			}
		case err, ok := <-s.l.newErr:
			{
				if ok {
					qlog.Get().Logger().Error(err)
				}
			}
		}
	}
}

func (s *Server) Close() error {
	return s.l.Close()
}
