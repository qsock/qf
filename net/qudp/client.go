package qudp

import (
	"net"
)

type Session struct {
	conn *net.UDPConn
	addr *net.UDPAddr

	recvChan chan []byte
	recvErr  chan error

	sendChan chan []byte
	sendErr  chan error

	s *Server
}

func (c *Session) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *Session) RemoteAddr() net.Addr {
	return c.addr
}

func newClient(conn *net.UDPConn, addr *net.UDPAddr, s *Server) *Session {
	c := &Session{conn: conn, addr: addr}
	c.recvChan = make(chan []byte, 1<<16)
	c.sendChan = make(chan []byte, 1<<16)

	c.recvErr = make(chan error, 1<<3)
	c.sendErr = make(chan error, 1<<3)
	c.s = s
	go c.run()
	return c
}

func (c *Session) run() {
	for {
		select {
		case data := <-c.recvChan:
			c.s.messageHandler(c, data)
		case err := <-c.recvErr:
			c.s.errorHandler(c, err)
		case data := <-c.sendChan:
			_, err := c.conn.WriteToUDP(data, c.addr)
			if err != nil {
				c.sendErr <- err
			} else {
				c.s.messageSentHandler(c, data)
			}
		case err := <-c.sendErr:
			c.s.errorHandler(c, err)
		}
	}
}

func (c *Session) Write(b []byte) (int, error) {
	c.sendChan <- b
	return len(b), nil
}

func (c *Session) Read(b []byte) (int, error) {
	select {
	case data := <-c.recvChan:
		copy(b, data)
		return len(data), nil
	case err := <-c.recvErr:
		return 0, err
	}
}

func (c *Session) ReadAll() (int, []byte, error) {
	select {
	case data := <-c.recvChan:
		return len(data), data, nil
	case err := <-c.recvErr:
		return 0, nil, err
	}
}

func (c *Session) Close() error {
	close(c.sendChan)
	close(c.sendErr)
	close(c.recvChan)
	close(c.sendErr)
	return nil
}
