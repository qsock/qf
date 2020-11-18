package bk

import (
	"github.com/kr/beanstalk"
	"strings"
	"sync"
	"time"
)

type Config struct {
	Addrs   []string `toml:"addrs"`
	Tube    string   `toml:"tube"`
	Workers int      `toml:"workers"`
}

type Consumer struct {
	c       *Config
	ch      chan []byte
	handler func([]byte)
	stop    chan struct{}
	isRun   bool
	wg      *sync.WaitGroup
}

func NewConsumer(c *Config, callback func([]byte)) *Consumer {
	consumer := new(Consumer)
	consumer.handler = callback
	consumer.ch = make(chan []byte)
	consumer.c = c
	consumer.stop = make(chan struct{})
	consumer.wg = new(sync.WaitGroup)
	consumer.isRun = true
	consumer.recv()
	return consumer
}

func (this *Consumer) recv() {
	for _, addr := range this.c.Addrs {
		go this.doRecv(strings.Split(addr, "|")[0])
	}
}

func (this *Consumer) Run() {
	for i := 0; i < this.c.Workers; i++ {
		go this.work()
	}
}

func (this *Consumer) Stop() {
	this.isRun = false
	close(this.stop)
	close(this.ch)
	this.wg.Wait()
}

func (this *Consumer) work() {
	this.wg.Add(1)
	for {
		msg, isOpen := <-this.ch
		if !isOpen {
			this.wg.Done()
			return
		}
		this.handler(msg)
	}
}

func (this *Consumer) doRecv(addr string) {
	connTube := connectBk(addr, this.c.Tube, nil)

	for this.isRun {
		if connTube == nil {
			connTube = connectBk(addr, this.c.Tube, connTube)
			time.Sleep(time.Second)
			continue
		}

		id, body, err := connTube.Reserve(time.Second)

		if err == beanstalk.ErrTimeout ||
			(err != nil && strings.Contains(err.Error(), "timeout") && !strings.Contains(err.Error(), " tcp ")) {
			continue
		}

		if err != nil {
			connTube = connectBk(addr, this.c.Tube, connTube)
			time.Sleep(time.Second)
			continue
		}

		select {
		case <-this.stop:
			connTube.Conn.Release(id, 1024, 0)
			return
		case this.ch <- body:
			if err := connTube.Conn.Delete(id); err != nil {

			}
		}
	}
}

func connectBk(addr string, tube string, connTube *beanstalk.TubeSet) *beanstalk.TubeSet {
	if connTube == nil {
		conn, err := beanstalk.Dial("tcp", addr)
		if err != nil {
			return nil
		}

		connTube = beanstalk.NewTubeSet(conn, tube)
		return connTube
	}

	if connTube.Conn != nil {
		connTube.Conn.Close()
		conn, err := beanstalk.Dial("tcp", addr)
		if err != nil {
			return nil
		}
		connTube = beanstalk.NewTubeSet(conn, tube)
		return connTube
	}

	return nil
}
