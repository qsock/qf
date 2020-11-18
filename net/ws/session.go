package ws

import (
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"sync"
	"time"
)

type Session struct {
	id string

	Request *http.Request
	Keys    map[string]interface{}

	//保持的连接
	conn *websocket.Conn

	// 输出
	output chan *Packet

	// 持有一个全局变量
	ws *Server

	// 读写锁
	lock *sync.RWMutex

	// 是否已关闭
	isClosed bool
}

func (s *Session) writeLoop() {
	defer s.close()
	for {
		select {
		case msg, ok := <-s.output:
			if !ok {
				return
			}
			err := s.writeSync(msg)
			if err != nil {
				s.ws.errorHandler(s, err)
				return
			}
			if msg.t == websocket.TextMessage {
				s.ws.messageSentHandler(s, msg.msg)
			}
			if msg.t == websocket.BinaryMessage {
				s.ws.messageSentHandlerBinary(s, msg.msg)
			}
			if msg.t == websocket.CloseMessage {
				s.ws.closeHandler(s, 0, "")
				return
			}
		}
	}
}

func (s *Session) readLoop() {
	s.conn.SetReadLimit(s.ws.c.MaxMessageSize)
	s.conn.SetCloseHandler(func(code int, text string) error {
		_ = s.Close()
		s.ws.closeHandler(s, code, text)
		return nil
	})

	for {
		if s.IsClosed() {
			return
		}
		_ = s.conn.SetReadDeadline(time.Now().Add(s.ws.c.TTLTimeout))
		t, msg, err := s.conn.ReadMessage()
		if err != nil {
			_ = s.Close()
			s.ws.errorHandler(s, err)
			return
		}

		if t == websocket.TextMessage {
			s.ws.messageHandler(s, msg)
		}
		if t == websocket.BinaryMessage {
			s.ws.messageHandlerBinary(s, msg)
		}
	}
}

func (s *Session) close() {
	if !s.IsClosed() {
		s.lock.Lock()
		s.isClosed = true
		// 关闭连接
		_ = s.conn.Close()
		close(s.output)
		s.lock.Unlock()
	}
}

func (s *Session) closed() bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.isClosed
}

func (s *Session) Close() error {
	return s.CloseWithMsg([]byte{})
}

func (s *Session) CloseWithMsg(msg []byte) error {
	if s.IsClosed() {
		return ErrSessionClosed
	}
	m := new(Packet)
	m.t = websocket.CloseMessage
	m.msg = msg
	_ = s.writeAsync(m)
	return nil
}

// session 是否关闭
func (s *Session) IsClosed() bool {
	return s.closed()
}

func (s *Session) writeAsync(msg *Packet) error {
	if s.IsClosed() {
		s.ws.errorHandler(s, ErrSessionClosed)
		return ErrSessionClosed
	}
	select {
	case s.output <- msg:
	default:
		// 消息发送太大
		s.ws.errorHandler(s, ErrSessionFulled)
	}
	return nil
}

func (s *Session) writeSync(msg *Packet) error {
	if s.IsClosed() {
		return ErrSessionClosed
	}
	// 设置写超时
	_ = s.conn.SetWriteDeadline(time.Now().Add(s.ws.c.WriteTimeout))
	// gorilla 写消息
	return s.conn.WriteMessage(msg.t, msg.msg)
}

func (s *Session) Ping() error {
	msg := new(Packet)
	msg.t = websocket.PingMessage
	msg.msg = []byte{}
	return s.writeAsync(msg)
}

func (s *Session) Pong() error {
	msg := new(Packet)
	msg.t = websocket.PongMessage
	msg.msg = []byte{}
	return s.writeAsync(msg)
}

func (s *Session) Write(msg []byte) error {
	if s.IsClosed() {
		return ErrSessionClosed
	}
	m := new(Packet)
	m.t = websocket.TextMessage
	m.msg = msg
	return s.writeAsync(m)
}

func (s *Session) WriteBinary(msg []byte) error {
	if s.IsClosed() {
		return ErrSessionClosed
	}
	m := new(Packet)
	m.t = websocket.BinaryMessage
	m.msg = msg
	return s.writeAsync(m)
}

// kv setting
func (s *Session) Set(key string, val interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.Keys == nil {
		s.Keys = make(map[string]interface{})
	}
	s.Keys[key] = val
}

func (s *Session) Get(key string) (interface{}, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	if s.Keys == nil {
		return nil, false
	}
	v, b := s.Keys[key]
	return v, b
}

func (s *Session) MustGet(key string) interface{} {
	if v, b := s.Get(key); b {
		return v
	}
	panic("qws:key not exists " + key)
}

func (s *Session) GetId() string {
	return s.id
}

func (s *Session) RemoteAddr() net.Addr {
	return s.conn.RemoteAddr()
}

func (s *Session) LocalAddr() net.Addr {
	return s.conn.LocalAddr()
}
