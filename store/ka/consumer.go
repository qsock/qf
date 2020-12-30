package ka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"math/rand"
	"sync"
)

type ConsumerConfig struct {
	// broker的集群地址
	Brokers []string `toml:"brokers"`
	// topic 的名称
	Topic string `toml:"topic"`
	// 消费组名称
	Group string `toml:"group"`

	// 多少个协程
	Workers int `toml:"workers"`
	// 是否从最老的开始消费
	Oldest bool `toml:"oldest"`
}

func (c *ConsumerConfig) Check() bool {
	if len(c.Topic) == 0 ||
		len(c.Brokers) == 0 {
		return false
	}
	if c.Group == "" {
		c.Group = fmt.Sprintf("%s%.8x", TestPrefix, rand.Int63())
	}
	if c.Workers == 0 {
		c.Workers = 10
	}
	return true
}

func NewConsumer(cfg *ConsumerConfig, name ...string) error {
	return NewConsumerWithConfig(cfg, nil, name...)
}

func NewConsumerWithInterceptor(cfg1 *ConsumerConfig, interceptor sarama.ConsumerInterceptor, name ...string) error {
	if !cfg1.Check() {
		return errors.New("config error")
	}
	cfg2 := sarama.NewConfig()
	cfg2.Consumer.Return.Errors = true
	cfg2.Version = sarama.V0_10_2_0
	if cfg1.Oldest {
		cfg2.Consumer.Offsets.Initial = sarama.OffsetOldest
	} else {
		cfg2.Consumer.Offsets.Initial = sarama.OffsetNewest
	}
	cfg2.Consumer.Interceptors = []sarama.ConsumerInterceptor{interceptor}
	return NewConsumerWithConfig(cfg1, cfg2, name...)
}

func NewConsumerWithConfig(cfg1 *ConsumerConfig, cfg2 *sarama.Config, name ...string) error {
	if !cfg1.Check() {
		return errors.New("config error")
	}
	if cfg2 == nil {
		cfg2 = sarama.NewConfig()
		cfg2.Consumer.Return.Errors = true
		cfg2.Version = sarama.V0_10_2_0
		if cfg1.Oldest {
			cfg2.Consumer.Offsets.Initial = sarama.OffsetOldest
		} else {
			cfg2.Consumer.Offsets.Initial = sarama.OffsetNewest
		}
	}
	if err := cfg2.Validate(); err != nil {
		return err
	}
	// Start with a client
	client, err := sarama.NewClient(cfg1.Brokers, cfg2)
	if err != nil {
		return err
	}
	// Start a new consumer group
	consumerGroup, err := sarama.NewConsumerGroupFromClient(cfg1.Group, client)
	if err != nil {
		return err
	}
	c := new(Consumer)
	c.cfg = cfg1
	c.config = cfg2
	c.group = consumerGroup
	c.wg = &sync.WaitGroup{}
	c.handlers = make(map[string]HandleConsumerFunc)
	c.ctx, c.cancel = context.WithCancel(context.Background())

	consumerName := gDefaultName
	if len(name) > 0 {
		consumerName = name[0]
	}
	_, ok := consumers[consumerName]
	if ok {
		return errors.New("consumer exists")
	}
	consumers[consumerName] = c
	return nil
}

type consumerGroupHandler struct {
	consumer *Consumer
}

func (consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case msg, ok := <-claim.Messages():
			if !ok {
				return nil
			}
			e := &Event{}
			if err := json.Unmarshal(msg.Value, e); err != nil {
				return err
			}
			h.consumer.handlerLock.RLock()
			handler, ok := h.consumer.handlers[e.Type]
			h.consumer.handlerLock.RUnlock()

			if ok {
				if err := handler(e); err != nil {
					return err
				}
			}
			sess.MarkMessage(msg, gMsgDone)
		}
	}
}

func (c *Consumer) HandleError(e HandleErrorFunc) {
	c.e = e
}

func (c *Consumer) HandleSucceed(t string, cc HandleConsumerFunc) {
	c.handlerLock.Lock()
	defer c.handlerLock.Unlock()
	c.handlers[t] = cc
}

func (c *Consumer) Run() {
	go c.Handle()
	for i := 0; i < c.cfg.Workers; i++ {
		go c.Worker()
	}
}

func (c *Consumer) Close() error {
	c.cancel()
	c.wg.Wait()
	return c.group.Close()
}

func (c *Consumer) Worker() {
	c.wg.Add(1)
	defer c.wg.Done()
	topics := []string{c.cfg.Topic}
	handler := consumerGroupHandler{c}
	if err := c.group.Consume(c.ctx, topics, handler); err != nil {
		if c.e != nil {
			c.e(err)
		}
	}
}

func (c *Consumer) Handle() {
	for {
		select {
		case err, ok := <-c.group.Errors():
			{
				if !ok {
					return
				}
				if c.e != nil {
					c.e(err)
				}
			}
		}
	}

}
