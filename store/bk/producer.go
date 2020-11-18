package bk

import (
	"errors"
	"github.com/kr/beanstalk"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	p              map[string]*Producer = make(map[string]*Producer)
	ProducerExists error                = errors.New("producer has exists")
	NoLiveConn     error                = errors.New("no conn alive")
)

type Producer struct {
	c  *Config
	ch chan []string
	sync.RWMutex
	m map[string]*beanstalk.Tube
}

func AddProducer(name string, c *Config) error {
	_, ok := p[name]
	if ok {
		return ProducerExists
	}

	producer := new(Producer)
	producer.c = c
	producer.m = make(map[string]*beanstalk.Tube)

	sm, l := ParseWeight(c.Addrs)

	producer.ch = make(chan []string, len(l))

	for i := 0; i < len(l); i++ {
		list := make([]string, len(l))
		for j := 0; j < len(l); j++ {
			list[j] = l[(i+j)%len(l)]
		}
		producer.ch <- list
	}

	for addr, _ := range sm {
		conn, err := beanstalk.Dial("tcp", addr)
		if err != nil {
			producer.m[addr] = nil
		} else {
			producer.m[addr] = &beanstalk.Tube{Name: c.Tube, Conn: conn}
		}
	}

	p[name] = producer
	go producer.ping()
	return nil
}

func ParseWeight(addrs []string) (map[string]int, []string) {
	//127.0.0.1:11300|1 权重值只仅支持0/1/2/3
	ret := make(map[string]int)
	for _, addr := range addrs {
		ss := strings.Split(addr, "|")
		if len(ss) == 1 {
			ret[ss[0]] = 1
		}

		if len(ss) == 2 {
			weight, err := strconv.Atoi(ss[1])
			if err == nil && weight <= 3 && weight >= 1 {
				ret[ss[0]] = weight
			}
		}
	}

	l := []string{}
	for a, w := range ret {
		for i := 0; i < w; i++ {
			l = append(l, a)
		}
	}

	return ret, l
}

func GetProducer(name string) *Producer {
	return p[name]
}

func (this *Producer) GetList() []string {
	l := <-this.ch
	this.ch <- l
	return l
}

func (this *Producer) Put(body []byte, pri uint32, delay, ttr time.Duration) (id uint64, err error) {
	for _, addr := range this.GetList() {
		conn, ok := this.m[addr]
		if ok && conn != nil {
			this.RLock()
			id, err = conn.Put(body, pri, delay, ttr)
			this.RUnlock()
			if err == nil {
				return
			} else {
				mlog.Info(err)
				this.Lock()
				conn.Conn.Close()
				this.m[addr] = nil
				this.Unlock()
			}
		}
	}

	return 0, NoLiveConn
}

func (this *Producer) ping() {
	for range time.NewTicker(time.Second).C {
		this.doping()
	}
}

func (this *Producer) doping() {
	for addr, conn := range this.m {
		if conn == nil {
			c, err := beanstalk.Dial("tcp", addr)
			if err == nil {
				this.Lock()
				this.m[addr] = &beanstalk.Tube{Name: this.c.Tube, Conn: c}
				this.Unlock()
			}
		}
	}
}
