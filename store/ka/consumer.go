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

func NewConsumer(cfg *ConsumerConfig) error {
	return NewConsumerWithConfig(cfg, nil)
}

func NewConsumerWithInterceptor(cfg1 *ConsumerConfig, interceptor sarama.ConsumerInterceptor) error {
	if !cfg1.Check() {
		return errors.New("config error")
	}
	cfg2 := sarama.NewConfig()
	cfg2.Consumer.Return.Errors = true
	cfg2.Version = version
	if cfg1.Oldest {
		cfg2.Consumer.Offsets.Initial = sarama.OffsetOldest
	} else {
		cfg2.Consumer.Offsets.Initial = sarama.OffsetNewest
	}
	cfg2.Consumer.Interceptors = []sarama.ConsumerInterceptor{interceptor}
	return NewConsumerWithConfig(cfg1, cfg2)
}

func NewConsumerWithConfig(cfg1 *ConsumerConfig, cfg2 *sarama.Config) error {
	if !cfg1.Check() {
		return errors.New("config error")
	}
	if cfg2 == nil {
		cfg2 = sarama.NewConfig()
		cfg2.Consumer.Return.Errors = true
		cfg2.Version = version
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
	c.ctx, c.cancel = context.WithCancel(context.Background())
	consumerName := cfg1.Topic
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

func (h consumerGroupHandler) getTopic() string {
	return h.consumer.GetTopic()
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
				sess.MarkMessage(msg, gMsgDone)
				return err
			}
			if handler := getConsumerHandler(h.getTopic(), e.Type); handler != nil {
				if err := handler(e); err != nil {
					return nil
				}
			}
			sess.MarkMessage(msg, gMsgDone)
		}
	}
}

func (c *Consumer) HandleError(e HandleErrorFunc) {
	c.e = e
}

func (c *Consumer) GetTopic() string {
	return c.cfg.Topic
}

func (c *Consumer) HandleSucceed(eventType string, consumerFunc HandleConsumerFunc) {
	setConsumerHandler(c.GetTopic(), eventType, consumerFunc)
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
	for {
		if err := c.group.Consume(c.ctx, topics, handler); err != nil {
			if c.e != nil {
				c.e(err)
			}
		}

		if c.ctx.Err() != nil {
			return
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
