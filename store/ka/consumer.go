package ka

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/qsock/qf/qlog"
	"math/rand"
	"strings"
	"sync"
)

type Config struct {
	Brokers []string `toml:"brokers"`
	Topic   string   `toml:"topic"`
	Group   string   `toml:"group"`
	Workers int      `toml:"workers"`
	Oldest  bool     `toml:"oldest"`
}

type Consumer struct {
	test    bool
	group   sarama.ConsumerGroup
	workers int
	handler func([]byte)
	stop    chan bool
	wg      *sync.WaitGroup
	cfg     *Config
}

const (
	TEST_PREFIX = "test_consumer_"
)

type consumerGroupHandler struct {
	consumer *Consumer
}

func (consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case msg := <-claim.Messages():
			h.consumer.handler(msg.Value)
			// 测试用的group不标记处理完毕，防止破坏集群数据
			if !h.consumer.test {
				sess.MarkMessage(msg, "done")
			}
		case <-h.consumer.stop:
			return nil

		}
	}
}

func NewConsumer(cfg *Config, hdl func([]byte)) *Consumer {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Version = sarama.V2_3_0_0

	if cfg.Oldest {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	} else {
		config.Consumer.Offsets.Initial = sarama.OffsetNewest
	}

	if err := config.Validate(); err != nil {
		qlog.Get().Logger().Error("kafka||err:%#v", err)
		return nil
	}

	group := cfg.Group
	if group == "" {
		group = fmt.Sprintf("%s%.8x", TEST_PREFIX, rand.Int63())
	}

	// Start with a client
	client, err := sarama.NewClient(cfg.Brokers, config)
	if err != nil {
		qlog.Get().Logger().Error("kafka||err:%#v", err)
		return nil
	}

	// Start a new consumer group
	consumerGroup, err := sarama.NewConsumerGroupFromClient(group, client)
	if err != nil {
		qlog.Get().Logger().Error("kafka||err:%#v", err)
		return nil
	}

	// Track errors
	go func() {
		for err := range consumerGroup.Errors() {
			qlog.Get().Logger().Error("kafka||err:%#v", err)
		}
	}()

	return &Consumer{
		test:    strings.HasPrefix(group, TEST_PREFIX),
		group:   consumerGroup,
		workers: cfg.Workers,
		handler: hdl,
		stop:    make(chan bool),
		wg:      &sync.WaitGroup{},
		cfg:     cfg,
	}
}

func (self *Consumer) Run() {
	for i := 0; i < self.workers; i++ {
		go self.Worker()
	}
}

func (self *Consumer) Stop() {
	close(self.stop)
	self.wg.Wait()
	_ = self.group.Close()
}

func (self *Consumer) Wait() {
	self.wg.Wait()
	_ = self.group.Close()
}

func (self *Consumer) Worker() {
	self.wg.Add(1)
	defer self.wg.Done()
	ctx := context.Background()
	topics := []string{self.cfg.Topic}
	handler := consumerGroupHandler{self}
	err := self.group.Consume(ctx, topics, handler)
	if err != nil {
		qlog.Get().Logger().Error("kafka||err:%#v", err)
		return
	}
}
