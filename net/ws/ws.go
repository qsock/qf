package ws

import (
	"github.com/gorilla/websocket"
	"github.com/qsock/qf/util/uuid"
	"net/http"
	"sync"
)

type Server struct {
	c        *Config
	upGrader *websocket.Upgrader

	//回调方法
	messageHandler           HandleMessageFunc
	messageHandlerBinary     HandleMessageFunc
	messageSentHandler       HandleMessageFunc
	messageSentHandlerBinary HandleMessageFunc

	errorHandler      HandleErrorFunc
	closeHandler      HandleCloseFunc
	connectHandler    HandleSessionFunc
	disconnectHandler HandleSessionFunc

	shutdownHandler HandleEmptyFunc

	g *G
}

func New(conf ...*Config) *Server {
	upGrader := new(websocket.Upgrader)
	upGrader.ReadBufferSize = 1024
	upGrader.WriteBufferSize = 1024
	upGrader.CheckOrigin = func(r *http.Request) bool { return true }

	g := newG()

	go g.loop()

	w := new(Server)
	if len(conf) != 0 {
		w.c = conf[0]
	} else {
		w.c = DefaultConfig()
	}
	w.upGrader = upGrader
	w.messageHandler = func(*Session, []byte) {}
	w.messageHandlerBinary = func(*Session, []byte) {}
	w.messageSentHandler = func(*Session, []byte) {}
	w.messageSentHandlerBinary = func(*Session, []byte) {}

	w.errorHandler = func(*Session, error) {}
	w.closeHandler = func(*Session, int, string) {}
	w.connectHandler = func(*Session) {}
	w.disconnectHandler = func(*Session) {}

	w.shutdownHandler = func() {}

	w.g = g
	return w
}

func (w *Server) HandleRequest(writer http.ResponseWriter, r *http.Request) error {
	return w.HandleRequestWithKeys(writer, r, nil)
}

func (w *Server) HandleRequestWithKeys(writer http.ResponseWriter, r *http.Request, keys map[string]interface{}) error {
	if w.g.closed() {
		w.shutdownHandler()
		return nil
	}
	conn, err := w.upGrader.Upgrade(writer, r, writer.Header())
	if err != nil {
		return err
	}

	sess := new(Session)
	sess.id = uuid.NewString()
	sess.Request = r
	sess.Keys = keys
	sess.conn = conn
	sess.output = make(chan *Packet, w.c.MessageBufferSize)
	sess.ws = w
	sess.lock = new(sync.RWMutex)

	w.g.register <- sess

	w.connectHandler(sess)

	go sess.writeLoop()
	// 死循环了
	sess.readLoop()
	// 解除注册
	if !w.g.closed() {
		w.g.unregister <- sess
	}

	w.disconnectHandler(sess)
	return nil
}

// 广播
func (w *Server) BroadcastFilter(msg []byte, fn func(*Session) bool) error {
	m := new(Packet)
	m.t = websocket.TextMessage
	m.msg = msg
	m.filter = fn

	w.g.broadcast <- m
	return nil
}

func (w *Server) Broadcast(msg []byte) error {
	return w.BroadcastFilter(msg, nil)
}

func (w *Server) BroadcastExcept(msg []byte, s *Session) error {
	return w.BroadcastFilter(msg, func(sess *Session) bool {
		return sess.id != s.id
	})
}

func (w *Server) BroadcastMultiple(msg []byte, ss []*Session) error {
	for _, sess := range ss {
		if err := sess.Write(msg); err != nil {
			return err
		}
	}
	return nil
}

// binary broadcast
func (w *Server) BroadcastBinaryFilter(msg []byte, fn func(*Session) bool) error {
	m := new(Packet)
	m.t = websocket.BinaryMessage
	m.msg = msg
	m.filter = fn
	w.g.broadcast <- m
	return nil
}

func (w *Server) BroadcastBinary(msg []byte) error {
	return w.BroadcastBinaryFilter(msg, nil)
}

func (w *Server) BroadcastBinaryExcept(msg []byte, s *Session) error {
	return w.BroadcastBinaryFilter(msg, func(sess *Session) bool {
		return sess.id != s.id
	})
}

func (w *Server) BroadcastBinaryMultiple(msg []byte, ss []*Session) error {
	for _, sess := range ss {
		if err := sess.WriteBinary(msg); err != nil {
			return err
		}
	}
	return nil
}

// close
func (w *Server) Close() error {
	return w.CloseWithMsg([]byte{})
}

func (w *Server) CloseWithMsg(msg []byte) error {
	m := new(Packet)
	m.t = websocket.CloseMessage
	m.msg = msg
	w.g.exit <- m

	return nil
}

func (w *Server) Len() int {
	return w.g.len()
}

func (w *Server) IsClosed() bool {
	return w.g.closed()
}

func (w *Server) GetSessionById(id string) *Session {
	return w.g.GetSessionById(id)
}

func (w *Server) ReplaceSessionById(session *Session, id string) {
	w.g.ReplaceSessionById(session, id)
}

func (w *Server) FmtCloseMsg(code int, text string) []byte {
	return websocket.FormatCloseMessage(code, text)
}
