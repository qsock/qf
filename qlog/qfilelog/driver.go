package qfilelog

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/qsock/qf/qlog"
	"github.com/qsock/qf/qlog/types"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"
	"time"
)

var (
	drv *driver
)

type driver struct {
	level    types.LEVEL
	cfg      Config
	host     string
	ch       chan *Atom
	bytePool *sync.Pool
	f        *os.File
	w        *bufio.Writer
	key      string
}

// 打开日志
func init() {
	drv = &driver{
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

	// 注册进去
	qlog.Register(Name(), drv)
}

func Name() string {
	return "file"
}

func (d *driver) Open(kv map[string]interface{}) error {
	b, _ := json.Marshal(kv)
	cfg := Config{}
	if err := json.Unmarshal(b, &cfg); err != nil {
		return err
	}
	cfg.check()
	d.cfg = cfg

	err := os.MkdirAll(d.cfg.Dir, 0755)
	if err != nil {
		return err
	}
	drv.f, err = os.OpenFile(d.cfg.fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	drv.w = bufio.NewWriterSize(drv.f, 1024*1024)
	d.run()
	return nil
}

func (d *driver) Close() error {
	if err := d.w.Flush(); err != nil {
		return err
	}
	close(d.ch)
	return d.f.Close()
}

func (d *driver) SetLevel(l types.LEVEL) types.IDriver {
	d.level = l
	return d
}

func (d *driver) Ctx(ctx context.Context) types.ILog {
	l := d.Logger()
	return l.Ctx(ctx)
}

func (d *driver) Logger(depth ...int) types.ILog {
	if len(depth) > 0 {
		return &Atom{line: depth[0]}
	}
	return &Atom{}
}

func (d *driver) CtxKey(key string) types.IDriver {
	d.key = key
	return d
}

func (d *driver) run() {
	go d.flush()
	go d.start()
}

func (d *driver) flush() {
	for range time.NewTicker(time.Second).C {
		d.ch <- nil
	}
}

func (d *driver) logname() string {
	t := fmt.Sprintf("%s", time.Now())[:19]
	tt := strings.Replace(
		strings.Replace(
			strings.Replace(t, "-", "", -1),
			" ", "", -1),
		":", "", -1)
	return fmt.Sprintf("%s.%s", d.cfg.fileName, tt)
}

func (d *driver) start() {
	for {
		a := <-d.ch
		if a != nil {
			b := d.bytes(a)
			_, _ = d.w.Write(b)
			continue
		}

		_ = d.w.Flush()
		fileInfo, err := os.Stat(d.cfg.fileName)
		if err != nil && os.IsNotExist(err) {
			_ = d.f.Close()
			d.f, _ = os.OpenFile(d.cfg.fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			d.w.Reset(d.f)
		}
		if fileInfo.Size() > d.cfg.fileSize {
			_ = d.f.Close()
			ofile := d.logname()
			_ = os.Rename(d.cfg.fileName, ofile)
			if d.cfg.UseGzip {
				go exec.Command("gzip", ofile).Output()
			}
			d.f, _ = os.OpenFile(d.cfg.fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			d.w.Reset(d.f)
			d.rm()
		}
		continue
	}
}

func (d *driver) rm() {
	if out, err := exec.Command("ls", d.cfg.Dir).Output(); err == nil {
		files := bytes.Split(out, []byte("\n"))
		total, idx := len(files)-1, 0
		for i := total; i >= 0; i-- {
			file := path.Join(d.cfg.Dir, string(files[i]))
			if strings.HasPrefix(file, d.cfg.fileName) && file != d.cfg.fileName {
				idx++
				if idx > d.cfg.FileNum {
					_ = os.Remove(file)
				}
			}
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
