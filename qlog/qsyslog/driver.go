package qsyslog

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/qsock/qf/qlog"
	"github.com/qsock/qf/qlog/types"
	"log/syslog"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	drv *driver
	mu  = new(sync.Mutex)
)

type driver struct {
	level    types.LEVEL
	cfg      Config
	host     string
	ch       chan *Atom
	bytePool *sync.Pool
	syslogW  *syslog.Writer
	key      string
	wg       *sync.WaitGroup
	closed   bool
}

func init() {
	drv = &driver{
		wg:       new(sync.WaitGroup),
		ch:       make(chan *Atom, 102400),
		bytePool: &sync.Pool{New: func() interface{} { return new(bytes.Buffer) }},
	}

	host, _ := os.Hostname()
	ss := strings.Split(host, "-")
	if len(ss) < 2 {
		drv.host = host
	} else {
		drv.host = ss[len(ss)-2] + ss[len(ss)-1]
	}
	qlog.Register(Name(), drv)
}

func Name() string {
	return "syslog"
}

func (d *driver) Open(kv map[string]interface{}) error {
	b, _ := json.Marshal(kv)
	cfg := Config{}
	err := json.Unmarshal(b, &cfg)
	if err != nil {
		return err
	}
	d.cfg = cfg

	d.syslogW, err = syslog.Dial(d.cfg.Network, d.cfg.Addr, syslog.LOG_USER|syslog.LOG_INFO, d.cfg.Tag)
	if err != nil {
		return err
	}
	d.run()
	return nil
}

func (d *driver) Close() error {
	mu.Lock()
	d.closed = true
	mu.Unlock()
	d.wg.Wait()
	close(d.ch)
	return d.syslogW.Close()
}

func (d *driver) SetLevel(l types.LEVEL) types.IDriver {
	d.level = l
	return d
}

func (d *driver) Ctx(ctx context.Context) types.ILog {
	l := d.Logger()
	return l.Ctx(ctx)
}

func (d *driver) Logger() types.ILog {
	return new(Atom)
}

func (d *driver) CtxKey(key string) types.IDriver {
	d.key = key
	return d
}

func (d *driver) run() {
	go d.start()
}

func (d *driver) start() {
	for {
		a, ok := <-d.ch
		if !ok {
			return
		}
		d.wg.Done()

		b := d.bytes(a)
		switch a.level {
		case types.DEBUG:
			_ = d.syslogW.Debug(string(b))
		case types.INFO:
			_ = d.syslogW.Info(string(b))
		case types.WARNING:
			_ = d.syslogW.Warning(string(b))
		case types.ERROR:
			_ = d.syslogW.Err(string(b))
		case types.FATAL:
			_ = d.syslogW.Crit(string(b))
		}
	}
}

func (d *driver) genTime() []byte {
	now := time.Now()
	_, month, day := now.Date()
	hour, minute, second := now.Clock()
	return []byte{byte(month/10) + 48, byte(month%10) + 48, '-', byte(day/10) + 48, byte(day%10) + 48, ' ',
		byte(hour/10) + 48, byte(hour%10) + 48, ':', byte(minute/10) + 48, byte(minute%10) + 48, ':',
		byte(second/10) + 48, byte(second%10) + 48, ' '}
}

func (d *driver) bytes(a *Atom) []byte {
	w := d.bytePool.Get().(*bytes.Buffer)
	defer func() {
		recover()
		w.Reset()
		d.bytePool.Put(w)
	}()
	ctxStr := ""

	if a.ctx != nil {
		if len(d.key) == 0 {
			d.key = types.MetaKey
		}
		if val, ok := a.ctx.Value(d.key).(map[string]string); ok && val != nil {
			for k, v := range val {
				ctxStr += fmt.Sprintf("%s:%s ", k, v)
			}
		}
	}

	w.Write(d.genTime())
	_, _ = fmt.Fprintf(w, "%s %s %s:%d %s", d.host, types.LevelText[a.level], a.file, a.line, ctxStr)
	if len(a.format) < 1 {
		for _, arg := range a.args {
			w.WriteByte(' ')
			_, _ = fmt.Fprint(w, arg)
		}
	} else {
		_, _ = fmt.Fprintf(w, a.format, a.args...)
	}
	w.WriteByte(10)
	b := make([]byte, w.Len())
	copy(b, w.Bytes())
	return b
}
