package ws

import (
	"github.com/qsock/qf/concurrent"
	"sync"
)

type G struct {
	sessions  *concurrent.IdMap
	broadcast chan *Packet
	exit      chan *Packet
	// 注册进去
	register chan *Session
	// 被解除注册
	unregister chan *Session
	isClosed   bool
	lock       *sync.RWMutex
}

func newG() *G {
	g := new(G)
	g.sessions = concurrent.NewIdMap()
	g.broadcast = make(chan *Packet, 32)
	g.exit = make(chan *Packet)
	g.register = make(chan *Session, 10)
	g.unregister = make(chan *Session, 10)
	g.lock = new(sync.RWMutex)
	return g
}

func (g *G) closed() bool {
	g.lock.RLock()
	defer g.lock.RUnlock()
	return g.isClosed
}

func (g *G) len() int {
	return int(g.sessions.Count())
}

func (g *G) loop() {
	for {
		select {
		case s := <-g.register:
			g.sessions.SetS(s.id, s)
		case s := <-g.unregister:
			g.sessions.DelS(s.id)
			_ = s.Close()
		case s := <-g.broadcast:
			for _, m := range g.sessions.Ms {
				m.Range(func(k, v interface{}) bool {
					sess, ok := v.(*Session)
					if !ok {
						return false
					}
					if s.filter != nil && !s.filter(sess) {
						return false
					}
					_ = sess.writeAsync(s)
					return true
				})
			}
		case s := <-g.exit:
			for _, m := range g.sessions.Ms {
				m.Range(func(k, v interface{}) bool {
					sess, ok := v.(*Session)
					if !ok {
						return false
					}
					_ = sess.writeAsync(s)
					g.unregister <- sess
					return true
				})
			}
			g.lock.Lock()
			g.isClosed = true
			g.lock.Unlock()
			return
		}
	}
}

func (g *G) GetSessionById(id string) *Session {
	sess, ok := g.sessions.GetS(id)
	if !ok {
		return nil
	}
	return sess.(*Session)
}

// 替换session
func (g *G) ReplaceSessionById(session *Session, id string) {
	{
		if session.id == id {
			return
		}
		sess := g.GetSessionById(id)
		if sess != nil {
			_ = sess.Close()
		}
	}

	g.unregister <- session
	session.id = id
	g.register <- session
}
